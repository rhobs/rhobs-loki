package push

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"

	"github.com/grafana/loki/v3/pkg/logproto"
	util_log "github.com/grafana/loki/v3/pkg/util/log"
)

// benchScale describes the number of streams and log lines used in each sub-benchmark.
type benchScale struct {
	streams        int
	linesPerStream int
}

var benchScales = []benchScale{
	{streams: 1, linesPerStream: 10},
	{streams: 10, linesPerStream: 100},
	{streams: 50, linesPerStream: 200},
}

const benchLogLine = "log line %d from service %d: something interesting happened here for context"

// --- Loki native push body builders ---

func buildLokiJSONBody(s benchScale) []byte {
	now := time.Now()
	var sb strings.Builder
	sb.WriteString(`{"streams":[`)
	for i := 0; i < s.streams; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, `{"stream":{"service_name":"svc-%d"},"values":[`, i)
		for j := 0; j < s.linesPerStream; j++ {
			if j > 0 {
				sb.WriteByte(',')
			}
			ts := now.Add(time.Duration(j) * time.Second).UnixNano()
			fmt.Fprintf(&sb, `["%d",%q]`, ts, fmt.Sprintf(benchLogLine, j, i))
		}
		sb.WriteString(`]}`)
	}
	sb.WriteString(`]}`)
	return []byte(sb.String())
}

func buildLokiProtoBody(s benchScale) []byte {
	now := time.Now()
	req := &logproto.PushRequest{
		Streams: make([]logproto.Stream, s.streams),
	}
	for i := 0; i < s.streams; i++ {
		entries := make([]logproto.Entry, s.linesPerStream)
		for j := 0; j < s.linesPerStream; j++ {
			entries[j] = logproto.Entry{
				Timestamp: now.Add(time.Duration(j) * time.Second),
				Line:      fmt.Sprintf(benchLogLine, j, i),
			}
		}
		req.Streams[i] = logproto.Stream{
			Labels:  fmt.Sprintf(`{service_name="svc-%d"}`, i),
			Entries: entries,
		}
	}
	raw, err := proto.Marshal(req)
	if err != nil {
		panic(err)
	}
	return snappy.Encode(nil, raw)
}

// --- OTLP body builders ---

func buildOTLPLogs(s benchScale) plog.Logs {
	now := time.Now()
	ld := plog.NewLogs()
	for i := 0; i < s.streams; i++ {
		rl := ld.ResourceLogs().AppendEmpty()
		rl.Resource().Attributes().PutStr("service.name", fmt.Sprintf("svc-%d", i))
		sl := rl.ScopeLogs().AppendEmpty()
		for j := 0; j < s.linesPerStream; j++ {
			lr := sl.LogRecords().AppendEmpty()
			lr.Body().SetStr(fmt.Sprintf(benchLogLine, j, i))
			lr.SetTimestamp(pcommon.Timestamp(now.Add(time.Duration(j) * time.Second).UnixNano()))
		}
	}
	return ld
}

func gzipBytes(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func buildOTLPProtoBody(s benchScale) []byte {
	exportReq := plogotlp.NewExportRequestFromLogs(buildOTLPLogs(s))
	data, err := exportReq.MarshalProto()
	if err != nil {
		panic(err)
	}
	return gzipBytes(data)
}

func buildOTLPJSONBody(s benchScale) []byte {
	exportReq := plogotlp.NewExportRequestFromLogs(buildOTLPLogs(s))
	data, err := exportReq.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return gzipBytes(data)
}

// --- Structured-metadata body builders ---
//
// Both paths produce the same output per entry:
//   stream label:        service_name="svc-X"
//   structured metadata: host.name="host-X", cloud.provider="aws", process.pid="1234"
//
// For Loki native push the metadata is carried explicitly on each entry.
// For OTLP the same three keys arrive as resource attributes; because none of them
// appear in the default index-label list they are all promoted to structured metadata
// and appended to every entry in that resource group — identical end result.

func buildLokiJSONBodyWithSM(s benchScale) []byte {
	now := time.Now()
	var sb strings.Builder
	sb.WriteString(`{"streams":[`)
	for i := 0; i < s.streams; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, `{"stream":{"service_name":"svc-%d"},"values":[`, i)
		for j := 0; j < s.linesPerStream; j++ {
			if j > 0 {
				sb.WriteByte(',')
			}
			ts := now.Add(time.Duration(j) * time.Second).UnixNano()
			smVal := fmt.Sprintf(`{"host.name":"host-%d","cloud.provider":"aws","process.pid":"1234"}`, i)
			fmt.Fprintf(&sb, `["%d",%q,%s]`, ts, fmt.Sprintf(benchLogLine, j, i), smVal)
		}
		sb.WriteString(`]}`)
	}
	sb.WriteString(`]}`)
	return []byte(sb.String())
}

func buildLokiProtoBodyWithSM(s benchScale) []byte {
	now := time.Now()
	req := &logproto.PushRequest{
		Streams: make([]logproto.Stream, s.streams),
	}
	for i := 0; i < s.streams; i++ {
		entries := make([]logproto.Entry, s.linesPerStream)
		for j := 0; j < s.linesPerStream; j++ {
			entries[j] = logproto.Entry{
				Timestamp: now.Add(time.Duration(j) * time.Second),
				Line:      fmt.Sprintf(benchLogLine, j, i),
				StructuredMetadata: []logproto.LabelAdapter{
					{Name: "host.name", Value: fmt.Sprintf("host-%d", i)},
					{Name: "cloud.provider", Value: "aws"},
					{Name: "process.pid", Value: "1234"},
				},
			}
		}
		req.Streams[i] = logproto.Stream{
			Labels:  fmt.Sprintf(`{service_name="svc-%d"}`, i),
			Entries: entries,
		}
	}
	raw, err := proto.Marshal(req)
	if err != nil {
		panic(err)
	}
	return snappy.Encode(nil, raw)
}

func buildOTLPLogsWithSM(s benchScale) plog.Logs {
	now := time.Now()
	ld := plog.NewLogs()
	for i := 0; i < s.streams; i++ {
		rl := ld.ResourceLogs().AppendEmpty()
		attrs := rl.Resource().Attributes()
		attrs.PutStr("service.name", fmt.Sprintf("svc-%d", i))
		attrs.PutStr("host.name", fmt.Sprintf("host-%d", i))
		attrs.PutStr("cloud.provider", "aws")
		attrs.PutStr("process.pid", "1234")
		sl := rl.ScopeLogs().AppendEmpty()
		for j := 0; j < s.linesPerStream; j++ {
			lr := sl.LogRecords().AppendEmpty()
			lr.Body().SetStr(fmt.Sprintf(benchLogLine, j, i))
			lr.SetTimestamp(pcommon.Timestamp(now.Add(time.Duration(j) * time.Second).UnixNano()))
		}
	}
	return ld
}

func buildOTLPProtoBodyWithSM(s benchScale) []byte {
	exportReq := plogotlp.NewExportRequestFromLogs(buildOTLPLogsWithSM(s))
	data, err := exportReq.MarshalProto()
	if err != nil {
		panic(err)
	}
	return gzipBytes(data)
}

func buildOTLPJSONBodyWithSM(s benchScale) []byte {
	exportReq := plogotlp.NewExportRequestFromLogs(buildOTLPLogsWithSM(s))
	data, err := exportReq.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return gzipBytes(data)
}

// BenchmarkIngestion is structured for benchstat comparison.
//
// Run and compare with:
//
//	go test ./pkg/loghttp/push/ -run='^$' -bench='^BenchmarkIngestion$' \
//	    -count=6 | tee results.txt
//	benchstat -col /push results.txt
//
// benchstat uses the key=value sub-benchmark names to pair results: /push selects
// the push dimension as the column axis so loki and otlp appear side-by-side for
// every (format, sm, streams, lines) combination.
func BenchmarkIngestion(b *testing.B) {
	limits := &fakeLimits{}
	resolver := newMockStreamResolver("tenant1", limits)

	type benchCase struct {
		format  string
		sm      bool
		push    string
		path    string
		ct      string
		ce      string
		buildFn func(benchScale) []byte
		parseFn RequestParser
	}

	cases := []benchCase{
		{"json", false, "loki", "/loki/api/v1/push", "application/json", "", buildLokiJSONBody, ParseLokiRequest},
		{"json", false, "otlp", "/otlp/v1/logs", "application/json", "gzip", buildOTLPJSONBody, ParseOTLPRequest},
		{"proto", false, "loki", "/loki/api/v1/push", "application/x-protobuf", "", buildLokiProtoBody, ParseLokiRequest},
		{"proto", false, "otlp", "/otlp/v1/logs", "application/x-protobuf", "gzip", buildOTLPProtoBody, ParseOTLPRequest},
		{"json", true, "loki", "/loki/api/v1/push", "application/json", "", buildLokiJSONBodyWithSM, ParseLokiRequest},
		{"json", true, "otlp", "/otlp/v1/logs", "application/json", "gzip", buildOTLPJSONBodyWithSM, ParseOTLPRequest},
		{"proto", true, "loki", "/loki/api/v1/push", "application/x-protobuf", "", buildLokiProtoBodyWithSM, ParseLokiRequest},
		{"proto", true, "otlp", "/otlp/v1/logs", "application/x-protobuf", "gzip", buildOTLPProtoBodyWithSM, ParseOTLPRequest},
	}

	for _, tc := range cases {
		smVal := "false"
		if tc.sm {
			smVal = "true"
		}
		// Key=value sub-benchmark names let benchstat treat each key as a
		// filterable/groupable dimension. /push becomes the column axis.
		prefix := fmt.Sprintf("format=%s/sm=%s/push=%s", tc.format, smVal, tc.push)
		b.Run(prefix, func(b *testing.B) {
			for _, s := range benchScales {
				body := tc.buildFn(s)
				b.Run(fmt.Sprintf("streams=%d/lines=%d", s.streams, s.linesPerStream), func(b *testing.B) {
					// Build the request template once: URL parsing, header map
					// allocation, and field setup all happen here outside the
					// timed region.
					tmpl := httptest.NewRequest("POST", tc.path, http.NoBody)
					tmpl.Header.Set("Content-Type", tc.ct)
					if tc.ce != "" {
						tmpl.Header.Set("Content-Encoding", tc.ce)
					}
					tmpl.ContentLength = int64(len(body))

					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						// Shallow-copy the template and attach a fresh reader.
						// bytes.NewReader is O(1): it wraps the pre-allocated
						// body slice without copying it.
						req := *tmpl
						req.Body = io.NopCloser(bytes.NewReader(body))
						_, _, err := tc.parseFn("tenant1", &req, limits, nil, 100<<20, 100<<20, nil, resolver, util_log.Logger)
						if err != nil {
							b.Fatal(err)
						}
					}
				})
			}
		})
	}
}

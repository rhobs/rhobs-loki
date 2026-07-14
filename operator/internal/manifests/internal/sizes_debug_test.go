package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	lokiv1 "github.com/grafana/loki/operator/api/loki/v1"
)

func TestResourceRequirementsForSize_WithOverrides(t *testing.T) {
	tests := []struct {
		name         string
		getResources func() ComponentResources
		validateFunc func(t *testing.T, baseResources ComponentResources, resources ComponentResources)
	}{
		{
			name: "no overrides - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, nil)
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Should match base t-shirt size exactly
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
				assert.Equal(t, baseResources.Distributor.Requests, resources.Distributor.Requests)
			},
		},
		{
			name: "no overrides - requests and limits set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, true, nil)
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Should match base t-shirt size
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
				// Limits should equal requests
				assert.Equal(t, resources.Ingester.Requests, resources.Ingester.Limits)
				assert.Equal(t, resources.Distributor.Requests, resources.Distributor.Limits)
			},
		},
		{
			name: "partial override - memory request only - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("8Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Memory request should be overridden
				memory := resources.Ingester.Requests[corev1.ResourceMemory]
				assert.Equal(t, "8Gi", memory.String())
				// CPU request should remain from base
				assert.Equal(t, baseResources.Ingester.Requests[corev1.ResourceCPU], resources.Ingester.Requests[corev1.ResourceCPU])
				// Limits should remain from base (nil)
				assert.Equal(t, baseResources.Ingester.Limits, resources.Ingester.Limits)
			},
		},
		{
			name: "partial override - cpu request only - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("500m"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// CPU request should be overridden
				cpu := resources.Ingester.Requests[corev1.ResourceCPU]
				assert.Equal(t, "500m", cpu.String())
				// Memory request should remain from base
				assert.Equal(t, baseResources.Ingester.Requests[corev1.ResourceMemory], resources.Ingester.Requests[corev1.ResourceMemory])
				// Limits should remain from base (nil)
				assert.Equal(t, baseResources.Ingester.Limits, resources.Ingester.Limits)
			},
		},
		{
			name: "partial override - cpu limit only - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("4"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Requests should remain from base
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
				// CPU limit should be overridden
				cpuLimit := resources.Ingester.Limits[corev1.ResourceCPU]
				assert.Equal(t, "4", cpuLimit.String())
				// Memory limit should NOT be set
				_, hasMemoryLimit := resources.Ingester.Limits[corev1.ResourceMemory]
				assert.False(t, hasMemoryLimit)
			},
		},
		{
			name: "partial override - requests only - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1"),
								corev1.ResourceMemory: resource.MustParse("4Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Requests should be overridden
				cpu := resources.Ingester.Requests[corev1.ResourceCPU]
				memory := resources.Ingester.Requests[corev1.ResourceMemory]
				assert.Equal(t, "1", cpu.String())
				assert.Equal(t, "4Gi", memory.String())
				// Limits should remain from base (nil for 1x.small ingester)
				assert.Equal(t, baseResources.Ingester.Limits, resources.Ingester.Limits)
			},
		},
		{
			name: "partial override - requests only - requests and limits set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, true, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1"),
								corev1.ResourceMemory: resource.MustParse("4Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Requests should be overridden
				cpu := resources.Ingester.Requests[corev1.ResourceCPU]
				memory := resources.Ingester.Requests[corev1.ResourceMemory]
				assert.Equal(t, "1", cpu.String())
				assert.Equal(t, "4Gi", memory.String())
				// Limits should be set from original base requests (useRequestsAsLimits applied first, then override)
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Limits)
			},
		},
		{
			name: "partial override - limits only - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("55Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Requests should be preserved from base t-shirt size
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
				// Limits should be overridden as specified
				memoryLimit := resources.Ingester.Limits[corev1.ResourceMemory]
				assert.Equal(t, "55Gi", memoryLimit.String())
				// CPU limit should NOT be set (only memory was specified)
				_, hasCPULimit := resources.Ingester.Limits[corev1.ResourceCPU]
				assert.False(t, hasCPULimit)
			},
		},
		{
			name: "partial override - limits only - requests and limits set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, true, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Ingester: &corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("55Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Requests should be preserved from base t-shirt size
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
				// Memory limit should be overridden as specified
				memoryLimit := resources.Ingester.Limits[corev1.ResourceMemory]
				assert.Equal(t, "55Gi", memoryLimit.String())
				// CPU limit should be set from requests (useRequestsAsLimits applies to non-overridden resources)
				cpuLimit := resources.Ingester.Limits[corev1.ResourceCPU]
				assert.Equal(t, resources.Ingester.Requests[corev1.ResourceCPU], cpuLimit)
			},
		},
		{
			name: "full override - both requests and limits - requests only set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Distributor: &corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("200m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("1Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				// Both requests and limits should be overridden
				reqCPU := resources.Distributor.Requests[corev1.ResourceCPU]
				reqMemory := resources.Distributor.Requests[corev1.ResourceMemory]
				limitCPU := resources.Distributor.Limits[corev1.ResourceCPU]
				limitMemory := resources.Distributor.Limits[corev1.ResourceMemory]
				assert.Equal(t, "200m", reqCPU.String())
				assert.Equal(t, "512Mi", reqMemory.String())
				assert.Equal(t, "500m", limitCPU.String())
				assert.Equal(t, "1Gi", limitMemory.String())
				// Non-overridden components should remain unchanged
				assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
			},
		},
		{
			name: "full override - both requests and limits - requests and limits set",
			getResources: func() ComponentResources {
				return ResourceRequirementsForSize(lokiv1.SizeOneXSmall, true, &lokiv1.DebugOptionsSpec{
					ResourceOverrides: &lokiv1.ComponentResourceOverrides{
						Distributor: &corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("200m"),
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("1Gi"),
							},
						},
					},
				})
			},
			validateFunc: func(t *testing.T, baseResources ComponentResources, resources ComponentResources) {
				reqCPU := resources.Distributor.Requests[corev1.ResourceCPU]
				reqMemory := resources.Distributor.Requests[corev1.ResourceMemory]
				limitCPU := resources.Distributor.Limits[corev1.ResourceCPU]
				limitMemory := resources.Distributor.Limits[corev1.ResourceMemory]
				assert.Equal(t, "200m", reqCPU.String())
				assert.Equal(t, "512Mi", reqMemory.String())
				assert.Equal(t, "500m", limitCPU.String())
				assert.Equal(t, "1Gi", limitMemory.String())
				assert.Equal(t, resources.Ingester.Requests, resources.Ingester.Limits)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseResources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, nil)
			resources := tt.getResources()
			tt.validateFunc(t, baseResources, resources)
		})
	}
}

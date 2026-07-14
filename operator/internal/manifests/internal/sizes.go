package internal

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	lokiv1 "github.com/grafana/loki/operator/api/loki/v1"
)

// ComponentResources is a map of component->requests/limits
type ComponentResources struct {
	IndexGateway ResourceRequirements
	Ingester     ResourceRequirements
	Compactor    ResourceRequirements
	Ruler        ResourceRequirements
	WALStorage   ResourceRequirements
	// these two don't need a PVCSize
	Querier       corev1.ResourceRequirements
	Distributor   corev1.ResourceRequirements
	QueryFrontend corev1.ResourceRequirements
	Gateway       corev1.ResourceRequirements
}

func (c ComponentResources) DeepCopy() ComponentResources {
	return ComponentResources{
		IndexGateway:  *c.IndexGateway.DeepCopy(),
		Ingester:      *c.Ingester.DeepCopy(),
		Compactor:     *c.Compactor.DeepCopy(),
		Ruler:         *c.Ruler.DeepCopy(),
		WALStorage:    *c.WALStorage.DeepCopy(),
		Querier:       *c.Querier.DeepCopy(),
		Distributor:   *c.Distributor.DeepCopy(),
		QueryFrontend: *c.QueryFrontend.DeepCopy(),
		Gateway:       *c.Gateway.DeepCopy(),
	}
}

// ResourceRequirements sets CPU, Memory, and PVC requirements for a component
type ResourceRequirements struct {
	Limits   corev1.ResourceList
	Requests corev1.ResourceList
	PVCSize  resource.Quantity
}

func (r *ResourceRequirements) DeepCopy() *ResourceRequirements {
	return &ResourceRequirements{
		Limits:   r.Limits.DeepCopy(),
		Requests: r.Requests.DeepCopy(),
		PVCSize:  r.PVCSize.DeepCopy(),
	}
}

// resourceRequirementsTable defines the default resource requests and limits for each size
var resourceRequirementsTable = map[lokiv1.LokiStackSizeType]ComponentResources{
	lokiv1.SizeOneXDemo: {
		Ruler: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
		},
		Ingester: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
		},
		Compactor: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
		},
		IndexGateway: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
		},
		WALStorage: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
		},
	},
	lokiv1.SizeOneXPico: {
		Querier: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("750m"),
				corev1.ResourceMemory: resource.MustParse("1.5Gi"),
			},
		},
		Ruler: ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
			PVCSize: resource.MustParse("10Gi"),
		},
		Ingester: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("3Gi"),
			},
		},
		Distributor: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
			},
		},
		QueryFrontend: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
			},
		},
		Compactor: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
			},
		},
		Gateway: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
			},
		},
		IndexGateway: ResourceRequirements{
			PVCSize: resource.MustParse("50Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("150m"),
				corev1.ResourceMemory: resource.MustParse("250Mi"),
			},
		},
		WALStorage: ResourceRequirements{
			PVCSize: resource.MustParse("150Gi"),
		},
	},
	lokiv1.SizeOneXExtraSmall: {
		Querier: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1.5"),
				corev1.ResourceMemory: resource.MustParse("3Gi"),
			},
		},
		Ruler: ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
			PVCSize: resource.MustParse("10Gi"),
		},
		Ingester: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("5.25Gi"),
			},
		},
		Distributor: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
		QueryFrontend: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
		Compactor: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		Gateway: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
			},
		},
		IndexGateway: ResourceRequirements{
			PVCSize: resource.MustParse("50Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
		WALStorage: ResourceRequirements{
			PVCSize: resource.MustParse("150Gi"),
		},
	},
	lokiv1.SizeOneXSmall: {
		Querier: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("4"),
				corev1.ResourceMemory: resource.MustParse("4Gi"),
			},
		},
		Ruler: ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("4"),
				corev1.ResourceMemory: resource.MustParse("8Gi"),
			},
			PVCSize: resource.MustParse("10Gi"),
		},
		Ingester: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("2.5"),
				corev1.ResourceMemory: resource.MustParse("13.25Gi"),
			},
		},
		Distributor: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		QueryFrontend: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("4"),
				corev1.ResourceMemory: resource.MustParse("2.5Gi"),
			},
		},
		Compactor: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("4Gi"),
			},
		},
		Gateway: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
		IndexGateway: ResourceRequirements{
			PVCSize: resource.MustParse("50Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		WALStorage: ResourceRequirements{
			PVCSize: resource.MustParse("150Gi"),
		},
	},
	lokiv1.SizeOneXMedium: {
		Querier: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("6"),
				corev1.ResourceMemory: resource.MustParse("10Gi"),
			},
		},
		Ruler: ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("8"),
				corev1.ResourceMemory: resource.MustParse("16Gi"),
			},
			PVCSize: resource.MustParse("10Gi"),
		},
		Ingester: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("6"),
				corev1.ResourceMemory: resource.MustParse("30Gi"),
			},
		},
		Distributor: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		QueryFrontend: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("4"),
				corev1.ResourceMemory: resource.MustParse("2.5Gi"),
			},
		},
		Compactor: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("4Gi"),
			},
		},
		Gateway: corev1.ResourceRequirements{
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
		IndexGateway: ResourceRequirements{
			PVCSize: resource.MustParse("50Gi"),
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		WALStorage: ResourceRequirements{
			PVCSize: resource.MustParse("150Gi"),
		},
	},
}

// ResourceRequirementsForSize returns the resource configuration for a specific LokiStack size.
// If debugOptions contains resource overrides, they will be applied on top of the base t-shirt size configuration.
// Only CPU and memory resources are overridden; PVC storage sizes remain from the base t-shirt size.
func ResourceRequirementsForSize(size lokiv1.LokiStackSizeType, useRequestsAsLimits bool, debugOptions *lokiv1.DebugOptionsSpec) ComponentResources {
	resources := resourceRequirementsTable[size].DeepCopy()
	if useRequestsAsLimits {
		resources.IndexGateway.Limits = resources.IndexGateway.Requests.DeepCopy()
		resources.Ingester.Limits = resources.Ingester.Requests.DeepCopy()
		resources.Compactor.Limits = resources.Compactor.Requests.DeepCopy()
		resources.Ruler.Limits = resources.Ruler.Requests.DeepCopy()
		resources.WALStorage.Limits = resources.WALStorage.Requests.DeepCopy()
		resources.Querier.Limits = resources.Querier.Requests.DeepCopy()
		resources.Distributor.Limits = resources.Distributor.Requests.DeepCopy()
		resources.QueryFrontend.Limits = resources.QueryFrontend.Requests.DeepCopy()
		resources.Gateway.Limits = resources.Gateway.Requests.DeepCopy()
	}

	// Apply resource overrides from debug options if specified
	if debugOptions != nil && debugOptions.ResourceOverrides != nil {
		resources = applyResourceOverrides(resources, debugOptions.ResourceOverrides)
	}

	return resources
}

// applyResourceOverrides applies per-component resource overrides on top of base resources.
// Only the specific resource fields (requests/limits) that are provided in overrides are changed.
// Unspecified fields preserve their t-shirt size values.
func applyResourceOverrides(base ComponentResources, overrides *lokiv1.ComponentResourceOverrides) ComponentResources {
	result := base.DeepCopy()

	// Apply overrides for components with PVC storage (using our custom ResourceRequirements type)
	// For these components, we selectively override only the specified requests/limits
	if overrides.IndexGateway != nil {
		applyResourceRequirementsOverride(&result.IndexGateway.Requests, &result.IndexGateway.Limits, overrides.IndexGateway)
	}
	if overrides.Ingester != nil {
		applyResourceRequirementsOverride(&result.Ingester.Requests, &result.Ingester.Limits, overrides.Ingester)
	}
	if overrides.Compactor != nil {
		applyResourceRequirementsOverride(&result.Compactor.Requests, &result.Compactor.Limits, overrides.Compactor)
	}
	if overrides.Ruler != nil {
		applyResourceRequirementsOverride(&result.Ruler.Requests, &result.Ruler.Limits, overrides.Ruler)
	}
	if overrides.Querier != nil {
		applyResourceRequirementsOverride(&result.Querier.Requests, &result.Querier.Limits, overrides.Querier)
	}
	if overrides.Distributor != nil {
		applyResourceRequirementsOverride(&result.Distributor.Requests, &result.Distributor.Limits, overrides.Distributor)
	}
	if overrides.QueryFrontend != nil {
		applyResourceRequirementsOverride(&result.QueryFrontend.Requests, &result.QueryFrontend.Limits, overrides.QueryFrontend)
	}

	return result
}

// applyResourceRequirementsOverride selectively applies resource overrides, preserving
// existing values for any requests/limits not specified in the override.
func applyResourceRequirementsOverride(baseRequests, baseLimits *corev1.ResourceList, override *corev1.ResourceRequirements) {
	// Only override requests if they are specified in the override
	if override.Requests != nil {
		if *baseRequests == nil {
			*baseRequests = make(corev1.ResourceList)
		}
		for resource, quantity := range override.Requests {
			(*baseRequests)[resource] = quantity.DeepCopy()
		}
	}

	// Only override limits if they are specified in the override
	if override.Limits != nil {
		if *baseLimits == nil {
			*baseLimits = make(corev1.ResourceList)
		}
		for resource, quantity := range override.Limits {
			(*baseLimits)[resource] = quantity.DeepCopy()
		}
	}
}

// StackSizeTable defines the default configurations for each size
var StackSizeTable = map[lokiv1.LokiStackSizeType]lokiv1.LokiStackSpec{
	lokiv1.SizeOneXDemo: {
		Size: lokiv1.SizeOneXDemo,
		Replication: &lokiv1.ReplicationSpec{
			Factor: 1,
		},
		Limits: &lokiv1.LimitsSpec{
			Global: &lokiv1.LimitsTemplateSpec{
				IngestionLimits: &lokiv1.IngestionLimitSpec{
					// Defaults from Loki docs
					IngestionRate:           4,
					IngestionBurstSize:      6,
					MaxLabelNameLength:      1024,
					MaxLabelValueLength:     2048,
					MaxLabelNamesPerSeries:  30,
					MaxLineSize:             256000,
					PerStreamDesiredRate:    3,
					PerStreamRateLimit:      5,
					PerStreamRateLimitBurst: 15,
				},
				QueryLimits: &lokiv1.QueryLimitSpec{
					// Defaults from Loki docs
					MaxEntriesLimitPerQuery: 5000,
					MaxChunksPerQuery:       2000000,
					MaxQuerySeries:          500,
					QueryTimeout:            "3m",
					CardinalityLimit:        100000,
					MaxVolumeSeries:         1000,
				},
			},
		},
		Template: &lokiv1.LokiTemplateSpec{
			Compactor: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Distributor: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Ingester: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Querier: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			QueryFrontend: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Gateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			IndexGateway: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Ruler: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
		},
	},

	lokiv1.SizeOneXPico: {
		Size: lokiv1.SizeOneXPico,
		Replication: &lokiv1.ReplicationSpec{
			Factor: 2,
		},
		Limits: &lokiv1.LimitsSpec{
			Global: &lokiv1.LimitsTemplateSpec{
				IngestionLimits: &lokiv1.IngestionLimitSpec{
					// Defaults from Loki docs
					IngestionRate:             4,
					IngestionBurstSize:        6,
					MaxGlobalStreamsPerTenant: 10000,
					MaxLabelNameLength:        1024,
					MaxLabelValueLength:       2048,
					MaxLabelNamesPerSeries:    30,
					MaxLineSize:               256000,
					PerStreamDesiredRate:      3,
					PerStreamRateLimit:        5,
					PerStreamRateLimitBurst:   15,
				},
				QueryLimits: &lokiv1.QueryLimitSpec{
					// Defaults from Loki docs
					MaxEntriesLimitPerQuery: 5000,
					MaxChunksPerQuery:       2000000,
					MaxQuerySeries:          500,
					QueryTimeout:            "3m",
					CardinalityLimit:        100000,
					MaxVolumeSeries:         1000,
				},
			},
		},
		Template: &lokiv1.LokiTemplateSpec{
			Compactor: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Distributor: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ingester: &lokiv1.LokiComponentSpec{
				Replicas: 3,
			},
			Querier: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			QueryFrontend: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Gateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			IndexGateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ruler: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
		},
	},

	lokiv1.SizeOneXExtraSmall: {
		Size: lokiv1.SizeOneXExtraSmall,
		Replication: &lokiv1.ReplicationSpec{
			Factor: 2,
		},
		Limits: &lokiv1.LimitsSpec{
			Global: &lokiv1.LimitsTemplateSpec{
				IngestionLimits: &lokiv1.IngestionLimitSpec{
					// Defaults from Loki docs
					IngestionRate:             4,
					IngestionBurstSize:        6,
					MaxGlobalStreamsPerTenant: 10000,
					MaxLabelNameLength:        1024,
					MaxLabelValueLength:       2048,
					MaxLabelNamesPerSeries:    30,
					MaxLineSize:               256000,
					PerStreamDesiredRate:      3,
					PerStreamRateLimit:        5,
					PerStreamRateLimitBurst:   15,
				},
				QueryLimits: &lokiv1.QueryLimitSpec{
					// Defaults from Loki docs
					MaxEntriesLimitPerQuery: 5000,
					MaxChunksPerQuery:       2000000,
					MaxQuerySeries:          500,
					QueryTimeout:            "3m",
					CardinalityLimit:        100000,
					MaxVolumeSeries:         1000,
				},
			},
		},
		Template: &lokiv1.LokiTemplateSpec{
			Compactor: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Distributor: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ingester: &lokiv1.LokiComponentSpec{
				Replicas: 3,
			},
			Querier: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			QueryFrontend: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Gateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			IndexGateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ruler: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
		},
	},

	lokiv1.SizeOneXSmall: {
		Size: lokiv1.SizeOneXSmall,
		Replication: &lokiv1.ReplicationSpec{
			Factor: 2,
		},
		Limits: &lokiv1.LimitsSpec{
			Global: &lokiv1.LimitsTemplateSpec{
				IngestionLimits: &lokiv1.IngestionLimitSpec{
					// Custom for 1x.small
					IngestionRate:             15,
					IngestionBurstSize:        20,
					MaxGlobalStreamsPerTenant: 10000,
					// Defaults from Loki docs
					MaxLabelNameLength:      1024,
					MaxLabelValueLength:     2048,
					MaxLabelNamesPerSeries:  30,
					MaxLineSize:             256000,
					PerStreamDesiredRate:    3,
					PerStreamRateLimit:      5,
					PerStreamRateLimitBurst: 15,
				},
				QueryLimits: &lokiv1.QueryLimitSpec{
					// Defaults from Loki docs
					MaxEntriesLimitPerQuery: 5000,
					MaxChunksPerQuery:       2000000,
					MaxQuerySeries:          500,
					QueryTimeout:            "3m",
					CardinalityLimit:        100000,
					MaxVolumeSeries:         1000,
				},
			},
		},
		Template: &lokiv1.LokiTemplateSpec{
			Compactor: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Distributor: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ingester: &lokiv1.LokiComponentSpec{
				Replicas: 3,
			},
			Querier: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			QueryFrontend: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Gateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			IndexGateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ruler: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
		},
	},

	lokiv1.SizeOneXMedium: {
		Size: lokiv1.SizeOneXMedium,
		Replication: &lokiv1.ReplicationSpec{
			Factor: 2,
		},
		Limits: &lokiv1.LimitsSpec{
			Global: &lokiv1.LimitsTemplateSpec{
				IngestionLimits: &lokiv1.IngestionLimitSpec{
					// Custom for 1x.medium
					IngestionRate:             50,
					IngestionBurstSize:        20,
					MaxGlobalStreamsPerTenant: 25000,
					// Defaults from Loki docs
					MaxLabelNameLength:      1024,
					MaxLabelValueLength:     2048,
					MaxLabelNamesPerSeries:  30,
					MaxLineSize:             256000,
					PerStreamDesiredRate:    3,
					PerStreamRateLimit:      5,
					PerStreamRateLimitBurst: 15,
				},
				QueryLimits: &lokiv1.QueryLimitSpec{
					// Defaults from Loki docs
					MaxEntriesLimitPerQuery: 5000,
					MaxChunksPerQuery:       2000000,
					MaxQuerySeries:          500,
					QueryTimeout:            "3m",
					CardinalityLimit:        100000,
					MaxVolumeSeries:         1000,
				},
			},
		},
		Template: &lokiv1.LokiTemplateSpec{
			Compactor: &lokiv1.LokiComponentSpec{
				Replicas: 1,
			},
			Distributor: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ingester: &lokiv1.LokiComponentSpec{
				Replicas: 3,
			},
			Querier: &lokiv1.LokiComponentSpec{
				Replicas: 3,
			},
			QueryFrontend: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Gateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			IndexGateway: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
			Ruler: &lokiv1.LokiComponentSpec{
				Replicas: 2,
			},
		},
	},
}

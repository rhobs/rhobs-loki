package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	lokiv1 "github.com/grafana/loki/operator/api/loki/v1"
)

func TestResourceRequirementsForSize_WithoutOverrides(t *testing.T) {
	// Test that base t-shirt sizing works without debug options
	resources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, nil)

	// Verify we get the expected base resources for 1x.small
	assert.NotNil(t, resources.Distributor.Requests)
	assert.NotNil(t, resources.Ingester.Requests)
	assert.NotNil(t, resources.Querier.Requests)

	// Check that PVC sizes are set from t-shirt sizing
	assert.False(t, resources.Ingester.PVCSize.IsZero())
	assert.False(t, resources.IndexGateway.PVCSize.IsZero())
}

func TestResourceRequirementsForSize_WithDistributorOverride(t *testing.T) {
	// Test overriding just the distributor resources
	debugOptions := &lokiv1.DebugOptionsSpec{
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
	}

	resources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, debugOptions)

	// Verify distributor uses override values
	cpu := resources.Distributor.Requests[corev1.ResourceCPU]
	memory := resources.Distributor.Requests[corev1.ResourceMemory]
	cpuLimit := resources.Distributor.Limits[corev1.ResourceCPU]
	memoryLimit := resources.Distributor.Limits[corev1.ResourceMemory]

	assert.Equal(t, "200m", cpu.String())
	assert.Equal(t, "512Mi", memory.String())
	assert.Equal(t, "500m", cpuLimit.String())
	assert.Equal(t, "1Gi", memoryLimit.String())

	// Verify other components still use base t-shirt size values
	baseResources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, nil)
	assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
	assert.Equal(t, baseResources.Querier.Requests, resources.Querier.Requests)
}

func TestResourceRequirementsForSize_WithIngesterOverride(t *testing.T) {
	// Test overriding ingester (component with PVC storage)
	debugOptions := &lokiv1.DebugOptionsSpec{
		ResourceOverrides: &lokiv1.ComponentResourceOverrides{
			Ingester: &corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1"),
					corev1.ResourceMemory: resource.MustParse("4Gi"),
				},
			},
		},
	}

	resources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, debugOptions)

	// Verify ingester uses override CPU/memory values
	cpu := resources.Ingester.Requests[corev1.ResourceCPU]
	memory := resources.Ingester.Requests[corev1.ResourceMemory]

	assert.Equal(t, "1", cpu.String())
	assert.Equal(t, "4Gi", memory.String())

	// Verify PVC size is preserved from t-shirt sizing (not overridden)
	baseResources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, nil)
	assert.Equal(t, baseResources.Ingester.PVCSize, resources.Ingester.PVCSize)
}

func TestResourceRequirementsForSize_WithMultipleOverrides(t *testing.T) {
	// Test overriding multiple components
	debugOptions := &lokiv1.DebugOptionsSpec{
		ResourceOverrides: &lokiv1.ComponentResourceOverrides{
			Distributor: &corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("300m"),
					corev1.ResourceMemory: resource.MustParse("256Mi"),
				},
			},
			Querier: &corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("2"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				},
			},
			QueryFrontend: &corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1"),
					corev1.ResourceMemory: resource.MustParse("1Gi"),
				},
			},
		},
	}

	resources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, debugOptions)

	// Verify all overridden components use new values
	distCPU := resources.Distributor.Requests[corev1.ResourceCPU]
	distMemory := resources.Distributor.Requests[corev1.ResourceMemory]
	assert.Equal(t, "300m", distCPU.String())
	assert.Equal(t, "256Mi", distMemory.String())

	querierCPU := resources.Querier.Requests[corev1.ResourceCPU]
	querierMemory := resources.Querier.Requests[corev1.ResourceMemory]
	assert.Equal(t, "2", querierCPU.String())
	assert.Equal(t, "2Gi", querierMemory.String())

	qfCPU := resources.QueryFrontend.Requests[corev1.ResourceCPU]
	qfMemory := resources.QueryFrontend.Requests[corev1.ResourceMemory]
	assert.Equal(t, "1", qfCPU.String())
	assert.Equal(t, "1Gi", qfMemory.String())

	// Verify non-overridden components still use base values
	baseResources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, false, nil)
	assert.Equal(t, baseResources.Ingester.Requests, resources.Ingester.Requests)
	assert.Equal(t, baseResources.Compactor.Requests, resources.Compactor.Requests)
}

func TestResourceRequirementsForSize_WithUseRequestsAsLimitsAndOverrides(t *testing.T) {
	// Test that overrides work correctly when useRequestsAsLimits is enabled
	debugOptions := &lokiv1.DebugOptionsSpec{
		ResourceOverrides: &lokiv1.ComponentResourceOverrides{
			Distributor: &corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("128Mi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("200m"),
					corev1.ResourceMemory: resource.MustParse("256Mi"),
				},
			},
		},
	}

	resources := ResourceRequirementsForSize(lokiv1.SizeOneXSmall, true, debugOptions)

	// Verify distributor override values are preserved (not overwritten by useRequestsAsLimits)
	reqCPU := resources.Distributor.Requests[corev1.ResourceCPU]
	reqMemory := resources.Distributor.Requests[corev1.ResourceMemory]
	limitCPU := resources.Distributor.Limits[corev1.ResourceCPU]
	limitMemory := resources.Distributor.Limits[corev1.ResourceMemory]

	assert.Equal(t, "100m", reqCPU.String())
	assert.Equal(t, "128Mi", reqMemory.String())
	assert.Equal(t, "200m", limitCPU.String())
	assert.Equal(t, "256Mi", limitMemory.String())

	// Verify non-overridden components get limits set from requests
	assert.Equal(t, resources.Ingester.Requests, resources.Ingester.Limits)
}

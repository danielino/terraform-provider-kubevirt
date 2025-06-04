package datavolume

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/pointer"
	"testing"

	"github.com/stretchr/testify/assert"
	api "k8s.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func TestExpandDataVolumeSpec_SourceRef(t *testing.T) {
	spec := []interface{}{map[string]interface{}{
		"source_ref": []interface{}{
			map[string]interface{}{
				"name": "test-source",
				"kind": "DataSource",
			},
		},
	}}
	out, err := ExpandDataVolumeSpec(spec)
	assert.NoError(t, err)
	assert.NotNil(t, out.SourceRef)
	assert.Equal(t, "test-source", out.SourceRef.Name)
	assert.Equal(t, "DataSource", out.SourceRef.Kind)
}

func TestExpandDataVolumeSpec_Storage(t *testing.T) {
	spec := []interface{}{map[string]interface{}{
		"storage": []interface{}{
			map[string]interface{}{
				"resources": []interface{}{
					map[string]interface{}{
						"requests": map[string]interface{}{
							"storage": "10Gi",
						},
					},
				},
			},
		},
	}}
	out, err := ExpandDataVolumeSpec(spec)
	assert.NoError(t, err)
	assert.NotNil(t, out.Storage)
	quantity := out.Storage.Resources.Requests[api.ResourceStorage]
	assert.Equal(t, "10Gi", quantity.String())
}

func TestFlattenPersistentVolumeClaimSpec(t *testing.T) {
	spec := api.PersistentVolumeClaimSpec{
		AccessModes: []api.PersistentVolumeAccessMode{api.ReadWriteOnce},
		Resources: api.ResourceRequirements{
			Requests: api.ResourceList{
				api.ResourceStorage: resource.MustParse("10Gi"),
			},
		},
	}
	flattened := FlattenPersistentVolumeClaimSpec(spec)
	assert.NotNil(t, flattened)
	assert.Equal(t, "10Gi", flattened[0].(map[string]interface{})["resources"].([]interface{})[0].(map[string]interface{})["requests"].(map[string]interface{})["storage"])
}

func TestFlattenDataVolumeSpec(t *testing.T) {
	spec := cdiv1.DataVolumeSpec{
		SourceRef: &cdiv1.DataVolumeSourceRef{
			Name:      "test-source",
			Kind:      "DataSource",
			Namespace: pointer.String(""),
		},
		Storage: &cdiv1.StorageSpec{
			Resources: api.ResourceRequirements{
				Requests: api.ResourceList{
					api.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
		},
	}
	flattened := FlattenDataVolumeSpec(spec)
	assert.NotNil(t, flattened)
	assert.Equal(t, "test-source", flattened[0].(map[string]interface{})["source_ref"].([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, "10Gi", flattened[0].(map[string]interface{})["storage"].([]interface{})[0].(map[string]interface{})["resources"].([]interface{})[0].(map[string]interface{})["requests"].(map[string]interface{})["storage"])
}

package datavolume

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	"github.com/stretchr/testify/assert"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func TestExpandDataVolumeStorage(t *testing.T) {
	tests := []struct {
		name           string
		input          []interface{}
		expectedOutput *cdiv1.StorageSpec
	}{
		{
			name:           "empty input",
			input:          []interface{}{},
			expectedOutput: nil,
		},
		{
			name: "valid input with requests",
			input: []interface{}{
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
			expectedOutput: &cdiv1.StorageSpec{
				Resources: api.ResourceRequirements{
					Requests: api.ResourceList{
						api.ResourceStorage: resource.MustParse("10Gi"),
					},
				},
			},
		},
		{
			name: "invalid quantity in requests",
			input: []interface{}{
				map[string]interface{}{
					"resources": []interface{}{
						map[string]interface{}{
							"requests": map[string]interface{}{
								"storage": "invalid",
							},
						},
					},
				},
			},
			expectedOutput: nil, // Invalid quantity should result in nil
		},
		{
			name: "valid input with access_modes",
			input: []interface{}{
				map[string]interface{}{
					"access_modes": schema.NewSet(schema.HashString, []interface{}{
						"ReadWriteOnce",
					}),
				},
			},
			expectedOutput: &cdiv1.StorageSpec{
				AccessModes: []api.PersistentVolumeAccessMode{
					api.ReadWriteOnce,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := expandDataVolumeStorage(tc.input)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

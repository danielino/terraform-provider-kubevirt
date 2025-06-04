package datavolume

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func dataVolumeSourceStorageFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_modes": {
			Type:        schema.TypeSet,
			Description: "AccessModes is a list of access modes supported by the storage.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Set: schema.HashString,
		},
		"resources": {
			Type:        schema.TypeList,
			Description: "The resources required by the storage (requests and limits).",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"requests": {
						Type:        schema.TypeMap,
						Description: "The minimum resources required.",
						Optional:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataVolumeStorageSchema() *schema.Schema {
	fields := dataVolumeSourceStorageFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "DataVolumeSourceStorage defines the Storage type specification.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandDataVolumeStorage(dataVolumeStorage []interface{}) *cdiv1.StorageSpec {
	if len(dataVolumeStorage) == 0 || dataVolumeStorage[0] == nil {
		return nil
	}

	result := &cdiv1.StorageSpec{}
	in := dataVolumeStorage[0].(map[string]interface{})

	// Process access_modes
	if v, ok := in["access_modes"].(*schema.Set); ok && v.Len() > 0 {
		result.AccessModes = expandPersistentVolumeAccessModes(v.List())
	}

	// Process resources
	if v, ok := in["resources"].([]interface{}); ok && len(v) > 0 {
		if res, ok := v[0].(map[string]interface{}); ok {
			if requests, ok := res["requests"].(map[string]interface{}); ok && len(requests) > 0 {
				result.Resources.Requests = api.ResourceList{}
				for rk, rv := range requests {
					if rvs, ok := rv.(string); ok && rvs != "" {
						q, err := resource.ParseQuantity(rvs)
						if err == nil {
							result.Resources.Requests[api.ResourceName(rk)] = q
						}
					}
				}
			}
		}
	}

	// Check if result is empty
	if len(result.AccessModes) == 0 && len(result.Resources.Requests) == 0 {
		return nil
	}

	return result
}

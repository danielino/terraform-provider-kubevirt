package datavolume

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "k8s.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func dataVolumeSourceFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"http": dataVolumeSourceHTTPSchema(),
		"pvc":  dataVolumeSourcePVCSchema(),
	}
}

func dataVolumeSourceSchema() *schema.Schema {
	fields := dataVolumeSourceFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("Source is the src of the data for the requested DataVolume."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func dataVolumeSourceHTTPFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Description: "url is the URL of the http source.",
			Optional:    true,
		},
		"secret_ref": {
			Type:        schema.TypeString,
			Description: "Secret_ref provides the secret reference needed to access the HTTP source.",
			Optional:    true,
		},
		"cert_config_map": {
			Type:        schema.TypeString,
			Description: "Cert_config_map provides a reference to the Registry certs.",
			Optional:    true,
		},
	}
}

func dataVolumeSourceHTTPSchema() *schema.Schema {
	fields := dataVolumeSourceHTTPFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "DataVolumeSourceHTTP provides the parameters to create a Data Volume from an HTTP source.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func dataVolumeSourcePVCFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"namespace": {
			Type:        schema.TypeString,
			Description: "The namespace which the PVC located in.",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the PVC.",
			Optional:    true,
		},
	}
}

func dataVolumeSourcePVCSchema() *schema.Schema {
	fields := dataVolumeSourcePVCFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "DataVolumeSourcePVC provides the parameters to create a Data Volume from an existing PVC.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func dataVolumeSourceRefFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"kind": {
			Type:        schema.TypeString,
			Description: "The kind of the source reference, currently only \"DataSource\" is supported",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the source reference.",
			Optional:    true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "The namespace of the source reference, defaults to the DataVolume namespace.",
			Optional:    true,
		},
	}
}

func dataVolumeSourceRefSchema() *schema.Schema {
	fields := dataVolumeSourceRefFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "DataVolumeSourceRef defines an indirect reference to the source of data for the DataVolume.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

// Expanders

func expandDataVolumeSource(dataVolumeSource []interface{}) *cdiv1.DataVolumeSource {
	result := &cdiv1.DataVolumeSource{}

	if len(dataVolumeSource) == 0 || dataVolumeSource[0] == nil {
		return result
	}

	in := dataVolumeSource[0].(map[string]interface{})

	result.HTTP = expandDataVolumeSourceHTTP(in["http"].([]interface{}))
	result.PVC = expandDataVolumeSourcePVC(in["pvc"].([]interface{}))

	return result
}

func expandDataVolumeSourceHTTP(dataVolumeSourceHTTP []interface{}) *cdiv1.DataVolumeSourceHTTP {
	if len(dataVolumeSourceHTTP) == 0 || dataVolumeSourceHTTP[0] == nil {
		return nil
	}

	result := &cdiv1.DataVolumeSourceHTTP{}

	in := dataVolumeSourceHTTP[0].(map[string]interface{})

	if v, ok := in["url"].(string); ok {
		result.URL = v
	}
	if v, ok := in["secret_ref"].(string); ok {
		result.SecretRef = v
	}
	if v, ok := in["cert_config_map"].(string); ok {
		result.CertConfigMap = v
	}

	return result
}

func expandDataVolumeSourcePVC(dataVolumeSourcePVC []interface{}) *cdiv1.DataVolumeSourcePVC {
	if len(dataVolumeSourcePVC) == 0 || dataVolumeSourcePVC[0] == nil {
		return nil
	}

	result := &cdiv1.DataVolumeSourcePVC{}

	in := dataVolumeSourcePVC[0].(map[string]interface{})

	if v, ok := in["namespace"].(string); ok {
		result.Namespace = v
	}
	if v, ok := in["name"].(string); ok {
		result.Name = v
	}

	return result
}

func expandDataVolumeSourceRef(dataVolumeSourceRef []interface{}) *cdiv1.DataVolumeSourceRef {
	if len(dataVolumeSourceRef) == 0 || dataVolumeSourceRef[0] == nil {
		return nil
	}

	result := &cdiv1.DataVolumeSourceRef{}

	in := dataVolumeSourceRef[0].(map[string]interface{})

	if v, ok := in["namespace"].(string); ok {
		result.Namespace = &v
	}
	if v, ok := in["name"].(string); ok {
		result.Name = v
	}
	if v, ok := in["kind"].(string); ok {
		result.Kind = v
	}

	return result
}

func expandPersistentVolumeAccessModes(s []interface{}) []api.PersistentVolumeAccessMode {
	out := make([]api.PersistentVolumeAccessMode, len(s))
	for i, v := range s {
		out[i] = api.PersistentVolumeAccessMode(v.(string))
	}
	return out
}

// Flatteners

func flattenDataVolumeSource(in *cdiv1.DataVolumeSource) []interface{} {
	att := make(map[string]interface{})

	if in.HTTP != nil {
		att["http"] = flattenDataVolumeSourceHTTP(*in.HTTP)
	}
	if in.PVC != nil {
		att["pvc"] = flattenDataVolumeSourcePVC(*in.PVC)
	}

	return []interface{}{att}
}

func flattenDataVolumeSourceHTTP(in cdiv1.DataVolumeSourceHTTP) []interface{} {
	att := map[string]interface{}{
		"url":             in.URL,
		"secret_ref":      in.SecretRef,
		"cert_config_map": in.CertConfigMap,
	}
	return []interface{}{att}
}

func flattenDataVolumeSourcePVC(in cdiv1.DataVolumeSourcePVC) []interface{} {
	att := map[string]interface{}{
		"namespace": in.Namespace,
		"name":      in.Name,
	}
	return []interface{}{att}
}

func flattenDataVolumeSourceRef(in cdiv1.DataVolumeSourceRef) []interface{} {
	att := map[string]interface{}{
		"namespace": *in.Namespace,
		"name":      in.Name,
		"kind":      in.Kind,
	}
	return []interface{}{att}
}

func flattenDataVolumeStorage(in cdiv1.StorageSpec) []interface{} {
	att := map[string]interface{}{}
	if in.Resources.Requests != nil && len(in.Resources.Requests) > 0 {
		att["resources"] = flattenResourceRequirements(in.Resources)
	}
	if len(att) == 0 {
		return nil
	}
	return []interface{}{att}
}

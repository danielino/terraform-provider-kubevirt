package datavolume

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func dataVolumeSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": dataVolumeSourceSchema(),
		"pvc":    k8s.PersistentVolumeClaimSpecSchema(),
		"content_type": {
			Type:        schema.TypeString,
			Description: "ContentType options: \"kubevirt\", \"archive\".",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"kubevirt",
				"archive",
			}, false),
		},
		"source_ref": dataVolumeSourceRefSchema(),
		"storage":    dataVolumeStorageSchema(),
	}
}

func DataVolumeSpecSchema() *schema.Schema {
	fields := dataVolumeSpecFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("DataVolumeSpec defines our specification for a DataVolume type"),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandDataVolumeSpec(dataVolumeSpec []interface{}) (cdiv1.DataVolumeSpec, error) {
	result := cdiv1.DataVolumeSpec{}

	if len(dataVolumeSpec) == 0 || dataVolumeSpec[0] == nil {
		return result, nil
	}

	in := dataVolumeSpec[0].(map[string]interface{})

	if in["source_ref"] == nil {
		result.Source = expandDataVolumeSource(in["source"].([]interface{}))
	} else {
		result.SourceRef = expandDataVolumeSourceRef(in["source_ref"].([]interface{}))
	}

	// DEBUG qui:

	if in["pvc"] != nil && len(in["pvc"].([]interface{})) > 0 {
		p, err := k8s.ExpandPersistentVolumeClaimSpec(in["pvc"].([]interface{}))
		if err != nil {
			return result, err
		}
		// PATCH: Solo se ha almeno un campo!
		if p != nil && (len(p.AccessModes) > 0 || len(p.Resources.Requests) > 0 || len(p.Resources.Limits) > 0 ||
			p.Selector != nil || p.VolumeName != "" || p.StorageClassName != nil) {
			result.PVC = p
		} else {
			result.PVC = nil
		}
	} else if in["storage"] != nil {
		arr, ok := in["storage"].([]interface{})
		if ok && len(arr) > 0 {
			s := expandDataVolumeStorage(arr)
			if s != nil {
				result.Storage = s
			}
		}
	}

	if v, ok := in["content_type"].(string); ok {
		result.ContentType = cdiv1.DataVolumeContentType(v)
	}

	return result, nil
}

func flattenResourceRequirements(in api.ResourceRequirements) []interface{} {
	m := map[string]interface{}{}
	if len(in.Requests) > 0 {
		requests := map[string]interface{}{}
		for k, v := range in.Requests {
			requests[string(k)] = v.String()
		}
		m["requests"] = requests
	}
	if len(m) == 0 {
		return nil
	}
	return []interface{}{m} // PATCH: dev’essere una slice!
}

func FlattenPersistentVolumeAccessModes(modes []api.PersistentVolumeAccessMode) []interface{} {
	if len(modes) == 0 {
		return nil
	}
	result := make([]interface{}, len(modes))
	for i, mode := range modes {
		result[i] = string(mode)
	}
	return result
}

func flattenLabelSelector(in *metav1.LabelSelector) map[string]interface{} {
	if in == nil {
		return nil
	}
	result := map[string]interface{}{}
	if len(in.MatchLabels) > 0 {
		result["match_labels"] = in.MatchLabels
	}
	if len(in.MatchExpressions) > 0 {
		expressions := make([]map[string]interface{}, len(in.MatchExpressions))
		for i, exp := range in.MatchExpressions {
			expressions[i] = map[string]interface{}{
				"key":      exp.Key,
				"operator": string(exp.Operator),
				"values":   exp.Values,
			}
		}
		result["match_expressions"] = expressions
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func FlattenPersistentVolumeClaimSpec(in api.PersistentVolumeClaimSpec) []interface{} {
	att := make(map[string]interface{})

	if am := FlattenPersistentVolumeAccessModes(in.AccessModes); len(am) > 0 {
		att["access_modes"] = am
	}

	res := flattenResourceRequirements(in.Resources)
	if res != nil && len(res) > 0 {
		att["resources"] = res
	}

	if in.Selector != nil {
		if sel := flattenLabelSelector(in.Selector); sel != nil && len(sel) > 0 {
			att["selector"] = sel
		}
	}
	if in.VolumeName != "" {
		att["volume_name"] = in.VolumeName
	}
	if in.StorageClassName != nil {
		att["storage_class_name"] = *in.StorageClassName
	}

	// PATCH: rimuovi chiavi con valori vuoti!
	for k, v := range att {
		empty := false
		switch vv := v.(type) {
		case nil:
			empty = true
		case string:
			empty = vv == ""
		case []interface{}:
			empty = len(vv) == 0
		case map[string]interface{}:
			empty = len(vv) == 0
		}
		if empty {
			delete(att, k)
		}
	}

	if len(att) == 0 {
		return nil
	}
	return []interface{}{att}
}

func FlattenDataVolumeSpec(spec cdiv1.DataVolumeSpec) []interface{} {
	var source interface{}
	if spec.Source != nil {
		source = flattenDataVolumeSource(spec.Source)
	}
	att := map[string]interface{}{}

	if source != nil {
		att["source"] = source
	}

	if spec.PVC != nil {
		pvcMap := FlattenPersistentVolumeClaimSpec(*spec.PVC)
		if pvcMap != nil && len(pvcMap) > 0 {
			isReallyEmpty := true
			for _, m := range pvcMap {
				if mm, ok := m.(map[string]interface{}); ok && len(mm) > 0 {
					isReallyEmpty = false
				}
			}
			if !isReallyEmpty {
				att["pvc"] = pvcMap
			}
		}
	}

	if spec.SourceRef != nil {
		ref := flattenDataVolumeSourceRef(*spec.SourceRef)
		if ref != nil {
			att["source_ref"] = ref
		}
	}

	if spec.ContentType != "" {
		att["content_type"] = string(spec.ContentType)
	}

	if spec.Storage != nil {
		st := flattenDataVolumeStorage(*spec.Storage)
		if st != nil {
			att["storage"] = st
		}
	}

	// Qui: se la mappa att è vuota, ritorna nil!
	if len(att) == 0 {
		return nil
	}
	return []interface{}{att}
}

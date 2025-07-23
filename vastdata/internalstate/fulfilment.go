// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"math/big"
)

func getNestedAttributes(attr any) map[string]any {
	switch a := attr.(type) {
	case dsschema.SingleNestedAttribute:
		return convertAttrMap(a.Attributes)
	case rschema.SingleNestedAttribute:
		return convertAttrMap(a.Attributes)

	case dsschema.ListNestedAttribute:
		return convertAttrMap(a.NestedObject.Attributes)
	case rschema.ListNestedAttribute:
		return convertAttrMap(a.NestedObject.Attributes)

	case dsschema.MapNestedAttribute:
		return convertAttrMap(a.NestedObject.Attributes)
	case rschema.MapNestedAttribute:
		return convertAttrMap(a.NestedObject.Attributes)

	case dsschema.SetNestedAttribute:
		return convertAttrMap(a.NestedObject.Attributes)
	case rschema.SetNestedAttribute:
		return convertAttrMap(a.NestedObject.Attributes)

	default:
		return nil
	}
}

func convertAttrMap[A any](in map[string]A) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func FillFrameworkValues(val tftypes.Value, schema any) (map[string]attr.Value, error) {
	if !val.IsKnown() || val.IsNull() {
		return nil, fmt.Errorf("value is null or unknown")
	}

	obj := map[string]tftypes.Value{}
	if err := val.As(&obj); err != nil {
		return nil, fmt.Errorf("decode tftypes.Value: %w", err)
	}

	var attrTypes map[string]attr.Type
	switch s := schema.(type) {
	case rschema.Schema:
		attrTypes = make(map[string]attr.Type, len(s.Attributes))
		for k, v := range s.Attributes {
			attrTypes[k] = v.GetType()
		}
	case dsschema.Schema:
		attrTypes = make(map[string]attr.Type, len(s.Attributes))
		for k, v := range s.Attributes {
			attrTypes[k] = v.GetType()
		}
	default:
		return nil, fmt.Errorf("unsupported schema type: %T", schema)
	}

	result := make(map[string]attr.Value, len(obj))
	for k, v := range obj {
		attrType, ok := attrTypes[k]
		if !ok {
			continue
		}
		av, err := BuildAttrValueFromAny(attrType, v)
		if err != nil {
			return nil, fmt.Errorf("BuildAttrValueFromAny failed for %q: %w\nInspected object: %v", k, err, obj)
		}
		result[k] = av
	}

	return result, nil
}

func BuildAttrValueFromAny(t attr.Type, val any) (attr.Value, error) {
	if tfVal, ok := val.(attr.Value); ok {
		return tfVal, nil
	}
	if _, ok := val.(tftypes.Value); ok {
		return tfTypeToAttrType(t, val.(tftypes.Value))
	}

	if IsNil(val) {
		switch t.String() {
		case types.StringType.String():
			return types.StringNull(), nil
		case types.Int64Type.String():
			return types.Int64Null(), nil
		case types.Float64Type.String():
			return types.Float64Null(), nil
		case types.BoolType.String():
			return types.BoolNull(), nil
		default:
			switch tt := t.(type) {
			case types.ListType:
				return types.ListNull(tt.ElemType), nil
			case types.SetType:
				return types.SetNull(tt.ElemType), nil
			case types.MapType:
				return types.MapNull(tt.ElemType), nil
			case types.ObjectType:
				return types.ObjectNull(tt.AttributeTypes()), nil
			default:
				return nil, fmt.Errorf("unsupported null type: %T", t)
			}
		}
	}

	switch t.String() {
	case types.StringType.String():
		return types.StringValue(fmt.Sprintf("%v", val)), nil
	case types.Int64Type.String():
		n, err := ToInt(val)
		if err != nil {
			return nil, err
		}
		return types.Int64Value(n), nil
	case types.Float64Type.String():
		f, err := ToFloat(val)
		if err != nil {
			return nil, err
		}
		return types.Float64Value(f), nil
	case types.BoolType.String():
		b, ok := val.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", val)
		}
		return types.BoolValue(b), nil
	}

	switch tt := t.(type) {
	case types.ListType:
		rawList, ok := val.([]any)
		if !ok {
			return nil, fmt.Errorf("expected []any for list, got %T, value = %v", val, val)
		}
		var elems []attr.Value
		for i, item := range rawList {
			elem, err := BuildAttrValueFromAny(tt.ElemType, item)
			if err != nil {
				return nil, fmt.Errorf("list[%d]: %w", i, err)
			}
			elems = append(elems, elem)
		}
		return types.ListValueMust(tt.ElemType, elems), nil

	case types.SetType:
		rawList, ok := val.([]any)
		if !ok {
			return nil, fmt.Errorf("expected []any for set, got %T. value = %v", val, val)
		}
		var elems []attr.Value
		for i, item := range rawList {
			elem, err := BuildAttrValueFromAny(tt.ElemType, item)
			if err != nil {
				return nil, fmt.Errorf("set[%d]: %w", i, err)
			}
			elems = append(elems, elem)
		}
		return types.SetValueMust(tt.ElemType, elems), nil

	case types.MapType:
		rawMap, ok := val.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map[string]any for map, got %T", val)
		}
		converted := make(map[string]attr.Value)
		for k, v := range rawMap {
			elem, err := BuildAttrValueFromAny(tt.ElemType, v)
			if err != nil {
				return nil, fmt.Errorf("map[%q]: %w", k, err)
			}
			converted[k] = elem
		}
		return types.MapValueMust(tt.ElemType, converted), nil

	case types.ObjectType:
		rawObj, ok := val.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map[string]any for object, got %T", val)
		}
		converted := make(map[string]attr.Value)
		for k, fieldType := range tt.AttributeTypes() {
			fieldVal, exists := rawObj[k]
			if !exists {
				converted[k], _ = BuildAttrValueFromAny(fieldType, nil)
				continue
			}
			elem, err := BuildAttrValueFromAny(fieldType, fieldVal)
			if err != nil {
				return nil, fmt.Errorf("object[%q]: %w", k, err)
			}
			converted[k] = elem
		}
		obj, diags := types.ObjectValue(tt.AttributeTypes(), converted)
		if diags.HasError() {
			return nil, fmt.Errorf("objectValue: %s", diags)
		}
		return obj, nil
	}

	return nil, fmt.Errorf("unsupported type: %T", t)
}

func tfTypeToAttrType(t attr.Type, val tftypes.Value) (attr.Value, error) {
	// Check if the value is null or unknown
	// We always set null values to the corresponding Terraform type
	// To avoid unnecessary fulfilment errors
	if val.IsNull() || !val.IsKnown() {
		switch t.String() {
		case types.StringType.String():
			return types.StringNull(), nil
		case types.Int64Type.String():
			return types.Int64Null(), nil
		case types.Float64Type.String():
			return types.Float64Null(), nil
		case types.BoolType.String():
			return types.BoolNull(), nil
		default:
			switch tt := t.(type) {
			case types.ListType:
				return types.ListNull(tt.ElemType), nil
			case types.SetType:
				return types.SetNull(tt.ElemType), nil
			case types.MapType:
				return types.MapNull(tt.ElemType), nil
			case types.ObjectType:
				ot, ok := t.(types.ObjectType)
				if !ok {
					return nil, fmt.Errorf("expected ObjectType, got %T", t)
				}
				return types.ObjectNull(ot.AttributeTypes()), nil
			default:
				return nil, fmt.Errorf("unsupported null type: %T", t)
			}
		}
	}

	switch t.String() {
	case types.StringType.String():
		var s string
		if err := val.As(&s); err != nil {
			return nil, fmt.Errorf("string.As: %w", err)
		}
		return types.StringValue(s), nil

	case types.Int64Type.String():
		var f *big.Float
		if err := val.As(&f); err != nil {
			return nil, fmt.Errorf("int64.As (from number): %w", err)
		}
		i, _ := f.Int64()
		return types.Int64Value(i), nil

	case types.Float64Type.String():
		var bf *big.Float
		if err := val.As(&bf); err != nil {
			return nil, fmt.Errorf("float64.As (from number): %w", err)
		}
		f, _ := bf.Float64()
		return types.Float64Value(f), nil

	case types.BoolType.String():
		var b bool
		if err := val.As(&b); err != nil {
			return nil, fmt.Errorf("bool.As: %w", err)
		}
		return types.BoolValue(b), nil
	}

	// Handle nested types
	switch tt := t.(type) {
	case types.ListType:
		var list []tftypes.Value
		if err := val.As(&list); err != nil {
			return nil, fmt.Errorf("list.As: %w", err)
		}
		var elems []attr.Value
		for _, item := range list {
			elem, err := tfTypeToAttrType(tt.ElemType, item)
			if err != nil {
				return nil, err
			}
			elems = append(elems, elem)
		}
		return types.ListValueMust(tt.ElemType, elems), nil

	case types.SetType:
		var set []tftypes.Value
		if err := val.As(&set); err != nil {
			return nil, fmt.Errorf("set.As: %w", err)
		}
		var elems []attr.Value
		for _, item := range set {
			elem, err := tfTypeToAttrType(tt.ElemType, item)
			if err != nil {
				return nil, err
			}
			elems = append(elems, elem)
		}
		return types.SetValueMust(tt.ElemType, elems), nil

	case types.MapType:
		var raw map[string]tftypes.Value
		if err := val.As(&raw); err != nil {
			return nil, fmt.Errorf("map.As: %w", err)
		}
		converted := make(map[string]attr.Value, len(raw))
		for k, v := range raw {
			elem, err := tfTypeToAttrType(tt.ElemType, v)
			if err != nil {
				return nil, fmt.Errorf("map[%q]: %w", k, err)
			}
			converted[k] = elem
		}
		return types.MapValueMust(tt.ElemType, converted), nil

	case types.ObjectType:
		ot, ok := t.(types.ObjectType)
		if !ok {
			return nil, fmt.Errorf("expected ObjectType, got %T", t)
		}

		var raw map[string]tftypes.Value
		if err := val.As(&raw); err != nil {
			return nil, fmt.Errorf("object.As: %w", err)
		}

		converted := make(map[string]attr.Value, len(raw))
		attrTypes := ot.AttributeTypes()

		for k, v := range raw {
			fieldType, ok := attrTypes[k]
			if !ok {
				return nil, fmt.Errorf("object field %q not found in type definition", k)
			}

			elem, err := tfTypeToAttrType(fieldType, v)
			if err != nil {
				return nil, fmt.Errorf("object[%q]: %s", k, err)
			}
			converted[k] = elem
		}

		obj, err := types.ObjectValue(attrTypes, converted)
		if err != nil {
			return nil, fmt.Errorf("objectValue: %s", err)
		}
		return obj, nil

	default:
		return nil, fmt.Errorf("unsupported type: %T", t)
	}
}

func ConvertAttrValueToRaw(val attr.Value, t attr.Type) any {
	if val.IsNull() || val.IsUnknown() {
		return nil
	}

	switch v := val.(type) {
	case types.String:
		return v.ValueString()
	case types.Int64:
		return v.ValueInt64()
	case types.Float64:
		return v.ValueFloat64()
	case types.Bool:
		return v.ValueBool()

	case types.List:
		listType := t.(types.ListType)
		items := v.Elements()
		var out = make([]any, 0)
		for _, e := range items {
			if e.IsNull() || e.IsUnknown() {
				continue
			}
			elem := ConvertAttrValueToRaw(e, listType.ElemType)
			if elem != nil {
				out = append(out, elem)
			}
		}
		return out

	case types.Set:
		setType := t.(types.SetType)
		items := v.Elements()
		var out = make([]any, 0)
		for _, e := range items {
			if e.IsNull() || e.IsUnknown() {
				continue
			}
			elem := ConvertAttrValueToRaw(e, setType.ElemType)
			if elem != nil {
				out = append(out, elem)
			}
		}
		return out

	case types.Map:
		mapType := t.(types.MapType)
		elems := v.Elements()
		out := make(map[string]any, 0)

		for k, e := range elems {
			if e.IsNull() || e.IsUnknown() {
				continue
			}
			elem := ConvertAttrValueToRaw(e, mapType.ElemType)
			if elem != nil {
				out[k] = elem
			}
		}
		return out

	case types.Object:
		objType := t.(types.ObjectType)
		attrs := v.Attributes()
		out := make(map[string]any, len(attrs))
		for k, a := range attrs {
			attrType, ok := objType.AttributeTypes()[k]
			if !ok {
				continue
			}
			out[k] = ConvertAttrValueToRaw(a, attrType)
		}
		return out
	}

	panic(fmt.Sprintf("<unsupported: %T>", val))
}

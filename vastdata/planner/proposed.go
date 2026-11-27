// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
//
// Port of opentofu internal/plans/objchange ProposedNew (attribute and nested-attribute
// rules) from cty/configschema to tftypes and terraform-plugin-framework resource schemas.

package planner

import (
	"fmt"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type nesting int

const (
	nestingSingle nesting = iota
	nestingList
	nestingSet
	nestingMap
)

// nestedOf returns the attributes and nesting mode of a nested attribute.
func nestedOf(a rschema.Attribute) (map[string]rschema.Attribute, nesting, bool) {
	switch n := a.(type) {
	case rschema.SingleNestedAttribute:
		return n.Attributes, nestingSingle, true
	case rschema.ListNestedAttribute:
		return n.NestedObject.Attributes, nestingList, true
	case rschema.SetNestedAttribute:
		return n.NestedObject.Attributes, nestingSet, true
	case rschema.MapNestedAttribute:
		return n.NestedObject.Attributes, nestingMap, true
	default:
		return nil, 0, false
	}
}

// proposedNew combines the computed attribute values from prior with the configured values
// from config (ProposedNew). prior must be wholly known.
func proposedNew(attrs map[string]rschema.Attribute, prior, config tftypes.Value) (tftypes.Value, error) {
	if config.IsNull() && prior.IsNull() {
		return prior, nil
	}
	if prior.IsNull() {
		prior = emptyObject(config.Type())
	}
	if config.IsNull() || !config.IsKnown() {
		return prior, nil
	}
	vals, err := proposedNewAttributes(attrs, prior, config)
	if err != nil {
		return tftypes.Value{}, err
	}
	return tftypes.NewValue(config.Type(), vals), nil
}

func proposedNewAttributes(attrs map[string]rschema.Attribute, prior, config tftypes.Value) (map[string]tftypes.Value, error) {
	cfg, err := objectValues(config)
	if err != nil {
		return nil, err
	}
	var pri map[string]tftypes.Value
	if !prior.IsNull() {
		if pri, err = objectValues(prior); err != nil {
			return nil, err
		}
	}
	out := make(map[string]tftypes.Value, len(attrs))
	for name, a := range attrs {
		configV, ok := cfg[name]
		if !ok {
			return nil, fmt.Errorf("config has no attribute %q", name)
		}
		priorV, ok := pri[name]
		if !ok {
			priorV = tftypes.NewValue(configV.Type(), nil)
		}

		var newV tftypes.Value
		_, _, isNested := nestedOf(a)
		switch {
		// Optional+computed attributes are only computed when there is no configuration.
		case a.IsComputed() && configV.IsNull():
			newV = priorV
			// A prior value holding non-computed nested values must have come from config.
			if optionalValueNotComputable(a, priorV) {
				newV = configV
			}
		case isNested:
			if newV, err = proposedNewNestedType(a, priorV, configV); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
		default:
			newV = configV
		}
		out[name] = newV
	}
	return out, nil
}

func proposedNewNestedType(a rschema.Attribute, prior, config tftypes.Value) (tftypes.Value, error) {
	if !config.IsKnown() {
		return config, nil
	}
	attrs, mode, _ := nestedOf(a)
	switch mode {
	case nestingSingle:
		if config.IsNull() {
			return config, nil
		}
		return proposedNewObjectAttributes(attrs, prior, config)
	case nestingList:
		return proposedNewNestingList(attrs, prior, config)
	case nestingMap:
		return proposedNewNestingMap(attrs, prior, config)
	case nestingSet:
		return proposedNewNestingSet(attrs, prior, config)
	}
	return config, nil
}

func proposedNewObjectAttributes(attrs map[string]rschema.Attribute, prior, config tftypes.Value) (tftypes.Value, error) {
	if config.IsNull() {
		return config, nil
	}
	vals, err := proposedNewAttributes(attrs, prior, config)
	if err != nil {
		return tftypes.Value{}, err
	}
	return tftypes.NewValue(config.Type(), vals), nil
}

// Nested lists are correlated by index.
func proposedNewNestingList(attrs map[string]rschema.Attribute, prior, config tftypes.Value) (tftypes.Value, error) {
	if config.IsNull() {
		return config, nil
	}
	cfg, err := sliceValues(config)
	if err != nil || len(cfg) == 0 {
		return config, err
	}
	var pri []tftypes.Value
	if prior.IsKnown() && !prior.IsNull() {
		if pri, err = sliceValues(prior); err != nil {
			return tftypes.Value{}, err
		}
	}
	out := make([]tftypes.Value, 0, len(cfg))
	for i, ce := range cfg {
		if i >= len(pri) {
			out = append(out, ce)
			continue
		}
		v, err := proposedNewObjectAttributes(attrs, pri[i], ce)
		if err != nil {
			return tftypes.Value{}, err
		}
		out = append(out, v)
	}
	return tftypes.NewValue(config.Type(), out), nil
}

// Nested maps are correlated by key.
func proposedNewNestingMap(attrs map[string]rschema.Attribute, prior, config tftypes.Value) (tftypes.Value, error) {
	if config.IsNull() || !config.IsKnown() {
		return config, nil
	}
	cfg := map[string]tftypes.Value{}
	if err := config.As(&cfg); err != nil {
		return tftypes.Value{}, err
	}
	if len(cfg) == 0 {
		return config, nil
	}
	pri := map[string]tftypes.Value{}
	if prior.IsKnown() && !prior.IsNull() {
		if err := prior.As(&pri); err != nil {
			return tftypes.Value{}, err
		}
	}
	out := make(map[string]tftypes.Value, len(cfg))
	for k, ce := range cfg {
		pe, ok := pri[k]
		if !ok {
			out[k] = ce
			continue
		}
		v, err := proposedNewObjectAttributes(attrs, pe, ce)
		if err != nil {
			return tftypes.Value{}, err
		}
		out[k] = v
	}
	return tftypes.NewValue(config.Type(), out), nil
}

// Nested sets are correlated by finding a prior element that could have been derived from
// the config element by adding only computed values (validPriorFromConfig).
func proposedNewNestingSet(attrs map[string]rschema.Attribute, prior, config tftypes.Value) (tftypes.Value, error) {
	if !config.IsKnown() || config.IsNull() {
		return config, nil
	}
	cfg, err := sliceValues(config)
	if err != nil || len(cfg) == 0 {
		return config, err
	}
	var pri []tftypes.Value
	if prior.IsKnown() && !prior.IsNull() {
		if pri, err = sliceValues(prior); err != nil {
			return tftypes.Value{}, err
		}
	}
	elemType := config.Type().(tftypes.Set).ElementType
	used := make([]bool, len(pri))
	out := make([]tftypes.Value, 0, len(cfg))
	for _, ce := range cfg {
		pe := tftypes.NewValue(elemType, nil)
		for i, cand := range pri {
			if used[i] {
				continue
			}
			if validPriorObject(attrs, cand, ce) {
				pe = cand
				used[i] = true
				break
			}
		}
		v, err := proposedNewObjectAttributes(attrs, pe, ce)
		if err != nil {
			return tftypes.Value{}, err
		}
		out = append(out, v)
	}
	return tftypes.NewValue(config.Type(), out), nil
}

// optionalValueNotComputable reports whether a prior nested value holds any non-null
// non-computed attribute, which means it must have come from configuration.
func optionalValueNotComputable(a rschema.Attribute, val tftypes.Value) bool {
	if !a.IsOptional() {
		return false
	}
	return hasNonComputedValue(a, val)
}

// hasNonComputedValue walks the objects held by nested attribute a and reports whether any
// non-null value sits at a non-computed attribute, at any depth.
func hasNonComputedValue(a rschema.Attribute, val tftypes.Value) bool {
	found := false
	_ = forEachNestedObject(a, val, func(attrs map[string]rschema.Attribute, obj tftypes.Value) error {
		if found || obj.IsNull() || !obj.IsKnown() {
			return nil
		}
		vals, err := objectValues(obj)
		if err != nil {
			return nil
		}
		for name, na := range attrs {
			v, ok := vals[name]
			if !ok || v.IsNull() {
				continue
			}
			if !na.IsComputed() || hasNonComputedValue(na, v) {
				found = true
				return nil
			}
		}
		return nil
	})
	return found
}

// validPriorObject reports whether prior could have been derived from config by adding only
// computed values (validPriorFromConfig for one object).
func validPriorObject(attrs map[string]rschema.Attribute, prior, config tftypes.Value) bool {
	if prior.Equal(config) {
		return true
	}
	if prior.IsNull() {
		return true
	}
	if config.IsNull() || !config.IsKnown() {
		return false
	}
	pv, err := objectValues(prior)
	if err != nil {
		return false
	}
	cv, err := objectValues(config)
	if err != nil {
		return false
	}
	for name, a := range attrs {
		p, c := pv[name], cv[name]
		if p.Equal(c) {
			continue
		}
		// Nested sets cannot be correlated further, so the values must be equal.
		if _, isSet := c.Type().(tftypes.Set); isSet {
			return false
		}
		if _, _, nested := nestedOf(a); nested {
			if !validPriorNested(a, p, c) {
				return false
			}
			continue
		}
		// A leaf may only differ from config when it is computed and not configured.
		if !a.IsComputed() || !c.IsNull() {
			return false
		}
	}
	return true
}

func validPriorNested(a rschema.Attribute, prior, config tftypes.Value) bool {
	if prior.IsNull() {
		return true
	}
	if config.IsNull() || !config.IsKnown() {
		return false
	}
	attrs, mode, _ := nestedOf(a)
	switch mode {
	case nestingSingle:
		return validPriorObject(attrs, prior, config)
	case nestingList:
		pri, err1 := sliceValues(prior)
		cfg, err2 := sliceValues(config)
		if err1 != nil || err2 != nil || len(cfg) < len(pri) {
			return false
		}
		for i := range pri {
			if !validPriorObject(attrs, pri[i], cfg[i]) {
				return false
			}
		}
		return true
	case nestingMap:
		pri, cfg := map[string]tftypes.Value{}, map[string]tftypes.Value{}
		if prior.As(&pri) != nil || config.As(&cfg) != nil {
			return false
		}
		for k, pe := range pri {
			ce, ok := cfg[k]
			if !ok || !validPriorObject(attrs, pe, ce) {
				return false
			}
		}
		return true
	}
	return false
}

// forEachNestedObject calls fn with every object value held by nested attribute a.
func forEachNestedObject(a rschema.Attribute, v tftypes.Value, fn func(map[string]rschema.Attribute, tftypes.Value) error) error {
	attrs, mode, ok := nestedOf(a)
	if !ok || v.IsNull() || !v.IsKnown() {
		return nil
	}
	switch mode {
	case nestingSingle:
		return fn(attrs, v)
	case nestingList, nestingSet:
		elems, err := sliceValues(v)
		if err != nil {
			return err
		}
		for _, e := range elems {
			if err := fn(attrs, e); err != nil {
				return err
			}
		}
	case nestingMap:
		m := map[string]tftypes.Value{}
		if err := v.As(&m); err != nil {
			return err
		}
		for _, e := range m {
			if err := fn(attrs, e); err != nil {
				return err
			}
		}
	}
	return nil
}

func emptyObject(t tftypes.Type) tftypes.Value {
	obj, ok := t.(tftypes.Object)
	if !ok {
		return tftypes.NewValue(t, nil)
	}
	vals := make(map[string]tftypes.Value, len(obj.AttributeTypes))
	for k, at := range obj.AttributeTypes {
		vals[k] = tftypes.NewValue(at, nil)
	}
	return tftypes.NewValue(t, vals)
}

func objectValues(v tftypes.Value) (map[string]tftypes.Value, error) {
	out := map[string]tftypes.Value{}
	if err := v.As(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func sliceValues(v tftypes.Value) ([]tftypes.Value, error) {
	var out []tftypes.Value
	if err := v.As(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Copyright (c) HashiCorp, Inc.

package planner

import (
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// After apply, planning the same configuration again changes nothing: every server-filled
// value is carried from state. This is the steady state of every resource.
func TestSweep_ReplanAfterApplyIsNoOp(t *testing.T) {
	for _, name := range sortedResourceNames() {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, name)
			config := f.minimalConfig()
			prior := f.mustApply(config)

			res := f.mustPlan(prior, config)
			assert.Equal(t, NoOp, res.Action, "changed=%v unknown=%v", res.Changed, res.Unknown)
			assert.Empty(t, res.Changed)
			assert.Empty(t, res.Unknown)
			assert.Empty(t, res.RequiresReplace)
		})
	}
}

// Changing any configurable attribute plans an update of exactly that attribute, or a
// replacement when the attribute requires one. Removing an optional attribute clears it;
// removing an optional+computed one keeps the value the server has.
func TestSweep_ChangingAnAttribute(t *testing.T) {
	for _, name := range sortedResourceNames() {
		f := newFixture(t, name)
		for _, k := range sortedAttributeNames(f.sch.Attributes) {
			a := f.sch.Attributes[k]
			if !(a.IsRequired() || a.IsOptional()) || f.isWriteOnly(k) {
				continue
			}
			t.Run(name+"/"+k, func(t *testing.T) {
				t.Parallel()
				f := newFixture(t, name)
				base := f.minimalConfig()
				if _, ok := base[k]; !ok {
					base[k] = f.value(k, 0)
				}
				prior, err := f.apply(base)
				if rejectedByValidation(err) {
					t.Skipf("the provider rejects the synthetic value: %v", err)
				}
				require.NoError(t, err)

				changed := copyConfig(base)
				changed[k] = f.value(k, 1)
				if changed[k].Equal(base[k]) {
					t.Skip("the attribute accepts a single value")
				}
				res, err := f.plan(prior, changed)
				if rejectedByValidation(err) {
					t.Skipf("the provider rejects the synthetic value: %v", err)
				}
				require.NoError(t, err)
				assertChanges(t, res, k)
				if _, nested := f.typ.AttributeTypes[k].(tftypes.Object); !nested {
					if _, _, isNested := nestedOf(a); !isNested {
						assert.True(t, attrOf(t, res.planned, k).Equal(changed[k]), "the plan holds the configured value")
					}
				}

				if a.IsRequired() {
					return
				}
				removed := copyConfig(base)
				delete(removed, k)
				res = f.mustPlan(prior, removed)
				priorV := attrOf(t, prior, k)
				switch def := attrOf(t, f.mustPlan(f.null(), removed).planned, k); {
				case def.IsKnown() && !def.IsNull():
					// A schema default replaces what the configuration leaves out.
					assert.True(t, attrOf(t, res.planned, k).Equal(def), "the plan holds the default")
					if def.Equal(priorV) {
						assert.Equal(t, NoOp, res.Action)
					} else {
						assertChanges(t, res, k)
					}
				case a.IsComputed() && optionalValueNotComputable(a, priorV):
					// The prior value holds configured nested values, so it cannot be the
					// server's: the server computes it again (known after apply), unless
					// UseStateForUnknown keeps the prior value.
					if slices.Contains(res.Unknown, k) {
						assert.Equal(t, Update, res.Action)
						assert.True(t, res.RawValues()[k].IsNull(), "unknown is sent as null")
					} else {
						assert.True(t, attrOf(t, res.planned, k).Equal(priorV), "UseStateForUnknown keeps the prior value")
						assert.NotContains(t, res.Changed, k)
					}
				case a.IsComputed():
					assert.Equal(t, NoOp, res.Action,
						"an optional+computed attribute left out keeps the server value (changed=%v unknown=%v)", res.Changed, res.Unknown)
				default:
					assertChanges(t, res, k)
					assert.True(t, attrOf(t, res.planned, k).IsNull(), "the plan clears the attribute")
				}
			})
		}
	}
}

// assertChanges checks a plan that changes attribute k: an update, or a replacement exactly
// when k requires one.
func assertChanges(t *testing.T, res *Result, k string) {
	t.Helper()
	assert.Contains(t, res.Changed, k)
	if slices.Contains(res.RequiresReplace, k) {
		assert.Equal(t, Replace, res.Action)
		return
	}
	assert.Equal(t, Update, res.Action, "requiresReplace=%v", res.RequiresReplace)
}

// knownChanges lists the changed attributes whose planned value is known.
func knownChanges(res *Result) []string {
	var out []string
	for _, k := range res.Changed {
		if !slices.Contains(res.Unknown, k) {
			out = append(out, k)
		}
	}
	return out
}

// Copyright (c) HashiCorp, Inc.

package planner

import (
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// Every write-only secret with a <field>_md5 companion: the secret is never stored, the
// fingerprint is; the same secret plans nothing, a rotated one re-sends it, and a cleared
// optional one nulls the fingerprint (TERF-320).
func TestSweep_WriteOnlyMD5Secrets(t *testing.T) {
	ran := 0
	for _, name := range sortedResourceNames() {
		for _, field := range newFixture(t, name).hints.WriteOnlyMD5Fields {
			ran++
			t.Run(name+"/"+field, func(t *testing.T) {
				f := newFixture(t, name)
				md5 := field + "_md5"
				base := f.minimalConfig()
				base[field] = f.value(field, 0)
				prior := f.mustApply(base)
				priorSt := f.state(prior)
				assert.True(t, attrOf(t, prior, field).IsNull(), "a write-only value is not stored")
				priorMD5 := attrOf(t, prior, md5)
				require.True(t, priorMD5.IsKnown() && !priorMD5.IsNull(), "the fingerprint is stored")

				same := f.mustPlan(prior, base)
				assert.Equal(t, NoOp, same.Action, "the same secret plans nothing (changed=%v)", same.Changed)

				rotated := copyConfig(base)
				rotated[field] = f.value(field, 1)
				res := f.mustPlan(prior, rotated)
				assertChanges(t, res, md5)
				assert.NotContains(t, res.Changed, field, "the secret itself is never in the plan")
				assert.True(t, res.IsNull(field), "the plan holds no write-only value")
				desired := res.Copy()
				is.OverlayWriteOnly(desired, f.state(f.object(rotated)))
				var want string
				require.NoError(t, rotated[field].As(&want))
				assert.Equal(t, want, desired.String(field), "Update takes the secret from config")
				switch {
				case res.Action != Update:
				case slices.Contains(f.hints.CreateOnlyFields, field):
					assert.NotContains(t, desired.GetUpdateParams(priorSt), field, "a create-only secret is never patched")
				default:
					assert.Contains(t, desired.GetUpdateParams(priorSt), field, "a rotated secret is sent")
				}

				cleared := copyConfig(base)
				delete(cleared, field)
				if f.sch.Attributes[field].IsRequired() {
					_, err := f.plan(prior, cleared)
					require.Error(t, err)
					assert.Contains(t, err.Error(), "missing required argument")
					return
				}
				res = f.mustPlan(prior, cleared)
				assertChanges(t, res, md5)
				assert.True(t, attrOf(t, res.planned, md5).IsNull(), "a cleared secret nulls its fingerprint")
				assert.True(t, res.RawValues()[field].IsNull())
				assert.NotContains(t, res.GetUpdateParams(priorSt), field,
					"a clear is applied by AfterUpdateResource, not sent in the PATCH body")
			})
		}
	}
	assert.Greater(t, ran, 30, "every write-only secret with a fingerprint is covered")
}

// A write-only value without a fingerprint of its own cannot be compared with anything the
// state holds, so Terraform plans no change when only that value changes.
func TestSweep_WriteOnlyWithoutFingerprintIsInvisible(t *testing.T) {
	ran := 0
	for _, name := range sortedResourceNames() {
		f := newFixture(t, name)
		for _, field := range f.hints.WriteOnlyFields {
			if slices.Contains(f.hints.WriteOnlyMD5Fields, field) {
				continue
			}
			ran++
			t.Run(name+"/"+field, func(t *testing.T) {
				f := newFixture(t, name)
				base := f.minimalConfig()
				base[field] = f.value(field, 0)
				prior := f.mustApply(base)

				rotated := copyConfig(base)
				rotated[field] = f.value(field, 1)
				res := f.mustPlan(prior, rotated)
				assert.Equal(t, NoOp, res.Action, "changed=%v", res.Changed)
				assert.True(t, res.IsNull(field), "the plan holds no write-only value")
			})
		}
	}
	assert.Positive(t, ran)
}

// EditOnly fields are never part of the main PATCH body, but changing one is still a change:
// Update runs and applies it through GetChangedEditOnlyParams.
func TestSweep_EditOnlyFields(t *testing.T) {
	ran := 0
	for _, name := range sortedResourceNames() {
		for _, field := range newFixture(t, name).hints.EditOnlyFields {
			ran++
			t.Run(name+"/"+field, func(t *testing.T) {
				f := newFixture(t, name)
				base := f.minimalConfig()
				base[field] = f.value(field, 0)
				prior := f.mustApply(base)
				priorSt := f.state(prior)

				assert.Equal(t, NoOp, f.mustPlan(prior, base).Action)

				changed := copyConfig(base)
				changed[field] = f.value(field, 1)
				res := f.mustPlan(prior, changed)
				assertChanges(t, res, field)
				assert.Equal(t, []string{field}, knownChanges(res))
				assert.NotContains(t, res.GetUpdateParams(priorSt), field)
				assert.Contains(t, res.GetChangedEditOnlyParams(priorSt), field)
			})
		}
	}
	assert.GreaterOrEqual(t, ran, 9)
}

// Sub-resource triggers (get_s3cors_configuration, get_views_count, ...) are configuration
// only: flipping one is a change Update has to see, but it is never sent in the PATCH body.
func TestSweep_SubResourceTriggers(t *testing.T) {
	ran := 0
	for _, name := range sortedResourceNames() {
		for _, sr := range newFixture(t, name).hints.SubResources {
			trigger := sr.FieldTrigger
			if trigger == "" {
				continue
			}
			ran++
			t.Run(name+"/"+trigger, func(t *testing.T) {
				f := newFixture(t, name)
				base := f.minimalConfig()
				base[trigger] = tftypes.NewValue(tftypes.Bool, false)
				prior := f.mustApply(base)
				priorSt := f.state(prior)

				assert.Equal(t, NoOp, f.mustPlan(prior, base).Action)

				on := copyConfig(base)
				on[trigger] = tftypes.NewValue(tftypes.Bool, true)
				res := f.mustPlan(prior, on)
				assertChanges(t, res, trigger)
				assert.Equal(t, []string{trigger}, knownChanges(res))
				assert.NotContains(t, res.GetUpdateParams(priorSt), trigger)
			})
		}
	}
	assert.GreaterOrEqual(t, ran, 9)
}

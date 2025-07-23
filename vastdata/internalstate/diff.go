// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// cmp.Options configured to ignore slice order for []string and [][]string comparisons
var opts = cmp.Options{
	// Sort individual string slices (e.g., []string)
	cmpopts.SortSlices(func(a, b string) bool {
		return a < b
	}),
	// Sort outer slice of []string slices (e.g., [][]string)
	cmpopts.SortSlices(func(a, b []string) bool {
		return fmt.Sprint(a) < fmt.Sprint(b)
	}),
}

// Equal compares expected and actual values, applying custom options to ignore slice order.
// It returns (true, "") if equal, or (false, diff string) if different.
func Equal(expected, actual any) (bool, string) {
	equal := cmp.Equal(expected, actual, opts...)
	if equal {
		return true, ""
	}
	diff := fmt.Sprintf("mismatch: %s", cmp.Diff(expected, actual, opts...))
	return false, diff
}

// DiffMap returns a new map containing only the keys from map1 that are
// missing or different in map2 (based on deep equality with slice sorting).
func DiffMap[T ~map[string]any](map1, map2 T) map[string]any {
	diff := make(map[string]any)

	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok {
			// Key is missing in map2
			diff[k] = v1
			continue
		}

		if equal, _ := Equal(v1, v2); !equal {
			diff[k] = v1
		}
	}

	return removeNilValues(diff).(map[string]any)
}

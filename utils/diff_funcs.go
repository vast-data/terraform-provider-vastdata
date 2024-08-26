package utils

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
   This package should contain a collection of diff functions to be reused when comparing attribute state
*/

var unit2seconds map[string]int64 = map[string]int64{
	"y": 60 * 60 * 24 * 365,
	"Y": 60 * 60 * 24 * 365,
	"M": 60 * 60 * 24 * 30,
	"w": 60 * 60 * 24 * 7,
	"W": 60 * 60 * 24 * 7,
	"d": 60 * 60 * 24,
	"D": 60 * 60 * 24,
	"h": 60 * 60,
	"H": 60 * 60,
	"m": 60,
	"s": 1,
	"S": 1,
}

func asStingsList(i []any) []string {
	s := []string{}
	for _, o := range i {
		s = append(s, fmt.Sprintf("%v", o))
	}
	return s
}

func compareStrings(x, y string) int {
	return cmp.Compare(x, y)
}

func ListsDiffSupress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	/*Due to unducumented terraform behaviour this will run on every element of the list
	  and not on th list itself.

	  The keys that are given ar at the format of k.0, k.1 .......
	  so if we use d.GetChange(k) we will get changes only for this specific key
	  but if we ask for the entire name of the list we will get the entire list
	  so the entire name will be used to compare , the downside is that will run for every element of the list and get the same results.
	  We convert any list to list of strings , sort it and now comapring made easy.
	*/
	key, _, _ := strings.Cut(k, ".") // k is the current name of the attribute compared , with the index , we simply need the attribute name so we stip anything after "."
	oldData, newData := d.GetChange(key)
	if oldData == nil || newData == nil { // if any of them is nil it means new data was set so there can be no diff
		return false
	}
	o := asStingsList(oldData.([]any))
	slices.SortFunc(o, compareStrings)
	n := asStingsList(newData.([]any))
	slices.SortFunc(n, compareStrings)
	return reflect.DeepEqual(o, n)
}

func FrameTimeDiff(k, oldValue, newValue string, d *schema.ResourceData) bool {
	oldData, newData := d.GetChange(k)
	if oldData == nil || newData == nil { // if any of them is nil it means new data was set so there can be no diff
		return false
	}
	old := fmt.Sprintf("%v", oldValue)
	new := fmt.Sprintf("%v", newValue)
	if len(old) == 0 || len(new) == 0 {
		return false
	}
	oldUnit := string(old[len(old)-1])
	newUnit := string(new[len(new)-1])
	oldNumber := string(old[:(len(old) - 1)])
	newNumber := string(new[:(len(new) - 1)])
	_oldUnit, _oldUnitExists := unit2seconds[oldUnit]
	_newUnit, _newUnitExists := unit2seconds[newUnit]
	if !(_oldUnitExists && _newUnitExists) {
		return false
	}
	o, e := strconv.Atoi(oldNumber)
	if e != nil {
		return false
	}

	n, e := strconv.Atoi(newNumber)
	if e != nil {
		return false
	}

	return int64(o)*_oldUnit == int64(n)*_newUnit
}

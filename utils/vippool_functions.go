package utils

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var cnode_ids_cache = map[string]bool{}

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
func VippoolCnodeIdsDiffSupress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	/*Due to unducumented terraform behaviour this will run on every element of the list
	  and not on th list itself.

	  The keys that are given ar at the format of cnode_ids.0, cnode_ids.1 .......
	  so if we use d.GetChange(k) we will get changes only for this specific key
	  but if we ask for the entire name of the list we will get the entire list
	  so we will use the cnode_ids attribute ignoring keys , sorting it and than comparing 2 lists.
	  This is not very efficient as it will run as the number of elelemnts in the list per VipPool object so
	  we create a cache holding the memory pointer as string of d with the result, if it was already calculated
	  we return the calculated value
	*/
	_d := fmt.Sprintf("%v", d)
	result, exists := cnode_ids_cache[_d]
	if exists {
		return result
	}

	oldData, newData := d.GetChange("cnode_ids")
	if oldData == nil || newData == nil { // if any of them is nil it means new data was set so there can be no diff
		return false
	}
	o := asStingsList(oldData.([]any))
	slices.SortFunc(o, compareStrings)
	n := asStingsList(newData.([]any))
	slices.SortFunc(n, compareStrings)
	cnode_ids_cache[_d] = reflect.DeepEqual(o, n)
	return cnode_ids_cache[_d]
}

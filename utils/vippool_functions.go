package utils

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var cnode_ids_cache = map[string]bool{}

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

func VipPoolBeforePostPatch(m map[string]interface{}, client interface{}, ctx context.Context, d *schema.ResourceData) (map[string]interface{}, error) {
	tenant_id, tenant_id_exists := d.GetOkExists("tenant_id")
	tflog.Debug(ctx, fmt.Sprintf("[VipPoolBeforePostPatch] Checking Tenant ID provided : Tenant ID Exists: %v , Teant ID Value: %v", tenant_id_exists, tenant_id))
	if tenant_id_exists {
		t := tenant_id.(int)
		tflog.Debug(ctx, fmt.Sprintf("[VipPoolBeforePostPatch] Tenant ID is number at the value of %v", t))
		if t == 0 { // 0 means all tenants to we are sending null
			m["tenant_id"] = nil
		}

	}
	return m, nil
}

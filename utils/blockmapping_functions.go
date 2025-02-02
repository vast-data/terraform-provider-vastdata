package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

type blockMappingObjectData struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type blockMappingObject struct {
	Volume     blockMappingObjectData `json:"volume,omitempty"`
	Block_host blockMappingObjectData `json:"block_host,omitempty"`
}
type blockMappingRequest struct {
	Snapshot_id     int              `json:"snapshot_id,omitempty"`
	Pairs_to_add    []map[string]int `json:"pairs_to_add"`
	Pairs_to_remove []map[string]int `json:"pairs_to_remove"`
}

func NewBlockMappingRequest(pairs_to_add, pairs_to_remove []map[string]int, snapshot_id int) blockMappingRequest {
	if snapshot_id <= 0 {
		return blockMappingRequest{Pairs_to_add: pairs_to_add, Pairs_to_remove: pairs_to_remove}
	}
	return blockMappingRequest{Pairs_to_add: pairs_to_add, Pairs_to_remove: pairs_to_remove, Snapshot_id: snapshot_id}
}

func BlockMappingCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	hosts_ids, hosts_ids_exist := data["hosts_ids"]
	pairs_to_add := []map[string]int{}
	volume_id := data["volume_id"].(int)
	snapshot_id := 0
	_snapshot_id, snapshot_id_exists := data["snapshot_id"]
	if snapshot_id_exists {
		snapshot_id = _snapshot_id.(int)
	}

	tflog.Debug(ctx, fmt.Sprintf("[BlockMappingCreateFunc] host IDs hosts_ids_exist:%v, host_ids: %v", hosts_ids_exist, hosts_ids))
	if !hosts_ids_exist {
		return nil, fmt.Errorf("hosts_ids,attribute was not found")
	}

	hosts_ids_list, is_list := hosts_ids.([]interface{})
	if !is_list {
		return nil, fmt.Errorf("hosts_ids,attribute is not a list of interface{} but from the type of %T", hosts_ids)

	}
	for _, r := range hosts_ids_list {
		if v, ok := r.(int); !ok {
			return nil, fmt.Errorf("Cannot convert %v into int as it from the type of %T", r, r)
		} else {
			pairs_to_add = append(pairs_to_add, map[string]int{"host_id": v, "volume_id": volume_id})
		}
	}
	blk := NewBlockMappingRequest(pairs_to_add, []map[string]int{}, snapshot_id)
	b, marshal_error := json.Marshal(blk)
	if marshal_error != nil {
		return nil, marshal_error
	}
	tflog.Debug(ctx, fmt.Sprintf("[BlockMappingCreateFunc] Calling PATCH with payload: %v", string(b)))
	h, err := client.Patch(ctx, GenPath("blockmappings/bulk"), "application/json", bytes.NewReader(b), map[string]string{})
	if err != nil {
		return h, err
	}
	data["id"] = fmt.Sprintf("blockmappings-volume-%v-snapshot-%v", volume_id, snapshot_id)
	return FakeHttpResponse(h, data)

}

func getCurrentBlockMapping(ctx context.Context, _client interface{}, volume_id, snapshot_id int) (*http.Response, []int, error) {
	i := []int{}
	client := _client.(vast_client.JwtSession)
	u := url.Values{}
	u.Add("volume__id__in", fmt.Sprintf("%v", volume_id))
	h, err := client.Get(ctx, GenPath("blockmappings"), u.Encode(), map[string]string{})
	if err != nil {
		return h, i, err
	}
	b, err := io.ReadAll(h.Body)
	if err != nil {
		return h, i, err
	}
	t := []blockMappingObject{}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, i, err
	}
	for _, q := range t {
		i = append(i, q.Block_host.Id)
	}
	return h, i, nil
}
func BlockMappingGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	volume_id := d.Get("volume_id").(int)
	_snapshot_id, snapshot_id_exists := d.GetOkExists("snapshot_id")
	snapshot_id := 0
	if snapshot_id_exists {
		snapshot_id = _snapshot_id.(int)
	}
	h, i, e := getCurrentBlockMapping(ctx, _client, volume_id, snapshot_id)
	if e != nil {
		return h, e
	}
	data := map[string]interface{}{}
	data["volume_id"] = volume_id
	data["hosts_ids"] = i
	data["id"] = d.Get("id")
	if snapshot_id != 0 {
		data["snapshot_id"] = snapshot_id
	}
	return FakeHttpResponse(h, data)

}

func BlockMappingUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	client := _client.(vast_client.JwtSession)
	old, new := d.GetChange("hosts_ids")
	volume_id := d.Get("volume_id").(int)
	_snapshot_id, snapshot_id_exists := d.GetOkExists("snapshot_id")
	snapshot_id := 0
	if snapshot_id_exists {
		snapshot_id = _snapshot_id.(int)
	}
	tflog.Debug(ctx, fmt.Sprintf("Old Hosts IDs: %v, New Hosts IDs:%v", old, new))
	old_ids := map[int]struct{}{}
	new_ids := map[int]struct{}{}
	pairs_to_remove := []map[string]int{}
	pairs_to_add := []map[string]int{}

	for _, o := range old.([]interface{}) {
		old_ids[o.(int)] = struct{}{}
	}
	for _, n := range new.([]interface{}) {
		_n := n.(int)
		new_ids[n.(int)] = struct{}{}
		_, exists := old_ids[_n]
		if exists {
			continue
		} else {
			pairs_to_add = append(pairs_to_add, map[string]int{"host_id": _n, "volume_id": volume_id})
		}

	}
	for _, o := range old.([]interface{}) {
		_o := o.(int)
		_, exists := new_ids[_o]
		if !exists {
			pairs_to_remove = append(pairs_to_remove, map[string]int{"host_id": _o, "volume_id": volume_id})

		}
	}
	blk := NewBlockMappingRequest(pairs_to_add, pairs_to_remove, snapshot_id)
	b, marshal_error := json.Marshal(blk)
	if marshal_error != nil {
		return nil, marshal_error
	}
	tflog.Debug(ctx, fmt.Sprintf("[BlockMappingUpdateFunc] Calling PATCH with payload: %v", string(b)))
	h, err := client.Patch(ctx, GenPath("blockmappings/bulk"), "application/json", bytes.NewReader(b), map[string]string{})
	return h, err

}

func BlockMappingDeleteFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	rgx := `blockmappings-volume-(?P<volume_id>\d+)-snapshot-(?P<snapshot_id>\d+)`
	client := _client.(vast_client.JwtSession)
	re := regexp.MustCompile(rgx)
	_id, id_exists := attr["id"]
	if !id_exists {
		return nil, fmt.Errorf("attributes provided does not hold the id")
	}
	m := re.FindStringSubmatch(fmt.Sprintf("%v", _id))
	if m == nil {
		return nil, fmt.Errorf("Id: %v , does not matches regular expression , %v", _id, rgx)
	}
	volume_id, _ := strconv.Atoi(m[1])
	snapshot_id, _ := strconv.Atoi(m[2])
	h, i, e := getCurrentBlockMapping(ctx, _client, volume_id, snapshot_id)
	if e != nil {
		return h, e
	}

	pairs_to_remove := []map[string]int{}
	for _, r := range i {
		pairs_to_remove = append(pairs_to_remove, map[string]int{"host_id": r, "volume_id": volume_id})
	}
	blk := NewBlockMappingRequest([]map[string]int{}, pairs_to_remove, snapshot_id)
	b, marshal_error := json.Marshal(blk)
	if marshal_error != nil {
		return nil, marshal_error
	}

	tflog.Debug(ctx, fmt.Sprintf("[BlockMappingDeleteFunc] Calling PATCH with payload: %v", string(b)))
	h, err := client.Patch(ctx, GenPath("blockmappings/bulk"), "application/json", bytes.NewReader(b), map[string]string{})
	if err != nil {
		return h, err
	}

	/*
	   We are goting to add some wait + check for those cases where we want to remove the volume just after
	   we remove the host_ids in some cases since this action is asynchronous actual removal might take time
	*/
	no_blocks := false
	for c := 0; (c <= 300) && (!no_blocks); c += 10 {
		h, i, e := getCurrentBlockMapping(ctx, _client, volume_id, snapshot_id)
		if e != nil {
			tflog.Error(ctx, fmt.Sprintf("[BlockMappingDeleteFunc], Error occured while wating for complease removal of all host_is from volume with id :%v", volume_id))
			return h, e
		}
		if len(i) > 0 {
			time.Sleep(10 * time.Second)
			continue
		} else {
			no_blocks = true
		}
	}
	if !no_blocks {
		return h, fmt.Errorf("Failed to remove all all blocks from volumeid :%v", volume_id)
	}

	return h, err

}

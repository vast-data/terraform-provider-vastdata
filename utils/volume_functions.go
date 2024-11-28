package utils

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var re = regexp.MustCompile(`\w+:\w+`)

func listToKVmap(ctx context.Context, n map[string]interface{}) (map[string]interface{}, error) {
	m := map[string]string{}
	volume_tags, volume_tags_exists := n["volume_tags"]
	if volume_tags_exists {
		tflog.Debug(ctx, fmt.Sprintf("[listToKVmap] volume_tags attribute exists , %v", volume_tags))
		l, volume_tags_is_list := volume_tags.([]interface{})
		if volume_tags_is_list {
			tflog.Debug(ctx, fmt.Sprintf("[listToKVmap] volume_tags is list , %v", volume_tags))
			for _, i := range l {
				j := fmt.Sprintf("%v", i)
				if re.Match([]byte(j)) {
					tflog.Debug(ctx, fmt.Sprintf(`[listToKVmap] volume tag %v matches \w+:\w+`, i))
					o := strings.SplitN(j, ":", 2)
					m[o[0]] = o[1]
				}

			}
		}
	}
	n["tags"] = m
	return n, nil
}

func kvMapToList(ctx context.Context, r *http.Response) (map[string]interface{}, error) {
	m := []string{}
	n := map[string]interface{}{}
	e := UnmarshelBodyToMap(r, &n)
	if e != nil {
		tflog.Debug(ctx, fmt.Sprintf("[kvMapToList] Error occured while unmarshling response body , %v", e))
		return n, e
	}
	tags, has_tags := n["tags"]
	tflog.Debug(ctx, fmt.Sprintf("[kvMapToList] Checking if the http resonse has tags unmarsheled data: %v", n))
	if has_tags {
		tflog.Debug(ctx, fmt.Sprintf("[kvMapToList] tags were found: %v, checking if tags are a map at the format of key:val", tags))
		//Tags should be at the format of { "key1":"value1", "key2":"value2" .....}
		tflog.Debug(ctx, fmt.Sprintf("[kvMapToList] The type of the tags is %T", tags))
		tags_map, is_tags_map := tags.(map[string]interface{})
		if is_tags_map {
			tflog.Debug(ctx, fmt.Sprintf("[kvMapToList] tags, %v  are atthe format of key:val building list of tags", tags))
			for k, v := range tags_map {
				tflog.Debug(ctx, fmt.Sprintf("[kvMapToList] Adding key value Key:%v , Value:%v", k, v))
				m = append(m, fmt.Sprintf("%v:%v", k, v))
			}
			tflog.Debug(ctx, fmt.Sprintf("[TagstoKVList] list of key:value maps created %v", m))
		}
	}
	n["volume_tags"] = m
	return n, nil
}

func VolumeCreateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, headers map[string]string) (*http.Response, error) {
	data, err := listToKVmap(ctx, data)
	if err != nil {
		return nil, err
	}
	return DefaultCreateFunc(ctx, _client, attr, data, headers)
}

func VolumeUpdateFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, data map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	data, err := listToKVmap(ctx, data)
	if err != nil {
		return nil, err
	}
	return DefaultUpdateFunc(ctx, _client, attr, data, d, headers)
}

func VolumeGetFunc(ctx context.Context, _client interface{}, attr map[string]interface{}, d *schema.ResourceData, headers map[string]string) (*http.Response, error) {
	response, err := DefaultGetFunc(ctx, _client, attr, d, headers)
	if err != nil {
		return response, err
	}
	data_with_tags, tags_err := kvMapToList(ctx, response)
	if tags_err != nil {
		return response, tags_err
	}
	tflog.Debug(ctx, fmt.Sprintf("[VolumeGetFunc] Data With Tags Returned %v", data_with_tags))
	return FakeHttpResponse(response, data_with_tags)

}

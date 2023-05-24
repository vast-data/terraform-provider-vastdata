package vast_client

import (
	"reflect"
	"fmt"
	"net/url"
	
)

//Build a query string from list of attributes supplied searching for attributes at in
func Query_builder[T any](in T ,attributes []string) string {
	//Convert to map for better performance 
	u:=url.Values{}
	attributes_map := make(map[string]int)
	for i,c := range attributes {
		attributes_map[c]=i
	}
	t:=reflect.TypeOf(in)
	q:=reflect.ValueOf(in)
	for _,n := range reflect.VisibleFields(t) {
		_,ok:= attributes_map[n.Name]
		if ok {
			u.Add(n.Name, fmt.Sprint(q.FieldByName(n.Name).Interface()))
		}
	}

	return u.Encode()
}

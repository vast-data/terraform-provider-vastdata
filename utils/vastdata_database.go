package utils

import "reflect"

// This struct will be filled by json Unmarshel
type key_value_types struct {
	Column_type string `json:"column_type"`
}
type db_field_type struct {
	Column_type string          `json:"column_type"`
	Fields      []db_field      `json:"fields"`
	Key_type    key_value_types `json:"key_type"`
	Value_type  key_value_types `json:"value_type"`
}

type db_field struct {
	Name  string        `json:"name"`
	Field db_field_type `json:"field"`
}

type db_fields_couple struct {
	Old  db_field
	New  db_field
	Path string
}

type deleted_db_field struct {
	Field db_field
	Path  string
}
type fields_change struct {
	Deleted []deleted_db_field
	Updated []db_fields_couple
}

func (f *fields_change) AddDeletedField(path string, o db_field) {
	f.Deleted = append(f.Deleted, deleted_db_field{Field: o, Path: path})
}

func compare_fields(old_fields, new_fields []db_field, chages *fields_change, path string) ([]db_field, error) {
	found := false

	if len(old_fields) > len(new_fields) { // Fields deleted
		found = false
		for _, f := range old_fields {
			for _, q := range new_fields {
				if reflect.DeepEqual(f, q) {
					//since we have found it be can break the loop
					found = true
					break
				}

			}
			if found { // A field was found so it was not deleted
				continue
			}
			chages.AddDeletedField(path, f)

		}

	}
}

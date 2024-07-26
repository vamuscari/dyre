package dyre

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

// TODO: Add a check for unkown fields in the json file
// check for unknow types in type fields
func Init(path string) (map[string]DyRe_Request, error) {
	m, err := readDyreJSON(path)
	if err != nil {
		return nil, err
	}

	re, err := parseDyreJSON(m)
	if err != nil {
		return nil, err
	}

	return re, nil
}

func readDyreJSON(path string) ([]map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseDyreJSON(m []map[string]interface{}) (map[string]DyRe_Request, error) {
	// required fields
	requests := make(map[string]DyRe_Request)
	for i, js_request := range m {
		if _, ok := js_request["name"]; !ok {
			return nil, errors.New(fmt.Sprintf("No field <name> on request index %d\n", i))
		}

		dy_request := DyRe_Request{
			name: js_request["name"].(string),
		}

		if fields, ok := js_request["fields"]; ok {
			dy_request.fields = parseDryeJSONFields(fields.([]interface{}))
		}

		if groups, ok := js_request["groups"]; ok {
			if groupsMap, ok := groups.([]interface{}); ok {
				dy_request.groups = parseDryeJSONGroups(groupsMap)
			} else {
				log.Printf("Field <groups> should be an array of objects on request  %s\nGivenType: %s", dy_request.name, reflect.TypeOf(groupsMap).String())
			}

		}

		if dy_request.fields == nil && dy_request.groups == nil {
			return nil, errors.New(fmt.Sprintf("No field <fields> on request  %s\n", dy_request.name))
		}

		if tableName, ok := js_request["tableName"]; ok {
			dy_request.sql = DyRe_SQL{tableName: tableName.(string)}
		}

		if sqlFile, ok := js_request["sqlFile"]; ok {
			dy_request.sql = DyRe_SQL{sqlFile: sqlFile.(string)}
		}

		dy_request.updateNames()

		requests[dy_request.name] = dy_request
	}

	return requests, nil
}

// check for string type,
// check for map/object type
// map [ name of field ] Dyre_Field.
// map is used for faster lookup times in large arrays
// TODO: check jsonMap to make sure field exists
func parseDryeJSONFields(a []any) map[string]DyRe_Field {

	dyre_fields := map[string]DyRe_Field{}

	for _, v := range a {

		field_type := reflect.TypeOf(v).String()

		if field_type == "string" {
			dyre_fields[v.(string)] = DyRe_Field{
				name:      v.(string),
				required:  false,
				typeName:  DefaultType,
				sqlSelect: v.(string),
			}
		}

		if field_type == "map[string]interface {}" {

			field_map := v.(map[string]interface{})

			new_field := DyRe_Field{}

			if name, ok := field_map["name"]; ok {
				if nameString, ok := name.(string); ok {
					new_field.name = nameString
				} else {
					log.Printf("Type <name> not string: %v,\n", name)
					continue
				}
			} else {
				continue
			}

			if required, ok := field_map["required"]; ok {
				if requiredBool, ok := required.(bool); ok {
					new_field.required = requiredBool
				} else {
					log.Printf("Type <required> not bool on field %s\n", new_field.name)
					new_field.required = false
				}
			} else {
				new_field.required = false
			}

			if typeName, ok := field_map["type"]; ok {
				if typeNameString, ok := typeName.(string); ok {
					new_field.typeName = typeNameString
				} else {
					log.Printf("Type <typeName> not string on field %s\n", new_field.name)
					new_field.typeName = DefaultType
				}
			} else {
				new_field.typeName = DefaultType
			}

			if querySelect, ok := field_map["querySelect"]; ok {
				if querySelectString, ok := querySelect.(string); ok {
					new_field.sqlSelect = querySelectString
				} else {
					log.Printf("Type <querySelect> not string on field %s\n", new_field.name)
					new_field.sqlSelect = new_field.name
				}
			} else {
				new_field.sqlSelect = new_field.name
			}

			dyre_fields[new_field.name] = new_field
		}
	}

	return dyre_fields
}

// look for defined group.
// infer infer group from json map if only name is given.
// not sure I want to "fill in the blank" on the map type
// If someone doesnt put all of the fields in on the map type is should not infer. This could lead to confusion otherwise
//
//	{
//	  "name": "group1",
//	  "fields": [
//	    "field5",
//	    "field6"
//	  ]
//	},
func parseDryeJSONGroups(a []interface{}) map[string]DyRe_Group {
	dyre_groups := map[string]DyRe_Group{}
	for _, v := range a {
		if _, ok := v.(map[string]interface{}); !ok {
			log.Printf("A dyre group must be a Object:\n%v", v)
			continue
		}
		group := v.(map[string]interface{})
		new_group := DyRe_Group{}

		if name, ok := group["name"]; ok {
			if nameString, ok := name.(string); ok {
				new_group.name = nameString
			} else {
				log.Printf("Type group <name> not string: %v,\n", name)
				continue
			}
		} else {
			log.Printf("Field <name> does not exit on group\n")
			continue
		}

		if group_required, ok := group["required"]; ok {
			new_group.required = group_required.(bool)
		} else {
			new_group.required = false
		}

		if group_fields, ok := group["fields"]; ok {
			new_group.fields = parseDryeJSONFields(group_fields.([]interface{}))
		} else {
			continue
		}

		dyre_groups[new_group.name] = new_group
	}

	return dyre_groups
}

func getMapKeys(m map[string]any) []string {
	count := len(m)
	keys := make([]string, count)
	i := 0
	for key := range m {
		keys[i] = key
		i += 1
	}
	return keys
}

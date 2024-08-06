package dyre

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Test parsing a known json file
func TestParseJSON(t *testing.T) {

	const parserJSON string = `
	[
		{
			"name": "ParserTest",
			"tableName": "Table",
			"fields": [
				{
					"name": "field1",
					"required": true
				},
				{
					"name": "field2",
					"required": true,
					"sqlSelect": "fieldX AS field2"
				},
				"field3"
			],
			"groups": [
				{
					"name": "group1",
					"fields": [
						"field5",
						"field6"
					]
				},
				{
					"name": "group2",
					"fields": [
						"field7",
						"field8"
					]
				}
			]
		}
	]
	`

	var m []map[string]interface{}
	err := json.Unmarshal([]byte(parserJSON), &m)
	if err != nil {
		t.Fatalf(`Failed Unmarshal %v`, err)
	}

	dyre_requests, err := parseDyreJSON(m)
	if err != nil {
		t.Fatalf(`error: %v`, err)
	}

	test_name := "ParserTest"
	re, ok := dyre_requests[test_name]
	if !ok {
		t.Error("Map key error for ParserTest")
	}

	t.Run("DyRe_Request.name", func(t *testing.T) {
		want := test_name
		got := re.name
		if want != got {
			t.Errorf("Want: %s, Got: %s", want, got)
		}
	})

	t.Run("DyRe_Request.FieldNames()", func(t *testing.T) {
		want := []string{"field1", "field2", "field3"}
		got := re.FieldNames()
		errors := deepEqualStringArray(want, got)
		if len(errors) > 0 {
			for _, err := range errors {
				t.Error(err)
			}
		}
	})

	t.Run("DyRe_Request.GroupNames()", func(t *testing.T) {
		want := []string{"group1", "group2"}
		got := re.GroupNames()
		errors := deepEqualStringArray(want, got)
		if len(errors) > 0 {
			for _, err := range errors {
				t.Error(err)
			}
		}
	})

	t.Run("DyRe_Request.ValidateRequest()", func(t *testing.T) {
		// has real field, has fake field, missing field
		given_fields := []string{"field1", "field100"}
		given_groups := []string{"group1", "group100"}

		valid, valid_err := re.ValidateRequest(given_fields, given_groups)
		if valid_err != nil {
			t.Error(valid_err)
		}

		valid_field_names := []string{}
		for _, f := range valid._fields {
			valid_field_names = append(valid_field_names, f.name)
		}

		valid_group_names := []string{}
		for _, g := range valid._groups {
			valid_group_names = append(valid_group_names, g.name)
		}

		expected_fields_names := []string{"field1", "field2"}
		expected_groups_names := []string{"group1"}

		field_errors := deepEqualStringArray(expected_fields_names, valid_field_names)
		if len(field_errors) > 0 {
			for _, err := range field_errors {
				t.Error(err)
			}
		}

		group_errors := deepEqualStringArray(expected_groups_names, valid_group_names)
		if len(group_errors) > 0 {
			for _, err := range group_errors {
				t.Error(err)
			}
		}

		expected_sql_field := []string{"field1", "fieldX AS field2", "field5", "field6"}

		sql_errors := deepEqualStringArray(expected_sql_field, valid._sqlFields)
		if len(sql_errors) > 0 {
			for _, err := range sql_errors {
				t.Error(err)
			}
		}

	})

}

func deepEqualStringArray(arr1 []string, arr2 []string) []error {
	errors := []error{}
	for _, k := range arr1 {
		if !contains(arr2, k) {
			errors = append(errors, fmt.Errorf("%s not found in %v\n", k, arr2))
		}
	}
	for _, k := range arr2 {
		if !contains(arr1, k) {
			errors = append(errors, fmt.Errorf("%s not found in %v\n", k, arr1))
		}
	}
	return errors
}

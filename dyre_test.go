package dyre

import (
	"testing"
)

func buildTestRequests() (map[string]DyRe_Request, error) {
	fileName := "./test_data/slqGen.json"
	js, err := readDyreJSON(fileName)
	if err != nil {
		return nil, err
	}

	dyre_requests, err := parseDyreJSON(js)
	if err != nil {
		return nil, err
	}

	return dyre_requests, nil
}

func TestValidation(t *testing.T) {
	requests, err := buildTestRequests()
	if err != nil {
		t.Fatalf("\nFailed to build test requests\nError:%v", err)
	}

	for name, request := range requests {
		name := name
		r := request
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if name != r.name {
				t.Errorf("Request %s map name %s does not match", r.name, name)
			}
			// Build many validated options
			valid := []DyRe_Validated{}
			fields := []string{}
			for _, f := range r.fields {
				fields = append(fields, f.name)
				groups := []string{}
				for _, g := range r.groups {
					groups = append(groups, g.name)
					f_valid, err := r.ValidateRequest(fields, fields)
					if err != nil {
						t.Errorf("Request:%s\nCould not validate %v and %v", name, fields, groups)
					} else {
						valid = append(valid, f_valid)
					}
				}
			}

			// Test Validated Options
			for _, v := range valid {
			}

		})
	}
}

package dyre

import (
	"fmt"
	"testing"
)

// TestReadJSON opens the dyre.json
func TestReadJSON(t *testing.T) {

	fileName := "./test_data/parser.json"
	_, err := readDyreJSON(fileName)
	if err != nil {
		t.Fatalf(`ReadJSON(%s), error: %v`, fileName, err)
	}
}

// Test for fake file reading json
func TestFakeFileReadJSON(t *testing.T) {
	fileName := "./test_data/fake.json"
	js, err := readJSON(fileName)
	if err == nil {
		t.Fatalf(`ReadJSON(%s) = %q, %v, want, error`, fileName, js, err)
	}
}

// Test parsing a known json file
func TestParseJSON(t *testing.T) {

	AddTypes(sqlTypes)

	fileName := "./test_data/parser.json"
	js, err := readDyreJSON(fileName)
	if err != nil {
		t.Fatalf(`ReadJSON(%s), error: %v`, fileName, err)
	}

	dyre_requests, err := parseDyreJSON(js)
	if err != nil {
		t.Fatalf(`PardeDyreJSON(%s), error: %v`, fileName, err)
	}

	for _, req := range dyre_requests {
		fmt.Printf("\n\nName: %s\n", req.name)
		// fmt.Printf("Fields: %v\n", req.fields)
		// fmt.Printf("FieldNames: %s\n", req._fieldNames)
		// fmt.Printf("Groups: %v\n", req.groups)
		// fmt.Printf("GroupNames: %s\n", req._groupNames)
		// fmt.Printf("JSONMap: %v\n", req.jsonMap)
	}

}

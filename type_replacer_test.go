package dyre

import (
	"fmt"
	"reflect"
	"testing"
)

// in a list of key values where the key is also the name of the declared type
//
//	 {
//			string:string
//			int:int
//	 }
//
// reflect should be the same as the output type
func CheckTypes(m map[string]interface{}, keyName string, printout bool) bool {
	badType := false
	for k, v := range m {
		refType := reflect.TypeOf(v).String()
		if refType == "map[string]interface {}" && CheckTypes(v.(map[string]interface{}), k, printout) {
			badType = true
		} else if refType != k && refType != "map[string]interface {}" {
			fmt.Printf("Failed Type Change (%s:%s) -> %s\n", keyName, k, refType)
			badType = true
		} else if printout {
			fmt.Printf("Successful Type change (%s:%s) -> %s\n", keyName, k, refType)
		}
	}
	return badType
}

func TestReplaceTypes(t *testing.T) {
	fileName := "./test_data/types.json"
	js, err := readJSON(fileName)
	if err != nil {
		t.Fatalf(`ReadJSON(%s) error: %v`, fileName, err)
	}

	replace_types, err := replaceTypes(js)
	failed := CheckTypes(replace_types, "base", false)

	if err != nil || failed == true {
		t.Fatalf(`ReplaceTypes() on %s, type name should match key name`, fileName)
	}
}

func TestBadReplaceTypes(t *testing.T) {
	fileName := "./test_data/types.json"
	js, err := readJSON(fileName)
	if err != nil {
		t.Fatalf(`ReadJSON(%s) error: %v`, fileName, err)
	}
	// inject invalid type then try to convert
	js["badType"] = "badType"
	replace_types, err := replaceTypes(js)
	failed := CheckTypes(replace_types, "base", false)
	if err == nil && failed == false {
		t.Fatal("Failed to handle bad type\n")
	}
}

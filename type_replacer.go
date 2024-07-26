package dyre

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
)

// For fields that type is not declared
// Needed is a map of types does not exist
var DefaultType = "string"

var Types = map[string]interface{}{
	"string":     string(""),
	"int":        int(0),
	"int8":       int8(0),
	"int16":      int16(0),
	"int32":      int32(0),
	"int64":      int64(0),
	"uint":       uint(0),
	"uint8":      uint8(0), // byte
	"uint16":     uint16(0),
	"uint32":     uint32(0),
	"uint64":     uint64(0),
	"uintptr":    uintptr(0), // rune
	"bool":       bool(false),
	"float32":    float32(0),
	"float64":    float64(0),
	"complex64":  complex64(0),
	"complex128": complex128(0),
}

var sqlTypes = map[string]interface{}{
	"sql.NullString":  sql.NullString{},
	"sql.NullBool":    sql.NullBool{},
	"sql.NullByte":    sql.NullByte{},
	"sql.NullTime":    sql.NullTime{},
	"sql.NullInt16":   sql.NullInt16{},
	"sql.NullInt32":   sql.NullInt32{},
	"sql.NullInt64":   sql.NullInt64{},
	"sql.NullFloat64": sql.NullFloat64{},
}

func readJSON(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func replaceTypes(m map[string]interface{}) (map[string]interface{}, error) {
	keys := getMapKeys(m)
	typeKeys := getMapKeys(Types)
	for _, key := range keys {
		if reflect.TypeOf(m[key]).String() == "string" {
			if slices.Contains(typeKeys, m[key].(string)) {
				m[key] = Types[m[key].(string)]
			} else {
				return nil, errors.New(fmt.Sprintf(`The declared type %s for %s is not known by dyre`, m[key].(string), key))
			}
		}
		if reflect.TypeOf(m[key]).String() == "map[string]interface {}" {
			var err error
			m[key], err = replaceTypes(m[key].(map[string]interface{}))
			if err != nil {
				return m, err
			}
		}
	}
	return m, nil
}

// if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
// 	masterData[v.Name()] = z.Bool
// 	continue
// }

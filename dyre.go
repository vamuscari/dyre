package dyre

import (
	"errors"
	"fmt"
	"maps"
)

// TODO: sql tools

type DyRe_Field struct {
	name      string
	typeName  string
	required  bool
	sqlSelect string
	tag       map[string]string
}

type DyRe_Group struct {
	_request *DyRe_Request
	name     string
	required bool
	fields   map[string]DyRe_Field
}

type DyRe_SQL struct {
	_request  *DyRe_Request
	tableName string
	sqlFile   string
}

type DyRe_Validated struct {
	_request   *DyRe_Request
	_headers   []string
	_sqlFields []string
	_sqlTypes  []string
	_fields    []DyRe_Field
	_groups    []DyRe_Group
}

type DyRe_Request struct {
	name        string
	requestType string
	fields      map[string]DyRe_Field
	fieldNames  []string
	groups      map[string]DyRe_Group
	groupNames  []string
	sql         DyRe_SQL
}

// TODO: Sub request into sql file

// Validates incoming fields and groups.
// Returns a validated struct for making sql queries.
func (re *DyRe_Request) ValidateRequest(fields []string, groups []string) (DyRe_Validated, error) {
	var selected DyRe_Validated
	for _, re_field := range re.fields {
		if re_field.required == true || contains(fields, re_field.name) {
			selected._fields = append(selected._fields, re_field)
			selected._sqlFields = append(selected._sqlFields, re_field.sqlSelect)
			selected._headers = append(selected._headers, re_field.name)
			selected._sqlTypes = append(selected._sqlTypes, re_field.typeName)
		}
	}
	for _, re_group := range re.groups {
		if re_group.required == true || contains(groups, re_group.name) {
			selected._groups = append(selected._groups, re_group)
			for _, re_gfield := range re_group.fields {
				selected._sqlFields = append(selected._sqlFields, re_gfield.sqlSelect)
				selected._headers = append(selected._headers, re_gfield.name)
				selected._sqlTypes = append(selected._sqlTypes, re_gfield.typeName)
			}
		}
	}

	if len(selected._fields) == 0 && len(selected._groups) == 0 {
		return selected, errors.New(fmt.Sprintf("No valid fields or groups selected for %s", re.name))
	}

	selected._request = re

	return selected, nil
}

func (re *DyRe_Request) FieldNames() []string {
	return re.fieldNames
}

func (re *DyRe_Request) GroupNames() []string {
	return re.groupNames
}

// Fields for SQL Select.
//
// SELECT  {Fields} FROM {Table}
func (valid *DyRe_Validated) SQLFields() []string {
	sqlFields := []string{}
	for _, v := range valid._sqlFields {
		sqlFields = append(sqlFields, v)
	}
	return sqlFields
}

// Returns the list of names for the fields that were calles
func (valid *DyRe_Validated) Headers() []string {
	headers := []string{}
	for _, v := range valid._headers {
		headers = append(headers, v)
	}
	return headers
}

func (re *DyRe_Request) TableName() string {
	return re.sql.tableName
}

func (re *DyRe_Request) SQLFile() string {
	return re.sql.sqlFile
}

func (re *DyRe_Request) Fields() map[string]DyRe_Field {
	return maps.Clone(re.fields)
}

func (re *DyRe_Request) Groups() map[string]DyRe_Group {
	return maps.Clone(re.groups)
}

func (field *DyRe_Field) Name() string {
	return field.name
}

func (field *DyRe_Field) Required() bool {
	return field.required
}

func (field *DyRe_Field) SQLSelect() string {
	return field.sqlSelect
}

func (group *DyRe_Group) Name() string {
	return group.name
}

func (group *DyRe_Group) Required() bool {
	return group.required
}

func (group *DyRe_Group) Fields() map[string]DyRe_Field {
	return maps.Clone(group.fields)
}

func contains(a []string, l string) bool {
	for _, v := range a {
		if l == v {
			return true
		}
	}
	return false
}

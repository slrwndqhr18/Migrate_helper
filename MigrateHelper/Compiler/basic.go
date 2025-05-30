package compiler

import (
	"strings"
)

type Compiler interface {
	Parse_block(string) string
	Analyze_block(string, bool)
	Generate_block(QueryBlock)
}

type TableInfo struct {
	Name string
	Idx  int
}
type ColumnInfo struct {
	Name string
	Type string
	IsPK bool
}
type ConstraintInfo struct {
	Name  string
	Table string
}

type QueryBlock struct {
	QueryTokens    []string
	Action         string
	MainEntityName string
	Constraints    []ConstraintInfo
	Tables         []TableInfo
	Columns        []ColumnInfo
}

type Model struct {
	Blocks     [][]QueryBlock
	EntityName CONFIG_map_names
}
type CONFIG_map_names map[string]NameInfo

type NameInfo struct {
	Schema string `yaml:"schema"`
	Orm    string `yaml:"orm"`
	Ignore bool   `yaml:"ignore"`
}

func (_schma CONFIG_map_names) _get_schema(_tableName string) string {
	if names, isIn := _schma[_tableName]; isIn && names.Schema != "" {
		return names.Schema
	} else if _, isIn = _schma["default"]; !isIn && names.Schema == "" {
		return ""
	} else {
		return _schma["default"].Schema
	}
}

func (_schma CONFIG_map_names) _get_orm(_tableName string) string {
	var ormName string
	if names, isIn := _schma[_tableName]; isIn && names.Orm != "" {
		ormName = names.Orm
	} else if _, isIn = _schma["default"]; !isIn && names.Orm == "" {
		ormName = "SNAKE"
	} else {
		ormName = _schma["default"].Orm
	}
	switch ormName {
	case "SNAKE":
		parts := strings.Split(_tableName, "_")
		for i := 1; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
		return strings.Join(parts, "")
	case "CAMEL":
		return _tableName
	default:
		return ormName
	}
}

func (_schma CONFIG_map_names) Add_schema(_tableName string) string {
	if schema := _schma._get_schema(_tableName); schema != "" {
		return schema + "." + _tableName
	} else {
		return _tableName
	}
}

package compiler

import (
	"fmt"
	"strings"
)

type MSSQL struct {
	Model
}

func (_m *MSSQL) Parse_and_analyze_blocks(_sql string, _nameInfo CONFIG_map_names, _onlyTableName bool) {
	_m.Model = Model{
		Blocks:     [][]QueryBlock{{}, {}},
		EntityName: _nameInfo,
	}
	blocks := _m.Parse_block(_sql)
	for _, block := range blocks {
		_m.Analyze_block(block, _onlyTableName)
	}
}

func (_m MSSQL) Parse_block(_sql string) []string {
	blocks := strings.Split(_sql, "GO")
	cleanBlocks := []string{}
	for _, s := range blocks {
		if s = strings.TrimSpace(s); s != "" {
			cleanBlocks = append(cleanBlocks, s)
		}
	}
	return cleanBlocks
}

func (_m *MSSQL) Analyze_block(_block string, _onlyTableName bool) {
	tokens := strings.Fields(_block)
	blockInfo := QueryBlock{
		Action:      tokens[0] + " " + tokens[1],
		QueryTokens: tokens,
	}
	switch blockInfo.Action {
	case "CREATE TABLE":
		if !_onlyTableName {
			blockInfo.Columns = Get_columns(tokens[2], _block)
		}
		fallthrough
	case "ALTER TABLE":
		blockInfo.MainEntityName = tokens[2]
		blockInfo.Tables, blockInfo.Constraints = _get_table_info(tokens[2], tokens)
		for _, t := range blockInfo.Tables {
			if _m.EntityName[t.Name].Ignore {
				return
			}
		}
		blockInfo.Tables = append(blockInfo.Tables, TableInfo{Name: tokens[2], Idx: 2})
	default:
		return
	}

	if blockInfo.Action == "CREATE TABLE" {
		_m.Model.Blocks[0] = append(_m.Model.Blocks[0], blockInfo)
	} else {
		_m.Model.Blocks[1] = append(_m.Model.Blocks[1], blockInfo)
	}
}

func (_m MSSQL) Generate_block(_block QueryBlock) string {
	for _, table := range _block.Tables {
		_block.QueryTokens[table.Idx] = _m.Model.EntityName.Add_schema(table.Name)
	}
	return strings.Join(_block.QueryTokens, " ")
}

// MSSQL specific functions =======================================
func Get_columns(_tableName string, _sql string) []ColumnInfo {
	_sql = strings.SplitN(_sql, "(", 2)[1]
	lines := strings.Split(_sql, ",\n")
	colList := []ColumnInfo{}
	var lastCheckedLine int
	for _, l := range lines {
		tokens := strings.Fields(strings.TrimSpace(l))
		if tokens[0] == "CONSTRAINT" {
			break
		} else if len(tokens) >= 2 {
			colList = append(colList, ColumnInfo{
				Name: tokens[0],
				Type: tokens[1],
				IsPK: false,
			})
		}
		lastCheckedLine += 1
	}

	PKName := "None"
	for i := lastCheckedLine; i < len(lines); i++ {
		if startIdx := strings.Index(lines[i], "PRIMARY KEY"); startIdx != -1 {
			PKName = lines[i][startIdx+11 : len(lines[i])]
			startIdx = strings.Index(PKName, "(")
			endIdx := strings.Index(PKName, ")")
			PKName = PKName[startIdx+1 : endIdx]
			for i, col := range colList {
				if col.Name == PKName {
					colList[i].IsPK = true
					return colList
				}
			}
			break
		}
	}
	if PKName == "None" {
		panic(fmt.Errorf("no primary key found at table \"%s\"", _tableName))
	} else {
		panic(fmt.Errorf("there are no column named \"%s\" found at table \"%s\"", PKName, _tableName))
	}
}

func _get_table_info(_table string, _tokens []string) ([]TableInfo, []ConstraintInfo) {
	constList := []ConstraintInfo{}
	tableList := []TableInfo{}
	for i, t := range _tokens {
		switch t {
		case "TABLE":
			fallthrough
		case "REFERENCES":
			tableList = append(tableList, TableInfo{
				Name: _tokens[i+1],
				Idx:  i + 1,
			})
		case "CONSTRAINT":
			constList = append(constList, ConstraintInfo{
				Name:  _tokens[i+1],
				Table: _table,
			})

		}
	}
	return tableList, constList
}

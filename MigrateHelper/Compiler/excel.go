package compiler

import (
	handlers "MigrationPG/MigrateHelper/Handlers"
	"fmt"
	"strings"
)

type CONFIG_table_format struct {
	Format []__tableFormat `yaml:"format"`
}

type __tableFormat struct {
	Name   string            `yaml:"name"`
	Src    string            `yaml:"src"`
	Struct map[string]string `yaml:"struct"`
}

type __excelData struct {
	ColIdxSeq map[string]int
	Data      [][]string
}

func (_e *__excelData) Set(_data [][]string) {
	_e.Data = _data[1:]
	_e.ColIdxSeq = make(map[string]int)
	for i, c := range _data[0] {
		_e.ColIdxSeq[c] = i
	}
}
func (_e __excelData) Get_col(_rowIdx int, _colName string) string {
	if i, isIn := _e.ColIdxSeq[_colName]; isIn {
		if len(_e.Data[_rowIdx]) < i {
			return ""
		}
		return _e.Data[_rowIdx][i]
	} else {
		panic(fmt.Errorf("there is no column named \"%s\"", _colName))
	}
}

type Format struct {
	ExcelFilePath string
	ExcelColSeq   []string
	ColTypSeq     []string
	EntityName    CONFIG_map_names
	QueryInsert   string
	Data          map[string]__excelData
}

func (_m *Format) Analyze_block(_block __tableFormat) bool {
	if _m.EntityName[_block.Name].Ignore {
		return false
	} else {
		_m.ExcelFilePath = _block.Src
		_m.ExcelColSeq = []string{}
		_m.ColTypSeq = []string{}
		if _, isIn := _m.Data[_block.Src]; !isIn {
			data := __excelData{}
			data.Set(handlers.Get_data_from_excel(_m.ExcelFilePath))
			_m.Data[_block.Src] = data
		}
		_m.Data[_block.Src] = __excelData{}
		TableColSeq := []string{}
		queryValues := []string{}
		for tbCol, srcCol := range _block.Struct {
			TableColSeq = append(TableColSeq, tbCol)
			queryValues = append(queryValues, "?")
			temp := strings.Split(srcCol, "|")
			_m.ExcelColSeq = append(_m.ExcelColSeq, temp[0])
			_m.ColTypSeq = append(_m.ColTypSeq, temp[1])
		}
		_m.QueryInsert = fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES (%s);",
			_m.EntityName.Add_schema(_block.Name),
			strings.Join(TableColSeq, ","),
			strings.Join(queryValues, ","))
		return true
	}
}

func (_m *Format) Generate_block(_rowIdx int) []any {
	values := []any{}
	for i, col := range _m.ExcelColSeq {
		values = append(values, __get_insert_value(_m.Data[_m.ExcelFilePath].Get_col(_rowIdx, col), _m.ExcelColSeq[i]))
	}
	return values
}

func __get_insert_value(_v string, _typ string) string {
	if _v != "" {
		return _v
	} else {
		switch _typ {
		case "int":
			return "0"
		case "date":
			return ""
		default:
			return "NULL"
		}
	}
}

// func (_m Format) Generate_block(_tableName string, _data [][]string) (query_ string, values_ []interface{}) {
// 	var queryRow []string
// 	var rowSize int
// 	var builder strings.Builder
// 	builder.WriteString(fmt.Sprintf("INSERT INTO %s(%s) VALUES", _m.EntityName.Add_schema(_tableName), strings.Join(_m.Block.TableColSeq, ",")))
// 	for _, row := range _data {
// 		queryRow = []string{}
// 		rowSize = len(row)
// 		for i, colIdx := range _m.Block.DataIdxSeq {
// 			if rowSize > colIdx {
// 				values_ = append(values_, __get_insert_value(row[colIdx], _m.Block.ColTyp[i]))
// 			} else {
// 				values_ = append(values_, __get_insert_value("", _m.Block.ColTyp[i]))
// 			}
// 			queryRow = append(queryRow, "?")
// 		}
// 		builder.WriteString(fmt.Sprintf("(%s),", strings.Join(queryRow, ",")))
// 	}
// 	query_ = strings.TrimSuffix(builder.String(), ",") + ";"
// 	return
// }

// func Split_partitions(_massData [][]string, _partisionSize int) [][][]string {
// 	partitions := [][][]string{}
// 	for i := 0; i < len(_massData); i += _partisionSize {
// 		lastI := i + _partisionSize
// 		if lastI > len(_massData) {
// 			lastI = len(_massData)
// 		}
// 		partitions = append(partitions, _massData[i:lastI])
// 	}
// 	return partitions
// }

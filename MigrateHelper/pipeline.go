package migratehelper

import (
	compiler "MigrationPG/MigrateHelper/Compiler"
	handlers "MigrationPG/MigrateHelper/Handlers"
	"fmt"
	"strings"
)

// ==============================================
// 1. sqlalchemey용

func (_c *Controller) Make_orm_model_file() {
	sqlFile := handlers.Read_sql_file(_c.Conf.FilePathList.SqlFile)
	model := compiler.MSSQL{}
	model.Parse_and_analyze_blocks(sqlFile, _c.Conf.MapTableAndSchema, false)
	orm := compiler.SQLALCHEMY{}
	orm.Model = model.Model

	modelFile, initFile, formatFile := orm.Generate_block_all()

	handlers.Write_to_file(_c.Conf.FilePathList.NewFormatFile, formatFile, true)
	handlers.Write_to_file(_c.Conf.FilePathList.NewPyModelFile, modelFile, true)
	handlers.Write_to_file(_c.Conf.FilePathList.NewPyInitFile, initFile, true)
}

// ==============================================
// 2. 엑셀일 용

func (_c *Controller) Flush_env_on_sql_file() {
	var err error
	var query string
	sqlFile := handlers.Read_sql_file(_c.Conf.FilePathList.SqlFile)
	model := compiler.MSSQL{}
	model.Parse_and_analyze_blocks(sqlFile, _c.Conf.MapTableAndSchema, true)

	_c.DBHandler.Start()
	for _, block := range model.Model.Blocks[1] {
		for _, c := range block.Constraints {
			query = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", model.EntityName.Add_schema(c.Table), c.Name)
			if err = _c.DBHandler.Run_sql(query, nil); err != nil {
				_c.Logger.Add_log("ERROR", "Failed to run query", query, err.Error())
			}
		}
	}
	for _, block := range model.Model.Blocks[0] {
		query = fmt.Sprintf("DROP TABLE %s", model.EntityName.Add_schema(block.MainEntityName))
		if err = _c.DBHandler.Run_sql(query, nil); err != nil {
			_c.Logger.Add_log("ERROR", "Failed to run query", query, err.Error())
		} else {
			_c.Logger.Add_log("INFO", "SUCC to run query", query)
		}
	}
}
func (_c *Controller) Flush_env_on_sql_file_dry_run() {
	var query string
	sqlFile := handlers.Read_sql_file(_c.Conf.FilePathList.SqlFile)
	model := compiler.MSSQL{}
	model.Parse_and_analyze_blocks(sqlFile, _c.Conf.MapTableAndSchema, true)
	for _, block := range model.Model.Blocks[1] {
		for _, c := range block.Constraints {
			query = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", model.EntityName.Add_schema(c.Table), c.Name)
			fmt.Println("Query:", query)
		}
	}
	for _, block := range model.Model.Blocks[0] {
		query = fmt.Sprintf("DROP TABLE %s", model.EntityName.Add_schema(block.MainEntityName))
		fmt.Println("Query:", query)
	}
	fmt.Printf("Total executed queries: %d (Constraints: %d, tables: %d)\n", len(model.Model.Blocks[1])+len(model.Model.Blocks[0]), len(model.Model.Blocks[1]), len(model.Model.Blocks[0]))
}

func (_c *Controller) Migrate_by_sql_file() {
	var query string
	sqlFile := handlers.Read_sql_file(_c.Conf.FilePathList.SqlFile)
	model := compiler.MSSQL{}
	model.Parse_and_analyze_blocks(sqlFile, _c.Conf.MapTableAndSchema, true)
	_c.DBHandler.Start()
	var TCLres string
	var err error
	for _, block := range model.Model.Blocks[0] {
		query = model.Generate_block(block)
		if err = _c.DBHandler.Run_sql(query, nil); err != nil {
			_c.Logger.Add_log("ERROR", "Failed to run query", query, err.Error())
		}
	}
	TCLres = _c.DBHandler.End_tx(_c.Logger.Is_ok(), true)
	_c.Logger.Add_log("INFO", fmt.Sprintf("TCL result after First quries: %s", TCLres))
	if _c.Logger.Is_ok() {
		for _, block := range model.Model.Blocks[1] {
			query = model.Generate_block(block)
			if err = _c.DBHandler.Run_sql(query, nil); err != nil {
				_c.Logger.Add_log("ERROR", "Failed to run query", query, err.Error())
			}
		}
		TCLres = _c.DBHandler.End_tx(_c.Logger.Is_ok(), false)
		_c.Logger.Add_log("INFO", fmt.Sprintf("TCL result after Second quries: %s", TCLres))
	}
}

func (_c *Controller) Migrate_by_sql_file_dry_run() {
	sqlFile := handlers.Read_sql_file(_c.Conf.FilePathList.SqlFile)
	model := compiler.MSSQL{}
	model.Parse_and_analyze_blocks(sqlFile, _c.Conf.MapTableAndSchema, true)

	exeCnt := 0
	for _, block := range model.Model.Blocks[0] {
		fmt.Println("Query:", model.Generate_block(block))
		exeCnt += 1
	}
	for _, block := range model.Model.Blocks[1] {
		fmt.Println("Query:", model.Generate_block(block))
		exeCnt += 1
	}
	fmt.Printf("Total executed queries: %d (First: %d, Second: %d)\n", exeCnt, len(model.Model.Blocks[0]), len(model.Model.Blocks[1]))
}

// ==============================================
// 3. 엑셀일 용

func (_c *Controller) Migrate_by_excel_file() {
	var formatBlocks compiler.CONFIG_table_format = compiler.CONFIG_table_format{}
	handlers.Read_yaml_file(_c.Conf.FilePathList.FormatFile, &formatBlocks)
	_c.DBHandler.Start()
	Model := compiler.Format{
		EntityName: _c.Conf.MapTableAndSchema,
	}
	for _, tableFormat := range formatBlocks.Format {
		tableName := _c.Conf.MapTableAndSchema.Add_schema(tableFormat.Name)
		if isSkip, err := _c.DBHandler.Check_is_data(tableName); isSkip {
			if err != nil {
				_c.Logger.Add_log("ERROR", fmt.Sprintf("Skip inserting Table '%s' - Cant't check is data", tableName), err.Error())
			} else {
				_c.Logger.Add_log("INFO", fmt.Sprintf("Skip inserting Table '%s' - alread has data", tableName))
			}
		} else {
			if isProgress := Model.Analyze_block(tableFormat); !isProgress {
				continue
			}
			fmt.Printf("[INFO] Inserting data to table \"%s\"\n", tableName)
			for rowIdx, _ := range Model.Data[""].Data {
				queryValues := Model.Generate_block(rowIdx)
				if err = _c.DBHandler.Run_sql(Model.QueryInsert, queryValues); err != nil {
					_c.Logger.Add_log("ERROR", "Failed to run query", err.Error(), Model.QueryInsert, _conv_arr_to_string(queryValues))
					break
				}
			}
			_c.DBHandler.End_tx(_c.Logger.Is_ok(), true)
		}
	}
}

func _conv_arr_to_string(_arr []interface{}) string {
	temp := make([]string, len(_arr))
	for i, v := range _arr {
		temp[i] = fmt.Sprint(v)
	}
	return strings.Join(temp, ", ")
}

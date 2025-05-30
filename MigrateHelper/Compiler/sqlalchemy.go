package compiler

import (
	"fmt"
	"strings"
)

type SQLALCHEMY struct {
	Model
}

func (_m SQLALCHEMY) Generate_block_all() (string, string, string) {
	var modelBuilder, initFileBuilder, formatBuilder strings.Builder
	modelBuilder.WriteString("from sqlalchemy import Table, Column, Integer, String, MetaData, create_engine\nmetadata = MetaData()\n")
	initFileBuilder.WriteString("from import (\n")
	formatBuilder.WriteString("format:\n")
	for _, block := range _m.Model.Blocks[0] {
		modelBuilder.WriteString(_m.Generate_block(block))
		initFileBuilder.WriteString(fmt.Sprintf("%s,\n", _m.EntityName._get_orm(block.MainEntityName)))
		formatBuilder.WriteString(_m.Generate_format_block(block))
	}
	initFileBuilder.WriteString(")")
	return modelBuilder.String(), initFileBuilder.String(), formatBuilder.String()
}

func (_m SQLALCHEMY) Generate_block(_block QueryBlock) string {
	if _block.Action == "CREATE TABLE" {
		var modelBuilder strings.Builder
		modelBuilder.WriteString(fmt.Sprintf("%s = Table(\n\"%s\",\nmetadata,\n", _m.EntityName._get_orm(_block.MainEntityName), _block.MainEntityName))
		for _, property := range _block.Columns {
			if property.IsPK {
				modelBuilder.WriteString(fmt.Sprintf("Column(\"%s\", %s, primary_key=True),\n", property.Name, __conv_col_typ_to_orm(property.Type)))
			} else {
				modelBuilder.WriteString(fmt.Sprintf("Column(\"%s\", %s),\n", property.Name, __conv_col_typ_to_orm(property.Type)))
			}
		}
		if schema := _m.EntityName._get_schema(_block.MainEntityName); schema != "" {
			modelBuilder.WriteString(fmt.Sprintf("schema=\"%s\"\n)\n", schema))
			return modelBuilder.String()
		} else {
			return strings.TrimSuffix(modelBuilder.String(), ",\n") + "\n)\n"
		}
	}
	return ""
}

func (_m SQLALCHEMY) Generate_format_block(_block QueryBlock) string {
	if _block.Action == "CREATE TABLE" {
		var modelBuilder strings.Builder
		modelBuilder.WriteString(fmt.Sprintf("- name: %s\n    src: .xlsx\n    struct:\n", _m.Model.EntityName.Add_schema(_block.MainEntityName)))
		for _, col := range _block.Columns {
			modelBuilder.WriteString(fmt.Sprintf("      %s: %s|%s\n", col.Name, col.Name, __conv_col_typ_to_format(col.Type)))
		}
		return modelBuilder.String()
	}
	return ""
}

func __conv_col_typ_to_orm(_typ string) string {
	switch strings.ToLower(_typ) {
	case "int":
		fallthrough
	case "bit":
		fallthrough
	case "numeric":
		return "Integer"
	default:
		return "String"
	}
}

func __conv_col_typ_to_format(_typ string) string {
	switch strings.ToLower(_typ) {
	case "int":
		fallthrough
	case "bit":
		fallthrough
	case "numeric":
		return "int"
	case "datetime":
		return "date"
	default:
		return "str"
	}
}

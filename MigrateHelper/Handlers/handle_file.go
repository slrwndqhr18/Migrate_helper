package handlers

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

func Read_yaml_file[T any](_filePath string, _formatData *T) {
	yamlFile, err := os.ReadFile(_filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, _formatData)
	if err != nil {
		panic(err)
	}
}

func Get_data_from_excel(_filePath string) [][]string {
	f, err := excelize.OpenFile(_filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sheets := f.GetSheetList()
	rows, _ := f.GetRows(sheets[0])
	return rows
}

func Read_sql_file(_filePath string) string {
	byteData, err := os.ReadFile(_filePath)
	if err != nil {
		panic(err)
	}
	queries := string(byteData)
	if queries[0] == '\n' {
		queries = queries[1:]
	}
	return queries
}

func Write_to_file(_filePath string, _str string, _override bool) {
	var flag int
	if _override {
		flag = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(_filePath, flag, 0644)
	defer func() {
		if err != nil {
			fmt.Println(_str)
			fmt.Println("[ERROR] Failed to write log at file\n", err.Error())
			fmt.Println("\tㄴ", err.Error())
		}
	}()
	if err == nil {
		defer f.Close()
		_, err = f.WriteString(_str)
	}
}

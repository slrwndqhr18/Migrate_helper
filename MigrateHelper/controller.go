package migratehelper

import (
	compiler "MigrationPG/MigrateHelper/Compiler"
	handlers "MigrationPG/MigrateHelper/Handlers"
	"fmt"
)

type TableColMappingInfo map[string]map[string]string

type _CONFIG_file_path struct {
	LogDir         string `yaml:"logDir"`
	MigrateData    string `yaml:"migrateData"`
	FormatFile     string `yaml:"formatFile"`
	SqlFile        string `yaml:"sqlFile"`
	NewFormatFile  string `yaml:"newFormatFile"`
	NewPyModelFile string `yaml:"newPyModelFile"`
	NewPyInitFile  string `yaml:"newPyInitFile"`
}

type CONFIG struct {
	FilePathList      _CONFIG_file_path         `yaml:"filePathList"`
	DatabaseInfo      handlers.CONFIG_db_info   `yaml:"databaseInfo"`
	MapTableAndSchema compiler.CONFIG_map_names `yaml:"mapTableAndSchema"`
	AllowedQueryError []string                  `yaml:"allowedQueryError"`
}

type Controller struct {
	Conf      CONFIG
	Logger    handlers.LOG_CTRL
	DBHandler handlers.DB_CTRL
}

func Init(_configFilePath string) Controller {
	fmt.Println("[START] System initiate")
	c_ := CONFIG{}
	handlers.Read_yaml_file("./files/config.yaml", &c_)
	o := Controller{
		Conf: c_,
	}
	o.Logger = handlers.Init_logger(o.Conf.FilePathList.LogDir)
	o.DBHandler = handlers.Init_database_controller(
		o.Conf.DatabaseInfo.Type,
		o.Conf.DatabaseInfo.ConnInfo,
		o.Conf.AllowedQueryError,
		o.Conf.DatabaseInfo.RunAs,
	)
	fmt.Println("[INTO] Starting sub process ...")
	return o
}

func (_c *Controller) Closure() {
	fmt.Println("[INFO] You got '", _c.Logger.ErrCnt, "' Errors")
	if _c.DBHandler.Conn.DB != nil {
		_c.DBHandler.Conn.DB.Close()
	}
	_c.Logger.Write_to_file()
	fmt.Print("[END] System has been successfully terminated, Have a nice day.\n\n")
	fmt.Println("--------------------------[Panic Messages]--------------------------")
	fmt.Print("If panic occured, the messages will be showen below.\n\n")
}

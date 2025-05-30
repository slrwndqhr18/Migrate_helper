package handlers

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

type CONFIG_db_info struct {
	Type          string `yaml:"type"`
	RunAs         string `yaml:"runAs"`
	Db_and_schema string `yaml:"db_and_schema"`
	ConnInfo      string `yaml:"connInfo"`
	MaxInsertRows int    `yaml:"maxInsertRows"`
}
type __DB_CONN_OBJ struct {
	DB *sql.DB
	Tx *sql.Tx
}
type __conn_info struct {
	DBType     string
	ConnString string
}
type DB_CTRL struct {
	Map_table_schema map[string]string
	Conn             __DB_CONN_OBJ
	AllowedError     []string
	CONNInfo         __conn_info
	IsTransaction    bool
}

func Init_database_controller(_dbType string, _connString string, _allowedErrList []string, _mode string) DB_CTRL {
	ctrl := DB_CTRL{
		Conn:          __DB_CONN_OBJ{DB: nil, Tx: nil},
		AllowedError:  _allowedErrList,
		CONNInfo:      __conn_info{DBType: _dbType, ConnString: _connString},
		IsTransaction: false,
	}
	dbMode := "Transaction"
	switch _mode {
	case "dryRun":
		dbMode = "DryRun"
	case "noTransaction":
		dbMode = "No Transaction"
	default:
		ctrl.IsTransaction = true
	}
	fmt.Printf("[INTO] Run as '%s' mode\n", dbMode)
	return ctrl
}

func (_Db *DB_CTRL) Start() {
	var err error
	if _Db.Conn.DB == nil {
		if _Db.Conn.DB, err = sql.Open(_Db.CONNInfo.DBType, _Db.CONNInfo.ConnString); err != nil {
			panic(err)
		}
		fmt.Println("[INFO] DB Connection opened")
	}
	if _Db.IsTransaction {
		if _Db.Conn.Tx, err = _Db.Conn.DB.Begin(); err != nil {
			panic(err)
		}
		fmt.Println("[INFO] DB Transaction started")
	}
}
func (_Db *DB_CTRL) End_tx(_isCommit bool, _isRestartTx bool) (finalAction_ string) {
	if _Db.IsTransaction {
		var err error
		if _isCommit {
			err = _Db.Conn.Tx.Commit()
			finalAction_ = "commit"
		} else {
			err = _Db.Conn.Tx.Rollback()
			finalAction_ = "rollback"
		}
		if err != nil {
			panic(err)
		}
		if _isRestartTx && _isCommit {
			_Db.Start()
		} else {
			_Db.Conn.Tx = nil
		}
	}
	return
}

func (_Db *DB_CTRL) Check_is_data(_tbName string) (bool, error) {
	var res int
	err := _Db.Conn.Tx.QueryRow(fmt.Sprintf("SELECT TOP 1 1 FROM %s;", _tbName)).Scan(&res)
	if err != nil && _Db.__is_allowed_err(err) {
		err = nil
	}
	return res == 1 || err != nil, err
}

func (_Db *DB_CTRL) Run_sql(_query string, _values []any) (err_ error) {
	if _Db.Conn.Tx == nil {
		if _values == nil {
			_, err_ = _Db.Conn.DB.Exec(_query)
		} else {
			_, err_ = _Db.Conn.DB.Exec(_query, _values...)
		}
	} else {
		if _values == nil {
			_, err_ = _Db.Conn.Tx.Exec(_query)
		} else {
			_, err_ = _Db.Conn.Tx.Exec(_query, _values...)
		}
	}
	if err_ != nil {
		if _Db.__is_allowed_err(err_) {
			err_ = nil
		}
	}
	return
}

func (_Db DB_CTRL) __is_allowed_err(_err error) bool {
	errMesg := _err.Error()
	for _, e := range _Db.AllowedError {
		if len(errMesg) >= len(e) && errMesg[:len(e)] == e {
			return true
		}
	}
	return false
}

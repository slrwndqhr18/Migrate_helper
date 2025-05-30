//go:debug x509negativeserial=1
package main

import (
	migratehelper "MigrationPG/MigrateHelper"
	"os"
)

type MainController migratehelper.Controller

func main() {
	o := migratehelper.Init("./files/config.yaml")
	defer o.Closure()

	argActionTyp := os.Args[1] // 0은 경로
	switch argActionTyp {
	case "MAKE_DB":
		if o.Conf.DatabaseInfo.RunAs == "dryRun" {
			o.Migrate_by_sql_file_dry_run()
		} else {
			o.Migrate_by_sql_file()
		}
	case "FLUSH_ENV":
		if o.Conf.DatabaseInfo.RunAs == "dryRun" {
			o.Flush_env_on_sql_file_dry_run()
		} else {
			o.Flush_env_on_sql_file()
		}
	case "MAKE_SQLALCHEMY":
		o.Make_orm_model_file()
	case "MIGRATE_EXCEL":
		o.Migrate_by_excel_file()
	}
	//Flush_env_on_sql_file(&o)
}

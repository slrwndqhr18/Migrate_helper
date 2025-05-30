# Migrate_helper
A Go-based Excel-to-Database & (MSSQL)SQL file-to-Database Migration Tool  

I originally created this for a flexible migration pipeline.
Currently, the code is designed to work with the "ERD Editor" extension for VS Code, and it can be quite useful when used together.

While "ERD Editor" is a great tool, it can be tricky at times—such as adding schema names in front of table names in MSSQL.

That's why I decided to develop this code: to support more use cases, even when "ERD Editor" isn't used.

## Features

In my project, I used MSSQL, excel, "ERD Editor".
It's not exclusive, but it is specialized to .SQL files made automatically by "ERD Editor".

There are now 4 main functions, defined in /MigrateHelper/pipeline.go
- Make_orm_model_file   : Make sqlalchemy.core "Table" definitions based on "CREATE TABLE" queries in .SQL file
- Flush_env_on_sql_file : Remove all DB objects in running DB, based on .SQL file
- Migrate_by_sql_file   : Make DB objects(mostly tables) based on .SQL file
- Migrate_by_excel_file : Migrate excel data to DB

## Project Structure
├── controller.go         // Main object <br/>
├── pipeline.go           // defines all main functions. You only need to use functions in this file.<br/>
├── Template_config.yaml  // example of pipeline configuration file<br/>
├── Template_config.yaml  // example of format file (Like mssql format file). Maps excel - DB<br/>
├── Handlers<br/>         // Utility <br/>
│ ├── handle_db.go        // handles only DB related functions (ex DB connection)<br/>
│ └── handle_files.go     // file I/O related functions<br/>
│ └── handle_logs.go      // logger (recored & make log files)<br/>
├── Compiler<br/>         // Defines model <br/>
│ └── excel.go            // Codes for handle excel file<br/>
│ └── mssql.go            // Codes for handle MSSQL SQL file<br/>
│ └── sqlalchemy.go       // Codes for generate sqlalchemy codes automatically

## How It Works

1. Define configurations by defining .yaml file (example "/Migrate_helper/Template_config.yaml")
   - If your using "Migrate_by_excel_file", you need to define format file (example "/Migrate_helper/Template_format.yaml")
2. run "go run main <OPTIONS>
   - OPTIONS: MAKE_DB, FLUSH_ENV, MAKE_SQLALCHEMY, MIGRATE_EXCEL
   - you can see the usage at "example.go"
3. To add more models like MySQL or sqlalchemy.orm, define codes at "/MigrateHelper/Compiler

## Getting Started

```bash
git clone https://github.com/slrwndqhr18/Migrate_helper.git
cd Migrate_helper
go run main.go <OPTION>
```
Feel free to customize. Feel free to share any feedback or report issues anytime !!! <br/>
Feel free to [open an issue](https://github.com/slrwndqhr18/Migrate_helper/issues) <br/>
or start a [discussion](https://github.com/slrwndqhr18/Migrate_helper/discussions) <br/>
You don't need write access to participate — everyone is welcome!

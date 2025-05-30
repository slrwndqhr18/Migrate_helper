# Migrate_helper
A Go-based Excel-to-Database & SQL-to-Database Migration Tool  
Migrate_helper is a lightweight utility that reads structured Excel files or SQL files and inserts their contents into a relational database based on predefined configurations.

I made this for flexible migration pipeline. 
Now, this code is for VScode extention "ERD Editor". It will be usefull with this extention.
But I will make this code for every cases(in case for not using "ERD Editor")

## Features

- Parses Excel (.xlsx) files
- Dynamically generates `INSERT` queries based on column order and types
- Uses predefined mapping to associate Excel columns with database fields
- Supports multiple table/entity mapping
- Block-based parsing for repeated data sections
- Uses YAML-based configuration for flexible mapping

## Project Structure
├── handlers
│ ├── excel.go // Excel parsing logic
│ └── utils.go // Utility functions
├── format
│ └── format.go // Core data structures (Format, __excelData, etc.)
├── config
│ └── map_config.yaml // Column-to-entity mapping configuration

## How It Works

1. Define column-to-entity mappings in `config/map_config.yaml`.
2. Prepare Excel files matching the expected format.
3. Run the tool with Go, and it will:
   - Read the Excel data
   - Match it to the defined entities
   - Generate and execute appropriate SQL insert statements

## Getting Started

```bash
git clone https://github.com/slrwndqhr18/Migrate_helper.git
cd Migrate_helper
go run main.go

## Requirements
Go 
Excel file with a known structure

SQL driver/configuration (implementation depends on usage context)

Feel free to customize or expand based on the features and usage patterns you implement!

package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func connectToPostgres() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file, using environment values")
	}

	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DB")
	port := os.Getenv("PG_PORT")
	sslmode := os.Getenv("PG_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("✅ Successfully connected to PostgreSQL")
	return db, nil
}

func getColumns(db *sql.DB, tableName string) ([]dbItem, error) {
	query := `
SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = $1`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []dbItem
	for rows.Next() {
		var colName, colType string
		if err := rows.Scan(&colName, &colType); err != nil {
			return nil, err
		}
		columns = append(columns, dbItem{
			name: colName + "\n" + fmt.Sprintf("(%s)", colType),
			kind: "column",
		})
	}

	return columns, nil
}

func getTables(db *sql.DB) ([]dbItem, error) {
	query := `SELECT tablename
                FROM pg_tables
              WHERE schemaname NOT IN ('pg_catalog', 'information_schema')`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []dbItem
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, dbItem{name: tableName, kind: "tables"})
	}
	return tables, nil
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgresql://acm-vit:8rKbDJw-ZkiBqbiHaNEBiQ@rc-2024-3929.7s5.aws-ap-south-1.cockroachlabs.cloud:26257/senior-citizen?sslmode=verify-full"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Fetch all tables in the public schema
	rows, err := db.Query(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public';
	`)
	if err != nil {
		log.Fatalf("Failed to fetch tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			log.Fatalf("Failed to scan table: %v", err)
		}
		tables = append(tables, table)
	}

	if len(tables) == 0 {
		fmt.Println("No tables found in public schema.")
		return
	}

	// Drop all tables
	for _, table := range tables {
		dropStmt := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", table)
		if _, err := db.Exec(dropStmt); err != nil {
			log.Fatalf("Failed to drop table %s: %v", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	fmt.Println("All tables in public schema have been dropped.")
}

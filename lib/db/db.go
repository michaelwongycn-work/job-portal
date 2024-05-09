package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func tableExists(db *sql.DB, tableName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pg_tables WHERE tablename = $1", tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createTableIfNotExists(db *sql.DB, tableName string, schema string) error {
	exists, err := tableExists(db, tableName)
	if err != nil {
		return err
	}
	if !exists {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
		log.Printf("Created Table %s\n", tableName)
	}
	return nil
}

func Connect(timeout time.Duration, dbname, host, port, user, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB\n", err)
		return nil, err
	}

	log.Printf("Connected to DB %s successfully\n", dbname)

	tables := []struct {
		name   string
		schema string
	}{
		{usersTable, usersTableSchema},
		{userTokensTable, userTokensTableSchema},
		{jobsTable, jobsTableSchema},
		{applicationsTable, applicationsTableSchema},
	}

	for _, table := range tables {
		err = createTableIfNotExists(db, table.name, table.schema)
		if err != nil {
			log.Printf("Error creating table %s: %s\n", table.name, err)
			return nil, err
		}
	}

	return db, nil
}

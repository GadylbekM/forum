package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Migrator()
}

func Migrator() {
	fmt.Println("Starting Migrations")
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = Migrate("./forum/migrator/up.sql", db)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Migrate(filepath string, db *sql.DB) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	scanner := bufio.NewScanner(f)
	statement := strings.Builder{}
	for scanner.Scan() {
		line := scanner.Text()
		statement.WriteString(line)
		if strings.Contains(line, ";") {
			_, err := db.Exec(statement.String())
			if err != nil {
				return err
			}
			statement.Reset()
		}
	}
	return nil
}

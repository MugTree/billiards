package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MugTree/billiards/app/db"

	_ "github.com/mattn/go-sqlite3"
)

// call this like
// go run ./app/db/seed/generate.go --urls=./app/db/seed/seed.csv --db=./billiards.db

func main() {

	filePtr := flag.String("urls", "", "the file to get the urls from - needs to be broken over lines")
	dbPtr := flag.String("db", "", "path to the db")

	flag.Parse()
	fmt.Println("urls:", *filePtr)

	f, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	people := []db.Person{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")

		fi := db.Person{
			Name: parts[0],
		}

		people = append(people, fi)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	dbhandle, err := sql.Open("sqlite3", *dbPtr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer dbhandle.Close()
	queries := db.New(dbhandle)

	ctx := context.Background()

	for _, p := range people {
		err := queries.InsertPerson(ctx, p.Name)
		if err != nil {
			log.Fatalf("error running log:  %v", err)
		}
	}
}

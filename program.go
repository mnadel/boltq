package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/boltdb/bolt"
	"github.com/mnadel/boltq/boltq"
)

var dbFile = flag.String("db", "", "boltq database file")
var createTest = flag.Bool("create-test", false, "create a test database")
var verbose = flag.Bool("verbose", false, "enable verbosity")
var debug = flag.Bool("debug", false, "enable extreme verbosity")

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	flag.Parse()

	if *verbose {
		log.SetLevel(log.InfoLevel)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if *dbFile == "" {
		log.Fatalln("missing: -db")
	}

	if *createTest {
		db, err := bolt.Open(*dbFile, 0600, nil)
		if err != nil {
			log.Fatalln(err)
		}
		db.Close()
		os.Exit(0)
	}

	query := flag.Args()
	if len(query) == 0 || query[0] == "" {
		log.Fatalln("usage: boltq -db [db file] query...")
	}

	db, err := bolt.Open(*dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = executeQuery(db, strings.Join(flag.Args(), " "))
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func executeQuery(db *bolt.DB, query string) error {
	parser := boltq.NewParser(strings.NewReader(query))

	selectStatement, err := parser.ParseSelect()
	if err == nil {
		log.Debugf("parsed select: %v", selectStatement)
		return executeSelect(selectStatement, db)
	}

	return fmt.Errorf("cannot parse: %s", query)
}

func executeSelect(stmt *boltq.SelectStatement, db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket

		for _, name := range stmt.BucketPath {
			log.Debugln("navigating to bucket", name)
			bucket = tx.Bucket([]byte(name))

			if bucket == nil {
				return fmt.Errorf("cannot find bucket %s", name)
			}
		}

		return nil
	})
}

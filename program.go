package main

import (
	"bytes"
	"encoding/gob"
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
	selectStatement, err := boltq.NewParser(strings.NewReader(query)).ParseSelect()
	if err == nil {
		log.Debugf("parsed select: %v", selectStatement)
		return executeSelect(selectStatement, db)
	} else {
		log.Debugln(err.Error())
	}

	updateStatement, err := boltq.NewParser(strings.NewReader(query)).ParseUpdate()
	if err == nil {
		log.Debugf("parsed update: %v", updateStatement)
		return executeUpdate(updateStatement, db)
	} else {
		log.Debugln(err.Error())
	}

	deleteStatement, err := boltq.NewParser(strings.NewReader(query)).ParseDelete()
	if err == nil {
		log.Debugf("parsed delete: %v", deleteStatement)
		return executeDelete(deleteStatement, db)
	} else {
		log.Debugln(err.Error())
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

		if containsAsterisk(stmt.Fields) {
			log.Debugln("interating keys")
			cursor := bucket.Cursor()

			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				emitKeypair(k, v)
			}
		} else {
			for _, k := range stmt.Fields {
				keyBytes := []byte(k)
				v := bucket.Get(keyBytes)
				emitKeypair(keyBytes, v)
			}
		}

		return nil
	})
}

func executeUpdate(stmt *boltq.UpdateStatement, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		var err error

		for _, name := range stmt.BucketPath {
			log.Debugln("navigating to bucket", name)
			bucket, err = tx.CreateBucketIfNotExists([]byte(name))

			if err != nil {
				return err
			}

			if bucket == nil {
				return fmt.Errorf("cannot find bucket %s", name)
			}
		}

		for k, v := range stmt.Fields {
			log.Debugf("putting %s -> %v", k, v)

			b, err := encode(v)
			if err != nil {
				return err
			}

			if err = bucket.Put([]byte(k), b); err != nil {
				return err
			}
		}

		return nil
	})
}

func executeDelete(stmt *boltq.DeleteStatement, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
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

func containsAsterisk(vals []string) bool {
	for _, v := range vals {
		if v == "*" {
			return true
		}
	}

	return false
}

func emitKeypair(k []byte, v []byte) {
	key := string(k)
	val, err := decode(v)
	if err != nil {
		fmt.Printf("%v -> cannot decode (%s)", key, err.Error())
	} else {
		fmt.Printf("%v -> %v\n", key, val)
	}
}

func decode(b []byte) (interface{}, error) {
	dec := gob.NewDecoder(bytes.NewReader(b))

	var intVal int
	err := dec.Decode(&intVal)
	if err == nil {
		return intVal, nil
	}

	var floatVal float64
	err = dec.Decode(&floatVal)
	if err == nil {
		return floatVal, nil
	}

	return strings.TrimSpace(string(b)), nil
}

func encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package boltq

import (
	"flag"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/medcl/mysqlbeat/beater"
)

func main() {
	err := beat.Run("mysqlbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}

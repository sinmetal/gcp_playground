package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sinmetal/gcp_playground/bigtable"
)

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fmt.Fprintf(w, "Hello, Google Cloud Platform Playground")
	end := time.Now()

	fmt.Printf("duration=%dns\n", end.Sub(start).Nanoseconds())
}

func main() {
	project := flag.String("project", "", "The Google Cloud Platform project ID. Required.")
	bigtableInstance := flag.String("bigtableInstance", "", "The Google Cloud Bigtable instance ID. Required.")
	flag.Parse()

	for _, f := range []string{"project"} {
		if flag.Lookup(f).Value.String() == "" {
			log.Fatalf("The %s flag is required.", f)
		}
	}

	bigtable.SetUp(*project, *bigtableInstance)

	http.HandleFunc("/", handler)
	http.HandleFunc("/bigtable", bigtable.HandlerBigtable)

	http.ListenAndServe(":8080", nil)
}

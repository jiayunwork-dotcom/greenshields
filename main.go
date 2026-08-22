// Command greenshields serves the Greenshields fundamental-diagram calculator
// over HTTP and can also print a single capacity figure from an example file
// without starting a server.
//
// Examples:
//
//	# Start the web server and JSON API on port 8080.
//	go run . -http :8080
//
//	# Print the capacity point derived from example/freeway.json.
//	go run . -print
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"greenshields/internal/api"
	"greenshields/internal/core"
	"greenshields/internal/example"
)

func main() {
	httpAddr := flag.String("http", ":8080", "HTTP listen address for the API and web UI")
	exampleDir := flag.String("example", "example", "directory containing freeway.json")
	webDir := flag.String("web", "web", "directory with static web assets")
	printExample := flag.Bool("print", false, "load example/freeway.json, print capacity, then exit")
	flag.Parse()

	if *printExample {
		if err := runPrint(*exampleDir); err != nil {
			log.Fatalf("greenshields: %v", err)
		}
		return
	}

	srv := api.NewServer(*exampleDir)
	mux := http.NewServeMux()
	srv.Register(mux, *webDir)

	fmt.Printf("greenshields: server listening on %s\n", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, mux); err != nil {
		log.Fatalf("greenshields: server error: %v", err)
	}
}

// runPrint loads the example file, builds the model, and prints the capacity
// point. It is the offline, no-server entry point used in documentation.
func runPrint(exampleDir string) error {
	path := filepath.Join(exampleDir, "freeway.json")
	p, err := example.LoadFile(path)
	if err != nil {
		return fmt.Errorf("cannot load example %s: %w", path, err)
	}
	m, err := core.New(p.Vf, p.Kj)
	if err != nil {
		return err
	}
	qm, km := m.Capacity()
	fmt.Printf("example %q\n", p.Name)
	fmt.Printf("  vf = %.4g  kj = %.4g\n", p.Vf, p.Kj)
	fmt.Printf("  capacity: qm = %.4g  at km = %.4g  (vm = %.4g)\n", qm, km, m.SpeedAtCapacity())
	fmt.Printf("  roots identity: kf + kc = %.4g (= kj)\n", km*2)
	return nil
}

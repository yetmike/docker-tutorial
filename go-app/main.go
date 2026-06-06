// go-app — Docker Tutorial sample (Multi-stage build + Compose pattern)
//
// Demonstrates:
//   - A two-stage Dockerfile: golang builder → distroless runtime
//   - The final image contains ONLY the compiled binary (no shell, no Go toolchain)
//   - Docker Compose wiring the Go app to a Postgres database
//   - Service-name DNS: connects to "db" (the Postgres service name in compose.yaml)
//
// Run with Docker Compose:
//
//	docker compose up --build
//	curl http://localhost:8080
//	curl http://localhost:8080/health
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func dbURL() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:secret@localhost:5432/app?sslmode=disable"
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", dbURL())
		if err != nil {
			http.Error(w, "DB open error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var version string
		if err := db.QueryRow("SELECT version()").Scan(&version); err != nil {
			fmt.Fprintf(w, "<h1>Hello from Docker (Go)!</h1><p>Could not reach Postgres: %v</p>", err)
			return
		}
		fmt.Fprintf(w,
			"<h1>Hello from Docker (Go)!</h1>"+
				"<p>Connected to Postgres successfully.</p>"+
				"<p><em>%s</em></p>"+
				"<p>This binary lives in a <strong>distroless</strong> image — no shell, no package manager.</p>",
			version,
		)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", dbURL())
		dbStatus := "ok"
		if err == nil {
			if pingErr := db.Ping(); pingErr != nil {
				dbStatus = "error: " + pingErr.Error()
			}
			db.Close()
		} else {
			dbStatus = "error: " + err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "ok",
			"database": dbStatus,
		})
	})

	addr := ":8080"
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

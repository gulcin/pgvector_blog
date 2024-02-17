package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to PostgreSQL database
	connStr := "user=gulcin.jelinek dbname=postgres host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Fetch extensions that are similar to pgvector based on their descriptions
	query := `
		WITH pgv AS (
			SELECT embedding
			  FROM document_embeddings JOIN documents USING (id)
			 WHERE title = 'pgvector'
		)
		SELECT title, content
		  FROM document_embeddings
		  JOIN documents USING (id)
		 WHERE embedding <-> (SELECT embedding FROM pgv) < 0.5;
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print results in a nice format
	for rows.Next() {
		var docTitle, docContent string
		err := rows.Scan(&docTitle, &docContent)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Document title: %s\n", docTitle)
		fmt.Printf("Document text: %s\n", docContent)
		fmt.Println()
	}
}

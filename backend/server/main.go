package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	api "social-network/app/router"
	middlewares "social-network/middleware"
	"social-network/pkg/db/sqlite"
)

var Db *sql.DB

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	if err := sqlite.ApplyMigrations(db); err != nil {
		panic("Migration failed: " + err.Error())
	}

	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	userIDs := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		userIDs = append(userIDs, id)
	}

	if len(userIDs) == 0 {
		log.Fatal("No users found in DB")
	}

	// 2️⃣ Shuffle users to get random different owners
	rand.Shuffle(len(userIDs), func(i, j int) { userIDs[i], userIDs[j] = userIDs[j], userIDs[i] })

	// 3️⃣ Insert 50 posts with random images (like stories)
	postCount := 0
	for i := 0; i < 50 && i < len(userIDs); i++ {
		userID := userIDs[i]
		postID := fmt.Sprintf("post-%03d", i+1)
		title := fmt.Sprintf("Post Title %d", i+1)
		content := fmt.Sprintf("This is the content for post %d", i+1)

		// Random image like stories
		imageURL := fmt.Sprintf("https://picsum.photos/400/600?random=%d", rand.Intn(1000))

		visibility := "public"

		_, err := db.Exec(`
			INSERT INTO posts (id, user_id, title, content, visibility, image_path)
			VALUES (?, ?, ?, ?, ?, ?)`,
			postID, userID, title, content, visibility, imageURL)
		if err != nil {
			log.Println("Error inserting post:", postID, err)
		} else {
			postCount++
		}
	}

	fmt.Printf("✅ %d posts inserted successfully with random images!\n", postCount)

	baseHandler := api.Routes()

	// Wrap the API routes with CORS
	handler := enableCORS(middlewares.SessionMiddleware(sqlite.Db, baseHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error in starting of server:", err)
		db.Close()
		return
	}
}

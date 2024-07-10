package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var db *sql.DB

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Isbn     string    `json:"isbn"`
	Director *Director `json:"director"`
}

func initDB() {
	var err error
	connStr := "user=postgres password=mimi123 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE DATABASE moviedb")
	if err != nil {
		if err.Error() == "pq: database \"moviedb\" already exists" {
			fmt.Println("Database already exists, skipping creation")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Database created successfully")
	}

	db.Close()

	connStr = "user=postgres password=mimi123 dbname=moviedb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	fmt.Println("Successfully connected to the database")
}

func createTables() {
	createDirectorsTable := `
	CREATE TABLE IF NOT EXISTS directors (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50)
	);`

	createMoviesTable := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100),
		isbn VARCHAR(20),
		director_id INT,
		CONSTRAINT fk_director
		FOREIGN KEY(director_id) 
		REFERENCES directors(id)
	);`

	_, err := db.Exec(createDirectorsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createMoviesTable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables created successfully or already exist")
}

func main() {
	r := mux.NewRouter()

	initDB()

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	fmt.Println("Starting server on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT movies.id, movies.title, movies.isbn, directors.first_name, directors.last_name FROM movies JOIN directors ON movies.director_id = directors.id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		var director Director
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Isbn, &director.FirstName, &director.LastName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		movie.Director = &director
		movies = append(movies, movie)
	}
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var directorID int
	err = db.QueryRow("INSERT INTO directors (first_name, last_name) VALUES ($1, $2) RETURNING id", movie.Director.FirstName, movie.Director.LastName).Scan(&directorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var movieID int
	err = db.QueryRow("INSERT INTO movies (title, isbn, director_id) VALUES ($1, $2, $3) RETURNING id", movie.Title, movie.Isbn, directorID).Scan(&movieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	movie.ID = strconv.Itoa(movieID)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movieID, _ := strconv.Atoi(params["id"])
	_, err = db.Exec("UPDATE movies SET title=$1, isbn=$2 WHERE id=$3", movie.Title, movie.Isbn, movieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("UPDATE directors SET first_name=$1, last_name=$2 WHERE id=(SELECT director_id FROM movies WHERE id=$3)", movie.Director.FirstName, movie.Director.LastName, movieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	movieID, _ := strconv.Atoi(params["id"])

	var directorID int
	err := db.QueryRow("DELETE FROM movies WHERE id=$1 RETURNING director_id", movieID).Scan(&directorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("DELETE FROM directors WHERE id=$1", directorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	movieID, _ := strconv.Atoi(params["id"])

	var movie Movie
	var director Director
	err := db.QueryRow("SELECT movies.id, movies.title, movies.isbn, directors.first_name, directors.last_name FROM movies JOIN directors ON movies.director_id = directors.id WHERE movies.id=$1", movieID).Scan(&movie.ID, &movie.Title, &movie.Isbn, &director.FirstName, &director.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	movie.Director = &director
	json.NewEncoder(w).Encode(movie)
}

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

/* Using RESTey api:
 *  Get all: 	curl localhost:8000/albums
 *  Get by ID: 	curl localhost:8000/album/1
 *  Delete:     curl localhost:8000/album/1 -X DELETE
 *  Create: 	curl localhost:8000/album/ -X POST -d '{"title":"foo","artist":"bar","price":5.99}'
 */

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type dbHandler struct {
	db *sql.DB
}

func main() {
	filename := "albums.db"

	//Check if the data file needs to be created
	log.Printf("Checking if file %v exists", filename)
	fileStatus, err := checkDataFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if !fileStatus {
		log.Printf("Creating %v...", filename)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		log.Printf("Created %v successfully.", filename)
	}

	//connect to the sqlite database
	db, err := initDB(filename)
	if err != nil {
		log.Fatal("Could not initialize database")
	}

	//Set up the database tables & structures
	createTable(db)

	// hand off data control to the database handler
	dbh := dbHandler{db: db}

	// The webserver and its routing channels
	log.Println("Starting webserver on localhost:8000")
	mux := mux.NewRouter()
	mux.HandleFunc("/albums/", dbh.getAlbums).Methods("Get")
	mux.HandleFunc("/album/", dbh.addAlbum).Methods("Post")
	mux.HandleFunc("/album/{id}", dbh.getAlbumByID).Methods("Get")
	mux.HandleFunc("/album/{id}", dbh.deleteAlbumByID).Methods("Delete")

	log.Fatal(http.ListenAndServe(":8000", mux))
}

//Check if file exists:
//  true if file exists
//  false if it must be created
func checkDataFile(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func initDB(filename string) (*sql.DB, error) {
	//open DB connection
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

//Create the table if needed
func createTable(db *sql.DB) {
	tbl := `
	CREATE TABLE IF NOT EXISTS album(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Title TEXT NOT NULL,
		Artist TEXT NOT NULL,
		Price FLOAT
	);`

	log.Println("Creating table...")
	statement, err := db.Prepare(tbl)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Created table successfully.")
}

//Get all albums
func (dbh dbHandler) getAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var albums []Album
	row, err := dbh.db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		item := Album{}
		err := row.Scan(&item.ID, &item.Title, &item.Artist, &item.Price)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, item)
	}
	json.NewEncoder(w).Encode(albums)
}

//Get an album for a given ID
func (dbh dbHandler) getAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var alb Album
	row := dbh.db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		}
		//should something go here?
	}
	json.NewEncoder(w).Encode(alb)
}

//Insert data into the database
func (dbh dbHandler) addAlbum(w http.ResponseWriter, r *http.Request) {
	album := Album{}
	// Parse the JSON post body
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := dbh.db.Exec("INSERT INTO ALBUM (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//If it cannot return the ID, something bad happened
	id, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	//Outputs the ID of the created item
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

//Delete a record by ID
//Should probably check if the record exists first
func (dbh dbHandler) deleteAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//Don't check if it exists, just delete from the DB
	_, err := dbh.db.Exec("DELETE FROM ALBUM WHERE ID=?", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	status := "ok"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

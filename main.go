package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Url struct {
	ID       uint   `gorm:"primaryKey"` // used for shortUrl index
	Url      string `gorm:"unique"`     // prevent duplicate url
	ExpireAt string
	ShortUrl string
}

type APIUrl struct {
	ID       uint
	ShortUrl string
}

var db *gorm.DB
var err error

func main() {
	// gain access to database by getting .env
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dpPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// database connection string
	dbURI_postgres := fmt.Sprintf("host=%s user=%s dbname=postgres sslmode=disable password =%s port=%s",
		host, user, password, dpPort)
	dbURI_urlShortener := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password =%s port=%s",
		host, user, dbName, password, dpPort)

	// connect to database or create a new one
	db, err = gorm.Open(dialect, dbURI_urlShortener)
	if err != nil {
		db, err = gorm.Open(dialect, dbURI_postgres)
		if err != nil {
			fmt.Printf("Fail to connect to postgresql database")
		} else {
			db = db.Exec("CREATE DATABASE urlshortener;")
			db, err = gorm.Open(dialect, dbURI_urlShortener)
			if err != nil {
				fmt.Printf("Fail to create a new database")
			} else {
				fmt.Printf("Create a new data base and connect to it")
			}
		}
	} else {
		fmt.Println("Success connect to db!")
	}

	// make migrations to the dbif they have not already been created
	db.AutoMigrate(&Url{})

	// API routes
	router := mux.NewRouter()

	router.HandleFunc("/urls", getURLs).Methods("GET")
	router.HandleFunc("/{id}", getURL).Methods("GET")

	router.HandleFunc("/create/url", createURL).Methods("POST")

	// Listener
	http.ListenAndServe(":8080", router)

	// close connection to db when main func finishes
	defer db.Close()
}

//API controllers
func getURLs(w http.ResponseWriter, r *http.Request) {
	// define a type
	var urls []Url
	// find and match
	db.Find(&urls)

	json.NewEncoder(w).Encode(&urls)
}

func getURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var url Url

	err = db.Find(&url, params["id"]).Error
	if err != nil {
		errorHandler(w, r, http.StatusNotFound)
	} else {
		// check if expired
		// RFC3339
		timeT, _ := time.Parse(time.RFC3339, url.ExpireAt)
		now := time.Now()
		expired := timeT.Before(now)

		if expired {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(url.Url)
		}

	}

}

func createURL(w http.ResponseWriter, r *http.Request) {
	var url []Url
	json.NewDecoder(r.Body).Decode(&url)

	// using loop for not only one POST
	for idx := range url {
		createPerson := db.Create(&url[idx])
		err = createPerson.Error
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			// Response here
			// update before response
			db.Model(&url[idx]).Update(Url{ShortUrl: "localhost/" + fmt.Sprint(url[idx].ID)})

			// scale down the return value
			rv := APIUrl{ID: url[idx].ID, ShortUrl: url[idx].ShortUrl}
			json.NewEncoder(w).Encode(&rv)
		}
	}

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "status 404\n")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	host := os.Getenv("HOST")
	dpPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// database connection string
	dsn_postgres := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dpPort)
	dsn_urlShortener := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbName, dpPort)

	// connect to database or create a new one
	db, err = gorm.Open(postgres.Open(dsn_urlShortener), &gorm.Config{})
	if err != nil {
		db, err = gorm.Open(postgres.Open(dsn_postgres), &gorm.Config{})
		if err != nil {
			fmt.Println("Fail to connect to postgresql database")
			panic("Stop process")
		} else {
			db = db.Exec("CREATE DATABASE urlshortener;")
			db, err = gorm.Open(postgres.Open(dsn_urlShortener), &gorm.Config{})
			if err != nil {
				fmt.Printf("Fail to create a new database")
				panic("Stop process")
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

	router.HandleFunc("/api/v1/urls", createURL).Methods("POST")
	router.HandleFunc("/create/urls", createURLs).Methods("POST")

	// Listener
	http.ListenAndServe(":80", router)

}

//API controllers
func getURLs(w http.ResponseWriter, r *http.Request) {

	var urls []Url
	// stmt := db.Session(&gorm.Session{DryRun: true}).Find(&urls).Statement
	// syntax := stmt.SQL.String() //returns SQL query string without the param value
	// fmt.Println(syntax)

	err := db.Find(&urls).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(&urls)
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var url Url

	err := db.Find(&url, params["id"]).Error
	if err != nil {
		fmt.Println("status not found")
		w.WriteHeader(http.StatusNotFound)
	} else {
		// check if expired
		// RFC3339
		timeT, _ := time.Parse(time.RFC3339, url.ExpireAt)
		now := time.Now()
		expired := timeT.Before(now)

		if expired {
			// w.WriteHeader(http.StatusNotFound)
			errorHandler(w, r, http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(url.Url)
		}
	}

}

func createURL(w http.ResponseWriter, r *http.Request) {
	var url Url
	json.NewDecoder(r.Body).Decode(&url)

	// using loop for not only one POST
	createPerson := db.Create(&url)
	err = createPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		// Response here
		// update before response
		// db.Model(&url).Update(Url{ShortUrl: "localhost/" + fmt.Sprint(url.ID)})
		db.Model(&url).Update("ShortUrl", "localhost/"+fmt.Sprint(url.ID))

		// scale down the return value
		rv := APIUrl{ID: url.ID, ShortUrl: url.ShortUrl}
		json.NewEncoder(w).Encode(&rv)
	}

}

func createURLs(w http.ResponseWriter, r *http.Request) {
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
			db.Model(&url[idx]).Update("ShortUrl", "localhost/"+fmt.Sprint(url[idx].ID))

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

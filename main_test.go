package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// test the sql func , assuming http request is OK
func TestGetURL(t *testing.T) {
	//set const answer for this test
	testQuery := "SELECT * FROM `urls` WHERE `id` = $1"
	id := 1

	//create response writer and request for testing
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	//set up the mock sql connection
	testDB, mock, err := sqlmock.New()
	//handle error
	if err != nil {
		panic("error")
	}

	// uses "gorm.io/driver/postgres" library
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 testDB,
		PreferSimpleProtocol: true,
	})
	db, err = gorm.Open(dialector, &gorm.Config{})
	//handle error
	if err != nil {
		panic("error")
	}

	//mock the db.Find function
	rows := sqlmock.NewRows([]string{"id", "url", "expire_at", "short_url"}).
		AddRow(1, "http://somelongurl.com", "some_date", "http://shorturl.com")
	mock.ExpectQuery(regexp.QuoteMeta(testQuery)).
		WillReturnRows(rows).WithArgs(id)

	getURL(w, r)

	fmt.Println("IPASDJOAJSDOIAJSDOIJASDOIAJDS")
	fmt.Println(rows)
	fmt.Println(w.Body)

	//check values in mockedWriter using assert

}

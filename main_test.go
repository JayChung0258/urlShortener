package main

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// test the sql func , assuming http request is OK
func TestGetURL(t *testing.T) {
	//set const answer for this test
	testQuery := `SELECT * FROM "urls" WHERE "urls"."id" = $1`
	id := "1"

	//set up the mock sql connection
	testDB, mock, err := sqlmock.New()
	if err != nil {
		panic("sqlmock.New() occurs an error")
	}

	// uses "gorm.io/driver/postgres" library
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 testDB,
		PreferSimpleProtocol: true,
	})
	db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("Cannot open stub database")
	}

	//mock the db.Find function
	rows := sqlmock.NewRows([]string{"id", "url", "expire_at", "short_url"}).
		AddRow(1, "url", "date", "shorurl")

	//try to match the real SQL syntax we get and testQuery
	mock.ExpectQuery(regexp.QuoteMeta(testQuery)).WillReturnRows(rows).WithArgs(id)

	//set the value send into the function
	vars := map[string]string{
		"id": "1",
	}

	//create response writer and request for testing
	mockedWriter := httptest.NewRecorder()
	mockedRequest := httptest.NewRequest("GET", "/{id}", nil)
	mockedRequest = mux.SetURLVars(mockedRequest, vars)

	//call getURL()
	getURL(mockedWriter, mockedRequest)

	//check result in mockedWriter mocksql built function
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL syntax is not match: %s", err)
	}

}

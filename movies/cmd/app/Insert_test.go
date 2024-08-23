// ********RoostGPT********
/*
Test generated by RoostGPT for test go-test-assistent using AI Type  and AI Model 

ROOST_METHOD_HASH=insert_6526688601
ROOST_METHOD_SIG_HASH=insert_f90e997c71

Scenario 1: Successful Movie Insertion

Details:
  Description: This test is meant to check the successful operation of the `insert` function, ensuring that a movie can be correctly inserted into a database.
Execution:
  Arrange: Setup a mock for the `Decode` method of `json.NewDecoder` to return a valid `Movie` instance, and a mock for the `Insert` method of `MovieModel` to produce a success result (and not to return an error).
  Act: Invoke `insert` function using these mocks.
  Assert: Use Go's testing facilities to verify that no errors occurred during the process, and that the `Printf` method of `infoLog` was called with the expected parameters.
Validation:
  We are asserting that the `Printf` method of `infoLog` was called correctly because this call is the indirect output of the method, indicating a successful insertion. This test is important for confirming the function's capability to insert a movie into the database correctly.

Scenario 2: Failure of Decoding JSON to Movie Object

Details:
    Description: This test is designed to check the `insert` function's behavior when it encounters invalid JSON for a movie.
Execution:
    Arrange: Setup a mock for the `Decode` method of `json.NewDecoder` to return an error.
    Act: Invoke the `insert` function.
    Assert: Validate that the `serverError` method was called with the expected parameters reflecting the decoding error.
Validation:
    We are asserting that the `serverError` method was called because the `Decode` method failed and, according to the function logic, it should lead to an invocation of `serverError`. This test is essential to verify that the function can gracefully handle an error during the JSON decoding process.

Scenario 3: Failure of Movie Insertion into the Database

Details:
    Description: This test checks how the `insert` function behaves when it cannot insert the provided movie into the database.
Execution:
    Arrange: Setup a mock for the `Insert` method of `MovieModel` to return an error.
    Act: Invoke the `insert` function.
    Assert: Verify that the `serverError` method was called, reflecting the insertion error.
Validation:
    We expect the `serverError` method to be called due to the failing `Insert` method. This test is crucial to confirm that the function behaves correctly when encountering issues with inserting a movie into the database.
*/

// ********RoostGPT********
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mmorejon/microservices-docker-go-mongodb/movies/pkg/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockMovieModel struct {
	Movie        models.Movie
	Error        error
	InsertedID   string
}

func (m *mockMovieModel) Insert(movie models.Movie) (insertResult InsertResult, err error) {
	return InsertResult{InsertedID: m.InsertedID}, m.Error
}

func TestInsert(t *testing.T) {
	// Scenario 1: Successful Movie Insertion
	t.Run("Successful Movie Insertion", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/movie", bytes.NewBufferString(`{"title":"Avatar","rating":"PG-13","year":2009}`))
		response := httptest.NewRecorder()

		var infoLog bytes.Buffer
		info := log.New(&infoLog, "INFO\t", log.Ldate|log.Ltime)
		app := application{
			infoLog: info,
			movies:  &mockMovieModel{InsertedID: "SomeID"},
		}
		r := mux.NewRouter()
		r.HandleFunc("/movie", app.insert)
		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, infoLog.String(), "New movie have been created, id=SomeID")
	})

	// Scenario 2: Failure of Decoding JSON to Movie Object
	t.Run("Failure of Decoding JSON to Movie Object", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/movie", bytes.NewBufferString(`{"title":}`))
		response := httptest.NewRecorder()

		var errorLog bytes.Buffer
		err := log.New(&errorLog, "ERROR\t", log.Ldate|log.Ltime)
		app := application{
			errorLog: err,
			movies:   &mockMovieModel{},
		}
		r := mux.NewRouter()
		r.HandleFunc("/movie", app.insert)
		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Contains(t, errorLog.String(), "Failed to decode JSON body")
	})

	// Scenario 3: Failure of Movie Insertion into the Database
	t.Run("Failure of Movie Insertion into the Database", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/movie", bytes.NewBufferString(`{"title":"Avatar","rating":"PG-13","year":2009}`))
		response := httptest.NewRecorder()

		var errorLog bytes.Buffer
		err := log.New(&errorLog, "ERROR\t", log.Ldate|log.Ltime)
		app := application{
			errorLog: err,
			movies:   &mockMovieModel{Error: errors.New("Insertion failed")},
		}
		r := mux.NewRouter()
		r.HandleFunc("/movie", app.insert)
		r.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Contains(t, errorLog.String(), "Insertion failed")
	})
}

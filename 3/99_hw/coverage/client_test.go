package main

// тут писать код тестов
import (
	// "io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindUsers(t *testing.T) {

	initDB()
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	client2 := SearchClient{
		AccessToken: "authorization",
		URL:         "http://127.0.0.1:8080/badurl",
	}
	testCases := []SearchRequest{
		SearchRequest{
			Limit:      10,
			Offset:     0,
			Query:      "EXOSIS",
			OrderField: "id",
			OrderBy:    0,
		},
		SearchRequest{
			Limit:      27,
			Offset:     -1,
			Query:      "EXOSIS",
			OrderField: "id",
			OrderBy:    1,
		},
		SearchRequest{
			Limit:      -1,
			Offset:     0,
			Query:      "EXOSIS",
			OrderField: "id",
			OrderBy:    -1,
		},
		SearchRequest{
			Limit:      25,
			Offset:     0,
			Query:      "male",
			OrderField: "Name",
			OrderBy:    0,
		},
	}
	client.FindUsers(testCases[0])
	client.FindUsers(testCases[1])
	client.FindUsers(testCases[2])
	client.FindUsers(testCases[3])
	client2.FindUsers(testCases[0])
	ts.Close()
}

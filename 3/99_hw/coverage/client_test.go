package main

// тут писать код тестов
import (
	// "io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindUsers(t *testing.T) {
	client := SearchClient{
		AccessToken: "authorization",
		URL:         "http://127.0.0.1:8080/server",
	}
	initDB()
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	req := SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "id",
		OrderBy:    0,
	}
	client.FindUsers(req)
	ts.Close()
}

package main

// тут писать код тестов
import (
	// "io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimitAndOffsetFindUsers(t *testing.T) {
	client := SearchClient{
		AccessToken: "authorization",
		URL:         "serverUrl",
	}
	err1, _ := client.FindUsers(SearchRequest{
		Limit: -1,
	})
	err2, _ := client.FindUsers(SearchRequest{
		Offset: -1,
	})
	if err1 != nil && err2 != nil {
		t.Errorf("limit and offset test fail")
	}
}

func TestTimeoutServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	//	client.Do()
	t.Errorf("fail")
	ts.Close()
}
func TestNilRequest(t *testing.T) {
	client := SearchClient{
		AccessToken: "authorization",
		URL:         "",
	}
	nilReq := SearchRequest{}
	err, _ := client.FindUsers(nilReq)
	if err != nil {
		t.Errorf("Return not nil error")
	}
}

func TestBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL + "/yahoo/",
	}
	err, _ := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err != nil {
		t.Errorf("fail BadRequest")
	}
	client.URL = ts.URL
	err2, _ := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Bad field",
		OrderBy:    OrderByAsIs,
	})
	if err2 != nil {
		t.Errorf("fail BadRequest BadOrderField")
	}
	client.URL = ts.URL + "/google/"
	err3, _ := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err3 != nil {
		t.Errorf("fail BadRequest unknown error")
	}
}
func TestServerAndAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	filename = "badxml.xml"
	err, _ := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err != nil {
		t.Errorf("fail ServerError")
	}
	filename = "dataset.xml"
	client.AccessToken = "bad token"
	err2, _ := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err2 != nil {
		t.Errorf("fail BadAccessToken test")
	}
	ts.Close()
}

func TestBadJsonResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(GiveBadJsonServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	err, _ := client.FindUsers(SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "",
		OrderField: "Name",
		OrderBy:    1,
	})
	if err != nil {
		t.Errorf("fail error Unmarshal result json")
	}
	ts.Close()
}
func TestDataFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	err, _ := client.FindUsers(SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    0,
	})
	err2, _ := client.FindUsers(SearchRequest{
		Limit:      30,
		Offset:     0,
		Query:      "male",
		OrderField: "Name",
		OrderBy:    -1,
	})
	if err == nil {
		t.Errorf("NextPage true fail")
	}
	if err2 == nil {
		t.Errorf("Normal test with next page false failed")
	}
	ts.Close()
}

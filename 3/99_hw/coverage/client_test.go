package main

// тут писать код тестов
import (
	// "io"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLimitAndOffsetFindUsers(t *testing.T) {
	client := SearchClient{
		AccessToken: "authorization",
		URL:         "serverUrl",
	}
	_, err := client.FindUsers(SearchRequest{
		Limit: -1,
	})
	_, err2 := client.FindUsers(SearchRequest{
		Offset: -1,
	})
	if err == nil && err2 == nil {
		t.Errorf("limit and offset test fail")
	}
}

func TestTimeoutServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	_, err := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "long request",
		OrderField: "Name",
		OrderBy:    OrderByDesc,
	})

	if err == nil {
		t.Errorf("timeout test failed")
	}
	ts.Close()
}
func TestNilRequest(t *testing.T) {
	client := SearchClient{
		AccessToken: "authorization",
		URL:         "",
	}
	nilReq := SearchRequest{}
	_, err := client.FindUsers(nilReq)
	if err == nil {
		t.Errorf("fail nil SearchRequest")
	}
}

func TestBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL + "/yahoo/",
	}
	_, err := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err == nil {
		t.Errorf("fail BadRequest")
	}
	client.URL = ts.URL
	_, err2 := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Bad field",
		OrderBy:    OrderByAsIs,
	})
	if err2 == nil {
		t.Errorf("fail BadRequest BadOrderField")
	}
	client.URL = ts.URL + "/google/"
	_, err3 := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err3 == nil {
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
	_, err := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err == nil {
		t.Errorf("fail ServerError")
	}
	filename = "dataset.xml"
	client.AccessToken = "bad token"
	_, err2 := client.FindUsers(SearchRequest{
		Limit:      300,
		Offset:     0,
		Query:      "EXOSIS",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err2 == nil {
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
	_, err := client.FindUsers(SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      "",
		OrderField: "Name",
		OrderBy:    1,
	})
	if err == nil {
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
	res, err := client.FindUsers(SearchRequest{
		Limit:      30,
		Offset:     0,
		Query:      "male",
		OrderField: "Name",
		OrderBy:    OrderByAsc,
	})
	expectedId := 34
	res2, err2 := client.FindUsers(SearchRequest{
		Limit:      10,
		Offset:     3,
		Query:      "Dillard",
		OrderField: "Name",
		OrderBy:    OrderByDesc,
	})
	expected2Id := 3
	if err != nil || res.Users[0].Id != expectedId { // сравниваем первые ID
		t.Errorf("NextPage true fail")
	}
	if err2 != nil || res2.Users[0].Id != expected2Id { // сравниваем первые ID
		t.Errorf("Normal test without next page failed")
	}
	ts.Close()
}

func TestServerAtoi(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	searcherParams := url.Values{}
	searcherParams.Add("limit", "blabla")
	searcherParams.Add("offset", "0")
	searcherParams.Add("query", "male")
	searcherParams.Add("order_field", "Name")
	searcherParams.Add("order_by", "1")
	searcherReq, errq := http.NewRequest("GET", ts.URL+"?"+searcherParams.Encode(), nil)
	searcherReq.Header.Add("AccessToken", "authorization")
	resp, err := client.Do(searcherReq)
	if resp.StatusCode != http.StatusBadRequest && err == nil && errq != nil {
		t.Errorf("fail bad limit test")
	}
	searcherParams.Add("offset", "eEe")
	searcherReq2, errq2 := http.NewRequest("GET", ts.URL+"?"+searcherParams.Encode(), nil)
	searcherReq.Header.Add("AccessToken", "authorization")
	resp2, err2 := client.Do(searcherReq2)
	if resp2.StatusCode != http.StatusBadRequest && err2 == nil && errq2 != nil {
		t.Errorf("fail bad offset test")
	}
	ts.Close()
}

func TestServerUnMarshal(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	filename = "bad_dataset.xml"
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	_, err := client.FindUsers(SearchRequest{
		Limit:      30,
		Offset:     0,
		Query:      "male",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err == nil {
		fmt.Errorf("fail another xml file")
	}
	ts.Close()
}

func TestError500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	filename = "bad_dataset.xml"
	client := SearchClient{
		AccessToken: "authorization",
		URL:         ts.URL,
	}
	_, err := client.FindUsers(SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "panic server",
		OrderField: "Name",
		OrderBy:    OrderByAsIs,
	})
	if err == nil {
		fmt.Errorf("fail another xml file")
	}
	ts.Close()
}

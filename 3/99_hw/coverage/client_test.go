package main

// тут писать код тестов
import (
	// "io"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type Row struct {
	XMLName       xml.Name `xml:"row"`
	Id            int      `xml:"id"  `
	Guid          string   `xml:"guid"`
	IsActive      bool     `xml:"isActive"`
	Balance       string   `xml:"balance"`
	Picture       string   `xml:"picture"`
	Age           int      `xml:"age"`
	EyeColor      string   `xml:"eyeColor"`
	First_name    string   `xml:"first_name"`
	Last_name     string   `xml:"last_name"`
	Gender        string   `xml:"gender"`
	Company       string   `xml:"company"`
	Email         string   `xml:"email"`
	Phone         string   `xml:"phone"`
	Address       string   `xml:"address"`
	About         string   `xml:"about"`
	Registered    string   `xml:"registere"`
	FavoriteFruit string   `xml:"favoriteFruit"`
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	Version string   `xml:"version,attr"`
	Rows    []Row    `xml:"row"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFromFile(filename string) (xmlData []byte) {
	xml, err := os.Open(filename)
	defer xml.Close()
	check(err)
	xmlData, _ = ioutil.ReadAll(xml)
	return
}

func TestGetUser(t *testing.T) {
	xmlData := readFromFile("dataset.xml")
	rows := &Root{}
	err := xml.Unmarshal(xmlData, rows)
	check(err)
	fmt.Println(rows)
	/*for caseNum, item := range rows {
		url := "http://example.com/api/user?id=" + item.ID
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		GetUser(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}*/
}

func TestFindUsers(t *testing.T) {

}

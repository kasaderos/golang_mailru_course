package main

// тут писать код тестов
import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var filename string = "dataset.xml"

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

func isInRow(q string, r Row) (res bool) {
	if strings.Contains(r.Guid, q) ||
		strings.Contains(r.Balance, q) ||
		strings.Contains(r.Picture, q) ||
		strings.Contains(strconv.Itoa(r.Age), q) ||
		strings.Contains(r.EyeColor, q) ||
		strings.Contains(r.First_name, q) ||
		strings.Contains(r.Last_name, q) ||
		strings.Contains(r.Company, q) ||
		strings.Contains(r.Email, q) ||
		strings.Contains(r.Phone, q) ||
		strings.Contains(r.Address, q) ||
		strings.Contains(r.About, q) ||
		strings.Contains(r.Registered, q) ||
		strings.Contains(r.FavoriteFruit, q) ||
		strings.Contains(r.Gender, q) {
		return true
	}
	return false
}

func sortUsers(users []User, order string) {
	sort.Slice(users, func(i, j int) bool {
		if order == "-1" {
			return users[i].Id > users[j].Id
		}
		return users[i].Id < users[j].Id
	})
}

func getRows() ([]Row, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open file %v", err)
	}
	xmlData, _ := ioutil.ReadAll(file)
	file.Close()
	root := &Root{}
	errunmarshal := xml.Unmarshal(xmlData, root)
	if errunmarshal != nil {
		return nil, fmt.Errorf("can't unmarshal %v", errunmarshal)
	}
	return root.Rows, nil
}

func appendUser(resQuery []User, r Row) []User {
	resQuery = append(resQuery,
		User{
			Id:     r.Id,
			Gender: r.Gender,
			About:  r.About,
			Name:   r.First_name + r.Last_name,
			Age:    r.Age,
		},
	)
	return resQuery
}

func findQuery(query string) ([]User, error) {
	var resQuery []User
	rows, err := getRows()
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if isInRow(query, r) {
			resQuery = appendUser(resQuery, r)
		}
	}
	return resQuery, nil
}

func GiveBadJsonServer(w http.ResponseWriter, r *http.Request) {
	data := `{"id": 1, "name":"Marshal me please"}`
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	if r.Header.Get("AccessToken") != "authorization" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(r.URL.String(), "/yahoo/") {
		w.WriteHeader(http.StatusBadRequest)
		errjson := `{"Error": "error YAHOO"}`
		io.WriteString(w, errjson)
		return
	}
	if strings.HasPrefix(r.URL.String(), "/google/") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	if query == "long request" {
		time.Sleep(time.Second) // ищем этот текст
	}
	if orderField != "Name" { // только по имени
		w.WriteHeader(http.StatusBadRequest)
		errjson := `{"Error": "ErrorBadOrderField"}`
		io.WriteString(w, errjson)
		return
	}
	order := r.URL.Query().Get("order_by")
	limit, errconv := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, errconv2 := strconv.Atoi(r.URL.Query().Get("offset"))
	if errconv != nil || errconv2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, errquery := findQuery(query)
	if errquery != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sortUsers(results, order)
	if limit > 25 {
		results = results[offset:limit+offset]
	}
	data, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

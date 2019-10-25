package main

// тут писать код тестов
import (
	// "io"

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
			return users[i].Name > users[j].Name
		}
		return users[i].Name < users[j].Name
	})
}

func getRows() []Row {
	file, err := os.Open(filename)
	if err != nil {
		panic("not found file " + filename)
	}
	xmlData, _ := ioutil.ReadAll(file)
	file.Close()
	root := &Root{}
	xml.Unmarshal(xmlData, root)
	return root.Rows
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

func findQuery(query string) []User {
	var resQuery []User
	for _, r := range getRows() {
		if isInRow(query, r) {
			resQuery = appendUser(resQuery, r)
		}
	}
	return resQuery
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
		fmt.Println(r.URL.String())
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
	if orderField != "Name" {
		w.WriteHeader(http.StatusBadRequest)
		errjson := `{"Error": "ErrorBadOrderField"}`
		io.WriteString(w, errjson)
		return
	}
	order := r.URL.Query().Get("order_by")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	results := findQuery(query)
	sortUsers(results, order)
	data, _ := json.Marshal(results[:limit])
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}

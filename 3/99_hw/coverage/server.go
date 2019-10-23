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

type ResponseData struct {
	Users []User
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

var root *Root

func initDB() {
	xmlData := readFromFile("dataset.xml")
	root = &Root{}
	err := xml.Unmarshal(xmlData, root)
	check(err)
}
func isInRow(q string, r Row) (res bool) {
	b1 := strings.Contains(r.Guid, q)
	b2 := strings.Contains(r.Balance, q)
	b3 := strings.Contains(r.Picture, q)
	b4 := strings.Contains(strconv.Itoa(r.Age), q)
	b5 := strings.Contains(r.EyeColor, q)
	b6 := strings.Contains(r.First_name, q)
	b7 := strings.Contains(r.Last_name, q)
	b8 := strings.Contains(r.Company, q)
	b9 := strings.Contains(r.Email, q)
	b10 := strings.Contains(r.Phone, q)
	b11 := strings.Contains(r.Address, q)
	b12 := strings.Contains(r.About, q)
	b13 := strings.Contains(r.Registered, q)
	b14 := strings.Contains(r.FavoriteFruit, q)
	b15 := strings.Contains(r.Gender, q)
	res = b1 || b2 || b3 || b4 || b5 || b6 || b7 ||
		b8 || b9 || b10 || b11 || b12 || b13 || b14 || b15
	return
}

func sortUsers(users []User, orderFieldNum int, order int) {
	sort.Slice(users, func(i, j int) bool {
		switch orderFieldNum {
		case 0:
			if order == 1 {
				return users[i].Name < users[j].Name
			} else if order == -1 {
				return users[i].Name > users[j].Name
			}
		case 1:
			if order == 1 {
				return users[i].Age < users[j].Age
			} else if order == -1 {
				return users[i].Age > users[j].Age
			}
		case 2:
			if order == 1 {
				return users[i].Gender < users[j].Gender
			} else if order == -1 {
				return users[i].Gender > users[j].Gender
			}
		}
		if order == -1 {
			return users[i].Id > users[j].Id
		}
		return users[i].Id < users[j].Id
	})
}

func findQuery(query string, order int) []User {
	var resQuery []User
	for _, r := range root.Rows {
		if isInRow(query, r) {
			resQuery = append(resQuery,
				User{
					Id:     r.Id,
					Gender: r.Gender,
					About:  r.About,
					Name:   r.First_name + r.Last_name,
					Age:    r.Age,
				},
			)
		}
	}
	return resQuery
}

var orderFields map[string]int = map[string]int{
	"Name":   0,
	"Age":    1,
	"Gender": 2,
	"Id":     3,
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	check(err)
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	order, errconv := strconv.Atoi(r.URL.Query().Get("order_by"))
	if errconv != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "order is not int")
		return
	}

	results := findQuery(query, order)
	fmt.Println("results len", len(results))
	sortUsers(results, orderFields[orderField], order)

	data, err := json.Marshal(results)
	check(err)
	sdata := string(data)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, sdata)
}

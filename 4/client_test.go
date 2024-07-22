package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type DataRowXml struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type DataRow struct {
	Id     int
	Name   string
	Age    int
	About  string
	Gender string
}

func SearchServerLoadRows() ([]DataRow, error) {
	file, err := os.Open("dataset.xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]DataRow, 0)

	d := xml.NewDecoder(file)
	for t, _ := d.Token(); t != nil; t, _ = d.Token() {
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "row" {
				rowXml := DataRowXml{}
				err = d.DecodeElement(&rowXml, &se)
				if err != nil {
					continue
				}
				result = append(result, DataRow{Id: rowXml.Id, Name: rowXml.FirstName + " " + rowXml.LastName, Age: rowXml.Age, About: rowXml.About, Gender: rowXml.Gender})
			}
		}
	}

	return result, nil
}

type ServerRequestParams struct {
	Limit      int
	Offset     int
	Query      string
	OrderField string
	OrderBy    int
}

var ServerRequestAcceptedOrderFields = []string{"Id", "Age", "Name"}

type SearchServerErrorResponse struct {
	Error string
}

func SearchServerWriteJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, _ = w.Write([]byte(jsonData))
	return nil
}

func SearchServerHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	params := ServerRequestParams{25, 0, "", "Name", 0}

	limit := r.URL.Query().Get("limit")

	if len(limit) > 0 {
		params.Limit, err = strconv.Atoi(limit)
		if err != nil || params.Limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_ = SearchServerWriteJson(w, SearchServerErrorResponse{`Bad param "limit"`})
			return
		}
	}

	offset := r.URL.Query().Get("offset")
	if len(offset) > 0 {
		params.Offset, err = strconv.Atoi(offset)
		if err != nil || params.Offset < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_ = SearchServerWriteJson(w, SearchServerErrorResponse{`Bad param "offset"`})
			return
		}
	}

	params.Query = r.URL.Query().Get("query")

	orderField := r.URL.Query().Get("order_field")
	if len(orderField) > 0 {
		if slices.Contains(ServerRequestAcceptedOrderFields, orderField) {
			params.OrderField = orderField
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_ = SearchServerWriteJson(w, SearchServerErrorResponse{`ErrorBadOrderField`})
			return
		}
	}

	orderBy := r.URL.Query().Get("order_by")
	if len(orderBy) > 0 {
		params.OrderBy, err = strconv.Atoi(orderBy)
		if err != nil || params.OrderBy < -1 || params.OrderBy > 1 {
			w.WriteHeader(http.StatusBadRequest)
			_ = SearchServerWriteJson(w, SearchServerErrorResponse{`Bad param "order_by"`})
			return
		}
	}

	AccessToken := r.Header.Get("AccessToken")
	if len(AccessToken) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rows := make([]DataRow, len(dbRows))
	copy(rows, dbRows)

	if params.OrderBy != 0 {
		sort.Slice(rows, func(i, j int) bool {
			switch params.OrderField {
			case "Id":
				if params.OrderBy == 1 {
					return rows[i].Id < rows[j].Id
				} else {
					return rows[i].Id > rows[j].Id
				}
			case "Age":
				if params.OrderBy == 1 {
					return rows[i].Age < rows[j].Age
				} else {
					return rows[i].Age > rows[j].Age
				}
			default:
				if params.OrderBy == 1 {
					return rows[i].Name < rows[j].Name
				} else {
					return rows[i].Name > rows[j].Name
				}
			}
		})
	}

	users := make([]DataRow, 0)

	i := 0
	f := 0
	for _, row := range rows {
		if f >= params.Limit {
			break
		}
		if len(params.Query) > 0 && !strings.Contains(row.Name, params.Query) && !strings.Contains(row.About, params.Query) {
			continue
		}
		i++
		if i <= params.Offset {
			continue
		}
		users = append(users, row)
		f++
	}

	_ = SearchServerWriteJson(w, users)

}

func SearchServerHandlerErrorTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1100 * time.Millisecond)
	_ = SearchServerWriteJson(w, true)
}

func SearchServerHandlerErrorInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func SearchServerHandlerErrorBadJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write([]byte("----bad-json"))
}

func SearchServerHandlerErrorBadJsonOnBadRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("----bad-json"))
}

var dbRows []DataRow

func init() {
	var err error
	dbRows, err = SearchServerLoadRows()
	if err != nil {
		panic(err)
	}
}

type TestCase struct {
	Params  SearchRequest
	IsError bool
	Result  *SearchResponse
}

type TestServerErrorCase struct {
	Token   string
	Handler func(w http.ResponseWriter, r *http.Request)
}

func (testcase *TestServerErrorCase) Client() *SearchClient {
	server := httptest.NewServer(http.HandlerFunc(testcase.Handler))
	client := &SearchClient{
		URL:         server.URL,
		AccessToken: "token",
	}
	if testcase.Token == "__empty" {
		client.AccessToken = ""
	} else if len(testcase.Token) > 0 {
		client.AccessToken = testcase.Token
	}
	return client
}

func TestFindUsers(t *testing.T) {

	allUsers := make([]User, 0, len(dbRows))
	for _, row := range dbRows {
		allUsers = append(allUsers, User(row))
	}

	testcases := []TestCase{
		{
			IsError: true,
			Params:  SearchRequest{Limit: -1},
		},
		{
			IsError: true,
			Params:  SearchRequest{Offset: -1},
		},
		{
			IsError: true,
			Params:  SearchRequest{OrderField: "Another"},
		},
		{
			IsError: true,
			Params:  SearchRequest{OrderBy: 10},
		},
		{
			IsError: false,
			Params:  SearchRequest{},
			Result:  &SearchResponse{Users: []User{}, NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 2, Offset: 1},
			Result:  &SearchResponse{Users: []User{allUsers[1], allUsers[2]}, NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Offset: 1000},
			Result:  &SearchResponse{Users: []User{}, NextPage: false},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 1, Query: "Hilda"},
			Result:  &SearchResponse{Users: []User{allUsers[1]}, NextPage: false},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 3, Offset: 3, Query: "A"},
			Result:  &SearchResponse{Users: []User{allUsers[9], allUsers[11], allUsers[15]}, NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 100},
			Result:  &SearchResponse{Users: allUsers[0:25], NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 2, Offset: 1, OrderField: "Id", OrderBy: 1},
			Result:  &SearchResponse{Users: []User{allUsers[1], allUsers[2]}, NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 2, Offset: 1, OrderField: "Id", OrderBy: -1},
			Result:  &SearchResponse{Users: []User{allUsers[33], allUsers[32]}, NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 2, Offset: 1, OrderField: "Age", OrderBy: 1},
			Result:  &SearchResponse{Users: []User{allUsers[15], allUsers[23]}, NextPage: true},
		},
		{
			IsError: false,
			Params:  SearchRequest{Limit: 2, Offset: 1, OrderField: "Age", OrderBy: -1},
			Result:  &SearchResponse{Users: []User{allUsers[13], allUsers[6]}, NextPage: true},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServerHandler))

	for caseNum, testcase := range testcases {

		c := &SearchClient{
			URL:         ts.URL,
			AccessToken: "token",
		}
		result, err := c.FindUsers(testcase.Params)

		if err != nil && !testcase.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && testcase.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(testcase.Result, result) {
			t.Errorf("[%d] wrong result, expected: \n %#v\n\ngot %#v\n", caseNum, testcase.Result, result)
		}
	}

	serverErrorsCases := []TestServerErrorCase{
		{
			Handler: SearchServerHandler,
			Token:   "__empty",
		},
		{
			Handler: SearchServerHandlerErrorTimeout,
		},
		{
			Handler: SearchServerHandlerErrorInternalServerError,
		},
		{
			Handler: SearchServerHandlerErrorBadJson,
		},
		{
			Handler: SearchServerHandlerErrorBadJsonOnBadRequest,
		},
	}

	for caseNum, errorCase := range serverErrorsCases {
		c := errorCase.Client()
		_, err := c.FindUsers(SearchRequest{})
		if err == nil {
			t.Errorf("Server error testcase [%d] expected error, got nil", caseNum)
		}
	}

	c := &SearchClient{
		URL:         "bad_url",
		AccessToken: "token",
	}
	_, err := c.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Unknown NET error expected error, got nil")
	}

}

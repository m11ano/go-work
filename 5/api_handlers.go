package main

import (
	"strconv"
	"net/http"
	"encoding/json"
	"slices"
	"unicode/utf8"
)

type JsonResponse struct {
	Error string `json:"error"`
	Response interface{} `json:"response,omitempty"`
}


// OtherApi
func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
	
	case r.URL.Path == "/user/create":
		
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotAcceptable)
			json, _ := json.Marshal(JsonResponse{Error: "bad method"})
			w.Write(json)
			return 
		}
		
        h.CreatePOST(w, r)
	
    default:
        w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(JsonResponse{Error: "unknown method"})
		w.Write(json)
    }
}


func (h *OtherApi) CreatePOST(w http.ResponseWriter, r *http.Request) {
	params := OtherCreateParams{}
	ctx := r.Context()
	
	authToken := r.Header.Get("X-Auth")
	if authToken != "100500" {
		w.WriteHeader(http.StatusForbidden)
		json, _ := json.Marshal(JsonResponse{Error: "unauthorized"})
		w.Write(json)
		return
	}
	
	 
	params.Username = r.FormValue("username")
	
	if len(params.Username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "username must me not empty"})
		w.Write(json)
		return
	}
	
	
	
	
	if utf8.RuneCountInString(params.Username) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "username len must be >= 3"})
		w.Write(json)
		return
	}
	
	  
	params.Name = r.FormValue("account_name")
	
	
	
	
	  
	params.Class = r.FormValue("class")
	
	
	if len(params.Class) == 0 {
		params.Class = "warrior"
	}
	
	
	ClassEnums := []string{ "warrior", "sorcerer", "rouge" }
	if !slices.Contains(ClassEnums, params.Class) {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "class must be one of [warrior, sorcerer, rouge]"})
		w.Write(json)
		return
	}
	
	
	  
	LevelValue, err := strconv.Atoi(r.FormValue("level"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "level must be int"})
		w.Write(json)
		return
	}
	params.Level = LevelValue
	
	
	
	if params.Level < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "level must be >= 1"})
		w.Write(json)
		return
	}
	
	
	if params.Level > 50 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "level must be <= 50"})
		w.Write(json)
		return
	}
	 
	res, err := h.Create(ctx, params)
	if err != nil {
		statuseCode := http.StatusInternalServerError
		if apiErr, ok := err.(ApiError); ok {
			statuseCode = apiErr.HTTPStatus
		}
		w.WriteHeader(statuseCode)
		json, _ := json.Marshal(JsonResponse{Error: err.Error()})
		w.Write(json)
		return
	}
	json, err := json.Marshal(JsonResponse{Response: res})
	if err != nil {
		panic(err)
	}
	w.Write(json)
}


// MyApi
func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
	
	case r.URL.Path == "/user/profile":
		
        h.ProfileALL(w, r)
	
	case r.URL.Path == "/user/create":
		
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotAcceptable)
			json, _ := json.Marshal(JsonResponse{Error: "bad method"})
			w.Write(json)
			return 
		}
		
        h.CreatePOST(w, r)
	
    default:
        w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(JsonResponse{Error: "unknown method"})
		w.Write(json)
    }
}


func (h *MyApi) ProfileALL(w http.ResponseWriter, r *http.Request) {
	params := ProfileParams{}
	ctx := r.Context()
	
	 
	params.Login = r.FormValue("login")
	
	if len(params.Login) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "login must me not empty"})
		w.Write(json)
		return
	}
	
	
	
	
	 
	res, err := h.Profile(ctx, params)
	if err != nil {
		statuseCode := http.StatusInternalServerError
		if apiErr, ok := err.(ApiError); ok {
			statuseCode = apiErr.HTTPStatus
		}
		w.WriteHeader(statuseCode)
		json, _ := json.Marshal(JsonResponse{Error: err.Error()})
		w.Write(json)
		return
	}
	json, err := json.Marshal(JsonResponse{Response: res})
	if err != nil {
		panic(err)
	}
	w.Write(json)
}

func (h *MyApi) CreatePOST(w http.ResponseWriter, r *http.Request) {
	params := CreateParams{}
	ctx := r.Context()
	
	authToken := r.Header.Get("X-Auth")
	if authToken != "100500" {
		w.WriteHeader(http.StatusForbidden)
		json, _ := json.Marshal(JsonResponse{Error: "unauthorized"})
		w.Write(json)
		return
	}
	
	 
	params.Login = r.FormValue("login")
	
	if len(params.Login) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "login must me not empty"})
		w.Write(json)
		return
	}
	
	
	
	
	if utf8.RuneCountInString(params.Login) < 10 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "login len must be >= 10"})
		w.Write(json)
		return
	}
	
	  
	params.Name = r.FormValue("full_name")
	
	
	
	
	  
	params.Status = r.FormValue("status")
	
	
	if len(params.Status) == 0 {
		params.Status = "user"
	}
	
	
	StatusEnums := []string{ "user", "moderator", "admin" }
	if !slices.Contains(StatusEnums, params.Status) {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "status must be one of [user, moderator, admin]"})
		w.Write(json)
		return
	}
	
	
	  
	AgeValue, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "age must be int"})
		w.Write(json)
		return
	}
	params.Age = AgeValue
	
	
	
	if params.Age < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "age must be >= 0"})
		w.Write(json)
		return
	}
	
	
	if params.Age > 128 {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(JsonResponse{Error: "age must be <= 128"})
		w.Write(json)
		return
	}
	 
	res, err := h.Create(ctx, params)
	if err != nil {
		statuseCode := http.StatusInternalServerError
		if apiErr, ok := err.(ApiError); ok {
			statuseCode = apiErr.HTTPStatus
		}
		w.WriteHeader(statuseCode)
		json, _ := json.Marshal(JsonResponse{Error: err.Error()})
		w.Write(json)
		return
	}
	json, err := json.Marshal(JsonResponse{Response: res})
	if err != nil {
		panic(err)
	}
	w.Write(json)
}


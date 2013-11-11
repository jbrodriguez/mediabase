package helper

import (
	"encoding/json"
	"github.com/gorilla/mux"
	// "io/ioutil"
	"log"
	"net/http"
)

func ReadForm(resp http.ResponseWriter, req *http.Request) bool {
	err := req.ParseForm()
	if err != nil {
		WriteJson(resp, 400, &StringMap{"message": "Invalid body"})
		return false
	}
	return true
}

func ReadJson(resp http.ResponseWriter, req *http.Request, reqD interface{}) bool {
	err := json.NewDecoder(req.Body).Decode(reqD)

	// hah, err2 := ioutil.ReadAll(req.Body)

	// if err2 != nil {
	// 	log.Printf("%s", err2)
	// }

	// log.Printf("%s", string(hah))

	if err != nil {
		WriteJson(resp, 400, &StringMap{"message": "Invalid body"})
		return false
	}
	return true
}

func WriteJson(resp http.ResponseWriter, status int, respD interface{}) {
	b, err := json.Marshal(respD)
	if err != nil {
		WriteJsonErr(resp, err)
	} else {
		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp.WriteHeader(status)
		resp.Write(b)
		resp.Write([]byte("\n"))
	}
}

func WriteJsonErr(resp http.ResponseWriter, err error) {
	log.Println("error:", err)
	WriteJson(resp, 500, &StringMap{"message": "Internal server error"})
}

func uvar(req *http.Request, name string) string {
	return mux.Vars(req)["id"]
}

func fvar(req *http.Request, name string) string {
	return req.FormValue(name)
}

func qvar(req *http.Request, name string) string {
	return req.FormValue(name)
}

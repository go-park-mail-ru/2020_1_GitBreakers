package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

var profilMap = map[string]Profile{
	"antonelagin": Profile{
		Login:     "AntonElagin",
		Avatar:    "https://avatars.mds.yandex.net/get-pdb/1964870/cfaa9b42-d74b-40f7-93cc-1d777bb5654c/s1200",
		Name:      "ANton ELagin",
		Bio:       "kek lol opps, long sfddsfdsfd dfdsfs dsfds sdfdsf sdfsd sdfd",
		URL:       "google.com",
		Followers: 1,
		Following: 2,
	},
	"keklol": Profile{
		Login:     "keeeekLooool",
		Avatar:    "https://avatars.mds.yandex.net/get-pdb/1964870/cfaa9b42-d74b-40f7-93cc-1d777bb5654c/s1200",
		Followers: 1,
		Following: 150,
	}}

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type Profile struct {
	Avatar    string `json:"avatar"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Bio       string `json:"Bio"`
	URL       string `json:"url"`
	Followers uint   `json:"followers"`
	Following uint   `json:"following"`
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	// TODO : проверка авторизации

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	prof := Profile{}
	body, errB := ioutil.ReadAll(r.Body)

	if errB != nil {
		json.NewEncoder(w).Encode(&Result{
			Err: "Server Error",
		})
		return
	}

	err := json.Unmarshal(body, &prof)
	if err != nil {
		json.NewEncoder(w).Encode(&Result{
			Err: "bad body",
		})
		return
	}
	login := strings.ToLower(prof.Login)
	if _, ok := profilMap[login]; !ok {
		json.NewEncoder(w).Encode(&Result{
			Err: "unknown login",
		})
		return
	}
	profilMap[login] = prof

	json.NewEncoder(w).Encode(&Result{Body: "success"})

}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	login := vars["login"]

	profile, ok := profilMap[strings.ToLower(login)]
	if !ok {
		json.NewEncoder(w).Encode(&Result{
			Err: "Not found"})
		return
	}

	json.NewEncoder(w).Encode(&Result{Body: profile})
}

// не работает
func UploadAvatar(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(5 * 1024 * 1025)
	image, header, err := r.FormFile("avatar")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer image.Close()

	byteImage, err := ioutil.ReadAll(image)
	if err != nil {
		fmt.Println(err)
		return
	}

	filePath := "./static/image/" + header.Filename
	err = ioutil.WriteFile(filePath, byteImage, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	webPath := "localhost:8080/static/image/" + header.Filename

	elem, er := profilMap["antonelagin"]
	if !er {
		json.NewEncoder(w).Encode(&Result{
			Err: "unknown login",
		})
		return
	}

	elem.Avatar = webPath
	profilMap["antonelagin"] = elem
}

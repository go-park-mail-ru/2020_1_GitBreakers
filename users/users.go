package users

import (
	"../form"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]
	profile, ok := form.UserSlice[strings.ToLower(login)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not issue this user!"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	prof := &form.SignupForm{}
	err := json.NewDecoder(r.Body).Decode(prof)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid json!"})
		return
	}

	login := strings.ToLower(prof.Login)
	if _, ok := form.UserSlice[login]; !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found this user"})
		return
	}
	form.UserSlice[login] = *prof
	json.NewEncoder(w).Encode(map[string]string{"info": "User data is changed"})
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 * 1024 * 1025)
	image, header, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found this user"})
		return
	}
	defer image.Close()
	//byteImage сама картинка
	byteImage, err := ioutil.ReadAll(image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		return
	}
	//засейвили картинку
	filePath := "./static/image/" + header.Filename
	err = ioutil.WriteFile(filePath, byteImage, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		return
	}

	webPath := "http://localhost:8080/static/image/" + header.Filename
	//TODO получить логин юзера
	elem, ok := form.UserSlice["antonelagin"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found this user"})
		return
	}

	elem.Avatar = webPath
	//TODO получить логин юзера
	form.UserSlice["antonelagin"] = elem
	json.NewEncoder(w).Encode(map[string]string{"info": "User data is changed"})
}

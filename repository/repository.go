package repository

import (
	"../form"
	"encoding/json"
	"fmt"
	"net/http"
)

var repositorySlice = []form.Repository{
	{
		Name:          "kekmda",
		Fullname:      "Deiklov/kekmda",
		Private:       false,
		Description:   "somerepository",
		Fork:          false,
		Owner:         "Deiklov",
		DefaultBranch: "master",
	},
	{
		Name:          "somerepo",
		Fullname:      "Dmitry/somerepo",
		Private:       false,
		Description:   "some text",
		Fork:          false,
		Owner:         "Deiklov",
		DefaultBranch: "dev",
	},
}

func NewRepository(w http.ResponseWriter, r *http.Request) {
	repoForm := &form.Repository{}
	err := json.NewDecoder(r.Body).Decode(repoForm)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid json must be field name, private(bool), description, fork(bool), owner"}`))
		return
	}

	session, _ := form.StoreSession.Get(r, "session_id")

	//вытащили из сессии юзерлогин
	repoForm.Owner = fmt.Sprintf("%v", session.Values["login"])
	repoForm.Fullname = repoForm.Name + "/" + repoForm.Owner

	//не позволит создать херню с одинаковым именем репозитория
	for _, v := range repositorySlice {
		if v.Fullname == repoForm.Fullname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Repo with this full_name already exsist!"})
			return
		}
	}
	repositorySlice = append(repositorySlice, *repoForm) //save to db
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"info": "Repo are created", "full_name": repoForm.Fullname})
	return
}

//TODO: Поиск не по name а по full_name
func GetRepository(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Your url must to be like ?name=<reponame>"})
		return
	}
	repoForm := &form.Repository{}

	for _, v := range repositorySlice {
		if v.Name == name {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(repoForm)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": "Not found repo with this name!"})
}

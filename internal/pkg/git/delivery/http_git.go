package delivery

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

type GitDelivery struct {
	UC     git.UseCase
	Logger *logger.SimpleLogger
	UserUC user.UCUser
}

//создать репак(id,name,description,private,owner)
func (GD *GitDelivery) CreateRepo(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		GD.Logger.HttpInfo(r.Context(), "already authorized", http.StatusBadRequest)
		return
	}
	userID := res.(int)
	newRepo := &gitmodels.Repository{IsPublic: true}
	if err := json.NewDecoder(r.Body).Decode(newRepo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		GD.Logger.HttpInfo(r.Context(), "json invalid", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(newRepo); err != nil {
		GD.Logger.HttpInfo(r.Context(), err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := GD.UC.Create(userID, newRepo); err != nil {
		GD.Logger.HttpInfo(r.Context(), "repo not created", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

//данные репока(модельку скинуть(id,name,private,owner)
func (GD *GitDelivery) GetRepo(w http.ResponseWriter, r *http.Request) {
	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]
	Repo, err := GD.UC.GetRepo(userName, repoName)
	if err != nil {
		GD.Logger.HttpInfo(r.Context(), "not got repo", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(&Repo); err != nil {
		GD.Logger.HttpInfo(r.Context(), "not encode json", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//
////все репозитории юзера
func (GD *GitDelivery) GetRepoList(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
		GD.Logger.HttpInfo(r.Context(), "already authorized", http.StatusBadRequest)
		return
	}
	userID := res.(int)
	userName := mux.Vars(r)["username"]
	if userName == "" {
		userModel, err := GD.UserUC.GetByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		userName = userModel.Login
	}
	repo, err := GD.UC.GetRepoList(userName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(&repo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//
////ветки репака(просто названия и ссылки)
func (GD *GitDelivery) GetBranchList(w http.ResponseWriter, r *http.Request) {
	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]

	res := r.Context().Value("UserID")
	var branches []gitmodels.Branch
	var err error

	if res == nil {
		branches, err = GD.UC.GetBranchList(nil, userName, repoName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		userID := res.(int)
		branches, err = GD.UC.GetBranchList(&userID, userName, repoName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	if err := json.NewEncoder(w).Encode(&branches); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

////cписок коммитов для ветки
//func (GD *GitDelivery) GetCommitsList(w http.ResponseWriter, r *http.Request) {
//
//}

package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"net/http"
)

type GitDelivery struct {
	UC     git.UseCase
	Logger *logger.SimpleLogger
}

//создать репак(id,name,description,private,owner)
func (GD *GitDelivery) CreateRepo(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res != nil {
		w.WriteHeader(http.StatusBadRequest)
		GD.Logger.HttpInfo(r.Context(), "already authorized", http.StatusBadRequest)
		return
	}
	userID := res.(int)
	newRepo := &gitmodels.Repository{}
	if err := json.NewDecoder(r.Body).Decode(newRepo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		GD.Logger.HttpInfo(r.Context(), "json invalid", http.StatusBadRequest)
		return
	}
	if err := GD.UC.Create(userID, newRepo); err != nil {
		GD.Logger.HttpInfo(r.Context(), "repo not created", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

////данные репока(модельку скинуть(id,name,private,owner)
//func (GD *GitDelivery) GetRepo(w http.ResponseWriter, r *http.Request) {
//	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]
//	Repo := &gitmodels.Repository{}
//	Repo = GD.UC.GetRepo(userName, repoName)
//}
//
////все репозитории юзера
//func (GD *GitDelivery) GetRepoList(w http.ResponseWriter, r *http.Request) {
//	userName := mux.Vars(r)["username"]
//	///вернет слайс репаков
//	Repo := GD.UC.GetRepoList(userName)
//}
//
////ветки репака(просто названия и ссылки)
//func (GD *GitDelivery) GetBranchList(w http.ResponseWriter, r *http.Request) {
//	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]
//	GD.UC.
//}
//
////cписок коммитов для ветки
//func (GD *GitDelivery) GetCommitsList(w http.ResponseWriter, r *http.Request) {
//
//}

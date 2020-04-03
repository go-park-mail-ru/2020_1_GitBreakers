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
	if res := r.Context().Value("UserID"); res != nil {
		w.WriteHeader(http.StatusBadRequest)
		GD.Logger.HttpInfo(r.Context(), "already authorized", http.StatusBadRequest)
		return
	}
	newRepo := &gitmodels.Repository{}
	if err := json.NewDecoder(r.Body).Decode(newRepo); err != nil {

	}

}

//данные репока(модельку скинуть(id,name,private,owner)
func (GD *GitDelivery) GetRepo(w http.ResponseWriter, r *http.Request) {

}

//все репозитории юзера
func (GD *GitDelivery) GetRepoList(w http.ResponseWriter, r *http.Request) {

}

//ветки репака
func (GD *GitDelivery) GetBranchList(w http.ResponseWriter, r *http.Request) {

}

//cписок коммитов для ветки
func (GD *GitDelivery) GetCommitsList(w http.ResponseWriter, r *http.Request) {

}

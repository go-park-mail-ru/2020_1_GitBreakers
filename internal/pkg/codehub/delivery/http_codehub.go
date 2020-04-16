package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"net/http"
)

type Http_Codehub struct {
	Logger    *logger.SimpleLogger
	UserUC    user.UCUser
	CodeHubUC codehub.UCCodeHub
}

func (GD *Http_Codehub) AddStar(w http.ResponseWriter, r *http.Request) {

}

func (GD *Http_Codehub) DelStar(w http.ResponseWriter, r *http.Request) {

}

func (GD *Http_Codehub) StarredRepos(w http.ResponseWriter, r *http.Request) {

}

func (GD *Http_Codehub) NewIssue(w http.ResponseWriter, r *http.Request) {

}

func (GD *Http_Codehub) UpdateIssue(w http.ResponseWriter, r *http.Request) {

}

func (GD *Http_Codehub) GetIssues(w http.ResponseWriter, r *http.Request) {

}

func (GD *Http_Codehub) CloseIssue(w http.ResponseWriter, r *http.Request) {

}

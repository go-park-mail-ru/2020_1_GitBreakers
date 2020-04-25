package delivery

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

type HttpCodehub struct {
	Logger    *logger.SimpleLogger
	CodeHubUC codehub.UCCodeHub
	Ws        websocket.Upgrader
}

func (GD *HttpCodehub) ModifyStar(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newStar := models.Star{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newStar); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newStar.AuthorID = userID

	if err := GD.CodeHubUC.ModifyStar(newStar); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "star modify success")

}

func (GD *HttpCodehub) StarredRepos(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repolist, err := GD.CodeHubUC.GetStarredRepo(userID)
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(repolist, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "repolist got success")
}

func (GD *HttpCodehub) NewIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newIssue := models.Issue{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newIssue.AuthorID = userID

	err := GD.CodeHubUC.CreateIssue(newIssue)

	switch {
	case err == entityerrors.AccessDenied():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case err == entityerrors.DoesNotExist():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	GD.Logger.HttpLogInfo(r.Context(), "issues created success")
}

func (GD *HttpCodehub) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newIssue := models.Issue{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newIssue.AuthorID = userID //чтобы автора не подменяли

	oldIssue, err := GD.CodeHubUC.GetIssue(newIssue.ID, newIssue.AuthorID)

	switch {
	case err == entityerrors.AccessDenied():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case err == entityerrors.DoesNotExist():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if govalidator.IsByteLength(newIssue.Message, 1, 1024) {
		oldIssue.Message = newIssue.Message
	}
	if govalidator.IsByteLength(newIssue.Title, 1, 256) {
		oldIssue.Message = newIssue.Title
	}
	if govalidator.IsByteLength(newIssue.Label, 0, 100) {
		oldIssue.Message = newIssue.Label
	}

	err = GD.CodeHubUC.UpdateIssue(oldIssue)

	switch {
	case err == entityerrors.AccessDenied():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case err == entityerrors.DoesNotExist():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "issues updated")
}

func (GD *HttpCodehub) GetIssues(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //скинули строку, а не число
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		return
	}

	userID, ok := res.(int64)
	if !ok {
		userID = -1 //соответствует неавторизованному юзеру
	}

	issueslist, err := GD.CodeHubUC.GetIssuesList(int64(repoID), userID)

	switch {
	case err == entityerrors.AccessDenied():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case err == entityerrors.DoesNotExist():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(issueslist, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(issueslist, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "issues returned success")
}

func (GD *HttpCodehub) CloseIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	oldIssue := models.Issue{} //достаточно скинуть просто id

	if err := easyjson.UnmarshalFromReader(r.Body, &oldIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := GD.CodeHubUC.CloseIssue(oldIssue.ID, userID)

	switch {
	case err == entityerrors.AccessDenied():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case err == entityerrors.DoesNotExist():
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "issues closed success")
}
func (GD *HttpCodehub) GetNews(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}


	GD.Logger.HttpLogInfo(r.Context(), "news getted success")
}

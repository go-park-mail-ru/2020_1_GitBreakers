package delivery

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Http_Codehub struct {
	Logger    *logger.SimpleLogger
	UserUC    user.UCUser
	CodeHubUC codehub.UCCodeHub
}

func (GD *Http_Codehub) ModifyStar(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newStar := models.Star{AuthorID: userID}

	if err := json.NewDecoder(r.Body).Decode(&newStar); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := GD.CodeHubUC.ModifyStar(newStar); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "star modify success")

}

func (GD *Http_Codehub) StarredRepos(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int)
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

	if err := json.NewEncoder(w).Encode(repolist); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "repolist got success")
}

func (GD *Http_Codehub) NewIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newIssue := models.Issue{AuthorID: userID}

	if err := json.NewDecoder(r.Body).Decode(&newIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//todo switch err(not found,access denied...)
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

func (GD *Http_Codehub) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newIssue := models.Issue{}

	if err := json.NewDecoder(r.Body).Decode(&newIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newIssue.AuthorID = userID //чтобы автора не подменяли

	oldIssue, err := GD.CodeHubUC.GetIssue(newIssue.AuthorID, newIssue.AuthorID)

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

func (GD *Http_Codehub) GetIssues(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //скинули строку, а не число
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		return
	}

	userID, ok := res.(int)
	if !ok {
		userID = -1 //соответствует неавторизованному юзеру
	}

	issueslist, err := GD.CodeHubUC.GetIssuesList(repoID, userID)

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

	if err := json.NewEncoder(w).Encode(issueslist); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "issues returned success")
}

func (GD *Http_Codehub) CloseIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID, ok := res.(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	oldIssue := models.Issue{}//достаточно скинуть просто id

	if err := json.NewDecoder(r.Body).Decode(&oldIssue); err != nil {
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

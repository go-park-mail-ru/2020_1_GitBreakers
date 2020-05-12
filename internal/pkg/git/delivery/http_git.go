package delivery

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"net/http"
)

type GitDelivery struct {
	UC     git.GitUseCaseI
	Logger *logger.SimpleLogger
	UserUC user.UCUser
}

//создать репак(id,name,description,private,owner)
func (GD *GitDelivery) CreateRepo(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := res.(int64)
	newRepo := &gitmodels.Repository{IsPublic: true}
	if err := easyjson.UnmarshalFromReader(r.Body, newRepo); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := govalidator.ValidateStruct(newRepo); err != nil {
		GD.Logger.HttpLogError(r.Context(), "govalidator", "ValidateStruct", errors.Cause(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := GD.UC.Create(userID, newRepo)
	switch {
	case errors.Is(err, entityerrors.AlreadyExist()):
		GD.Logger.HttpLogError(r.Context(), "repository", "Create", errors.Cause(err))
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		GD.Logger.HttpLogError(r.Context(), "repository", "Create", errors.Cause(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	GD.Logger.HttpInfo(r.Context(), "repo created", http.StatusCreated)
}

//данные репока(модельку скинуть(id,name,private,owner)
func (GD *GitDelivery) GetRepo(w http.ResponseWriter, r *http.Request) {
	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]
	userID := r.Context().Value("UserID")

	Repo, err := GD.UC.GetRepo(userName, repoName, GD.idToIntPointer(userID))

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), "access denied to repo", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(r.Context(), "repo does not exsist", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpInfo(r.Context(), err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(Repo, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "repo received", http.StatusOK)
}

//
////все репозитории юзера
func (GD *GitDelivery) GetRepoList(w http.ResponseWriter, r *http.Request) {
	userIDFromContext := r.Context().Value("UserID")
	userIDPtr := GD.idToIntPointer(userIDFromContext)

	userName := mux.Vars(r)["username"]

	if userName == "" {
		if userIDPtr == nil {
			GD.Logger.HttpInfo(r.Context(), "user doesn't exist", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		userModel, err := GD.UserUC.GetByID(*userIDPtr)
		if err != nil {
			GD.Logger.HttpInfo(r.Context(), "user doesn't exist", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		userName = userModel.Name
	}

	repo, err := GD.UC.GetRepoList(userName, userIDPtr)
	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), fmt.Sprintf("access denied for user=%s", userName), http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(repo, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "repolist received", http.StatusOK)
}

//
////ветки репака(просто названия и ссылки)
func (GD *GitDelivery) GetBranchList(w http.ResponseWriter, r *http.Request) {
	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]

	res := r.Context().Value("UserID")
	branches, err := GD.UC.GetBranchList(GD.idToIntPointer(res), userName, repoName)

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), "access denied", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(r.Context(), "not found repo with that fullname ", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogError(r.Context(), "", "", errors.Cause(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(branches, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "branches received", http.StatusOK)
}

//cписок коммитов для ветки
func (GD *GitDelivery) GetCommitsList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("UserID")
	//сжирает два параметра
	vars := mux.Vars(r)
	commitParams := &gitmodels.CommitRequest{
		UserLogin:  vars["username"],
		RepoName:   vars["reponame"],
		CommitHash: vars["hash"],
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(commitParams, r.URL.Query())

	if err != nil {
		GD.Logger.HttpLogError(r.Context(), "schema", "decoder.Decode", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := GD.UC.GetCommitsByCommitHash(*commitParams, GD.idToIntPointer(userID))

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), "access denied to this commits", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(r.Context(), "not found repo with that fullname ", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogError(r.Context(), "", "", errors.Cause(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "commits list received", http.StatusOK)
}

func (GD *GitDelivery) ShowFiles(w http.ResponseWriter, r *http.Request) {
	userIDpointer := GD.idToIntPointer(r.Context().Value("UserID"))

	vars := mux.Vars(r)
	showParams := gitmodels.FilesCommitRequest{
		UserName:    vars["username"],
		Reponame:    vars["reponame"],
		HashCommits: vars["hashcommits"],
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	err := decoder.Decode(&showParams, r.URL.Query())
	if err != nil {
		GD.Logger.HttpLogError(r.Context(), "schema", "decoder.Decode", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := GD.UC.FilesInCommitByPath(showParams, userIDpointer)

	if err != nil {
		res, err := GD.UC.GetFileByPath(showParams, userIDpointer)

		switch {
		case errors.Is(err, entityerrors.AccessDenied()):
			GD.Logger.HttpInfo(r.Context(), "access denied to this files", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		case errors.Is(err, entityerrors.DoesNotExist()):
			GD.Logger.HttpInfo(r.Context(), "not found this repo or path", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		case err != nil:
			GD.Logger.HttpLogError(r.Context(), "", "", errors.Cause(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, _, err := easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "files returned", http.StatusOK)
}

func (GD *GitDelivery) GetCommitsByBranchName(w http.ResponseWriter, r *http.Request) {
	userIDpointer := GD.idToIntPointer(r.Context().Value("UserID"))

	vars := mux.Vars(r)
	userName, repoName, branchName := vars["username"], vars["reponame"], vars["branchname"]

	res, err := GD.UC.GetCommitsByBranchName(userName, repoName, branchName, 0, 100, userIDpointer)

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), "access denied to this files", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(r.Context(), "not found this repo or path", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogError(r.Context(), "", "", errors.Cause(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "commits returned", http.StatusOK)
}

func (GD *GitDelivery) idToIntPointer(id interface{}) *int64 {
	intID, ok := id.(int64)
	if !ok {
		return nil
	} else {
		return &intID
	}
}

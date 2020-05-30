package delivery

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/http/helpers"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type GitDelivery struct {
	UC     git.GitUseCaseI
	Logger *logger.SimpleLogger
	UserUC user.UCUser
}

//создать репак(id,name,description,private,owner)
func (GD *GitDelivery) CreateRepo(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
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
	userID := r.Context().Value(models.UserIDKey)

	repo, err := GD.UC.GetRepo(userName, repoName, GD.idToIntPointer(userID))

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), "access denied to repo", http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(r.Context(), "repo does not exist", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpInfo(r.Context(), err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err = easyjson.MarshalToHTTPResponseWriter(repo, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "repo received", http.StatusOK)
}

func (GD *GitDelivery) DeleteRepo(w http.ResponseWriter, r *http.Request) {
	userIDFromContext := r.Context().Value(models.UserIDKey)
	userIDPtr := GD.idToIntPointer(userIDFromContext)
	if userIDPtr == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized repository deletion",
			http.StatusUnauthorized)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ownerID := *userIDPtr

	var repoForDelete gitmodels.Repository
	if err := easyjson.UnmarshalFromReader(r.Body, &repoForDelete); err != nil {
		GD.Logger.HttpInfo(
			r.Context(),
			fmt.Sprintf("user with id=%d send bad request", ownerID),
			http.StatusBadRequest,
		)
		return
	}

	repoName := repoForDelete.Name
	if repoName == "" {
		GD.Logger.HttpInfo(
			r.Context(),
			fmt.Sprintf("user with id=%d send bad request", ownerID),
			http.StatusBadRequest,
		)
		return
	}

	err := GD.UC.DeleteByOwnerID(ownerID, repoName)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(
			r.Context(),
			fmt.Sprintf(
				"user with id=%d try delete repo=%s, which not exists",
				ownerID,
				repoName,
			),
			http.StatusNotFound,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case errors.Is(err, entityerrors.Conflict()):
		GD.Logger.HttpInfo(
			r.Context(),
			fmt.Sprintf(
				"user with id=%d try delete repo=%s, but repo have opened pull requests",
				ownerID,
				repoName,
			),
			http.StatusConflict,
		)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	GD.Logger.HttpInfo(
		r.Context(),
		fmt.Sprintf(
			"user with id=%d deleted repo=%s",
			ownerID,
			repoName,
		),
		http.StatusOK,
	)
}

//
////все репозитории юзера
func (GD *GitDelivery) GetRepoList(w http.ResponseWriter, r *http.Request) {
	userIDFromContext := r.Context().Value(models.UserIDKey)
	userIDPtr := GD.idToIntPointer(userIDFromContext)

	offset, limit, err := helpers.ParseLimitAndOffset(r.URL.Query())
	if err != nil {
		offset, limit = helpers.DefaultOffset, helpers.DefaultLimit
	}

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

		userName = userModel.Login
	}

	repo, err := GD.UC.GetRepoList(userName, offset, limit, userIDPtr)
	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), fmt.Sprintf("access denied for user=%s", userName), http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if _, _, err = easyjson.MarshalToHTTPResponseWriter(repo, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "repolist received", http.StatusOK)
}

//
////все репозитории юзера
func (GD *GitDelivery) GetRepoListByUserID(w http.ResponseWriter, r *http.Request) {
	userIDFromContext := r.Context().Value(models.UserIDKey)
	userIDPtr := GD.idToIntPointer(userIDFromContext)

	slugID := mux.Vars(r)["id"]
	reposOwnerID, err := strconv.ParseInt(slugID, 10, 64)
	if err != nil {
		GD.Logger.HttpLogInfo(r.Context(), fmt.Sprintf("bad request: %v", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	offset, limit, err := helpers.ParseLimitAndOffset(r.URL.Query())
	if err != nil {
		offset, limit = helpers.DefaultOffset, helpers.DefaultLimit
	}

	reposOwnerModel, err := GD.UserUC.GetByID(reposOwnerID)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(r.Context(), "user doesn't exist", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	repo, err := GD.UC.GetRepoList(reposOwnerModel.Login, offset, limit, userIDPtr)
	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(r.Context(), fmt.Sprintf("access denied for userIDPtr=%v",
			userIDPtr), http.StatusForbidden)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(repo, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "repolist received", http.StatusOK)
}

//
////ветки репака(просто названия и ссылки)
func (GD *GitDelivery) GetBranchList(w http.ResponseWriter, r *http.Request) {
	userName, repoName := mux.Vars(r)["username"], mux.Vars(r)["reponame"]

	res := r.Context().Value(models.UserIDKey)
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

	if _, _, err = easyjson.MarshalToHTTPResponseWriter(branches, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "branches received", http.StatusOK)
}

//cписок коммитов для ветки
func (GD *GitDelivery) GetCommitsList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.UserIDKey)
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

	if _, _, err = easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "commits list received", http.StatusOK)
}

func (GD *GitDelivery) ShowFiles(w http.ResponseWriter, r *http.Request) {
	userIDpointer := GD.idToIntPointer(r.Context().Value(models.UserIDKey))

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
		// FIXME(nickeskov): what is this??? ignoring error, fix it
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
		case errors.Is(err, entityerrors.Invalid()):
			GD.Logger.HttpInfo(r.Context(), "trying get invalid entity", http.StatusNotAcceptable)
			w.WriteHeader(http.StatusNotAcceptable)
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

	if _, _, err = easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "files returned", http.StatusOK)
}

func (GD *GitDelivery) GetCommitsByBranchName(w http.ResponseWriter, r *http.Request) {
	userIDpointer := GD.idToIntPointer(r.Context().Value(models.UserIDKey))

	vars := mux.Vars(r)
	userName, repoName, branchName := vars["username"], vars["reponame"], vars["branchname"]

	offset, limit, err := helpers.ParseLimitAndOffset(r.URL.Query())
	if err != nil {
		offset, limit = helpers.DefaultOffset, helpers.DefaultLimit
	}

	res, err := GD.UC.GetCommitsByBranchName(userName, repoName, branchName, offset, limit, userIDpointer)

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

	if _, _, err = easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(r.Context(), "commits returned", http.StatusOK)
}

func (GD *GitDelivery) GetRepoHead(w http.ResponseWriter, r *http.Request) {
	userIDPointer := GD.idToIntPointer(r.Context().Value(models.UserIDKey))

	vars := mux.Vars(r)

	userName := vars["username"]
	repoName := vars["reponame"]

	res, err := GD.UC.GetRepoHead(userName, repoName, userIDPointer)
	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpInfo(
			r.Context(),
			fmt.Sprintf("access denied to repository=%s/%s", userName, repoName),
			http.StatusForbidden,
		)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpInfo(
			r.Context(),
			fmt.Sprintf("not found repository=%s/%s", userName, repoName),
			http.StatusNotFound,
		)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	case errors.Is(err, entityerrors.ContentEmpty()):
		w.WriteHeader(http.StatusNoContent)

	case err != nil:
		GD.Logger.HttpLogError(
			r.Context(),
			"git/delivery/http_git",
			"GetRepoHead",
			err,
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	default:
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(res, w); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			GD.Logger.HttpInfo(
				r.Context(),
				fmt.Sprintf("successfully returned HEAD=%s for repository=%s/%s",
					res.Name, userName, repoName),
				http.StatusOK,
			)
		}
	}
}

func (GD *GitDelivery) Fork(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int64)
	forkData := gitmodels.RepoFork{}
	if err := easyjson.UnmarshalFromReader(r.Body, &forkData); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	valid, err := govalidator.ValidateStruct(forkData)
	if !valid || err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if forkData.FromRepoID <= 0 {
		forkData.FromRepoID = -1
	}

	err = GD.UC.Fork(
		forkData.FromRepoID,
		forkData.FromAuthorName,
		forkData.FromRepoName,
		forkData.NewName,
		userID,
	)

	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.AlreadyExist()) || errors.Is(err, entityerrors.Conflict()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusConflict)
		return
	case err == nil:
		GD.Logger.HttpLogInfo(r.Context(), fmt.Sprintf("fork successfully created, forkData=%+v", forkData))
		w.WriteHeader(http.StatusCreated)
	default:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (GD *GitDelivery) GetBranchInfoByNames(w http.ResponseWriter, r *http.Request) {
	userIDPointer := GD.idToIntPointer(r.Context().Value(models.UserIDKey))

	vars := mux.Vars(r)

	userName, repoName, branchName := vars["username"], vars["reponame"], vars["branchname"]

	branch, err := GD.UC.GetBranchInfoByNames(userName, repoName, branchName, userIDPointer)
	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("access denied to repo=%s/%s branch=%s for user with id=%v",
				userName, repoName, branchName, userIDPointer))
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("not found repo=%s/%s branch=%s for user with id=%v",
				userName, repoName, branchName, userIDPointer))
		w.WriteHeader(http.StatusNotFound)

	case err == nil:
		GD.Logger.HttpLogInfo(r.Context(), fmt.Sprintf(
			"successfully get branchInfo repo=%s/%s branch=%s for user with id=%v",
			userName, repoName, branchName, userIDPointer))
		if _, _, err = easyjson.MarshalToHTTPResponseWriter(branch, w); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	default:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (GD *GitDelivery) GetFileContentByBranch(w http.ResponseWriter, r *http.Request) {
	userIDPointer := GD.idToIntPointer(r.Context().Value(models.UserIDKey))

	vars := mux.Vars(r)

	userName, repoName, branchName := vars["username"], vars["reponame"], vars["branchname"]
	filePath := vars["path"]

	content, err := GD.UC.GetFileContentByBranch(userName, repoName, branchName, filePath, userIDPointer)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("not dound in repo=%s/%s in branch=%s file=%s for user with id=%v",
				userName, repoName, branchName, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("access denied in repo=%s/%s in branch=%s file=%s for user with id=%v",
				userName, repoName, branchName, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

	case errors.Is(err, entityerrors.Invalid()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("bad request in repo=%s/%s in branch=%s file=%s for user with id=%v",
				userName, repoName, branchName, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

	case errors.Is(err, entityerrors.TooLarge()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("too large file in repo=%s/%s in branch=%s file=%s for user with id=%v",
				userName, repoName, branchName, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)

	case err == nil:
		contentType := http.DetectContentType(content)
		w.Header().Set("Content-Type", contentType)

		if _, err := w.Write(content); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)

	default:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (GD *GitDelivery) GetFileContentByCommitHash(w http.ResponseWriter, r *http.Request) {
	userIDPointer := GD.idToIntPointer(r.Context().Value(models.UserIDKey))

	vars := mux.Vars(r)

	userName, repoName, commitHash := vars["username"], vars["reponame"], vars["hash"]
	filePath := vars["path"]

	content, err := GD.UC.GetFileContentByCommitHash(userName, repoName, commitHash, filePath, userIDPointer)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("not dound in repo=%s/%s in hash=%s file=%s for user with id=%v",
				userName, repoName, commitHash, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("access denied in repo=%s/%s in hash=%s file=%s for user with id=%v",
				userName, repoName, commitHash, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

	case errors.Is(err, entityerrors.Invalid()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("bad request in repo=%s/%s in hash=%s file=%s for user with id=%v",
				userName, repoName, commitHash, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

	case errors.Is(err, entityerrors.TooLarge()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("too large file in repo=%s/%s in hash=%s file=%s for user with id=%v",
				userName, repoName, commitHash, filePath, userIDPointer))

		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)

	case err == nil:
		contentType := http.DetectContentType(content)
		w.Header().Set("Content-Type", contentType)

		if _, err := w.Write(content); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)

	default:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (GD *GitDelivery) idToIntPointer(id interface{}) *int64 {
	intID, ok := id.(int64)
	if !ok {
		return nil
	} else {
		return &intID
	}
}

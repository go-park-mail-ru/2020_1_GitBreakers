package http

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

type HttpCodehub struct {
	Logger     *logger.SimpleLogger
	CodeHubUC  codehub.UCCodeHubI
	NewsClient interfaces.NewsClientI
	UserClient interfaces.UserClientI
}

func (GD *HttpCodehub) ModifyStar(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := res.(int64)

	newStar := models.Star{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newStar); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newStar.AuthorID = userID
	newStar.RepoID = int64(repoID)

	err = GD.CodeHubUC.ModifyStar(newStar)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()) || errors.Is(err, entityerrors.AlreadyExist()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "star modify success")

}

func (GD *HttpCodehub) StarredRepos(w http.ResponseWriter, r *http.Request) {
	userLogin := mux.Vars(r)["login"]

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	user, err := GD.UserClient.GetByLogin(userLogin)
	//todo чуть лучше обработку сделать
	switch {
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	repolist, err := GD.CodeHubUC.GetStarredRepos(user.ID, int64(limit), int64(offset))
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
func (GD *HttpCodehub) UserWithStar(w http.ResponseWriter, r *http.Request) {
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	userlist, err := GD.CodeHubUC.GetUserStaredList(int64(repoID), int64(limit), int64(offset))

	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(userlist, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "userlist got success")
}

func (GD *HttpCodehub) NewIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := res.(int64)
	newIssue := models.Issue{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newIssue.AuthorID = userID
	newIssue.RepoID = int64(repoID)

	err = GD.CodeHubUC.CreateIssue(newIssue)

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.AlreadyExist()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusConflict)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
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
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	userID := res.(int64)

	newIssue := models.Issue{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newIssue); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldIssue, err := GD.CodeHubUC.GetIssue(newIssue.ID, userID)

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
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
		oldIssue.Title = newIssue.Title
	}
	if govalidator.IsByteLength(newIssue.Label, 0, 100) {
		oldIssue.Label = newIssue.Label
	}
	oldIssue.AuthorID = userID
	err = GD.CodeHubUC.UpdateIssue(oldIssue)

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
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
	res := r.Context().Value(models.UserIDKey)
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest) //скинули строку, а не число
		return
	}

	userID, ok := res.(int64)
	if !ok {
		userID = -1 //соответствует неавторизованному юзеру
	}

	issueslist, err := GD.CodeHubUC.GetIssuesList(int64(repoID), userID, 100, 0)

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
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

	GD.Logger.HttpLogInfo(r.Context(), "issues returned success")
}

func (GD *HttpCodehub) CloseIssue(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
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
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
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
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	news, err := GD.NewsClient.GetNews(int64(repoID), res.(int64), int64(limit), int64(offset))

	switch {
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogInfo(r.Context(), "news access denied")
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogInfo(r.Context(), "news does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(news, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpLogInfo(r.Context(), "news got success")
}

func (GD *HttpCodehub) Search(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int64)

	query := r.URL.Query().Get("query")
	params := mux.Vars(r)["params"]
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	data, err := GD.CodeHubUC.Search(query, params, int64(limit), int64(offset), userID)
	switch {
	case errors.Is(err, entityerrors.Invalid()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repolist, ok := data.(models.RepoSet)
	if ok {
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(repolist, w); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	userlist, ok := data.(models.UserSet)
	if ok {
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(userlist, w); err != nil {
			GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (GD *HttpCodehub) CreatePullReq(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int64)
	plModel := models.PullRequest{}

	if err := easyjson.UnmarshalFromReader(r.Body, &plModel); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	plModel.AuthorId = &userID
	isCorrect, err := govalidator.ValidateStruct(plModel)
	if !isCorrect || err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, err := GD.CodeHubUC.CreatePL(plModel)
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
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(pr, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(
		r.Context(),
		fmt.Sprintf("successfully created pr=%+v", pr),
		http.StatusCreated,
	)

	w.WriteHeader(http.StatusCreated)
}

func (GD *HttpCodehub) GetPullReqList(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	repoID, err := strconv.Atoi(mux.Vars(r)["repoID"])
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	direction := mux.Vars(r)["direction"]
	repoModel := gitmodels.Repository{ID: int64(repoID)}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	var PLlist models.PullReqSet
	switch direction {
	case "in":
		PLlist, err = GD.CodeHubUC.GetPLIn(repoModel, int64(limit), int64(offset))
	case "out":
		PLlist, err = GD.CodeHubUC.GetPLOut(repoModel, int64(limit), int64(offset))
	default:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(PLlist, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (GD *HttpCodehub) ApproveMerge(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int64)
	plModel := models.PullRequest{}
	if err := easyjson.UnmarshalFromReader(r.Body, &plModel); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	plModel.AuthorId = &userID
	err := GD.CodeHubUC.ApprovePL(plModel, userID)

	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case errors.Is(err, entityerrors.Conflict()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
func (GD *HttpCodehub) UndoPullReq(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int64)

	plModel := models.PullRequest{}
	if err := easyjson.UnmarshalFromReader(r.Body, &plModel); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	plModel.AuthorId = &userID

	err := GD.CodeHubUC.ClosePL(plModel, userID)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case errors.Is(err, entityerrors.AccessDenied()):
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusForbidden)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
func (GD *HttpCodehub) GetAllPLFromUser(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(models.UserIDKey)
	if res == nil {
		GD.Logger.HttpInfo(r.Context(), "unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int64)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	pllist, err := GD.CodeHubUC.GetAllMRUser(userID, int64(limit), int64(offset))
	if err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(pllist, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (GD *HttpCodehub) GetMRByID(w http.ResponseWriter, r *http.Request) {
	strMRID := mux.Vars(r)["id"]

	mrID, err := strconv.ParseInt(strMRID, 10, 64)
	if err != nil {
		GD.Logger.HttpLogInfo(r.Context(), fmt.Sprintf("bad request: %v", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pr, err := GD.CodeHubUC.GetMRByID(mrID)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("merge request with id=%d does not exist", mrID))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(pr, w); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	GD.Logger.HttpInfo(
		r.Context(),
		fmt.Sprintf("successfully get mr with id=%d", mrID),
		http.StatusOK,
	)
}

func (GD *HttpCodehub) GetMRDiffByID(w http.ResponseWriter, r *http.Request) {
	strMRID := mux.Vars(r)["id"]

	mrID, err := strconv.ParseInt(strMRID, 10, 64)
	if err != nil {
		GD.Logger.HttpLogInfo(r.Context(), fmt.Sprintf("bad request: %v", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	diff, err := GD.CodeHubUC.GetMRDiffByID(mrID)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		GD.Logger.HttpLogInfo(r.Context(),
			fmt.Sprintf("merge request with id=%d does not exist", mrID))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case err != nil:
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(diff)); err != nil {
		GD.Logger.HttpLogCallerError(r.Context(), *GD, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set plain text, no json
	w.Header().Set("Content-Type", "plain/text")

	GD.Logger.HttpInfo(
		r.Context(),
		fmt.Sprintf("successfully get diff for mr with id=%d", mrID),
		http.StatusOK,
	)
}

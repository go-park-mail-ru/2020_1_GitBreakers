package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/gitserver"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	permTypes "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sosedoff/gitkit"
	"net/http"
)

const (
	GitReceivePackService      = "git-receive-pack"
	GitUploadPackService       = "git-upload-pack"
	OwnerLoginMuxParameter     = "OwnerLogin"
	RepositoryNameMuxParameter = "RepositoryName"
)

type GitServerDelivery struct {
	UseCase gitserver.UseCase
	Logger  logger.Logger
}

type repoBasicInfo struct {
	ownerLogin string
	repoName   string
}

func CreateGitIfoRefsMiddleware(delivery GitServerDelivery) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodGet {
				http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
				delivery.Logger.HttpLogWarning(r.Context(), "gitserver/delivery",
					"GitInfoRefsMiddleware", "called not implemented method")
				return
			}

			needNextHandler, err := gitInfoRefsHandler(w, r, delivery)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				delivery.Logger.HttpLogCallerError(r.Context(), delivery, err)
				return
			}

			if needNextHandler {
				next.ServeHTTP(w, r)
			}

		})
	}
}

func CreateGitUploadPackMiddleware(delivery GitServerDelivery) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
				delivery.Logger.HttpLogWarning(r.Context(), "gitserver/delivery",
					"GitUploadPackMiddleware", "called not implemented method")
				return
			}

			needNextHandler, err := gitUploadPackHandler(w, r, delivery)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				delivery.Logger.HttpLogCallerError(r.Context(), delivery, err)
				return
			}

			if needNextHandler {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func CreateGitReceivePackMiddleware(delivery GitServerDelivery) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
				delivery.Logger.HttpLogWarning(r.Context(), "gitserver/delivery",
					"GitReceivePackMiddleware", "called not implemented method")
				return
			}

			needNextHandler, err := gitReceivePackHandler(w, r, delivery)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				delivery.Logger.HttpLogCallerError(r.Context(), delivery, err)
				return
			}

			if needNextHandler {
				next.ServeHTTP(w, r)

				// Update merge requests associated with this repo
				repoInfo := getRepoInfo(r)
				if err := updateMrStatuses(repoInfo, delivery); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					delivery.Logger.HttpLogCallerError(r.Context(), delivery, err)
				}
			}
		})
	}
}

func gitInfoRefsHandler(w http.ResponseWriter, r *http.Request, delivery GitServerDelivery) (bool, error) {
	service := r.URL.Query().Get("service")
	if service != GitUploadPackService && service != GitReceivePackService {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false, nil
	}

	switch service {
	case GitUploadPackService:
		return gitUploadPackHandler(w, r, delivery)
	case GitReceivePackService:
		return gitReceivePackHandler(w, r, delivery)
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return false, nil
}

func gitUploadPackHandler(w http.ResponseWriter, r *http.Request, delivery GitServerDelivery) (bool, error) {
	repoInfo := getRepoInfo(r)

	haveReadAccess, err := delivery.UseCase.CheckGitRepositoryReadAccess(nil,
		repoInfo.ownerLogin, repoInfo.repoName)
	if err != nil {
		return false, err
	}
	if haveReadAccess {
		return true, nil
	}

	perm, needNextHandle, err := getPermissionsForRepo(w, r, repoInfo, delivery)
	if err != nil {
		return false, err
	}
	if !needNextHandle {
		return false, nil
	}

	if perm == permTypes.NoAccess() {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return false, nil
	}

	return needNextHandle, nil
}

func gitReceivePackHandler(w http.ResponseWriter, r *http.Request, delivery GitServerDelivery) (bool, error) {
	repoInfo := getRepoInfo(r)

	perm, needNextHandle, err := getPermissionsForRepo(w, r, repoInfo, delivery)
	if err != nil {
		return false, err
	}
	if !needNextHandle {
		return false, nil
	}

	if perm == permTypes.NoAccess() || perm == permTypes.ReadAccess() {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return false, nil
	}

	return needNextHandle, nil
}

func updateMrStatuses(repoInfo repoBasicInfo, delivery GitServerDelivery) error {
	return delivery.UseCase.UpdateMergeRequestsStatuses(repoInfo.ownerLogin, repoInfo.repoName)
}

func getPermissionsForRepo(w http.ResponseWriter, r *http.Request, repoInfo repoBasicInfo,
	delivery GitServerDelivery) (permTypes.Permission, bool, error) {

	cred, ok := processCredentials(w, r)
	if !ok {
		delivery.Logger.HttpLogInfo(r.Context(), "user not authenticated")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return permTypes.NoAccess(), false, nil
	}

	// check user credentials
	isValidCred, err := delivery.UseCase.CheckUserPassword(cred.Username, cred.Password)

	switch {
	case errors.Cause(err) == entityerrors.DoesNotExist():
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return permTypes.NoAccess(), false, nil
	case err != nil:
		return permTypes.NoAccess(), false, err
	}

	// if credentials is invalid in any case send http.StatusForbidden
	if !isValidCred {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return permTypes.NoAccess(), false, nil
	}

	// get user
	user, err := delivery.UseCase.GetUser(cred.Username)
	if err != nil {
		return permTypes.NoAccess(), false, err
	}

	// get permissions for user
	perm, err := delivery.UseCase.GetGitRepositoryPermission(&user.ID, repoInfo.ownerLogin,
		repoInfo.repoName)
	if err != nil {
		return permTypes.NoAccess(), false, err
	}

	return perm, true, nil
}

func processCredentials(w http.ResponseWriter, r *http.Request) (gitkit.Credential, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Send auth request header
		w.Header().Set("WWW-Authenticate", `Basic realm=""`)
		return gitkit.Credential{}, false
	}

	return getCredential(r)
}

func getCredential(request *http.Request) (cred gitkit.Credential, ok bool) {
	cred.Username, cred.Password, ok = request.BasicAuth()
	return cred, ok
}

func getRepoInfo(request *http.Request) repoBasicInfo {
	ownerLogin := mux.Vars(request)[OwnerLoginMuxParameter]
	repositoryName := mux.Vars(request)[RepositoryNameMuxParameter]

	return repoBasicInfo{
		ownerLogin: ownerLogin,
		repoName:   repositoryName,
	}
}

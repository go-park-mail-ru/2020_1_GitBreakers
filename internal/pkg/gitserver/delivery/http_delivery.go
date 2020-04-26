package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/gitserver"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	permTypes "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/pkg/errors"
	"github.com/sosedoff/gitkit"
	"net/http"
	"strings"
)

const (
	pkgAnchor      = 0
	gitReceivePack = "git-receive-pack"
	gitUploadPack  = "git-upload-pack"
)

type GitServerDelivery struct {
	UseCase gitserver.UseCase
	Logger  logger.Logger
}

type repoBasicInfo struct {
	ownerLogin string
	repoName   string
}

func CreateMainAuthMiddleware(delivery GitServerDelivery) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var needNextHandler bool
			var err error

			switch r.Method {
			case http.MethodGet:
				if strings.HasSuffix(r.URL.Path, "/info/refs") {
					needNextHandler, err = gitInfoRefsHandler(w, r, delivery)
				}
			case http.MethodPost:

				switch true {
				case strings.HasSuffix(r.URL.Path, "/"+gitUploadPack):
					needNextHandler, err = gitUploadPackHandler(w, r, delivery)
				case strings.HasSuffix(r.URL.Path, "/"+gitReceivePack):
					needNextHandler, err = gitReceivePackHandler(w, r, delivery)
				}
			default:
				http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
				delivery.Logger.HttpLogWarning(r.Context(), "gitserver/delivery",
					"MainAuthMiddleware", "called not implemented method")
				return
			}

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				delivery.Logger.HttpLogCallerError(r.Context(), pkgAnchor, err)
				return
			}

			if needNextHandler {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func gitInfoRefsHandler(w http.ResponseWriter, r *http.Request, delivery GitServerDelivery) (bool, error) {
	service := r.Header.Get("service")
	if service != gitUploadPack && service != gitReceivePack {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false, nil
	}

	switch service {
	case gitUploadPack:
		return gitUploadPackHandler(w, r, delivery)
	case gitReceivePack:
		return gitReceivePackHandler(w, r, delivery)
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return false, nil
}

func gitUploadPackHandler(w http.ResponseWriter, r *http.Request, delivery GitServerDelivery) (bool, error) {
	repoInfo := getRepoInfo(gitUploadPack, r)

	perm, needNextHandle, err := GetPermissionsForRepo(w, r, repoInfo, delivery)
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

func gitReceivePackHandler(w http.ResponseWriter, r *http.Request, delivery GitServerDelivery) (bool, error) {
	repoInfo := getRepoInfo(gitReceivePack, r)

	haveReadAccess, err := delivery.UseCase.CheckGitRepositoryReadAccess(nil,
		repoInfo.ownerLogin, repoInfo.repoName)
	if err != nil {
		return false, err
	}
	if haveReadAccess {
		return true, nil
	}

	perm, needNextHandle, err := GetPermissionsForRepo(w, r, repoInfo, delivery)
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

func GetPermissionsForRepo(w http.ResponseWriter, r *http.Request, repoInfo repoBasicInfo,
	delivery GitServerDelivery) (permTypes.Permission, bool, error) {

	cred, ok := processCredentials(w, r)
	if !ok {
		delivery.Logger.HttpLogInfo(r.Context(), "user not authenticated")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return permTypes.NoAccess(), false, nil
	}

	// check user credentials
	isValidCred, err := delivery.UseCase.CheckUserPassword(cred.Username, cred.Password)

	switch true {
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

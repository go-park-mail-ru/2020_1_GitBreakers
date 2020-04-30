package http

import (
	"bytes"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

type UserHttp struct {
	SessHttp session.SessDelivery
	Logger   *logger.SimpleLogger
	UClient  interfaces.UserClientI
}

func (UsHttp *UserHttp) Create(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		UsHttp.Logger.HttpInfo(r.Context(), "already authorized", http.StatusNotAcceptable)
		return
	}

	User := &models.User{}

	if err := easyjson.UnmarshalFromReader(r.Body, User); err != nil {
		UsHttp.Logger.HttpLogError(r.Context(), "json", "decode", errors.Cause(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := govalidator.ValidateStruct(User); err != nil {
		UsHttp.Logger.HttpLogError(r.Context(), "validator", "validate struct", errors.Cause(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := UsHttp.UClient.Create(*User)
	switch {
	case err == entityerrors.AlreadyExist():
		UsHttp.Logger.HttpLogError(r.Context(), "user", " Create", errors.Cause(err))
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		UsHttp.Logger.HttpLogError(r.Context(), "user", " Create", errors.Cause(err))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	UsHttp.Logger.HttpLogInfo(r.Context(), "user created in postgres")

	UserFromDB, err := UsHttp.UClient.GetByLogin(User.Login)
	if err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := UsHttp.SessHttp.Create(UserFromDB.ID)
	if err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	UsHttp.Logger.HttpInfo(r.Context(), "user created successful", http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}

func (UsHttp *UserHttp) Update(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId := res.(int64)
	newUserData := models.User{}

	if err := easyjson.UnmarshalFromReader(r.Body, &newUserData); err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := UsHttp.UClient.Update(userId, newUserData)
	switch {
	case err == entityerrors.AlreadyExist():
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (UsHttp *UserHttp) Login(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res != nil {
		UsHttp.Logger.HttpInfo(r.Context(), "already authorized", http.StatusNotAcceptable)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	input := &models.SignInForm{}
	if err := easyjson.UnmarshalFromReader(r.Body, input); err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		UsHttp.Logger.HttpLogError(r.Context(), "govalidator", "validate struct", errors.Cause(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	User, err := UsHttp.UClient.GetByLogin(input.Login)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		w.WriteHeader(http.StatusNotFound)
		UsHttp.Logger.HttpLogError(r.Context(), "", "GetByLogin", errors.Cause(err))
		return
	case err != nil:
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isUser, err := UsHttp.UClient.CheckPass(User.Login, input.Password)
	if err != nil || !isUser {
		UsHttp.Logger.HttpLogWarning(r.Context(), " ", "CheckPass", errors.Cause(err).Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookie, err := UsHttp.SessHttp.Create(User.ID)
	if err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	UsHttp.Logger.HttpInfo(r.Context(), "user authorized", http.StatusOK)
}

func (UsHttp *UserHttp) Logout(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res == nil {
		UsHttp.Logger.HttpInfo(r.Context(), "user unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		UsHttp.Logger.HttpInfo(r.Context(), "cookies doesn't exsist", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := UsHttp.SessHttp.Delete(cookie.Value); err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	UsHttp.Logger.HttpInfo(r.Context(), "user logouted", http.StatusOK)
}

func (UsHttp *UserHttp) GetInfo(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		UsHttp.Logger.HttpInfo(r.Context(), "user unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := res.(int64)
	User, err := UsHttp.UClient.GetByID(userID)

	switch {
	case err == entityerrors.DoesNotExist():
		UsHttp.Logger.HttpInfo(r.Context(), "user not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(User, w); err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	UsHttp.Logger.HttpInfo(r.Context(), "info received", http.StatusOK)
}

func (UsHttp *UserHttp) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res == nil {
		UsHttp.Logger.HttpInfo(r.Context(), "user unauthorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//6MB max
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		UsHttp.Logger.HttpLogError(r.Context(), "http", "ParseMultipartForm", errors.Cause(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image, header, err := r.FormFile("avatar")
	if err != nil {
		UsHttp.Logger.HttpLogError(r.Context(), "http", "FormFile", errors.Cause(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if err := image.Close(); err != nil {
			UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	binaryImage := bytes.NewBuffer(nil)
	if _, err := io.Copy(binaryImage, image); err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	currUser := r.Context().Value("UserID").(int64)

	err = UsHttp.UClient.UploadAvatar(currUser, header.Filename, binaryImage.Bytes(), int64(binaryImage.Len()))
	if err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	UsHttp.Logger.HttpLogInfo(r.Context(), "new avatar loaded")
}

func (UsHttp *UserHttp) GetInfoByLogin(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["login"]
	userData, err := UsHttp.UClient.GetByLogin(slug)

	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		UsHttp.Logger.HttpLogWarning(r.Context(), "user", "GetByLogin", errors.Cause(err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(userData, w); err != nil {
		UsHttp.Logger.HttpLogCallerError(r.Context(), *UsHttp, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	UsHttp.Logger.HttpInfo(r.Context(), "info received", http.StatusOK)
}

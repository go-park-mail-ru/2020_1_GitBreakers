package delivery

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type UserHttp struct {
	SessHttp session.SessDelivery
	UserUC   user.UCUser
	Logger   *logger.SimpleLogger
}

func (UsHttp *UserHttp) Create(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res != nil {
		w.WriteHeader(http.StatusBadRequest)
		UsHttp.Logger.HttpInfo(r.Context(), "already authorized", http.StatusBadRequest)
		return
	}
	User := &models.User{}
	err := json.NewDecoder(r.Body).Decode(User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		UsHttp.Logger.HttpLogWarning(r.Context(), "validator", "", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := UsHttp.UserUC.Create(*User); err != nil {
		UsHttp.Logger.HttpLogWarning(r.Context(), "", "", "не создали юзера")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UsHttp.Logger.HttpLogInfo(r.Context(), "создали юзера в постгресе")
	UserFromDB, err := UsHttp.UserUC.GetByLogin(User.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie, err := UsHttp.SessHttp.Create(UserFromDB.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		UsHttp.Logger.HttpInfo(r.Context(), "не создали сессию", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	UsHttp.Logger.HttpInfo(r.Context(), "создали юзера", http.StatusCreated)
}

func (UsHttp *UserHttp) Update(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := res.(int)
	newUserData := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&newUserData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := UsHttp.UserUC.Update(userId, newUserData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (UsHttp *UserHttp) Login(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res != nil {
		UsHttp.Logger.HttpInfo(r.Context(), "уже авторизован", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input := &models.SignInForm{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		UsHttp.Logger.HttpInfo(r.Context(), "невалидный json", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		UsHttp.Logger.HttpInfo(r.Context(), err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	User, err := UsHttp.UserUC.GetByLogin(input.Login)
	if err != nil {
		UsHttp.Logger.HttpInfo(r.Context(), "не нашли такой логин", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	isUser, err := UsHttp.UserUC.CheckPass(User, input.Password)
	//выйдем если хреновый пароль
	if err != nil || !isUser {
		UsHttp.Logger.HttpInfo(r.Context(), "неверный пароль ", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//создали сессию челику
	cookie, err := UsHttp.SessHttp.Create(User.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	UsHttp.Logger.HttpInfo(r.Context(), "авторизовали юзера", http.StatusOK)
}

func (UsHttp *UserHttp) Logout(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie || cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := UsHttp.SessHttp.Delete(cookie.Value); err != nil {
		UsHttp.Logger.HttpInfo(r.Context(), "ошибка с удалением куки", http.StatusInternalServerError)
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
}

func (UsHttp *UserHttp) GetInfo(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value("UserID")
	if res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := res.(int)
	User, err := UsHttp.UserUC.GetByID(userID)
	if err != nil {
		UsHttp.Logger.HttpInfo(r.Context(), "ошибка при получении id юзера", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(User); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
func (UsHttp *UserHttp) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if res := r.Context().Value("UserID"); res == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//6MB max
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	image, header, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if err := image.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	currUser := r.Context().Value("UserID")
	User, err := UsHttp.UserUC.GetByID(currUser.(int))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := UsHttp.UserUC.UploadAvatar(User, header, image); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (UsHttp *UserHttp) GetInfoByLogin(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["login"]
	userData, err := UsHttp.UserUC.GetByLogin(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//todo не очень круто
	userData.Password = ""
	if err := json.NewEncoder(w).Encode(userData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

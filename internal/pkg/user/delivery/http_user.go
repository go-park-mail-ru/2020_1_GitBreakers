package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"net/http"
)

type UserHttp struct {
	SessHttp session.SessDelivery
	UserUC   user.UCUser
}

func (UsHttp *UserHttp) Create(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAuth").(bool) {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (UsHttp *UserHttp) Update(w http.ResponseWriter, r *http.Request) {
	//call userUC
}

func (UsHttp *UserHttp) Login(w http.ResponseWriter, r *http.Request) {
	//залогинить челика
}

func (UsHttp *UserHttp) Logout(w http.ResponseWriter, r *http.Request) {
	//разлогинить чела
	//call sesshttp.delete()
}

func (UsHttp *UserHttp) GetInfo(w http.ResponseWriter, r *http.Request) {
	//get данные юзера
	//userUC.getByLogin()
}

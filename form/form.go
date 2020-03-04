package form

import "github.com/gorilla/sessions"

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type Repository struct {
	Name          string `json:"name"`
	Fullname      string `json:"full_name"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	Fork          bool   `json:"fork"`
	Owner         string `json:"owner"`
	DefaultBranch string `json:"default_branch"`
}
type SignupForm struct {
	Avatar    string `json:"avatar"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
	Bio       string `json:"bio"`
	URL       string `json:"url"`
	Followers uint   `json:"followers"`
	Following uint   `json:"following"`
}

//хранилище сессий,
var StoreSession = sessions.NewCookieStore([]byte("top-secret"))

//
var UserSlice = map[string]SignupForm{
	"antonelagin": {
		Login:     "AntonElagin",
		Avatar:    "https://avatars.mds.yandex.net/get-pdb/1964870/cfaa9b42-d74b-40f7-93cc-1d777bb5654c/s1200",
		Email:     "ANton ELagin",
		Bio:       "kek lol opps, long sfddsfdsfd dfdsfs dsfds sdfdsf sdfsd sdfd",
		URL:       "google.com",
		Followers: 1,
		Following: 2,
		Password:  []byte("fsffds"),
	},
	"Kekmdasher": {
		Login:     "Kekmdasher",
		Avatar:    "https://avatars.mds.yandex.net/get-pdb/1964870/cfaa9b42-d74b-40f7-93cc-1d777bb5654c/s1200",
		Email:     "heheh@mail.ru",
		Bio:       "kek lolsf sdfsd sdfd",
		URL:       "google.com",
		Followers: 1,
		Following: 2,
		Password:  []byte("fsffds"),
	},
}

package routes

import (
	"../handlers"
	"../models"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(config *AppConfig, ctx *handlers.StoresContext) *mux.Router {
	r := mux.NewRouter()

	//csrf.Secure(!config.IsDebug)
	///csrfMiddleware := csrf.Protect(config.CSRFTAuthKey)

	authRouter := r.PathPrefix("/auth").Subrouter()
	//authRouter.Use(csrfMiddleware)
	//authRouter.Use(mux.CORSMethodMiddleware(authRouter))

	authRouter.HandleFunc("/login", ctx.Login).Methods("POST")
	authRouter.HandleFunc("/signup", ctx.Signup).Methods("POST")
	authRouter.HandleFunc("/logout", ctx.Logout).Methods("POST")
	userRouter := r.PathPrefix("/user").Subrouter()
	//userRouter.Use(csrfMiddleware)
	userRouter.Use(mux.CORSMethodMiddleware(userRouter))

	userRouter.HandleFunc("/whoami", ctx.GetUser).Methods("POST")

	r.HandleFunc("/new/repository", models.NewRepository).Methods("POST")
	r.HandleFunc("/settings/profile", models.UpdateProfile).Methods("POST")
	r.HandleFunc("/profile/{login}", models.GetProfile).Methods("GET")
	r.HandleFunc("/repository/{login}", models.GetRepositoryList).Methods("GET")
	r.HandleFunc("/settings/avatar", models.UploadAvatar).Methods("POST")
	r.HandleFunc("/settings/profile", models.UpdateProfileOption).Methods("OPTIONS")

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", fs))

	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler))

	r.HandleFunc("/repository", models.GetRepository).Methods(http.MethodGet, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))

	return r
}

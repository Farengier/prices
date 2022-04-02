package admin

import (
	"context"
	"github.com/Farengier/prices/src/pkg/server/deps"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type ctrl struct {
	l  zerolog.Logger
	db interface{}
	t  deps.Template
}

func Ctrl(l zerolog.Logger, db interface{}, t deps.Template) *ctrl {
	return &ctrl{l: l, db: db, t: t}
}

func (c ctrl) Init(r *mux.Router) {
	sr := r.PathPrefix("/admin/").Subrouter()
	sr.Use(authMiddleware)
	sr.HandleFunc("/", c.mainPage)
}

func matcher(r *http.Request, rm *mux.RouteMatch) bool {
	if c, err := r.Cookie("admin_secret"); err == nil {
		return c.Value == "superSecret"
	}
	return false
}

func (c ctrl) mainPage(w http.ResponseWriter, r *http.Request) {
	isAdmin := r.Context().Value("is_admin").(bool)
	p := c.t.NewPage("Admin")
	if isAdmin {
		p.AddContent([]byte("Welcome to admin panel"))
		p.Render(w)
		return
	}

	p.AddScript("document.addEventListener('DOMContentLoaded', function() {" +
		"document.getElementsByClassName(\"js-secret-btn\")[0].onclick = function(){" +
		"document.cookie = \"admin_secret=\"+document.getElementsByClassName(\"js-secret\")[0].value;" +
		"window.location.reload()" +
		"}" +
		"}, false);",
	)
	p.AddContent([]byte("Welcome to admin panel<br/>"))
	p.AddContent([]byte("<label>Secret code:<input type=\"password\" class=\"js-secret\"></input></label>"))
	p.AddContent([]byte("<button class=\"js-secret-btn\">Ok</button>"))
	p.Render(w)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := false
		if c, err := r.Cookie("admin_secret"); err == nil && c.Value == "superSecret" {
			isAdmin = true
		}
		nr := r.WithContext(context.WithValue(r.Context(), "is_admin", isAdmin))
		next.ServeHTTP(w, nr)
	})
}

//
//func Middleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token := r.Header.Get("X-Session-Token")
//
//		if user, found := amw.tokenUsers[token]; found {
//			// We found the token in our map
//			log.Printf("Authenticated user %s\n", user)
//			// Pass down the request to the next middleware (or final handler)
//			next.ServeHTTP(w, r)
//		} else {
//			// Write an error and stop the handler chain
//			http.Error(w, "Forbidden", http.StatusForbidden)
//		}
//	})
//}

package server

import (
	"gap/internal/userservice"
	"net/http"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	sessionKey, err := userservice.CreateUserSession(r.Context(), s.db.Q(), r.FormValue("username"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userservice.SetCookie(w, sessionKey)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

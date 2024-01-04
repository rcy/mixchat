package server

import "net/http"

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("TODO: login not implemented"))
}

package server

import (
	"net/http"
)

func (s *Server) nextHandler(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("/media/Icq_old_sound-6iCPIUGnHQ8.ogg"))
	w.Write([]byte("/media/Squid_-_The_Cleaner_Official_Audio-T3XY6FPrbuM.ogg"))
}

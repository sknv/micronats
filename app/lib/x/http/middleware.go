package http

import (
	"log"
	"net/http"
)

const (
	errSmoke = "something went wrong"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if IsHandlerAborted(rvr) {
					return // response is already flushed
				}
				log.Print("[ERROR] panic: ", rvr)
				http.Error(w, errSmoke, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

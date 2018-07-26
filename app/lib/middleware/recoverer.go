package middleware

import (
	"net/http"

	xhttp "github.com/sknv/micronats/app/lib/net/http"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if xhttp.IsHandlerAborted(rvr) {
					return // response is already flushed
				}
				panic(rvr) // throw the error
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

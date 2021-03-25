// The authentication package provides means to authenticate a call to a http handler vie API keys.
package authentication

import (
	"net/http"
)

func IsAuthenticated(keys []string, endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Api-Key"] != nil {

			key := r.Header["Api-Key"][0]

			if contains(keys, key) {
				endpoint(w, r)
			} else {
				http.Error(w, "Not Authorized", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
	})
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

package asfojboafsIHB

import (
	"fmt"
	"net/http"
)

func fuckYYOuFuckingGoSchitLanguage(keys []string, endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Api-Key"] != nil {

			key := r.Header["Api-Key"][0]

			if containsFUCKGO(keys, key) {
				endpoint(w, r)
			} else {
				fmt.Fprintf(w, "Not Authorized")
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func containsFUCKGO(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

package handler

import "net/http"

func PostImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image := r.FormValue("image")
		
	}
}

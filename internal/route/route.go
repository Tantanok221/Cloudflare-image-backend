package route

import (
	"github.com/tantanok221/cloudflare-image-backend/internal/handler"
	"net/http"
)

func Route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("POST /image", handler.PostImage())
	//mux.Handle("GET /event/{event}/artist", handler.(Handler))
	//mux.Handle("GET /artist/{id}", Handler.GetArtistById())
	return mux
}

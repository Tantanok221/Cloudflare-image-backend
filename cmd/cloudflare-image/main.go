package main

import (
	"flag"
	"github.com/tantanok221/cloudflare-image-backend/internal/route"
	"log"
	"net/http"
	"os"
)

func main() {
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	addr := flag.String("addr", ":3001", "HTTP Network address")
	print("\n Server Running \n")
	//database := db.Init()
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  route.Route(),
	}
	srv.ListenAndServe()
}

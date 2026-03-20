package main

import (
	"net/http"
	"time"
)

func CreateRouter() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func StartServer(router *http.ServeMux, port string) error {
   server := &http.Server{
	   Addr: ":" + port,
	   Handler: router,
	   IdleTimeout: time.Minute,
	   ReadHeaderTimeout: 5 * time.Second,
	   ReadTimeout: 10 * time.Second,
	   WriteTimeout: 10 * time.Second,
   }

   return server.ListenAndServe()

}
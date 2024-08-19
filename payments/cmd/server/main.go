package payments

import (
	"log"
	"net/http"
	"time"
)

func Serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.Write([]byte("payments"))
	})

	log.Println("running payments on 4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}

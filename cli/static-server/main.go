package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := newRouter()

	addr := ":3000"
	log.Println("Static server listen on " + addr)
	err := http.ListenAndServe(addr, router)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(static)
	router.HandleFunc("/generator.csv", generatorHandler)
	return router
}

func generatorHandler(w http.ResponseWriter, r *http.Request) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		count = 100
	}

	for i := 0; i < count; i++ {
		price := random.Float64() * float64(random.Intn(count))
		row := fmt.Sprintf("Product %d;%.2f\n", random.Intn(count), price)

		if r.FormValue("plain") != "true" {
			w.Header().Add("Content-type", "text/csv")
			w.Header().Add("Content-Disposition", "attachment; filename=generator.csv")
			w.Header().Add("Pragma", "no-cache")
			w.Header().Add("Expires", "0")
		}

		w.Write([]byte(row))
	}
}

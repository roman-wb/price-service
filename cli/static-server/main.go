package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/purini-to/zapmw"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var addr = flag.String("addr", "localhost:3000", "Listen on host:port")
var mode = flag.String("mode", "dev", "Run mode dev or prod")

func main() {
	flag.Parse()

	var logger *zap.Logger
	if *mode == "prod" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync() //nolint:errcheck

	router := newRouter(logger)

	logger.Sugar().Infof("Static server listen on %s", *addr)
	err := http.ListenAndServe(*addr, router)
	if err != nil && err != http.ErrServerClosed {
		logger.Sugar().Fatal(err)
	}
}

func newRouter(logger *zap.Logger) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(
		zapmw.WithZap(logger),
		zapmw.Request(zapcore.InfoLevel, "request"),
		zapmw.Recoverer(zapcore.ErrorLevel, "recover", zapmw.RecovererDefault),
	)
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

		w.Write([]byte(row)) //nolint:errcheck
	}
}

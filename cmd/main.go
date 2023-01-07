package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/KhasanOrsaev/logger-client"
	"github.com/KhasanOrsaev/orse/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	logger.NewLoggerDefault(map[string]interface{}{
		"module": "test",
		"format": "json",
		"level":  4,
	})
}

type Response struct {
	State int
}
type Data struct {
	Controllers []domain.Controller
}

var data Data

func main() {
	// HTTP server
	srv := &http.Server{Handler: service()}

	// Server run context
	srvCtx, srvStopCtx := context.WithCancel(context.Background())

	// syscall signal to interrupt
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ch

		// Shutdown signal with grace period of 20 seconds
		shutdownCtx, _ := context.WithTimeout(srvCtx, 20*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		srvStopCtx()
	}()

	// Run the server
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-srvCtx.Done()
}

func service() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	data.Controllers = []domain.Controller{
		{
			Address:  "192.168.0.21",
			ID:       1,
			State:    0,
			Name:     "Зал",
			IsActive: true,
			Topic:    "light1",
		},
		{
			Address:  "192.168.0.21",
			ID:       2,
			State:    0,
			Name:     "Зал2",
			IsActive: true,
			Topic:    "light2",
		},
		{
			Address:  "192.168.0.21",
			ID:       3,
			State:    0,
			Name:     "Зал3",
			IsActive: true,
			Topic:    "light3",
		},
		{
			Address:  "192.168.0.21",
			ID:       4,
			State:    0,
			Name:     "Зал4",
			IsActive: true,
			Topic:    "light4",
		},
		{
			Address:  "192.168.0.22",
			ID:       5,
			State:    0,
			Name:     "Спальня Хасан 1",
			IsActive: true,
			Topic:    "light1",
		},
		{
			Address:  "192.168.0.22",
			ID:       6,
			State:    0,
			Name:     "Спальня Хасан 2",
			IsActive: true,
			Topic:    "light2",
		},
		{
			Address:  "192.168.0.23",
			ID:       7,
			State:    0,
			Name:     "Спальня Аслан 1",
			IsActive: true,
			Topic:    "light1",
		},
		{
			Address:  "192.168.0.23",
			ID:       8,
			State:    0,
			Name:     "Спальня Аслан 2",
			IsActive: true,
			Topic:    "light2",
		},
	}
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmpl, _ := template.ParseFiles("templates/index.html")
		err := tmpl.Execute(writer, data)
		if err != nil {
			logger.Error("on parse html", "", nil, err)
		}
	})
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RemoteAddr)
	})
	r.HandleFunc("/execute", func(writer http.ResponseWriter, request *http.Request) {
		var sig string
		id, _ := strconv.Atoi(request.URL.Query()["id"][0])
		res := Response{}
		switch request.URL.Query()["state"][0] {
		case "0":
			res.State = 1
			sig = "on"
		case "1":
			res.State = 0
			sig = "off"
		}
		for _, v := range data.Controllers {
			if v.ID == id {
				http.DefaultClient.Timeout = 2 * time.Second
				resp, err := http.Get("http://" + v.Address + "/" + sig + "?pin=" + v.Topic)
				if err != nil {
					logger.Error("on send signal", "", nil, err)
					writer.WriteHeader(http.StatusBadGateway)
					return
				}
				if resp.StatusCode != http.StatusOK {
					logger.Error("on send signal", resp.Status, nil, nil)
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
				logger.Info("send signal", "", &map[string]interface{}{
					"topic":  v.Topic,
					"signal": sig,
				}, err)
			}
		}

		j, _ := json.Marshal(res)
		writer.Write(j)
	})
	return r
}

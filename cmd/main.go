package main

import (
	"context"
	"encoding/json"
	"github.com/KhasanOrsaev/logger-client"
	"github.com/KhasanOrsaev/orse/internal/domain"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	logger.NewLoggerDefault(map[string]interface{}{
		"module": "test",
		"format": "json",
		"level": 4,
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
	handler := http.NewServeMux()
	srv := http.Server{Handler: handler}
	data.Controllers = []domain.Controller{
		{
			Address: "192.168.0.21",
			ID: 1,
			State: 0,
			Name: "Зал",
			IsActive: true,
			Topic: "light1",
		},
		{
			Address: "192.168.0.21",
			ID: 2,
			State: 0,
			Name: "Зал2",
			IsActive: true,
			Topic: "light2",
		},
		{
			Address: "192.168.0.21",
			ID: 3,
			State: 0,
			Name: "Зал3",
			IsActive: true,
			Topic: "light3",
		},
		{
			Address: "192.168.0.21",
			ID: 4,
			State: 0,
			Name: "Зал4",
			IsActive: true,
			Topic: "light4",
		},
		{
			Address: "192.168.0.22",
			ID: 5,
			State: 0,
			Name: "Спальня Хасан 1",
			IsActive: true,
			Topic: "light1",
		},
		{
			Address: "192.168.0.22",
			ID: 6,
			State: 0,
			Name: "Спальня Хасан 2",
			IsActive: true,
			Topic: "light2",
		},
		{
			Address: "192.168.0.23",
			ID: 7,
			State: 0,
			Name: "Спальня Аслан 1",
			IsActive: true,
			Topic: "light1",
		},
		{
			Address: "192.168.0.23",
			ID: 8,
			State: 0,
			Name: "Спальня Аслан 2",
			IsActive: true,
			Topic: "light2",
		},
	}
	handler.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmpl,_ := template.ParseFiles("templates/index.html")
		err := tmpl.Execute(writer,data)
		if err != nil {
			logger.Error("on parse html","",nil,err)
		}
	})
	handler.HandleFunc("/execute", func(writer http.ResponseWriter, request *http.Request) {
		var sig string
		id,_ := strconv.Atoi(request.URL.Query()["id"][0])
		res := Response{}
		switch request.URL.Query()["state"][0] {
		case "0":
			res.State = 1
			sig = "on"
		case "1":
			res.State = 0
			sig = "off"
		}
		for _,v := range data.Controllers {
			if v.ID == id {
				http.DefaultClient.Timeout = 2*time.Second
				resp, err := http.Get("http://"+v.Address+"/"+sig+"?pin="+v.Topic)
				if err != nil {
					logger.Error("on send signal","",nil,err)
					writer.WriteHeader(http.StatusBadGateway)
					return
				}
				if resp.StatusCode != http.StatusOK {
					logger.Error("on send signal",resp.Status,nil,nil)
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
				logger.Info("send signal","",&map[string]interface{}{
					"topic": v.Topic,
					"signal": sig,
				},err)
			}
		}

		j,_ := json.Marshal(res)
		writer.Write(j)
	})

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatal("on start server","",nil,err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("on stop server","",nil,err)
	}
}


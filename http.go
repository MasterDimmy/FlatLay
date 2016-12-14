package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TProtoJSError struct {
	Success bool
	Message string
}

type TProtoJSSuccess struct {
	Success bool
	Items   interface{}
}

func ProtoError(w http.ResponseWriter, s string) {
	p := TProtoJSError{
		Success: false,
		Message: s,
	}
	buf, _ := json.MarshalIndent(p, "", " ")
	w.Write(buf)
}

func ProtoSuccess(w http.ResponseWriter, s interface{}) {
	p := TProtoJSSuccess{
		Success: true,
		Items:   s,
	}
	buf, _ := json.MarshalIndent(p, "", " ")
	w.Write(buf)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Println(r.URL.String())
	switch r.URL.String() {
	case "/":
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Println(r.URL.String())
	err := r.ParseForm()
	if err != nil {
		ProtoError(w, "Ошибка: "+err.Error())
		return
	}
	ws := r.Form.Get("width")
	hs := r.Form.Get("height")
	data, err := app.Cfg.Collager.getCollage(ws, hs)
	if err != nil {
		ProtoError(w, "Ошибка: "+err.Error())
	} else {
		ProtoSuccess(w, data)
	}
}

func (a *TApp) createWebServer() error {
	if len(a.Cfg.WebServer.Static) == 0 {
		return fmt.Errorf("Ошибка! Пусть путь к папке Веб-сервера static (задается в настройках)")
	}

	http.HandleFunc("/", root)
	http.HandleFunc("/get_data", getData)
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir(app.Cfg.WebServer.Static)))
	http.Handle("/static/", fileServer)

	fmt.Println("\nЗапуск Веб-сервера по адресу http://" + a.Cfg.WebServer.IpPort + "\nКаталог Веб-сервера: " + app.Cfg.WebServer.Static)

	err := http.ListenAndServe(a.Cfg.WebServer.IpPort, nil)
	return err
}

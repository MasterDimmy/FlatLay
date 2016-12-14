package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*
	Формирователь коллажа для Рината Бакиева
	Настройки: config.json
	База: database.json

	Настройки
*/

type TWebConfig struct {
	IpPort string
	Static string
}

type TConfig struct {
	WebServer TWebConfig
	Collager  TCollager
}

type TApp struct {
	Cfg TConfig
}

func ErrorCheck(err error, s string) bool {
	if err != nil {
		fmt.Println(s)
		fmt.Println("ERROR:" + err.Error())
		return true
	}
	return false
}

var app TApp

func main() {
	f, err := ioutil.ReadFile("config.json")
	if ErrorCheck(err, "Ошибка! Невозможно прочитать файл настроек config.json") {
		return
	}

	err = json.Unmarshal(f, &app.Cfg)
	if len(app.Cfg.WebServer.IpPort) == 0 {
		fmt.Printf("Ошибка! Порт демо Веб-интерфейса пуст!")
		return
	}

	if ErrorCheck(err, "Ошибка чтения файла config.json") {
		return

	}
	/*
		fmt.Println("Считанные настройки: ")
		b, err := json.MarshalIndent(app.Cfg, "", " ")
		if ErrorCheck(err, "ERROR") {
			return
		}
		fmt.Printf("%s", string(b))
	*/

	err = app.createWebServer()
	if ErrorCheck(err, "Ошибка запуска демо Веб-сервера!") {
		return
	}
}

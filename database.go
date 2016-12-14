package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TImage struct {
	Name   string //имя картинки
	Path   string //локальный путь для веб-сервера
	Width  int    //ширина
	Height int    //высота
	Weight int    //архимедов вес
	Group  int    //группа стиля одежды (категория?)
}

type TDatabase struct {
	Images []TImage //картинки
	groups int      //сколько всего групп в картинках?
}

func (t *TCollager) readDatabase() (*TDatabase, error) {
	d, err := ioutil.ReadFile(t.Database)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла базы картинок: " + err.Error())
	}
	db := TDatabase{}
	err = json.Unmarshal(d, &db)
	if err != nil {
		return nil, fmt.Errorf("Ошибка парсинга файла базы картинок: " + err.Error())
	}
	var m map[int]bool
	for _, i := range db.Images {
		_, ok := m[i.Group]
		if !ok {
			m[i.Group] = true
			db.groups++
		}
	}

	return &db, nil
}

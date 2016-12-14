package main

import (
	"fmt"
)

type TCollager struct {
	Items    int    //количество картинок на странице
	Database string //путь к базе json

	DB *TDatabase
}

//вернуть коллаж для JS в JSON
func (t *TCollager) getCollage(gs string, ws string, hs string) (*TCollage, error) {
	if len(ws) == 0 {
		return nil, fmt.Errorf("Ширина не указана")
	}
	if len(hs) == 0 {
		return nil, fmt.Errorf("Ширина не указана")
	}
	if t.DB == nil {
		return nil, fmt.Errorf("Отсутствует база картинок!")
	}
	if len(t.DB.Images) == 0 {
		return nil, fmt.Errorf("Отсутствует картинки в базе!")
	}
	//требуется сделать коллаж по размерам
	fmt.Println("gs=" + gs + " ws=" + ws + " hs=" + hs)
	w := 0
	h := 0
	g := 0
	_, err := fmt.Sscanf(gs+" "+ws+" "+hs+"\n", "%d %d %d", &g, &w, &h)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения размеров окна: " + err.Error())
	}

	if w <= 0 {
		return nil, fmt.Errorf("Задана нулевая ширина окна! ")
	}
	if h <= 0 {
		return nil, fmt.Errorf("Задана нулевая высота окна! ")
	}

	return t.create(g, w, h)
}

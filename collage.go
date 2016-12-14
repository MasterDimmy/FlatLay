package main

import (
	"fmt"
)

//данные для JS
type TCollageImage struct {
	PosX int
	PosY int
	Path string
}

type TCollage struct {
	Images []TCollageImage
}

type TCollager struct {
	Items int //количество картинок на странице
}

//вернуть коллаж для JS в JSON
func (t *TCollager) getCollage(ws string, hs string) (*TCollager, error) {
	if len(ws) == 0 {
		return nil, fmt.Errorf("Ширина не указана")
	}
	if len(hs) == 0 {
		return nil, fmt.Errorf("Ширина не указана")
	}
	fmt.Println("ws=" + ws + " hs=" + hs)
	w := 0
	h := 0
	_, err := fmt.Sscanf(ws+" "+hs+"\n", "%d %d", &w, &h)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения размеров окна: " + err.Error())
	}
	collage := TCollager{}
	return &collage, nil
}

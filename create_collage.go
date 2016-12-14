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

func (t *TCollager) create(w int, h int) (*TCollage, error) {
	collage := TCollage{}
	return &collage, nil
}

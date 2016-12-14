package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

var help string = `
	Вспомогательная программа для коллажей.
	Создает базу database.json на основе картинок каталога.
	
	Использование: this.exe <путь_к_каталогу>

	Формат базы:
	{
		"Images": {
		[
			"Name": "image1.jpg",
			"Path": "folder/image1.jpg",
			"Width": 123,
			"Height": 342,
			"Weight": 1
		],
		[
			"Name": "image1.jpg",
			"Path": "folder/image1.jpg",
			"Width": 123,
			"Height": 342,
			"Weight": 1
		]
		}
	}
`

type TImage struct {
	Name   string //имя картинки
	Path   string //локальный путь для веб-сервера
	Width  int    //ширина
	Height int    //высота
	Weight int    //архимедов вес
}

type TDatabase struct {
	Images []TImage
}

var db TDatabase

func walkFn(path string, info os.FileInfo, err error) error {
	if info == nil {
		return nil
	}
	if info.IsDir() {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка чтения файла: " + err.Error())
		return err
	}
	fmt.Println("Обработка: " + path)

	img, err := jpeg.Decode(f)
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return nil
	}

	image := TImage{
		Name:   filepath.Base(path),
		Path:   path,
		Width:  img.Bounds().Max.X,
		Height: img.Bounds().Max.Y,
		Weight: 1,
	}

	db.Images = append(db.Images, image)
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		return
	}

	ex, err := exists(os.Args[1])
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	if !ex {
		fmt.Println("Ошибка: путь " + os.Args[1] + " не существует!")
		return
	}

	err = filepath.Walk(os.Args[1], walkFn)
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	data, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	err = ioutil.WriteFile("database.json", data, 0777)
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	fmt.Printf("База database.json создана!\nЗаписей: %d", len(db.Images))
	return
}

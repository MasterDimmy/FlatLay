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
	Group  int    //группа стиля одежды (категория?)
}

type TDatabase struct {
	Images []TImage
}

var db TDatabase //вся база

var names = make(map[string]int) //номер  группы (категория картинки)
var min_name_len = 10000

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
		Group:  0, //значение по умолчанию
	}

	if min_name_len > len(image.Name) {
		min_name_len = len((image.Name))
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

	//обход каталога - построение модели базы в памяти
	err = filepath.Walk(os.Args[1], walkFn)
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	//формируем группы одежды в рамках одного стиля на основе имени файла
	/* из имен файлов
	23.11.2016  22:29            90 091 produktybermudy-kargo-loose-fit-sinij-v-kletku-923883stl.jpg
	23.11.2016  22:29            16 534 produktybermudy-kargo-loose-fit-sinij-v-kletku-923883stl_product1.jpg
	23.11.2016  22:29             5 292 produktybermudy-kargo-loose-fit-sinij-v-kletku-923883stl_product2.jpg
	23.11.2016  22:29            35 287 produktybermudy-kargo-loose-fit-sinij-v-kletku-923883stl_product3.jpg
	23.11.2016  22:29             6 272 produktybermudy-kargo-loose-fit-sinij-v-kletku-923883stl_product4.jpg
	23.11.2016  22:28            81 629 produktybrjuki-5-karmanov-regular-fit-straight-bezhevyj-907597stl.jpg
	23.11.2016  22:28            22 252 produktybrjuki-5-karmanov-regular-fit-straight-bezhevyj-907597stl_product1.jpg
	23.11.2016  22:28             5 830 produktybrjuki-5-karmanov-regular-fit-straight-bezhevyj-907597stl_product2.jpg
	23.11.2016  22:28            18 011 produktybrjuki-5-karmanov-regular-fit-straight-bezhevyj-907597stl_product3.jpg
	23.11.2016  22:28            16 156 produktybrjuki-5-karmanov-regular-fit-straight-bezhevyj-907597stl_product4.jpg
	23.11.2016  22:30           109 119 produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl.jpg
	23.11.2016  22:30            20 835 produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl_product1.jpg
	23.11.2016  22:30            19 700 produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl_product2.jpg
	23.11.2016  22:30             4 737 produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl_product3.jpg

	выделяем самое короткое имя produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl.jpg
	убираем все до .
	считаем эту длину "базовой"
	все имена, что входят в "базу" - принадлежат одной группе

	produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl
	produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl_product1
	produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl_product2
	produktybrjuki-bez-zastezhki-classic-fit-sinij-939493stl_product3
	*/

	//формируем базу имен
	for _, i := range db.Images {
		base := i.Name[:min_name_len]
		_, ok := names[base]
		if !ok {
			names[base] = len(names) + 1
		}
	}

	//формируем группы картинок на основе базы имен
	for n, i := range db.Images {
		db.Images[n].Group = names[i.Name[:min_name_len]]
	}

	//формируем модель базы json
	data, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	//пишем
	err = ioutil.WriteFile("database.json", data, 0777)
	if err != nil {
		fmt.Println("Ошибка: " + err.Error())
		return
	}

	fmt.Printf("База database.json создана!\nЗаписей: %d\nКатегорий: %d", len(db.Images), len(names))
	return
}

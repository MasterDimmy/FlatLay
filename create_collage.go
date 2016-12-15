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
	Images       []TCollageImage //картинки
	Square       int64           //занятая площадь
	MaxAvailable int64           //максимально возможная площадь
}

//ограничения поиска
type TLimits struct {
	maxX  int //размеры рамки X, W
	maxY  int //размеры рамки Y, H
	group int //группа (категория)
}

func (t *TCollager) create(group int, w int, h int) (*TCollage, error) {
	fmt.Println("create for ", group, w, h)
	collage := TCollage{MaxAvailable: int64(w) * int64(h)}

	//создаем поле для группы с ее ограничениями
	//может работать параллельно одновременно для разных пользователей с разными ограничениями
	b, f := t.gen(&TLimits{maxX: w, maxY: h, group: group}, 0, make(TField))

	fmt.Println("Создан коллаж из", len(f), " картинок")

	//преобразуем поле в коллаж для JS
	for n, v := range f {
		collage.Images = append(collage.Images, TCollageImage{
			PosX: v.x,
			PosY: v.y,
			Path: t.DB.Images[n].Path,
		})
		fmt.Println("x ", v.x, "y ", v.y)
	}

	//указываем занятую площадь
	collage.Square = b

	return &collage, nil
}

//элемент поля
type TFieldItem struct {
	x int //ее позиция на поле по Х
	y int //позиция по Y
}

//поле картинок [номер в базе] позиция
type TField map[int]*TFieldItem

//можно разместить картинку по указанным координатам?
//num - номер вставляемой по базе
//x,y - куда вставляем
func (t *TCollager) trypaste(limits *TLimits, num int, x int, y int, f TField) bool {
	//1. картинка должна вмещаться в рамки по X
	if (x + t.DB.Images[num].Width) > limits.maxX {
		return false
	}
	//2. картинка должна вмещаться в рамки по Y
	if (y + t.DB.Images[num].Height) > limits.maxY {
		return false
	}
	//	fmt.Printf("Test: num:%d x:%d y:%d w:%d h:%d\n", num, x, y, t.DB.Images[num].Width, t.DB.Images[num].Height)
	for i, v := range f {
		//3. проверяем правило размещения по весам:
		//не должно быть картинки БОЛЬШЕ весом, которая будет НАД текущей
		//т.е. нельзя чтобы any.y+any.H<=y
		if (t.DB.Images[i].Height+v.y > y) && (t.DB.Images[i].Weight > t.DB.Images[num].Weight) {
			return false
		}
		//4. не должно быть наложений картинок
		if v.x <= x && x <= v.x+t.DB.Images[i].Width { //левая внутри кого-то
			if v.y <= y && y <= v.y+t.DB.Images[i].Height { //верхняя наша
				return false
			}
			if y <= v.y && v.y <= y+t.DB.Images[num].Height { //верхняя чужая
				return false
			}
			if v.y <= y+t.DB.Images[num].Height && y+t.DB.Images[num].Height <= v.y+t.DB.Images[i].Height { //нижняя наша
				return false
			}
			if y <= v.y+t.DB.Images[i].Height && v.y+t.DB.Images[i].Height <= y+t.DB.Images[num].Height { //нижняя чужая
				return false
			}
		}
		if v.x <= x+t.DB.Images[num].Width && x+t.DB.Images[num].Width <= v.x+t.DB.Images[i].Width { //правая внутри кого-то
			if v.y <= y && y <= v.y+t.DB.Images[i].Height { //верхняя наша
				return false
			}
			if y <= v.y && v.y <= y+t.DB.Images[num].Height { //верхняя чужая
				return false
			}
			if v.y <= y+t.DB.Images[num].Height && y+t.DB.Images[num].Height <= v.y+t.DB.Images[i].Height { //нижняя наша
				return false
			}
		}
		//5. полное накрытие начал
		if v.x == x && v.y == y {
			return false
		}
	}
	return true
}

//генерируем поле!
func (t *TCollager) gen(limits *TLimits, used_square int64, used_field TField) (int64, TField) {
	best_square := used_square      //максимальная площадь
	best_field := used_field        //лучшее поле для максимальной площади
	for n, _ := range t.DB.Images { //перебираем каждую картинку и пробуем ее вставить
		//условие отсева: они приналежат одной группе с имеющимися на поле
		if limits.group == t.DB.Images[n].Group {
			_, ok := used_field[n]
			if ok {
				continue //если картинка уже есть на поле - пропускаем
			}

			varx := []int{0} //создаем варианты размещения по X
			vary := []int{0} //создаем варианты размещения по Y
			for i, v := range used_field {
				varx = append(varx, v.x+1)
				varx = append(varx, v.x+t.DB.Images[i].Width+1)
				vary = append(varx, v.y+1)
				vary = append(varx, v.y+t.DB.Images[i].Height+1)
			}
			for _, x := range varx { //перебираем варианты по X
				for _, y := range vary { //перебираем варианты по Y
					if t.trypaste(limits, n, x, y, used_field) { //если можно разместить
						used_field[n] = &TFieldItem{ //размещаем на поле
							x: x,
							y: y,
						}
						ns, nf := t.gen(limits, used_square+int64(t.DB.Images[n].Width)*int64(t.DB.Images[n].Height), used_field)
						if ns > best_square {
							best_square = ns
							best_field = make(TField) //копируем новое поле
							for k, v := range nf {
								best_field[k] = v
							}
						}
						delete(used_field, n) //удаляем вариант с _текущего_ поля
					}
				}
			}
		}
	}
	return best_square, best_field
}

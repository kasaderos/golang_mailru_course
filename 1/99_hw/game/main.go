package main

import (
	"strings"
)

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные перменные ( тут можно ), функции
*/
type arrStrFuncType func(*[]string)

var man Man
var rooms []Room

type RoomName int
type ItemName int

const (
	Kitchen RoomName = iota
	Hall
	Bedroom
	Street
)

const (
	Table ItemName = iota
	Chair
	Backpack
	Tea
	Keys
	Notes
	Door
)

func (d RoomName) String() string {
	return [...]string{"кухня", "коридор", "комната", "улица"}[d]
}

func (d ItemName) String() string {
	return [...]string{"стол", "стул", "рюкзак", "чай", "ключи", "конспекты", "дверь"}[d]
}

type Man struct {
	inventory    []string
	haveBackpack bool
	location     *Room
	mission      string
}

type Item struct {
	name   ItemName
	parent ItemName
}

type Room struct {
	name  RoomName
	items []Item
	paths []*Room
	info  string
}

// идти
func (man *Man) move(newLoc string) string {
	for _, loc := range man.location.paths {
		if loc.name.String() == newLoc {
			man.location = loc
			return loc.info
		}
	}
	return "нет пути в комната"
}

// осмотреться
func (man *Man) lookAround() (res string) {
	res = man.location.info
	if man.location.items != nil {
		return ""
	}
	return ""
}
func (man *Man) putOn(item string) string {
	if item == Backpack.String() {
		man.haveBackpack = true
		return "вы надели: " + item
	}
	return "нечего одеть"
}

func (man *Man) take(item string) string {
	if man.haveBackpack {
		man.inventory = append(man.inventory, item)
		return "предмет добавлен в инвентарь: " + item + "."
	}
	return "некуда класть."
}

func (man *Man) apply(item string, toItem string) string {
	//..
	return ""
}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
	initGame()
}

func initGame() {
	rooms = []Room{
		Room{
			name: Kitchen,
			items: []Item{
				Item{
					name:   Tea,
					parent: Table,
				},
			},
			info: "ты находишься на кухне.",
		},
		Room{
			name: Hall,
			items: []Item{
				Item{
					name: Door,
				},
			},
			info: "ничего интересного",
		},
		Room{
			name: Bedroom,
			items: []Item{
				Item{
					name:   Keys,
					parent: Table,
				},
				Item{
					name:   Notes,
					parent: Table,
				},
				Item{
					name:   Backpack,
					parent: Chair,
				},
			},
			info: "ты в своей комнате",
		},
		Room{
			name: Street,
			info: "на улице весна",
		},
	}
	rooms[0].paths = append(rooms[0].paths, &rooms[1]) // кухня    -> коридор
	rooms[1].paths = append(rooms[1].paths, &rooms[0]) // коридор  -> кухня
	rooms[1].paths = append(rooms[1].paths, &rooms[2]) // коридор  -> комнату
	rooms[1].paths = append(rooms[1].paths, &rooms[3]) // коридор  -> улицу
	rooms[2].paths = append(rooms[2].paths, &rooms[1]) // комната  -> коридор
	rooms[3].paths = append(rooms[3].paths, &rooms[1]) // улица    -> коридор
	//fmt.Println(rooms[0].paths[0].name)
	man = Man{
		location: &rooms[0],
		mission:  "надо собрать рюкзак и идти в универ",
	}
}

/*
	данная функция принимает команду от "пользователя"
	и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
*/
func handleCommand(command string) string {
	args := strings.Split(command, " ")
	if args == nil {
		panic("No commands")
	}
	if args[0] == "осмотреться" {
		return man.lookAround()
	} else if args[0] == "идти" {
		return man.move(args[1])
	} else if args[0] == "взять" {
		return man.take(args[1])
	} else if args[0] == "надеть" {
		return man.putOn(args[1])
	} else if args[0] == "применить" {
		return man.apply(args[1], args[2])
	}

	return "not implemented"
}

/*
actions: map[string]arrStrFuncType{
			//"осмотреться": lookAround,
			"идти":      man.move,
			"надеть":    man.putOn,
			"взять":     man.take,
			"применить": man.apply,
		},*/

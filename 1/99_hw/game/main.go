package main

import (
	"fmt"
	"strings"
)

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные перменные ( тут можно ), функции
*/
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

var mapOfItems map[string]ItemName = map[string]ItemName{
	"стол":      Table,
	"стул":      Chair,
	"рюкзак":    Backpack,
	"чай":       Tea,
	"ключи":     Keys,
	"конспекты": Notes,
	"дверь":     Door,
}

type Man struct {
	inventory      []*Item
	haveBackpack   bool
	location       *Room
	mission        []string
	curMessageRoom *string
}

type Item struct {
	name       ItemName
	parent     ItemName
	applying   ItemName
	f          bool
	afterApply string
}

type Room struct {
	name      RoomName
	items     []*Item
	paths     []*Room
	info      string
	afterMove string
}

// идти
func (man *Man) move(newLoc string) (string, bool) {
	for _, loc := range man.location.paths {
		if loc.name.String() == newLoc {
			door, _ := findItem(Door, &man.location.items)
			if door != nil && newLoc == Street.String() {
				if !(*door).f {
					return "дверь закрыта", false
				} else {
					man.location = loc
					return loc.afterMove + ".", true
				}
			}
			man.location = loc
			return loc.afterMove + ".", true
		}
	}
	return "нет пути в комната", false
}

// осмотреться
func (m *Man) lookAround() (res string) {
	if m.location.info != "" {
		res = m.location.info + ", "
	}
	onTheTable, onTheChair := "", ""
	for _, item := range m.location.items {
		if item != nil && (*item).parent == Table {
			onTheTable += item.name.String() + ", "
		} else if item != nil && (*item).parent == Chair {
			onTheChair += item.name.String() + ", "
		}
	}

	if onTheTable != "" {
		res += "на столе: " + onTheTable
	}
	if onTheChair != "" {
		res += "на стуле: " + onTheChair
	}
	if man.location.name == Kitchen {
		res += getMission()
	}
	if res == "" {
		res = "пустая комната."
	} else {
		res = res[:len(res)-2] + "."
	}
	return
}

func getMission() (res string) {
	res = "надо "
	for _, m := range man.mission {
		if m != "" {
			res += m + " и "
		}
	}
	return res[:len(res)-2]
}

func (m *Man) deleteMissionBackpack() {
	if m.haveBackpack {
		keys, _ := findItem(Keys, &m.inventory)
		notes, _ := findItem(Notes, &m.inventory)
		if keys != nil && notes != nil {
			for i, v := range m.mission {
				if v == "собрать рюкзак" {
					m.mission[i] = ""
				}
			}
		}
	}
}

func (man *Man) putOn(item string) string {
	if item == Backpack.String() {
		man.haveBackpack = true
		_, ind := findItem(Backpack, &man.location.items)
		man.location.items[ind] = nil
		return "вы надели: " + item
	}
	return "нечего одеть"
}

func (man *Man) take(item string) string {
	if man.haveBackpack {
		if itemName, ok := mapOfItems[item]; ok {
			pItem, ind := findItem(itemName, &man.location.items)
			if pItem != nil {
				man.inventory = append(man.inventory, pItem)
				man.deleteMissionBackpack()
				man.location.items[ind] = nil
				return "предмет добавлен в инвентарь: " + item
			} else {
				return "нет такого"
			}
		} else {
			return "нет такого"
		}
	}
	return "некуда класть"
}

func findItem(name ItemName, items *[]*Item) (itm *Item, ind int) {
	for i, item := range *items {
		if item != nil && (*item).name == name {
			return item, i
		}
	}
	return nil, -1
}

func (man *Man) apply(item string, toItem string) string {
	if itemName, exist := mapOfItems[item]; exist {
		pItem, _ := findItem(itemName, &man.inventory)
		if pItem == nil {
			return "нет предмета в инвентаре - " + item
		}
		pToItem, _ := findItem((*pItem).applying, &man.location.items)
		if pToItem == nil || (*pToItem).name.String() != toItem {
			return "не к чему применить"
		}
		(*pToItem).f = true

		return (*pToItem).afterApply
	}

	return "нет предмета в инвентаре - " + item
}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
	/*
		initGame()
		fmt.Println(handleCommand("осмотреться"))
		fmt.Println(handleCommand("идти коридор"))
		fmt.Println(handleCommand("идти комната"))
		fmt.Println(handleCommand("осмотреться"))
		fmt.Println(handleCommand("надеть рюкзак"))
		fmt.Println(handleCommand("взять ключи"))
		fmt.Println(handleCommand("взять конспекты"))
		fmt.Println(handleCommand("идти коридор"))
		fmt.Println(handleCommand("применить ключи дверь"))
		fmt.Println(handleCommand("идти улица"))
	*/
	initGame()
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("завтракать"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("надеть рюкзак"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println("12", handleCommand("взять телефон"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("осмотреться"))

	fmt.Println(handleCommand("взять конспекты"))
	fmt.Println(handleCommand("осмотреться"))

	fmt.Println(handleCommand("идти коридор"))

	fmt.Println(handleCommand("идти кухня"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти улица"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("применить телефон шкаф"))
	fmt.Println(handleCommand("применить ключи шкаф"))
	fmt.Println(handleCommand("идти улица"))
}

func initGame() {
	rooms = []Room{
		Room{
			name: Kitchen,
			items: []*Item{
				&Item{
					name:   Tea,
					parent: Table,
				},
			},
			info:      "ты находишься на кухне",
			afterMove: "кухня, ничего интересного",
		},
		Room{
			name: Hall,
			items: []*Item{
				&Item{
					name:       Door,
					afterApply: "дверь открыта",
				},
			},
			info:      "ничего интересного",
			afterMove: "ничего интересного",
		},
		Room{
			name: Bedroom,
			items: []*Item{
				&Item{
					name:     Keys,
					parent:   Table,
					applying: Door,
				},
				&Item{
					name:   Notes,
					parent: Table,
				},
				&Item{
					name:   Backpack,
					parent: Chair,
				},
			},
			afterMove: "ты в своей комнате",
		},
		Room{
			name:      Street,
			info:      "",
			afterMove: "на улице весна",
		},
	}
	rooms[Kitchen].paths = append(rooms[Kitchen].paths, &rooms[Hall]) // кухня    -> коридор
	rooms[Hall].paths = append(rooms[Hall].paths, &rooms[Kitchen])    // коридор  -> кухня
	rooms[Hall].paths = append(rooms[Hall].paths, &rooms[Bedroom])    // коридор  -> комнату
	rooms[Hall].paths = append(rooms[Hall].paths, &rooms[Street])     // коридор  -> улицу
	rooms[Bedroom].paths = append(rooms[Bedroom].paths, &rooms[Hall]) // комната  -> коридор
	rooms[Street].paths = append(rooms[Street].paths, &rooms[Hall])   // улица    -> коридор

	man = Man{
		location: &rooms[0],
		mission: []string{
			"собрать рюкзак",
			"идти в универ",
		},
	}
}

/*
	данная функция принимает команду от "пользователя"
	и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
*/
func whereToGo(f bool) string {
	if !f {
		return ""
	}
	res := " можно пройти - "
	for _, way := range man.location.paths {
		if man.location.name == Street && (*way).name == Hall {
			res += "домой  "
		} else {
			res += (*way).name.String() + ", "
		}
	}
	length := len(res)
	return res[:length-2]
}

func handleCommand(command string) (res string) {
	args := strings.Split(command, " ")
	if args == nil {
		panic("No commands")
	}
	if args[0] == "осмотреться" {
		res = man.lookAround()
		res += whereToGo(true)
	} else if args[0] == "идти" {
		var f bool
		res, f = man.move(args[1])
		res += whereToGo(f)
	} else if args[0] == "взять" {
		res = man.take(args[1])
	} else if args[0] == "надеть" {
		res = man.putOn(args[1])
	} else if args[0] == "применить" {
		res = man.apply(args[1], args[2])
	} else {
		return "неизвестная команда"
	}
	return
}

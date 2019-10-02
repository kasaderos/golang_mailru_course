package main

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные перменные ( тут можно ), функции
*/
type Backpack struct {
	inventory []string
}

type Man struct {
	backpack Backpack
}

// идти
func (man *Man) move() {
}

// осмотреться
func (man *Man) lookAround() {

}

// надеть
func (man *Man) putOn() {
}

type Item struct {
	name string
}

type Room struct {
	name  string
	items []Item
	paths []string
}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
}

func initGame() {
	var rooms = []Room {
		Room {
			name : "кухня",
			items : {
				{name : "стол"},
				{name : "чай"},
				{name : "стол"},
			}
		}
	}
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/
	return "not implemented"
}

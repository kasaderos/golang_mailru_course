package main

import (
	"fmt"
	"golang-2018-2/1/5_visibility/person"
)

func main() {
	p := person.NewPerson(1, "rvasily", "secret")

	// p.secret undefined (cannot refer to unexported field or method secret)
	// fmt.Printf("main.PrintPerson: %+v\n", p.secret)

	secret := person.GetSecret(p)
	fmt.Println("GetSecret", secret)
}

package main

import (
	"fmt"
)

func main () {
	arrosages := arrosages()

	aTraiter := arrosages.toDo().current()

	if len(aTraiter) > 0 {
		fmt.Println("Arrosage en cours : "+aTraiter[0].Event.Summary)
		// TODO arroser
		aTraiter[0].setDone().save()
	} else {
		fmt.Println("Aucun arrosage à effectuer")
	}
}
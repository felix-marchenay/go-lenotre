package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main () {
	arrosages := arrosages()

	aTraiter := arrosages.toDo().current()

	if len(aTraiter) > 0 {
		fmt.Println("Arrosage en cours : "+aTraiter[0].Event.Summary)

		cmdOpen := exec.Command("gpio", "write", "2", "1")
		cmdClose := exec.Command("gpio", "write", "2", "0")
		err := cmdOpen.Run()
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second*20)

		err = cmdClose.Run()
		if err != nil {
			log.Fatal(err)
		}

		aTraiter[0].setDone().save()
	} else {
		fmt.Println("Aucun arrosage à effectuer")
	}
}
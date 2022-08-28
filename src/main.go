package main

import (
	"fmt"
	"time"
)

var db FileDatabase[Arrosage]

func main() {

	db = FileDatabase[Arrosage]{
		filename: "arrosages",
	}

	all := AllArrosages()
	afaire := all.AFaire()

	fmt.Println(len(all), "arrosages à total, ", len(afaire), "à traiter")

	for _, arrosage := range afaire {

		if arrosage.Done {
			fmt.Println("Déjà fait -- ", arrosage.Event.Id, arrosage.Event.Summary)
			continue
		}

		now := time.Now()

		if now.Before(arrosage.Start) || now.After(arrosage.End) {
			fmt.Println("C'est pas l'heure -- ", arrosage.Event.Id, arrosage.Event.Summary)
			continue
		}

		arrosage.setDone()

		db.Set(arrosage.Event.Id, arrosage)
		err := arrosage.save()
		if err != nil {
			fmt.Println("Erreur enregistrement vers gcalendar event"+arrosage.Event.Id+" : ", err)
		}

		arrosage.arroser()
	}
}

func AllArrosages() (result ArrosageSlice) {

	arrosagesGoogle, err := ArrosagesFromG()

	if err != nil {
		return db.All()
	}

	for _, arrosage := range arrosagesGoogle {
		if db.Exist(arrosage.Event.Id) {
			arrosageLocal := db.Get(arrosage.Event.Id)
			if arrosageLocal.Done == true && arrosage.Done == false {
				arrosage.setDone()
				arrosage.save()
			}
		}

		db.Set(arrosage.Event.Id, arrosage)

		result = append(result, arrosage)
	}

	return result
}

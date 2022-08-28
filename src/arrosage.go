package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/calendar/v3"
)

var pinOut map[string]string = map[string]string{
	"pompe": "0",
	"out1":  "1",
	"out2":  "2",
	"out3":  "3",
	"out4":  "4",
}

type ArrosageSlice []Arrosage

func (arrosages ArrosageSlice) AFaire() (afaire []Arrosage) {

	for _, arrosage := range arrosages {

		if arrosage.Done {
			continue
		}

		now := time.Now()

		if now.Before(arrosage.Start) || now.After(arrosage.End) {
			continue
		}

		afaire = append(afaire, arrosage)
	}

	return afaire
}

type Arrosage struct {
	Done     bool
	Start    time.Time
	End      time.Time
	Event    calendar.Event
	Duration time.Duration
	Sortie   string
}

type summary string

func (arrosage *Arrosage) setDone() {
	arrosage.Done = true
	arrosage.Event.Summary = arrosage.Event.Summary + " -OK"
}

func (str summary) lastchars(length int) (result summary) {

	if len(str) < length {
		return str
	}

	result = str[len(str)-length:]
	return
}

func ArrosagesFromG() (arrosages []Arrosage, err error) {

	events, err := getEvents()
	if err != nil {
		return nil, err
	}

	for _, event := range events.Items {

		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)

		outNum := regexp.MustCompile(`out[0-9]`)
		outFound := outNum.FindString(event.Summary)

		sortie, found := pinOut[outFound]
		if found == false {
			sortie = "out1"
		}

		arrosages = append(arrosages, Arrosage{
			summary(event.Summary).lastchars(3) == "-OK",
			start,
			end,
			*event,
			durationFromString(event.Summary),
			sortie,
		})

	}

	return arrosages, nil
}

func durationFromString(str string) time.Duration {

	regSec := regexp.MustCompile(`[^\d]([0-9]+)s`)
	result := regSec.FindStringSubmatch(str)

	if len(result) > 1 {
		seconds, err := strconv.Atoi(result[1])
		if err != nil {
			panic(err)
		}
		return time.Second * time.Duration(seconds)
	}

	regMin := regexp.MustCompile(`[^\d]([0-9]+)m`)
	resultm := regMin.FindStringSubmatch(str)

	if len(resultm) > 1 {
		min, err := strconv.Atoi(resultm[1])
		if err != nil {
			panic(err)
		}
		return time.Minute * time.Duration(min)
	}

	return time.Second * 30
}

func (arrosage Arrosage) save() error {
	return saveEvent(&arrosage.Event)
}

func (arr *Arrosage) arroser() {

	exec.Command("gpio", "write", pinOut["pompe"], "0").Run()
	exec.Command("gpio", "write", pinOut["out1"], "1").Run()
	exec.Command("gpio", "write", pinOut["out2"], "1").Run()
	exec.Command("gpio", "write", pinOut["out3"], "1").Run()
	exec.Command("gpio", "write", pinOut["out4"], "1").Run()

	cmdOpenPompe := exec.Command("gpio", "write", "0", "1")
	cmdClosePompe := exec.Command("gpio", "write", "0", "0")
	cmdOpenVanne := exec.Command("gpio", "write", pinOut[arr.Sortie], "0")
	cmdCloseVanne := exec.Command("gpio", "write", pinOut[arr.Sortie], "1")

	cmdOpenVanne.Run()

	time.Sleep(time.Millisecond * 500)

	cmdOpenPompe.Run()

	time.Sleep(arr.Duration)

	cmdClosePompe.Run()
	time.Sleep(time.Millisecond * 500)
	cmdCloseVanne.Run()
}

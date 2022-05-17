package main

import (
	"fmt"
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

type Arrosage struct {
	Done     bool
	Start    time.Time
	End      time.Time
	Event    calendar.Event
	Duration time.Duration
	Sortie   string
}

type summary string

func (arrosage Arrosage) setDone() Arrosage {
	if arrosage.Done {
		return arrosage
	}

	arrosage.Done = true
	arrosage.Event.Summary = arrosage.Event.Summary + " -OK"

	return arrosage
}

func (str summary) lastchars(length int) (result summary) {

	if len(str) < length {
		return str
	}

	result = str[len(str)-length:]
	return
}

func arrosage() (arrosages []Arrosage) {

	events := getEvents()

	now := time.Now()

	for _, event := range events.Items {

		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)

		outNum := regexp.MustCompile(`out[0-9]`)
		outFound := outNum.FindString(event.Summary)

		if start.Before(now) && end.After(now) && summary(event.Summary).lastchars(3) != "-OK" {
			fmt.Println("out:", outFound)
			arrosages = append(arrosages, Arrosage{
				false,
				start,
				end,
				*event,
				durationFromString(event.Summary),
				outFound,
			})
		}
	}

	return arrosages
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

func (arrosage Arrosage) save() {
	saveEvent(&arrosage.Event)
}

func (arr Arrosage) arroser() {
	fmt.Println("Arrosage en cours : "+arr.Event.Summary, "temps : ", arr.Duration)

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

	arr.save()

	cmdClosePompe.Run()
	time.Sleep(time.Millisecond * 500)
	cmdCloseVanne.Run()
}

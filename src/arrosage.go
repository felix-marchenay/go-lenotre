package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/calendar/v3"
)

type Arrosage struct {
	Done     bool
	Start    time.Time
	End      time.Time
	Event    calendar.Event
	Duration time.Duration
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

func arrosage() Arrosage {

	events := getEvents()

	now := time.Now()

	for _, event := range events.Items {

		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)

		if start.Before(now) && end.After(now) && summary(event.Summary).lastchars(3) != "-OK" {
			return Arrosage{
				false,
				start,
				end,
				*event,
				durationFromString(event.Summary),
			}
		}
	}

	panic("Aucun arrosage")
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
	fmt.Println("Arrosage en cours : " + arr.Event.Summary)

	cmdOpen := exec.Command("gpio", "write", "2", "1")
	cmdClose := exec.Command("gpio", "write", "2", "0")
	err := cmdOpen.Run()

	if err != nil {
		panic(err)
	}

	time.Sleep(arr.Duration)

	arr.save()

	err = cmdClose.Run()
	if err != nil {
		panic(err)
	}
}

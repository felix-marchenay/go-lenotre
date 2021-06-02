package main

import (
	"google.golang.org/api/calendar/v3"
	"time"
)

type Arrosage struct {
	Done bool
	Start time.Time
	End time.Time
	Event calendar.Event
}

type Arrosages []Arrosage

type summary string

func (arrosages Arrosages) toDo() (result Arrosages) {
	for _, arr := range arrosages {
		if !arr.Done {
			result = append(result, arr)
		}
	}

	return
}

func (arrosages Arrosages) current() (result Arrosages) {
	now := time.Now()
	for _, arr := range arrosages {
		if arr.Start.Before(now) && arr.End.After(now) {
			result = append(result, arr)
		}
	}

	return
}

func (arrosage Arrosage) setDone() Arrosage{
	if arrosage.Done {
		return arrosage
	}

	arrosage.Done = true
	arrosage.Event.Summary = arrosage.Event.Summary+" -OK"

	return arrosage
}

func (str summary) lastchars(length int) (result summary) {

	if (len(str) < length) {
		return str
	}

	result = str[len(str)-length:]
	return
}

func arrosages() (arrosages Arrosages) {

	events := getEvents()

	for _, event := range events.Items {
		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)
		arrosages = append(arrosages,Arrosage{summary(event.Summary).lastchars(3) == "-OK", start, end, *event})
	}

	return
}

func (arrosage Arrosage) save() {
	saveEvent(&arrosage.Event)
}
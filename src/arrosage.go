package main

import (
	"google.golang.org/api/calendar/v3"
	"time"
	"fmt"
)

type Arrosage struct {
	Done bool
	Start time.Time
	End time.Time
	Event calendar.Event
}

type Arrosages []Arrosage

type summary string

func (arrosages Arrosages) aEffectuer() (result Arrosages) {
	now := time.Now()
	for _, arr := range arrosages {
		if !arr.Done && arr.Start.After(now) && arr.End.Before(now){
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
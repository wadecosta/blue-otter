package main

import (
	"time"
	"strconv"
	)

type Dash struct {
	ID		int
	Username 	string
	TimeOfDay	string
	Stickies	[]Sticky
}

func GetTimeOfDay() (string) {
	now := time.Now()
	TOD, err := strconv.Atoi(now.Format("15"))
	if err != nil {
		return "ERROR"
	}

	if TOD < 12 {
		return "Good Morning"
	} else if TOD < 18 {
		return "Good Afternoon"
	} else {
		return "Good Evening"
	}
}

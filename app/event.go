package main

import (
	)

type Event struct {
	ID		int
	When		At
	Description	string
}

type Events struct {
	Events []Event
}

type At struct {
	/* Day */
	Day		int
	Month		int
	Year		int
}

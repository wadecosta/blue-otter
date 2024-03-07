package main

import (
	)

type Sticky struct {
	ID		int
	Title		string
	Description	string
}

type Stickies struct {
	Sticky []Sticky
}

type DeleteRequest struct {
	ButtonID 	string `json:"button_id"`
}

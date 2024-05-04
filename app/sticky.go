package main

import (
	)

type Sticky struct {
	ID		int
	Title		string
	Description	string
	DashID		int
}

type Stickies struct {
	Sticky []Sticky
}

type AddRequest struct {
	Title		string `json:"title"`
	Description	string `json:"description"`
}

type DeleteRequest struct {
	ButtonID 	string `json:"button_id"`
}

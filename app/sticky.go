package main

import (
	)

type Sticky struct {
	ID		int
	Title		string
	Description	string
	DashID		int
}

type ModSticky struct {
	ID		int
	UserID		int
	Description	string
	Title		string
	AESKey		string
}

type AddRequest struct {
	Title		string `json:"title"`
	Description	string `json:"description"`
}

type ModStickyRequest struct {
	ID		string `json:"id"`
	Old_Description string `json:"old_sticky_description"`
	Old_Title	string `json:"old_sticky_title"`
	New_Description string `json:"new_sticky_description"`
	New_Title	string `json:"new_sticky_title"`
}

type DeleteRequest struct {
	ButtonID 	string `json:"button_id"`
}

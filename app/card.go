package main

import (
	)

type Card struct {
	ID		int
	CardBank	string
	CardName	string
	Balance		string
	DueDate 	string
	DashID          int
}


type AddCardRequest struct {
	CardBank    	string	`json:"card_bank"`
        CardName    	string	`json:"card_name"`
        Balance  	string	`json:"balance"`
        DueDate 	string	`json:"due_date"`
}

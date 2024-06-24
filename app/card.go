package main

import (
	"fmt"
	)

type Card struct {
	ID		int
	CardID		int
	CardBank	string
	CardName	string
	CardArtwork	string
	Balance		string
	DueDate 	string
	DashID          int
}

type CardDetails struct {
        ID              int
	CardBank	string
        CardName        string
	CardArtwork	string
}


type AddCardRequest struct {
	CardID    	string	`json:"card_id"`
        Balance  	string	`json:"balance"`
        DueDate 	string	`json:"due_date"`
}

func getListCards() (cards []CardDetails, err error) {

	stmt := "SELECT * FROM list_cards"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var temp CardDetails
		var ID	int
		var CardBank string
		var CardName string
		var CardArtwork string
		err = rows.Scan(&ID, &CardBank, &CardName, &CardArtwork)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		temp.ID = ID
		temp.CardBank = CardBank
		temp.CardName = CardName
		temp.CardArtwork = CardArtwork

		cards = append(cards, temp)
	}

	return cards, nil

}

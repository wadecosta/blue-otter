package main

import (
	)

type Dash struct {
	ID		int
	Username 	string
	AESKey		string
	Stickies	[]Sticky
	Cards		[]Card
	Banks		[]Bank
	BankAccounts	[]BankAccount
}

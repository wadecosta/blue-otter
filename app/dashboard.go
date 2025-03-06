package main

import (
	)

type Dash struct {
	ID		int
	Username 	string
	AESKey		string
	Admin		bool
	Stickies	[]Sticky
	Cards		[]Card
	Banks		[]Bank
	BankAccounts	[]BankAccount
	CDs		[]CD
}

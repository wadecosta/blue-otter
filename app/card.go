package main

import (
	"os"
	"io"
	"fmt"
	"context"
	"net/http"
	"path/filepath"
	//"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

func AddCardHandler(w http.ResponseWriter, r *http.Request) {
	/* TODO check if user is an admin */

	/* Ensure this is a POST request */
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	/* Parse the multipart form */
	err := r.ParseMultipartForm(10 << 20) /* 10 MB limit for file upload */
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	/* Get the name of the Bank */
	bankName := r.FormValue("bank_name")
	if bankName == "" {
		http.Error(w, "Bank name is required", http.StatusBadRequest)
		return
	}
	fmt.Println("Bank Name:", bankName)

	/* Get the name of the Card */
	cardName := r.FormValue("card_name")
	if cardName == "" {
		http.Error(w, "Card name is required", http.StatusBadRequest)
		return
	}
	fmt.Println("Card Name:", cardName)

	file, header, err := r.FormFile("card_image")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}

	defer file.Close()

	/* Extract the original file name */
	originalFileName := filepath.Base(header.Filename)

	/* Check to make sure data/card_images exists */
	if err := os.MkdirAll("data/card_images", 0755); err != nil {
		fmt.Println("Error creating data/card_images directory:", err)
		return
	}

	/* Create the file to save on the server */
	dstPath := filepath.Join("data/card_images", originalFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	/* Copy the uploaded file to the server */
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}

	/* Insert card name and image location into the database */
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO list_cards (bank_id, card_name, card_artwork, to_delete) VALUES (?,?,?,0);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, bankName, cardName, "data/card_images/" + originalFileName)
	if err != nil {
		fmt.Println("Inserting new card into database:", err)
	}
	defer insertStmt.Close()

	/* Respond with a success message */
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully!", "file_name": "%s"}`, originalFileName)))

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
		var Delete int
		err = rows.Scan(&ID, &CardBank, &CardName, &CardArtwork, &Delete)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if Delete == 0 {
			temp.ID = ID
			temp.CardBank = CardBank
			temp.CardName = CardName
			temp.CardArtwork = CardArtwork

			cards = append(cards, temp)
		}
	}

	return cards, nil

}

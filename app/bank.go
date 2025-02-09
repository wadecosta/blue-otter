package main

import (
	"os"
	"io"
	"fmt"
	"context"
	"net/http"
	"path/filepath"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Bank struct {
	ID	int
	Name	string
	Image	string
}

func checkifExists() {
	if err := os.MkdirAll("data/bank_images", 0755); err != nil {
		fmt.Println("Error creating data/bank_images directory:", err)
		return
	}
}

func getListBanks() (banks []Bank, err error) {
	stmt := "SELECT * FROM list_banks"
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var temp Bank
		var ID		int
		var Name	string
		var Image	string
		err = rows.Scan(&ID, &Name, &Image)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		temp.ID = ID
		temp.Name = Name
		temp.Image = Image

		banks = append(banks, temp)
	}

	return banks, nil
}

func AddBankHandler(w http.ResponseWriter, r *http.Request) {
	
	/* TODO check to see if person is an admin */

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
	bankName := r.FormValue("title")
	if bankName == "" {
		http.Error(w, "Bank name is required", http.StatusBadRequest)
		return
	}
	fmt.Println("Bank Name:", bankName)

	/* Get the file from the request */
	file, header, err := r.FormFile("bank_image")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	/* Extract the original file name */
	originalFileName := filepath.Base(header.Filename)

	/* Check to make sure data/bank_images exists */
	checkifExists()

	/* Create the file to save on the server */
	dstPath := filepath.Join("data/bank_images", originalFileName)
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

	if err != nil {
    		fmt.Println("Form Parsing Error:", err)
    		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
    		return
	}

	/* Insert bank name and image location into the database */
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO list_banks (bank_name, bank_artwork) VALUES (?,?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, bankName, "data/bank_images/" + originalFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer insertStmt.Close()
	
	/* Respond with a success message */
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully!", "file_name": "%s"}`, originalFileName)))
}

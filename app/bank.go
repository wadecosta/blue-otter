package main

import (
	"os"
	"io"
	"fmt"
	"context"
	"net/http"
	"path/filepath"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Bank struct {
	ID	int
	Name	string
	Image	string
	Delete	int
}

type BankAccount struct {
	ID		int
	BankID		int
	BankName	string
	Amount		string
	Delete		int
	DashID		int
}

type AddBankAccountRequest struct {
        Bank		string `json:"bank"`
        Amount		string `json:"amount"`
}

type ModBankAccountRequest struct {
	ID			string	`json:"id"`
	Old_Bank_Account_Amount	string	`json:"old_bank_account_amount"`
	New_Bank_Account_Amount	string	`json:"new_bank_account_amount"`
}

func checkifExists() {
	if err := os.MkdirAll("data/bank_images", 0755); err != nil {
		fmt.Println("Error creating data/bank_images directory:", err)
		return
	}
}

func GetBankName(id int) string {
	var BankName string
	
	stmt := "SELECT bank_name FROM list_banks WHERE id = ?"
	row := db.QueryRow(stmt, id);
	err := row.Scan(&BankName)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			BankName = "ERROR - NO FOUND"

		} else {
			fmt.Println(err)
			BankName = "ERROR - INVAILD"
		}
	}

	return BankName
}

func GetListBanks() (banks []Bank, err error) {
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
		var Delete	int
		err = rows.Scan(&ID, &Name, &Image, &Delete)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		/* Check to see if bank is not deleted */
		if Delete == 0 {
			temp.ID = ID
			temp.Name = Name
			temp.Image = Image
			banks = append(banks, temp)
		}
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
	insertStmt, err = db.Prepare("INSERT INTO list_banks (bank_name, bank_artwork, to_delete) VALUES (?,?,?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, bankName, "data/bank_images/" + originalFileName, 0)
	if err != nil {
		fmt.Println(err)
	}
	defer insertStmt.Close()
	
	/* Respond with a success message */
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "File uploaded successfully!", "file_name": "%s"}`, originalFileName)))
}

func DelBankHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
                fmt.Println(err)
                return
	}

	successMessage := "Delete request received successfully"
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": successMessage})

        fmt.Println(req)

        /* set bank to DELETE in the database*/
        var delStmt *sql.Stmt
        delStmt, err = db.Prepare("UPDATE list_banks SET to_delete = 1 WHERE id = ?")
        if (err != nil) {
                fmt.Println("error preparing statement", err)
                tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem registering this account")
                return
        }

        ctx := context.Background()
        _, err = delStmt.ExecContext(ctx, req.ButtonID)
        if (err != nil) {
                fmt.Println(err)
        }

        defer delStmt.Close()
}

func AddBankAccountHandler(w http.ResponseWriter, r *http.Request) {
	
	var req AddBankAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	bank := req.Bank
	amount := req.Amount

	fmt.Println("Bank: ", bank)
	fmt.Println("Amount: ", amount)

	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	/* insert user bank account data into database */
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO list_bank_accounts (user_id, bank_id, amount, to_delete) VALUES (?, ?, ?, 0);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem inserting this account into the database")
		return
	}

	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, userID, bank, amount)
	if err != nil {
		fmt.Println(err)
	}

	defer insertStmt.Close()

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func EditBankAccountHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var req ModBankAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("error decoding request body:", err)
		return
	}

	id := req.ID
	oldAmount := req.Old_Bank_Account_Amount
	newAmount := req.New_Bank_Account_Amount

	fmt.Println("Requested data:", id, oldAmount, newAmount)

	updateStmt, err := db.Prepare("UPDATE list_bank_accounts SET amount = ? WHERE id = ? AND user_id = ? AND amount = ?")
	
	if err != nil {
		fmt.Println("error preparing statement:", err)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	defer updateStmt.Close()

	ctx := context.Background()
	res, err := updateStmt.ExecContext(ctx, newAmount, id, userID, oldAmount)
	if err != nil {
		fmt.Println("error executing statement:", err)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
                return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("error fetching rows affected:", err)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
                return
	}
	if rowsAffected == 0 {
		fmt.Println("no rows updated")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
                return
	}
		
	http.Redirect(w, r, "/dashboard", http.StatusFound)
	
}

func DelBankAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	
	successMessage := "Delete request received successfully"
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": successMessage})

	fmt.Println(req)

	/* set bank accout to DELETE in the database*/
	var delStmt *sql.Stmt
	delStmt, err = db.Prepare("UPDATE list_bank_accounts SET to_delete = 1 WHERE id = ?")
	if (err != nil) {
		fmt.Println("error preparing statement", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem deleting this bank account!")
		return
	}

	ctx := context.Background()
	_, err = delStmt.ExecContext(ctx, req.ButtonID)
	if (err != nil) {
		fmt.Println(err)
	}
	
	defer delStmt.Close()
}

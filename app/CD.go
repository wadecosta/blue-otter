package main

import(
	"fmt"
	"context"
	"net/http"
	"database/sql"
	"encoding/json"
)

type CD struct {
	ID		int
	BankID		int
	BankName	string
	StartDate	string
	Deposit		string
	Term		string
	Apy		string
	Delete		int
	DashID		int
}

type AddCDRequest struct {
	Bank		string `json:"bank"`
	StartDate	string `json:"startDate"`
	Deposit		string `json:"deposit"`
	Term		string `json:"term"`
	Apy		string `json:"apy"`
}

func AddCDHandler(w http.ResponseWriter, r *http.Request) {

	var req AddCDRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	bank := req.Bank
	startDate := req.StartDate
	deposit := req.Deposit
	term := req.Term
	apy := req.Apy

	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	/* insert user CD into database */
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO CD (user_id, bank_id, start_date, deposit, term, apy, to_delete) VALUES (?, ?, ?, ?, ?, ?, 0);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, userID, bank, startDate, deposit, term, apy)
	if err != nil {
		fmt.Println(err)
	}

	defer insertStmt.Close()

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func DelCDHandler(w http.ResponseWriter, r *http.Request) {
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

        /* set CD to DELETE in the database*/
        var delStmt *sql.Stmt
        delStmt, err = db.Prepare("UPDATE CD SET to_delete = 1 WHERE id = ?")
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

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

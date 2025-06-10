package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"unicode"
	"context"
    	"github.com/gorilla/mux"
    	"github.com/gorilla/sessions"
    	_ "github.com/go-sql-driver/mysql"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

var tpl *template.Template
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:Password@tcp(localhost:3306)/yourdbname")
	if err != nil {
		fmt.Println("Is MySQL running? Please check MySQL!")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
        	panic(err.Error())
	}

	tpl, _ = template.ParseGlob("templates/*.html")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler).Methods("GET")

	/* Server serving Handlers */
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))

	/* Login Handlers */
	router.HandleFunc("/login", LoginHandler).Methods("GET")
	router.HandleFunc("/login", LoginSubmitHandler).Methods("POST")

	/* Vault Handler */
	router.HandleFunc("/vault", VaultHandler).Methods("GET")
	router.HandleFunc("/vault", VaultSubmitHandler).Methods("POST")

	/* User Creation Handlers */
	router.HandleFunc("/register", RegisterHandler).Methods("GET")
	router.HandleFunc("/register", RegisterSubmitHandler).Methods("POST")

	/* User update Handlers */
	router.HandleFunc("/updatePassword", UpdatePasswordHandler).Methods("GET")
	router.HandleFunc("/updatePassword", UpdatePasswordSubmitHandler).Methods("POST")

	/* User Dashboard */
	router.HandleFunc("/dashboard", DashboardHandler).Methods("GET")

	/* Admin Dashboard */
	router.HandleFunc("/admin", AdminHandler).Methods("GET")
	router.HandleFunc("/addBank", AddBankHandler).Methods("POST")
	router.HandleFunc("/delBank", DelBankHandler).Methods("POST")

	/* User Bank Account Handlers */
	router.HandleFunc("/addBankAccount", AddBankAccountHandler).Methods("POST")
	router.HandleFunc("/editBankAccount", EditBankAccountHandler).Methods("POST")
	router.HandleFunc("/delBankAccount", DelBankAccountHandler).Methods("POST")

	/* User CD Handlers */
	router.HandleFunc("/addCD", AddCDHandler).Methods("POST")
	router.HandleFunc("/delCD", DelCDHandler).Methods("POST")
	
	/* User Sticky Handlers */
	router.HandleFunc("/addSticky", AddStickyHandler).Methods("GET")
	router.HandleFunc("/addSticky", AddStickySubmitHandler).Methods("POST")
	router.HandleFunc("/editSticky", EditStickySubmitHandler).Methods("POST")
	router.HandleFunc("/delSticky", DelStickyHandler).Methods("POST")

	/* User Card Handlers */
	//router.HandleFunc("/addCard", AddCardHandler).Methods("GET")
	router.HandleFunc("/addCard", AddCardHandler).Methods("POST")
	router.HandleFunc("/delCard", DelCardHandler).Methods("POST")

	router.HandleFunc("/profile", ProfileHandler).Methods("GET")
	router.HandleFunc("/logout", LogoutHandler).Methods("GET")

	http.Handle("/", router)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "home.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := authenticateUser(username, password)
	if err != nil {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "login.html", "Unable to Login. Please try again.")
        	return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func VaultHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        user, err := getUserByID(userID)
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }

	tpl.ExecuteTemplate(w, "vault.html", user)
}

func VaultSubmitHandler(w http.ResponseWriter, r *http.Request) {
	
	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        user, err := getUserByID(userID)
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
	
	AESKey := r.FormValue("key")
	AESHash := user.AESKey

	fmt.Println("AESHash:", AESHash)
        fmt.Println("AESKey:", AESKey)
	
	err = authenticateVault(AESHash, AESKey)
	if (err != nil) {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "vaultFailure.html", user)
	}

	tpl.ExecuteTemplate(w, "vaultSuccess.html", user)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "register.html", nil)
}

func RegisterSubmitHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username")

	// check username for only alphaNumeric chars
	var nameAlphaNumeric = true
	for _, char := range username {
		if (unicode.IsLetter(char) == false) && (unicode.IsNumber(char) == false) {
			nameAlphaNumeric = false
		}
	}

	// check username length
	var nameLength bool
	if (5 <= len(username)) && (len(username) <= 50) {
		nameLength = true
	}

	if (!nameAlphaNumeric || !nameLength) {
		tpl.ExecuteTemplate(w, "register.html", "Please check username and password criteria")
		return
	}

	/* Check to make sure all values are lowercase */
	allLower := AllLower(username)
	if (!allLower) {
		tpl.ExecuteTemplate(w, "register.html", "Please only use lowercase for username")
		return
	}

	//TODO impliment email restrictions
	email := r.FormValue("email")
	allLower = AllLower(email)
        if (!allLower) {
                tpl.ExecuteTemplate(w, "register.html", "Please only use lowercase for email")
                return
        }

	password := r.FormValue("password")

	stmt := "SELECT id FROM users WHERE username = ? OR email = ?"
	row := db.QueryRow(stmt, username, email)

	var uID string
	err := row.Scan(&uID)

	if (err != sql.ErrNoRows) {
		tpl.ExecuteTemplate(w, "register.html", "username or email is already taken")
		return
	}

	// create hash from password
	var passwordHash []byte
	passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if (err != nil) {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
		return
	}
	fmt.Println("hash:", passwordHash)
	fmt.Println("string(passwordHash):", string(passwordHash))

	// Provide AES-256 Key
	key, err := GenerateAES256Key()
	if err != nil {
		fmt.Println("Error generating AES key:", err)
		return
	}
	fmt.Println("Generated AES-256 key:", key)

	/* Hash AES-256 Key into database */
	var keyHash []byte
	keyHash, err = bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if (err != nil) {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
		return
	}
	fmt.Println("hash:", keyHash)
	fmt.Println("string(keyHash):", string(keyHash))

	// insert user data into database
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users (username, password, email, AESkey, is_admin) VALUES (?, ?, ?, ?, 0);")
	if (err != nil) {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
		return
	}

	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, passwordHash, email, keyHash)
	
	if (err != nil) {
    		fmt.Println("error inserting new user:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
    		return
	}

	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)

	if (err != nil) {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
		return
	}

	user, err := authenticateUser(username, password)
    	if err != nil {
        	http.Redirect(w, r, "/login", http.StatusSeeOther)
        	return
    	}

    	session, _ := store.Get(r, "session-name")
    	session.Values["user_id"] = user.ID
    	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusSeeOther) 
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := getUserByID(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "profile.html", user)
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "updatePassword.html", nil)
}

func UpdatePasswordSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := getUserByID(userID)
    	if err != nil {
        	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        	return
    	}

	/* Check to see if old password is correct */
	r.ParseForm()
	oldPassword := r.FormValue("oldPassword")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
        
	if (err != nil) {
		tpl.ExecuteTemplate(w, "updatePassword.html", "Incorrect old password")
                return
        }

	newPassword := r.FormValue("newPassword")

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error generating hash for the new password", http.StatusInternalServerError)
		return
	}

	/* Update the password in the database */
	_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", newHash, userID)
	if err != nil {
		http.Error(w, "Error updating password in the database", http.StatusInternalServerError)
		return
	}

	/* Password updated successfully */
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AddStickyHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        user, err := getUserByID(userID)
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
	tpl.ExecuteTemplate(w, "addSticky.html", user)
}

func AddStickySubmitHandler(w http.ResponseWriter, r *http.Request) {
	
	var req AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	title := req.Title
	description := req.Description

	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

	/* insert user sticky data into database */
        var insertStmt *sql.Stmt
        insertStmt, err = db.Prepare("INSERT INTO stickies (user_id, sticky_description, sticky_title, to_delete) VALUES (?, ?, ?, 0);")
        if (err != nil) {
                fmt.Println("error preparing statement:", err)
                tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem registering this account")
                return
        }

	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, userID, description, title)
	if err != nil {
		fmt.Println(err)
	}

        defer insertStmt.Close()

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func EditStickySubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var req ModStickyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("error decoding request body:", err)
		return
    	}

	id := req.ID
	oldDescription := req.Old_Description
	oldTitle := req.Old_Title
	newDescription := req.New_Description
	newTitle := req.New_Title

	updateStmt, err := db.Prepare("UPDATE stickies SET sticky_description = ?, sticky_title = ? WHERE id = ? AND user_id = ? AND sticky_description = ? AND sticky_title = ?")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem updating this sticky")
		return
	}
	defer updateStmt.Close()

	ctx := context.Background()
	res, err := updateStmt.ExecContext(ctx, newDescription, newTitle, id, userID, oldDescription, oldTitle)
	if err != nil {
		fmt.Println("error executing statement:", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem updating this sticky")
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("error fetching rows affected:", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem updating this sticky")
		return
	}
	if rowsAffected == 0 {
		fmt.Println("no rows updated")
		tpl.ExecuteTemplate(w, "dashboard.html", "No sticky was updated. Please check the provided details.")
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func DelStickyHandler(w http.ResponseWriter, r *http.Request) {
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

	/* set sticky to DELETE in the database*/
	var delStmt *sql.Stmt
	delStmt, err = db.Prepare("UPDATE stickies SET to_delete = 1 WHERE id = ?")
	if (err != nil) {
		fmt.Println("error preparing statement", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem deleting this sticky!")
		return
	}

	ctx := context.Background()
	_, err = delStmt.ExecContext(ctx, req.ButtonID)
	if (err != nil) {
		fmt.Println(err)
	}
	
	defer delStmt.Close()
}

func DelCardHandler(w http.ResponseWriter, r *http.Request) {
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

        /* set sticky to DELETE in the database*/
        var delStmt *sql.Stmt
        delStmt, err = db.Prepare("UPDATE cards SET to_delete = 1 WHERE id = ?")
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

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var dash Dash
	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        user, err := getUserByID(userID)
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }

	dash.ID = user.ID
	dash.Username = user.Username
	dash.AESKey = user.AESKey
	dash.Admin = user.isAdmin


	/* Load Stickies */
	stmt := "SELECT id, sticky_description, sticky_title FROM stickies WHERE (user_id = ? AND to_delete = 0)"
        rows, err := db.Query(stmt, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	
	dashboard_num := 0

	for rows.Next() {
		var tempSticky Sticky
        	var stickyData string
        	var stickyTitle string
		var stickyID int
        	err = rows.Scan(&stickyID, &stickyData, &stickyTitle)
		if err != nil {
			fmt.Println(err)
			return
		}

		tempSticky.Title = stickyTitle
		tempSticky.Description = stickyData
		tempSticky.ID = stickyID
		tempSticky.DashID = dashboard_num

		dash.Stickies = append(dash.Stickies, tempSticky)
		dashboard_num++
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	/* Load Cards */
	stmt = "SELECT id, card_id, balance, due_date FROM cards WHERE (user_id = ? AND to_delete = 0)"
	rows, err = db.Query(stmt, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	dashboard_num = 0

	for rows.Next() {
		var tempCard Card
		var ID int
		var cardID int
		var Balance string
		var DueDate string
		err = rows.Scan(&ID, &cardID, &Balance, &DueDate)
		if err != nil {
			fmt.Println(err)
			return
		}

		tempCard.ID = ID
		tempCard.CardID = cardID
		tempCard.Balance = Balance
		tempCard.DueDate = DueDate
		tempCard.DashID = dashboard_num

		dash.Cards = append(dash.Cards, tempCard)
		dashboard_num++
	}

	/* Load Banks */
        banks, err := GetListBanks()
        if err != nil {
                fmt.Println(err)
                return
        }
        dash.Banks = banks
	
	dashboard_num = 0

	/* Load Bank Accounts */
	stmt = "SElECT id, bank_id, amount FROM list_bank_accounts WHERE (user_id = ? AND to_delete = 0)"
	rows, err = db.Query(stmt, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempBankAccount BankAccount
		var ID 		int
		var BankID 	int
		var BankName	string
		var Amount 	string
		err = rows.Scan(&ID, &BankID, &Amount)
		if err != nil {
			fmt.Println(err)
			return
		}

		BankName = GetBankName(BankID)

		tempBankAccount.ID = ID
		tempBankAccount.BankID = BankID
		tempBankAccount.BankName = BankName
		tempBankAccount.Amount = Amount
		tempBankAccount.DashID = dashboard_num

		dash.BankAccounts = append(dash.BankAccounts, tempBankAccount)
		dashboard_num++
	}

	dashboard_num = 0

	/* Load CD(s) */
	stmt = "SELECT id, bank_id, deposit, term, apy FROM CD WHERE (user_id = ? AND to_delete = 0)"
	rows, err = db.Query(stmt, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempCD CD
		var ID		int
		var BankID	int
		var BankName	string
		var Deposit	string
		var Term	string
		var Apy		string
		err = rows.Scan(&ID, &BankID, &Deposit, &Term, &Apy)
		if err != nil {
			fmt.Println(err)
			return
		}

		BankName = GetBankName(BankID)

		tempCD.ID = ID
		tempCD.BankID = BankID
		tempCD.BankName = BankName
		tempCD.Deposit = Deposit
		tempCD.Term = Term
		tempCD.Apy = Apy
		tempCD.DashID = dashboard_num

		dash.CDs = append(dash.CDs, tempCD)
		dashboard_num++
	}

	tpl.ExecuteTemplate(w, "dashboard.html", dash)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminDash AdminDash
	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        user, err := getUserByID(userID)
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }

	/* Check if user is really an Admin */
	if user.isAdmin != true {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	adminDash.User = *user

	cards, err := getListCards()
	if err != nil {
		fmt.Println(err)
		return
	}
	adminDash.Cards = cards

	/* Get the list of Banks */
	banks, err := GetListBanks()
	if err != nil {
		fmt.Println(err)
		return
	}
	adminDash.Banks = banks
	
	tpl.ExecuteTemplate(w, "admin.html", adminDash)
}

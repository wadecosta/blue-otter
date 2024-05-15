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
	
	/* User Sticky Handlers */
	router.HandleFunc("/addSticky", AddStickyHandler).Methods("GET")
	router.HandleFunc("/addSticky", AddStickySubmitHandler).Methods("POST")
	router.HandleFunc("/delSticky", DelStickyHandler).Methods("POST")

	/* User Card Handlers */
	router.HandleFunc("/addCard", AddCardHandler).Methods("GET")
	router.HandleFunc("/addCard", AddCardSubmitHandler).Methods("POST")
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

	fmt.Println(err)

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
	insertStmt, err = db.Prepare("INSERT INTO users (username, password, email, AESkey) VALUES (?, ?, ?, ?);")
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

	fmt.Println(userID)

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
	//r.ParseForm()
	//title := r.FormValue("encryptedTitle")
	//description := r.FormValue("encryptedDescription")
	
	var req AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	title := req.Title
	description := req.Description

	//TODO
	fmt.Println(title)
	fmt.Println(description)


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

	fmt.Println(req)

	/* set sticky to DELETE in the database*/
	var delStmt *sql.Stmt
	delStmt, err = db.Prepare("UPDATE stickies SET to_delete = 1 WHERE id = ?")
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

func AddCardHandler(w http.ResponseWriter, r *http.Request) {
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
        tpl.ExecuteTemplate(w, "addCard.html", user)
}

func AddCardSubmitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TODO")
	var req AddCardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	card_bank := req.CardBank
	fmt.Println("Card Bank:", card_bank)
	card_name := req.CardName
	balance := req.Balance
	due_date := req.DueDate

	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO cards (user_id, card_bank, card_name, balance, due_date, to_delete) VALUES (?, ?, ?, ?, ?, 0);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "dashboard.html", "There was a problem adding this card to the database!")
		return
	}
	
	ctx := context.Background()
	_, err = insertStmt.ExecContext(ctx, userID, card_bank, card_name, balance, due_date)
	if err != nil {
		fmt.Println(err)
	}

	defer insertStmt.Close()


	http.Redirect(w, r, "/dashboard", http.StatusFound)
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

        fmt.Println(req)

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
		fmt.Println("Data:", stickyData, "Title:", stickyTitle)

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

	fmt.Println(dash.Stickies)

	/* Load Cards */
	stmt = "SELECT id, card_bank, card_name, balance, due_date FROM cards WHERE (user_id = ? AND to_delete = 0)"
	rows, err = db.Query(stmt, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	dashboard_num = 0

	for rows.Next() {
		var tempCard Card
		var cardID int
		var cardBank string
		var cardName string
		var Balance string
		var DueDate string
		err = rows.Scan(&cardID, &cardBank, &cardName, &Balance, &DueDate)
		if err != nil {
			fmt.Println(err)
			return
		}

		tempCard.ID = cardID
		tempCard.CardBank = cardBank
		tempCard.CardName = cardName
		tempCard.Balance = Balance
		tempCard.DueDate = DueDate
		tempCard.DashID = dashboard_num

		dash.Cards = append(dash.Cards, tempCard)
		dashboard_num++
	}


	tpl.ExecuteTemplate(w, "dashboard.html", dash)
}

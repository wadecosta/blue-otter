package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"unicode"
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

	/* Login Handlers */
	router.HandleFunc("/login", LoginHandler).Methods("GET")
	router.HandleFunc("/login", LoginSubmitHandler).Methods("POST")

	/* User Creation Handlers */
	router.HandleFunc("/register", RegisterHandler).Methods("GET")
	router.HandleFunc("/register", RegisterSubmitHandler).Methods("POST")

	/* User update Handlers */
	router.HandleFunc("/updatePassword", UpdatePasswordHandler).Methods("GET")
	router.HandleFunc("/updatePassword", UpdatePasswordSubmitHandler).Methods("POST")

	/* User Dashboard */
	router.HandleFunc("/dashboard", DashboardHandler).Methods("GET")
	
	/* User Event Handlers */
	router.HandleFunc("/addEvent", AddEventHandler).Methods("GET")
	router.HandleFunc("/addEvent", AddEventSubmitHandler).Methods("POST")

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
		tpl.ExecuteTemplate(w, "login.html", "Unable to Login. Please try again.")
        	return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if (err != nil) {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))

	// insert user data into database
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users (username, password, email) VALUES (?, ?, ?);")
	if (err != nil) {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering this account")
		return
	}

	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, hash, email)
	
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

func AddEventHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "addEvent.html", nil)
}

func AddEventSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
        
	description := r.FormValue("description")
	when := r.FormValue("when")

	key, salt := deriveKey(r.FormValue("password"))

	session, _ := store.Get(r, "session-name")
        userID, ok := session.Values["user_id"].(int)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

	encryptedDescription, err := encrypt(description, key)
	if err != nil {
		http.Error(w, "Failed to encrypt description", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	encryptedWhen, err := encrypt(when, key)
	if err != nil {
		http.Error(w, "Failed to encrypt when", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	fmt.Println(userID, description, when, r.FormValue("password"), encryptedDescription, encryptedWhen)

	decryptedDescription, err := decrypt(encryptedDescription, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	decryptedWhen, err := decrypt(encryptedWhen, key)
        if err != nil {
                fmt.Println("Decryption error:", err)
                return
        }

	fmt.Println("key:", key)
	fmt.Println("salt:", salt)

	fmt.Println("Decrypted Description", decryptedDescription)
	fmt.Println("Decrypted When", decryptedWhen)
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

	tod := GetTimeOfDay()
	dash.TimeOfDay = tod
	
	var eventsData Events
	
	stmt := "SELECT event_description, event_when FROM events WHERE user_id = ?"
        rows, err := db.Query(stmt, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
        	var eventData string
        	var eventWhen string
        	err = rows.Scan(&eventData, &eventWhen)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Data:", eventData, "When:", eventWhen)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	dash.Events = eventsData.Events

	tpl.ExecuteTemplate(w, "dashboard.html", dash)
}

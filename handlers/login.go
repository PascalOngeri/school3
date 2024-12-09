package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type API struct {
	Name  string
	Icon  string
	IName string
}

type LoginData struct {
	Name     string
	Icon     string
	Username string
	Password string
	Remember bool
}

func getAPIDetails(db *sql.DB) (API, error) {
	var api API
	query := "SELECT name, icon, iname FROM api LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&api.Name, &api.Icon, &api.IName)
	if err != nil {
		log.Printf("Error fetching API details: %v", err)
		return api, err
	}
	return api, nil
}

func renderLoginPage(w http.ResponseWriter, api API) {
	loginData := LoginData{
		Name:     api.Name,
		Icon:     api.Icon,
		Username: "",
		Password: "",
		Remember: false,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, loginData)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		remember := r.FormValue("remember") == "on"

		var userID int
		var foundInAdmin bool
		var adm, phone string

		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
		if err == nil {
			foundInAdmin = true
		} else {
			queryRegistration := "SELECT id, adm, username, phone, password FROM registration WHERE username = ? AND password = ?"
			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password)
			if err != nil {
				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
				return
			}
		}

		session, _ := store.Get(r, "store")
		if foundInAdmin {
			session.Values["sturecmsaid"] = userID
			session.Values["username"] = username
		} else {
			session.Values["sturecmsaid"] = userID
			session.Values["adm"] = adm
			session.Values["username"] = username
			session.Values["phone"] = phone
			session.Values["password"] = password
		}

		session.Save(r, w)

		http.SetCookie(w, &http.Cookie{Name: "user_login", Value: username, Path: "/", MaxAge: 86400})
		if remember {
			http.SetCookie(w, &http.Cookie{Name: "userpassword", Value: password, Path: "/", MaxAge: 86400})
		}

		if foundInAdmin {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
		}
		return
	}

	api, _ := getAPIDetails(db)
	renderLoginPage(w, api)
}

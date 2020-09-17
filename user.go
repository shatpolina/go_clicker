package main

import (
     "fmt"
     "net/http"
     "io/ioutil"
)

type User struct {
    UserID int `json:"id"`
    Login string `json:"login"`
    Password string `json:"password"`
    Num int `json:"num"`
    Numps int `json:"numps"`
    Price_clicker int `json:"price_clicker"`
    Qty_clicker int `json:"qty_clicker"`
    Price_super_clicker int `json:"price_super_clicker"`
    Qty_super_clicker int `json:"qty_super_clicker"`
    Price_mega_clicker int `json:"price_mega_clicker"`
    Qty_mega_clicker int `json:"qty_mega_clicker"`
    Connected bool `json:"connected"`
}

func Registration(w http.ResponseWriter, r *http.Request, login string, password string) {
    var userID int
    db.QueryRow("insert into users (Login, Password, Num, Price_clicker, Qty_clicker, Price_super_clicker, Qty_super_clicker, Price_mega_clicker, Qty_mega_clicker, Numps) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning UserID", login, password, 0, 10, 0, 100, 0, 1000, 0, 0).Scan(&userID)
    c, _:= r.Cookie("session")
    db.Exec("update sessions set UserID = $1 where UUID = $2", userID, c.Value)
    http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func Authorization(w http.ResponseWriter, r *http.Request, login string, password string) {
    var userID int
    err := db.QueryRow("select UserID from users where Login = $1 and Password = $2", login, password).Scan(&userID)
    if err == nil {
        c, _:= r.Cookie("session")
        db.Exec("update sessions set UserID = $1 where UUID = $2", userID, c.Value)
        http.Redirect(w, r, "/home", http.StatusSeeOther)
    } else {
        fmt.Fprintf(w, string("Неверное имя пользователя или пароль"))
    }
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    if len(r.Form["login"]) == 1 {
        login := r.Form["login"][0]
        var value int
        db.QueryRow("select UserID from users where Login = $1", login).Scan(&value)
        if value > 0 {
            fmt.Fprintf(w, string("Login taken, return on register"))
        }
    }
}

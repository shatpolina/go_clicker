package main

import (
    "net/http"
)

func set_routes() {
    http.HandleFunc("/auth", Auth_page)
    http.HandleFunc("/reg", Reg_page)
    http.HandleFunc("/checklogin", CheckLogin)
    http.HandleFunc("/exit", Exit)
    http.HandleFunc("/home", HomePage)
    http.HandleFunc("/clicknum", Clicknum)
    http.HandleFunc("/ws", wsEndpoint)
    http.HandleFunc("/clickerbuy", ClickerBuy)
    http.HandleFunc("/superclickerbuy", SuperClickerBuy)
    http.HandleFunc("/megaclickerbuy", MegaClickerBuy)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    
    dat, _ := ioutil.ReadFile("./button.html")
    fmt.Fprintf(w, string(dat))
}

func Exit(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    c, _:= r.Cookie("session")
    db.Exec("update sessions set UserID = $1 where UUID = $2", nil, c.Value)
    fmt.Println("user exit")
    http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

func Auth_page(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "/home", "")
    r.ParseForm()

    if len(r.Form["login"]) == 1 && len(r.Form["password"]) == 1 {
        login := r.Form["login"][0]
        password := r.Form["password"][0]
        Authorization(w, r, login, password)
    } else {
        dat, _ := ioutil.ReadFile("./authorization.html")
        fmt.Fprintf(w, string(dat))
    }
}

func Reg_page(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "/home", "")
    r.ParseForm()
    
    if len(r.Form["login"]) == 1 && len(r.Form["password"]) == 1 {
        CheckLogin(w, r)
        login := r.Form["login"][0]
        password := r.Form["password"][0]
        Registration(w, r, login, password)
    } else {
        dat, _ := ioutil.ReadFile("./registration.html")
        fmt.Fprintf(w, string(dat))
    }
}

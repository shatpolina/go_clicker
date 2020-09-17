package main

import (
    "fmt"
    "net/http"
)

type Clicker struct {
    UserID int `json:"id"`
    Numps int `json:"numps"`
    Num int `json:"num"`
    Price_clicker int `json:"price_clicker"`
    Qty_clicker int `json:"qty_clicker"`
    Price_super_clicker int `json:"price_super_clicker"`
    Qty_super_clicker int `json:"qty_super_clicker"`
    Price_mega_clicker int `json:"price_mega_clicker"`
    Qty_mega_clicker int `json:"qty_mega_clicker"`
}

func Clicknum(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    c, _ := r.Cookie("session")
    var userID int
    db.QueryRow("select UserID from sessions where UUID = $1", c.Value).Scan(&userID)
    var num int
    db.QueryRow("select Num from users where UserID = $1", userID).Scan(&num)
    num += 1
    fmt.Fprintf(w, "%d", num)
    db.Exec("update users set Num = $1 where UserID = $2", num, userID)
}

func Clickers(userID int) {
    var num int
    var numps int
    db.QueryRow("select Num, Numps from users where UserID = $1", userID).Scan(&num, &numps)
    num += numps
    db.Exec("update users set Num = $1 where UserID = $2", num, userID)
}

func Buy(w http.ResponseWriter, r *http.Request, qty_name string, price_name string) {
    checkSession(w, r, "", "/auth")
    c, _ := r.Cookie("session")
    var userID int
    db.QueryRow("select UserID from sessions where UUID = $1", c.Value).Scan(&userID)
    var num int
    var numps int
    var qty_value int
    var price_value int
    var select_text string = "select " + qty_name + ", " + price_name + ", Num, Numps from users where UserID = $1"
    db.QueryRow(select_text, userID).Scan(&qty_value, &price_value, &num, &numps)
    num -= price_value
    if num > 0 {
        qty_value += 1
        price_value = int(float64(price_value) * 1.1)
        if qty_name == "Qty_clicker" {
            numps += qty_value * 1
        } else if qty_name == "Qty_super_clicker" {
            numps += qty_value * 10
        } else if qty_name == "Qty_mega_clicker" {
            numps += qty_value * 100
        }
        
        db.Exec("update users set " + qty_name + "= $1," + price_name + "= $2, Num = $3, Numps = $4 where UserID = $5", qty_value, price_value, num, numps, userID)
        fmt.Fprintf(w, "%d", qty_value)
    } else {
        fmt.Fprintf(w, "error")
    }  
}

func ClickerBuy (w http.ResponseWriter, r *http.Request) {
        Buy(w, r, "Qty_clicker", "Price_clicker")
    }
    
func SuperClickerBuy (w http.ResponseWriter, r *http.Request) {
        Buy(w, r, "Qty_super_clicker", "Price_super_clicker")
    }
    
func MegaClickerBuy (w http.ResponseWriter, r *http.Request) {
        Buy(w, r, "Qty_mega_clicker", "Price_mega_clicker")
    }

package main

import (
    "fmt"
    "net/http"
    "math/big"
)

type Clicker struct {
    UserID int `json:"id"`
    Numps string `json:"numps"`
    Num string `json:"num"`
    Price_clicker string `json:"price_clicker"`
    Qty_clicker int `json:"qty_clicker"`
    Price_super_clicker string `json:"price_super_clicker"`
    Qty_super_clicker int `json:"qty_super_clicker"`
    Price_mega_clicker string `json:"price_mega_clicker"`
    Qty_mega_clicker int `json:"qty_mega_clicker"`
}

func Clicknum(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    c, _ := r.Cookie("session")
    
    var userID int
    db.QueryRow("select UserID from sessions where UUID = $1", c.Value).Scan(&userID)
    
    var num string
    db.QueryRow("select Num from users where UserID = $1", userID).Scan(&num)
    
    var num_value big.Int
    num_value.SetString(num, 10)
    num_value.Add(&num_value, big.NewInt(int64(1)))
    num = num_value.String()
    fmt.Fprintf(w, num)
    
    db.Exec("update users set Num = $1 where UserID = $2", num, userID)
}

func Clickers(userID int) {
    var num string
    var numps string
    db.QueryRow("select Num, Numps from users where UserID = $1", userID).Scan(&num, &numps)
    var num_value big.Int
    num_value.SetString(num, 10)
    var numps_value big.Int
    numps_value.SetString(numps, 10)
    num_value.Add(&num_value, &numps_value)
    num = num_value.String()
    db.Exec("update users set Num = $1 where UserID = $2", num, userID)
}

func Buy(w http.ResponseWriter, r *http.Request, qty_name string, price_name string) {
    checkSession(w, r, "", "/auth")
    c, _ := r.Cookie("session")
    
    var userID int
    db.QueryRow("select UserID from sessions where UUID = $1", c.Value).Scan(&userID)
    var num string
    var numps string
    var qty int
    var price string
    var select_text string = "select " + qty_name + ", " + price_name + ", Num, Numps from users where UserID = $1"
    
    db.QueryRow(select_text, userID).Scan(&qty, &price, &num ,&numps)
    
    var num_value big.Int
    num_value.SetString(num, 10)
    var price_value big.Int
    price_value.SetString(price, 10)
    cmp := num_value.Cmp(&price_value)
    
    if cmp >= 0 {
        num_value.Sub(&num_value, &price_value)
        num = num_value.String()
        
        db.Exec("update users set Num = $1 where UserID = $2", num, userID)
        
        var price_float big.Float
        price_float.SetInt(&price_value)
        var k big.Float
        k.SetFloat64(1.2)
        price_float.Mul(&price_float, &k)
        price_float.Int(&price_value)
        
        db.Exec("update users set " + price_name + " = $1 where UserID = $2", price, userID)

        qty += 1
        var numps_value big.Int
        numps_value.SetString(numps, 10)

        if qty_name == "Qty_clicker" {
            numps_value.Add(&numps_value, big.NewInt(int64(1)))
        } else if qty_name == "Qty_super_clicker" {
            numps_value.Add(&numps_value, big.NewInt(int64(5)))
        } else if qty_name == "Qty_mega_clicker" {
            numps_value.Add(&numps_value, big.NewInt(int64(10)))
        }
        
        numps = numps_value.String()
        price = price_value.String()
        
        db.Exec("update users set " + qty_name + "= $1," + price_name + "= $2, Numps = $3 where UserID = $4", qty, price, numps, userID)
        
        fmt.Fprintf(w, "%d", qty)
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

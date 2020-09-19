package main

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "encoding/json"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func test_ws (conn *websocket.Conn) {
    for {
        _, value, err := conn.ReadMessage()
        if value != nil {
            fmt.Println(value)
        }
        if err != nil {
            fmt.Println(err)
            break
        }
    }
}

func wsWorker(conn *websocket.Conn, userID int) {
    for {
        Clickers(userID)
        
        var num string
        var numps string
        var price_clicker string
        var qty_clicker int
        var price_super_clicker string
        var qty_super_clicker int 
        var price_mega_clicker string
        var qty_mega_clicker int
        
        db.QueryRow("select Num, Numps, Price_clicker, Qty_clicker, Price_super_clicker, Qty_super_clicker, Price_mega_clicker, Qty_mega_clicker from users where UserID = $1", userID).Scan(&num, &numps, &price_clicker, &qty_clicker, &price_super_clicker, &qty_super_clicker, &price_mega_clicker, &qty_mega_clicker)
    
        var message Clicker
        message.UserID = userID
        message.Numps = numps
        message.Num = num
        message.Price_clicker = price_clicker
        message.Qty_clicker = qty_clicker
        message.Price_super_clicker = price_super_clicker
        message.Qty_super_clicker = qty_super_clicker
        message.Price_mega_clicker = price_mega_clicker
        message.Qty_mega_clicker = qty_mega_clicker
           
        data, _ := json.Marshal(message)
        msg := string(data)
            
        if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
            log.Println(err)
            return
        }
        time.Sleep(1 * time.Second)
        defer db.Exec("update users set Connected = $1 where UserID = $2", false, userID)
    }
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    if len(r.Form["session"]) != 1 {
        return
    }
    uuid := r.Form["session"][0]
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }
    log.Println("Client Connected")
    err = ws.WriteMessage(1, []byte("Hi Client!"))
    if err != nil {
        log.Println(err)
    }
    
    var userID int
    var connected bool
    db.QueryRow("select UserID from sessions where UUID = $1", uuid).Scan(&userID)
    db.QueryRow("select Connected from users where UserID = $1", userID).Scan(&connected)
    
    if !connected {
        go test_ws(ws)
        go wsWorker(ws, userID)
        db.Exec("update users set Connected = $1 where UserID = $2", true, userID)
    }
}

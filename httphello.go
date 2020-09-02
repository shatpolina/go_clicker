package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "time"
    "crypto/rand"
    _ "github.com/lib/pq"
    "database/sql"
)

var db *sql.DB

type Session struct {
    UUID string `json:"uuid"`
    UserID int `json:"id"`
}

type User struct {
    UserID int `json:"id"`
    Login string `json:"login"`
    Password string `json:"password"`
    Num int `json:"num"`
}

func main() {
    connStr := "user=postgres password=postgres dbname=httphello sslmode=disable"
    local_db, err := sql.Open("postgres", connStr)
    db = local_db
    if err != nil {
        panic(err)
    } 
    defer db.Close()
    
    http.HandleFunc("/auth", Authorization)
    http.HandleFunc("/reg", Registration)
    http.HandleFunc("/checklogin", CheckLogin)
    http.HandleFunc("/exit", Exit)
    http.HandleFunc("/hello", HelloServer)
    http.HandleFunc("/num", Givenum)
    http.ListenAndServe(":8080", nil)
}

func getUUID()(uuid string) {
    b := make([]byte, 16)
    rand.Read(b)
    uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return
}

func setCookie(w http.ResponseWriter, name string, value string, ttl time.Duration){        
    expire := time.Now().Add(ttl*time.Minute)
    cookie := http.Cookie{
        Name:    name,
        Value:   value,
        Expires: expire,
    }
    http.SetCookie(w, &cookie)
}

func setSession(w http.ResponseWriter) {
    uuid := getUUID()
    setCookie(w, "session", uuid, 60)
    _, err := db.Exec("insert into sessions (UUID, UserID) values ($1, $2)", uuid, nil)
    if err != nil {
        panic(err)
    }
    fmt.Println("Cookie set")
}

func checkSession(w http.ResponseWriter, r *http.Request, redirOnAuth string, redirOnNoAuth string) {
    c, err := r.Cookie("session")
    if err != nil {
        setSession(w)
        if len(redirOnNoAuth) > 0 {
            http.Redirect(w, r, redirOnNoAuth, http.StatusSeeOther)
        }
    } else {
        row := db.QueryRow("select * from sessions where UUID = $1", c.Value);
        var dbUserID sql.NullInt64
        var UUID string
        err = row.Scan(&UUID, &dbUserID)
        if err != nil {
            setSession(w)
        }
        if dbUserID.Valid {
            if len(redirOnAuth) > 0 {
                http.Redirect(w, r, redirOnAuth, http.StatusSeeOther)
            }
            fmt.Println("AutoAuth")
        } else {
            if len(redirOnNoAuth) > 0 {
                http.Redirect(w, r, redirOnNoAuth, http.StatusSeeOther)
            }
        }
    }
    fmt.Println("Session check")
}

func Authorization(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "/hello", "")
    r.ParseForm()

    if len(r.Form["login"]) == 1 && len(r.Form["password"]) == 1 {
        login := r.Form["login"][0]
        password := r.Form["password"][0]
        var userID int
        err := db.QueryRow("select UserID from users where Login = $1 and Password = $2", login, password).Scan(&userID)
        if err == nil {
            c, _:= r.Cookie("session")
            db.Exec("update sessions set UserID = $1 where UUID = $2", userID, c.Value)
            http.Redirect(w, r, "/hello", http.StatusSeeOther)
        } else {
            fmt.Fprintf(w, string("Неверное имя пользователя или пароль"))
        }
    } else {
        dat, _ := ioutil.ReadFile("./authorization.html")
        fmt.Fprintf(w, string(dat))
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

func Registration(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "/hello", "")
    r.ParseForm()
    
    if len(r.Form["login"]) == 1 && len(r.Form["password"]) == 1 {
        CheckLogin(w, r)
        login := r.Form["login"][0]
        password := r.Form["password"][0]
        var userID int
        db.QueryRow("insert into users (Login, Password, Num) values ($1, $2, $3) returning UserID", login, password, 0).Scan(&userID)
        c, _:= r.Cookie("session")
        db.Exec("update sessions set UserID = $1 where UUID = $2", userID, c.Value)
        http.Redirect(w, r, "/hello", http.StatusSeeOther)
    } else {
        dat, _ := ioutil.ReadFile("./registration.html")
        fmt.Fprintf(w, string(dat))
    }
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    dat, _ := ioutil.ReadFile("./button.html")
    fmt.Fprintf(w, string(dat))
}

func Givenum(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    c, _ := r.Cookie("session")
    var userID int
    db.QueryRow("select UserID from sessions where UUID = $1", c.Value).Scan(&userID)
    var num int
    db.QueryRow("select Num from users where UserID = $1", userID).Scan(&num)
    fmt.Fprintf(w, "%d", num)
    db.Exec("update users set Num = $1 where UserID = $2", num + 1, userID)
}

func Exit(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r, "", "/auth")
    c, _:= r.Cookie("session")
    db.Exec("update sessions set UserID = $1 where UUID = $2", nil, c.Value)
    fmt.Println("user exit")
    http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

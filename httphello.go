package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "time"
    "crypto/rand"
    "encoding/json"
)

type Session struct {
    UserID int `json:"id"`
}

type User struct {
    Login string `json:"login"`
    Password string `json:"password"`
    Num int `json:"num"`
}

var sessionDB (map[string]*Session)
var userDB (map[int]*User)

func readDB(file string, container interface{}) {
    rawDataIn, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Println("Cannot load " + file)
    }
    err = json.Unmarshal(rawDataIn, &container)
    if err != nil {
        fmt.Println("Invalid format " + file, err)
    }
}

func saveDB(file string, container interface{}) {
    data, _ := json.Marshal(container)
    ioutil.WriteFile(file, []byte(data), 0644)
    fmt.Println("DB " + file + " saved")
}

func main() {
    sessionDB = make(map[string]*Session)
    userDB = make(map[int]*User)
    readDB("sessions.json", &sessionDB)
    readDB("users.json", &userDB)
    
    http.HandleFunc("/auth", Authorization)
    http.HandleFunc("/reg", Registration)
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
    var newSession *Session = new(Session)
    newSession.UserID = -1
    sessionDB[uuid] = newSession
    fmt.Println("Cookie set")
}

func checkSession(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("session")
    if err != nil {
        setSession(w)
    } else {
        uuid := c.Value   
        if session, ok := sessionDB[uuid]; ok {
            if session.UserID != -1 {
                fmt.Println("AutoAuth")
                http.Redirect(w, r, "/hello", http.StatusSeeOther)
            }
        } else {
            setSession(w)
        }
    }
    fmt.Println("Session check")
}

func Authorization(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r)
    r.ParseForm()

    if len(r.Form["login"]) == 1 && len(r.Form["password"]) == 1 {
        auth := false
        for ID, user := range userDB {
            if user.Login == r.Form["login"][0] && user.Password == r.Form["password"][0] {
                c, _:= r.Cookie("session")
                var session *Session = sessionDB[c.Value]
                session.UserID = ID
                saveDB("sessions.json", sessionDB)
                fmt.Println(ID, session.UserID, c.Value)
                auth = true
                break
            }
        }
        if auth {
            http.Redirect(w, r, "/hello", http.StatusSeeOther)
        } else {
            fmt.Fprintf(w, string("Неверное имя пользователя или пароль"))
        }
    } else {
        dat, _ := ioutil.ReadFile("./authorization.html")
        fmt.Fprintf(w, string(dat))
    }
}

func Registration(w http.ResponseWriter, r *http.Request) {
    checkSession(w, r)
    r.ParseForm()
    
    if len(r.Form["login"]) == 1 && len(r.Form["password"]) == 1 {
        user_id := len(userDB)
        var user *User = new(User)
        user.Login = r.Form["login"][0]
        user.Password = r.Form["password"][0]
        user.Num = 0
        userDB[user_id] = user
        saveDB("users.json", userDB)
        c, _:= r.Cookie("session")
        var session *Session = sessionDB[c.Value] 
        session.UserID = user_id
        saveDB("sessions.json", sessionDB)
        fmt.Println(user_id, session.UserID, c.Value)
        http.Redirect(w, r, "/hello", http.StatusSeeOther)
    } else {
        dat, _ := ioutil.ReadFile("./registration.html")
        fmt.Fprintf(w, string(dat))
    }
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("session")
    if err != nil {
        http.Redirect(w, r, "/auth", http.StatusSeeOther)
    } else {
        uuid := c.Value   
        if session, ok := sessionDB[uuid]; ok {
            if session.UserID != -1 {
                fmt.Println("/hello AutoAuth")
            } else {
                http.Redirect(w, r, "/auth", http.StatusSeeOther)
            }
        } else {
            http.Redirect(w, r, "/auth", http.StatusSeeOther)
        }
    }

    saveDB("sessions.json", sessionDB)
    saveDB("users.json", userDB)
    
    dat, _ := ioutil.ReadFile("./button.html")
    fmt.Fprintf(w, string(dat))
}

func Givenum(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("session")
    if err != nil {
        http.Redirect(w, r, "/auth", http.StatusSeeOther)
    } else {
        uuid := c.Value
        if session, ok := sessionDB[uuid]; ok {
            var user *User = userDB[session.UserID]
            fmt.Fprintf(w, "%d", user.Num)
            user.Num += 1
            saveDB("users.json", userDB)
        } else {
            fmt.Fprintf(w, "ТЫ ЧО, МЕНЯ НАЕБАТЬ РЕШИЛ??!!!")
        }
    }
}

func Exit(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("session")
    if err != nil {
        http.Redirect(w, r, "/auth", http.StatusSeeOther)
    } else { 
        if _, ok := sessionDB[c.Value]; ok {
            var session *Session = sessionDB[c.Value]
            session.UserID = -1
            saveDB("sessions.json", sessionDB)
            fmt.Println("User exit, c.Value: " + c.Value)
            http.Redirect(w, r, "/auth", http.StatusSeeOther)
        } else {
            http.Redirect(w, r, "/auth", http.StatusSeeOther)
        }
    }
}

package main

import (
     "crypto/rand"
     "fmt"
     "net/http"
     "database/sql"
     "time"
)

type Session struct {
    UUID string `json:"uuid"`
    UserID int `json:"id"`
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

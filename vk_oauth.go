package main

import (
    "fmt"
    "net/http"
    "context"
    "log"
    "golang.org/x/oauth2"
    "github.com/go-vk-api/vk"
)

type User_vk struct {
    ID        int64  `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Photo     string `json:"photo_400_orig"`
    City      City   `json:"city"`
}

type City struct {
    Title string `json:"title"`
}

var ctx = context.Background()
var conf = &oauth2.Config{
    ClientID:     obj_vk.ClientID,
    ClientSecret: obj_vk.ClientSecret,
    Scopes:       []string{},
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://oauth.vk.com/authorize",
        TokenURL: "https://oauth.vk.com/access_token",
        },
    }

func getCurrentUser(api *vk.Client) User_vk {
    var users []User_vk

    api.CallMethod("users.get", vk.RequestParams{
        "fields": "photo_400_orig,city",
    }, &users)

    return users[0]
}

func VK_oauth(w http.ResponseWriter, r *http.Request) {
    url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
    fmt.Println("Visit the URL for the auth dialog: %v", url)
    fmt.Fprintf(w, "Visit the URL for the auth dialog: %v", url)
}

func VK_auth(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Oy e, auth it's ok!")
    var code string
    if _, err := fmt.Scan(&code); err != nil {
        log.Fatal(err)
    }
    
    tok, err := conf.Exchange(ctx, code)
    if err != nil {
        log.Fatal(err)
    }

    client, err := vk.NewClientWithOptions(vk.WithToken(tok.AccessToken))
    if err != nil {
        log.Fatal(err)
    }
    user := getCurrentUser(client)
    fmt.Println("Result: ", user)
}

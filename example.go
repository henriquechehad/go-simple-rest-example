package main

import (
    "encoding/json"
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "net/http"
)

type Conn struct {
    db gorm.DB
}

func (DB *Conn) Init() {
    DB.db, _ = gorm.Open("sqlite3", "messages.db")
    DB.db.LogMode(true)
    DB.db.CreateTable(&Message{})
    DB.db.AutoMigrate()
}

func (DB *Conn) ListAll() []Message {
    messages := []Message{}
    DB.db.Find(&messages)

    return messages
}

type Message struct {
    ID   int
    Text string `sql:"type:varchar(100)"`
}

func (DB *Conn) ListMessages(w http.ResponseWriter, r *http.Request) {
    messages := DB.ListAll()

    b, err := json.Marshal(messages)

    if err != nil {
        fmt.Fprintf(w, "err")
    }

    w.Write(b)
}

func (DB *Conn) CreateMessage(w http.ResponseWriter, r *http.Request) {
    msg := r.URL.Query().Get("message")

    if len(msg) > 1 {
        m := Message{Text: msg}

        // Save message
        DB.db.Save(&m)

        // Return created message
        messages := Message{}

        b, err := json.Marshal(DB.db.Last(&messages))

        if err != nil {
            fmt.Fprintf(w, "Error: %s", err)
        }

        w.Write(b)

    } else {
        fmt.Fprintf(w, "\"message\" param is required")
    }

}

func main() {
    fmt.Println("Starting")

    db := Conn{}
    db.Init()

    http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
        db.CreateMessage(w, r)
    })

    http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        db.ListMessages(w, r)
    })

    fmt.Println("Listening on port 80")
    http.ListenAndServe(":80", nil)

}

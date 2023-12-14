package main

import (
    "database/sql"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Chat struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Message struct {
	ChatID  int    `json:"chat_id"`
	Content string `json:"content"`
}


func main() {
    // Connect to MySQL database
    var err error
    db, err = sql.Open("mysql", "root:5256@tcp(localhost:3306)/chat_system_test")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    router := gin.Default()

    router.POST("/chats/create", CreateChatHandler)
    router.POST("/messages/create", AddMessageHandler)

    port := 8080
    fmt.Printf("Server is running on port %d...\n", port)
    http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func CreateChatHandler(c *gin.Context) {
    var newChat Chat

    // Bind JSON body to the Chat struct
    if err := c.ShouldBindJSON(&newChat); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validate chat name
    if newChat.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Chat name cannot be empty"})
        return
    }

    // Insert the chat item into the database
    _, err := db.Exec("INSERT INTO go (name) VALUES (?)", newChat.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat item"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully", "chat": newChat})
}


func AddMessageHandler(c *gin.Context) {
	var newMessage Message

    // Bind JSON body to the message struct
    if err := c.ShouldBindJSON(&newMessage); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    // Validate chat ID and message content
    if newMessage.ChatID <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
        return
    }
    if newMessage.Content == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Message content cannot be empty"})
        return
    }

    // Insert the message into the database
    _, err := db.Exec("INSERT INTO messages (chat_id, content) VALUES (?, ?)", newMessage.ChatID, newMessage.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add message"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Message added successfully"})
}
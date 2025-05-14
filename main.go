package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("SIGNUP_SECRET"))
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	router.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		phone := c.PostForm("phone")
		password := c.PostForm("password")

		if password != os.Getenv("SIGNUP_SECRET") {
			c.HTML(http.StatusUnauthorized, "register.html", gin.H{"error": "Incorrect password."})
			return
		}

		_, err := db.Exec("INSERT INTO users (username, phone) VALUES (?, ?)", username, phone)
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Username or phone already exists."})
			return
		}
		c.HTML(http.StatusOK, "success.html", gin.H{"username": username})
	})
	router.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})

	router.Run(":8080")
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	// gorm
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetDBConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "tester:secret@tcp(db:3306)/test")

	if err != nil {
		panic(err.Error())
	}
	return db
}

// 一覧取得
func dbGetAll() *gorm.DB {
	db := GetDBConnection()

	// defer db.Close()

	var users []User
	result := db.Find(&users)
	println("user:", result)

	return result
}

// func getUsers() []*User {
func getUsers() []*User {
	// Open up our database connection.
	db, err := sql.Open("mysql", "tester:secret@tcp(db:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	// Execute the query
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var users []*User
	for results.Next() {
		var u User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&u.ID, &u.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		users = append(users, &u)
	}

	return users
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func userPage(w http.ResponseWriter, r *http.Request) {
	users := getUsers()

	fmt.Println("Endpoint Hit: usersPage")
	json.NewEncoder(w).Encode(users)
}

// 一覧取得
// func dbGetAll() []User {
// 	db := GetDBConnection()

// 	defer db.Close()

// 	var users []User
// 	db.Order("created_at desc").Find(&users)
// 	return users
// }

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/index.html")

	// 一覧
	r.GET("/userlist", func(c *gin.Context) {
		// users := dbGetAll()
		users := getUsers()

		if users != nil {
			fmt.Printf("Cannot Load users!!")
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"user": users,
		})
	})

	// fmt.Print("Golang from docker")
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/users", userPage)
	// log.Fatal(http.ListenAndServe(":8080", nil))
	r.Run(":8080")
}

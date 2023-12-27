package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	host     = "localhost"
	port     = 5437
	user     = "host"
	password = "postgres"
	dbname   = "user"
)

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{})

	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusAccepted, gin.H{"Message": "Welcome to the Go App And Gin App"})
	})

	router.POST("/users", func(c *gin.Context) {
		createUser(c, db)
	})
	router.POST("/users/:id", func(context *gin.Context) {
		GetUserbyId(context, db)
	})

	router.GET("/users/all", func(context *gin.Context) {
		GetUsers(context, db)
	})

	router.PATCH("/users", func(context *gin.Context) {
		UpdateUser(context, db)
	})

	router.DELETE("/users", func(context *gin.Context) {
		DeleteUser(context, db)
	})

	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}

func createUser(c *gin.Context, db *gorm.DB) {
	var newUser User
	err := c.ShouldBindJSON(&newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, newUser)
	return
}

func GetUsers(c *gin.Context, db *gorm.DB) {
	var users []User
	result := db.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": result.Error})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserbyId(c *gin.Context, db *gorm.DB) {

	var user User
	id := c.Param("id")

	result := db.Find(user).Where("id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, user)
	return

}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	idStr := c.Query("id")
	name := c.Query("name")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := db.Model(&User{ID: uint(id)}).Where("id = ?", id).Update("name", name)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context, db *gorm.DB) {
	idStr := c.Query("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := db.Where("ID = ?", id).Delete(&User{})

	if result != nil {
		c.JSON(http.StatusOK, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Deleted Successfully"})
}

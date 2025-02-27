package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"

	"time"
)


var jwtKey = []byte("my_secret_key")

type User struct {
	Name string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Meals struct {
	Name string `json:"name"`
	Cost float32 `json:"cost"`
	Weight float32 `json:"weight"`
}

var meals = []Meals{
	{Name: "borscht", Cost: 300.5, Weight: 250},
	{Name: "tea", Cost: 200, Weight: 500},
}

var users = []User{
	{Name: "Иван", Password: "Ivan_1980"},
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Second)
	claim := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtKey)
}

func login(c *gin.Context) {
	for _, user := range users {
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		if user.Name == "Иван" && user.Password == "Ivan_1980" {
			token, err := generateToken(user.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(toke *jwt.Token) (interface{}, error){
			return jwtKey, nil})
		if !token.Valid || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}


func getMeals(c *gin.Context) {
	c.JSON(http.StatusOK, meals)
}

func getMealsByName(c *gin.Context) {
	name := c.Param("name")
	for _, meal := range meals {
		if meal.Name == name {
			c.JSON(http.StatusOK, meal)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Блюдо не найдено"})
}

func main() {
  r := gin.Default()
  r.POST("/login", login)
  protected := r.Group("/")
  protected.Use(authMiddleware())
  {
	protected.GET("/meals", getMeals)
  	protected.GET("/meals/:name", getMealsByName)
  }
  r.Run()
}
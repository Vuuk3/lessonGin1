package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type Meals struct {
	Name string `json:"name"`
	Cost float32 `json:"cost"`
	Weight float32 `json:"weight"`
}

var meals = []Meals{
	{Name: "borscht", Cost: 300.5, Weight: 250},
	{Name: "tea", Cost: 200, Weight: 500},
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
  r.GET("/meals", getMeals)
  r.GET("/meals/:name", getMealsByName)
  r.Run()
}
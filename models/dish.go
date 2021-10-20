package models

type Dish struct {
	Dish_id          int    `json:"dish_id"`
	Name             string `json:"name"`
	Preparation_time int    `json:"preparation_time"`
	Complexity       int    `json:"complexity"`
	Cooking_aparatus string `json:"cooking_aparatus"`
}

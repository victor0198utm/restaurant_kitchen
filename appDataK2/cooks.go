package appData

import "github.com/victor0198utm/restaurant_kitchen/models"

func GetCook(id int) models.Cook {
	cooks := []models.Cook{
		{1, 3, 3, "Gordon Ramsay", "Hey, panini head, are you listening to me?"},
		{2, 2, 2, "Julia Child", "Move, potato head!"},
		{3, 1, 1, "Rachael Ray", "Wait for me!"},
		{4, 2, 2, "Daniel Heart", "I have to moove faster!"},
		{5, 3, 3, "Maria Truman", "Give me a knife!"},
	}

	return cooks[id]
}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/victor0198utm/restaurant_kitchen/appData"
	"github.com/victor0198utm/restaurant_kitchen/models"
)

var w sync.WaitGroup

var m sync.Mutex

var receivedOrders = []models.Order{}
var dishesReady = []models.KitchenResponse{}

var activity = []int{}

var stoves = []models.Aparatus{}
var ovens = []models.Aparatus{}

func call_kitchen(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Kitchen")
	fmt.Fprintf(w, "Welcome to the Kitchen!")
}

func post_order(w http.ResponseWriter, r *http.Request) {
	// get order from request body

	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process the order in a separate thread
	fmt.Println("Received order with id: ", newOrder.Order_id, "| ", newOrder)

	go addOrder(newOrder)
}

func addOrder(newOrder models.Order) {
	m.Lock()
	receivedOrders = append(receivedOrders, newOrder)

	kr := models.KitchenResponse{
		newOrder.Order_id,
		newOrder.Table_id,
		newOrder.Waiter_id,
		[]int{},
		newOrder.Priority,
		newOrder.Max_wait,
		newOrder.Pick_up_time,
		0,
		[]models.CookingDetails{},
	}

	kr.Items = append(kr.Items, newOrder.Items...)

	dishesReady = append(dishesReady, kr)

	m.Unlock()
}

func removeOrder(s []models.Order, index int) []models.Order {
	return append(s[:index], s[index+1:]...)
}

func removeResponse(s []models.KitchenResponse, index int) []models.KitchenResponse {
	return append(s[:index], s[index+1:]...)
}

func removeDishId(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func cook(cook_id int) {
	// coordinates := []Table_Order{}
	me := appData.GetCook(cook_id)

	fmt.Println(me.Name, "started to work.\n")
	for {
		time.Sleep(10 * time.Millisecond)

		if me.Proficiency > activity[cook_id] {
			chosenOrderIdx := -1
			// highestPriorityValue := 0
			chosenDishIdx := 0
			timeWaited := 0
			m.Lock()

			// send completed orders to hall
			for dr_idx := 0; dr_idx < len(dishesReady); dr_idx++ {
				if len(dishesReady[dr_idx].Items) == len(dishesReady[dr_idx].Cooking_details) {

					return_order(dishesReady[dr_idx])

					for j := 0; j < len(receivedOrders); j++ {
						if receivedOrders[j].Order_id == dishesReady[dr_idx].Order_id {
							receivedOrders = removeOrder(receivedOrders, j)
						}
					}

					dishesReady = removeResponse(dishesReady, dr_idx)

				}
			}

			cause := ""

			// pick dish if it has waited for too long
			for j := 0; j < len(receivedOrders); j++ {
				if len(receivedOrders[j].Items) == 0 {
					continue
				}

				timeWaited = int(time.Now().Unix()) - receivedOrders[j].Pick_up_time
				allowWait := receivedOrders[j].Max_wait - int(float64(receivedOrders[j].Max_wait*100)/130)
				if timeWaited >= allowWait-2 {

					// fmt.Println("TW", timeWaited, "|AW", allowWait)
					chosenDishIdx = search_dish_to_make(receivedOrders[j], me.Rank)
					// fmt.Println("Time|", chosenDishIdx)

					if chosenDishIdx >= 0 {
						dish := appData.GetDish(receivedOrders[j].Items[chosenDishIdx] - 1)
						success := get_aparatus(dish.Cooking_aparatus)

						if success {
							chosenOrderIdx = j
							cause = "waited for too long" + strconv.Itoa(timeWaited) + "s"
							break
						}
					}
				}
			}

			// if there was no picked dish
			// pick dish which has shotest preparation time
			if chosenOrderIdx == -1 {
				smallestTimeOrderIdx := -1
				smallestTimeDishIdx := -1
				smallestTime := 120
				for j := 0; j < len(receivedOrders); j++ {
					dishes := receivedOrders[j].Items
					for dish_idx := 0; dish_idx < len(dishes); dish_idx++ {
						dish := appData.GetDish(dishes[dish_idx] - 1)
						if dish.Complexity == me.Rank || dish.Complexity == me.Rank-1 {
							if dish.Preparation_time < smallestTime {
								smallestTime = dish.Preparation_time
								smallestTimeDishIdx = dish_idx
								smallestTimeOrderIdx = j
							}
						}
					}
				}
				if smallestTimeOrderIdx > -1 {
					chosenOrderIdx = smallestTimeOrderIdx
					chosenDishIdx = smallestTimeDishIdx
					cause = "has smallest preparation time"
				}
			}

			if chosenDishIdx == -1 {
				fmt.Println("No dish in order ", receivedOrders[chosenOrderIdx].Order_id)
			} else if chosenDishIdx == -2 {
				fmt.Println("Too low or high rank. Cook:", me.Name)
			}

			if chosenOrderIdx > -1 && chosenDishIdx >= 0 {

				if chosenDishIdx < len(receivedOrders[chosenOrderIdx].Items) {

					ro := receivedOrders[chosenOrderIdx]
					fmt.Print(
						"Cook:",
						cook_id,
						" took oder:",
						ro)
					fmt.Print(
						"dish idx:",
						ro.Items[chosenDishIdx],
						". Cause:",
						cause,
						"\n")

					dishesInOrder := receivedOrders[chosenOrderIdx].Items
					dishToMake := dishesInOrder[chosenDishIdx]
					dishesInOrder = removeDishId(dishesInOrder, chosenDishIdx)
					receivedOrders[chosenOrderIdx].Items = dishesInOrder

					activity[cook_id] = activity[cook_id] + 1
					w.Add(1)
					go prepareDish(dishToMake, cook_id, receivedOrders[chosenOrderIdx].Order_id)
				}
			}

			m.Unlock()

		}
	}

	w.Done()
}

func get_aparatus(aparatus string) bool {

	if aparatus == "" {
		return true
	}

	if aparatus == "stove" {
		for i := 0; i < len(stoves); i++ {
			if stoves[i].State == 0 {
				stoves[i].State = 1
				return true
			}
		}
	}

	if aparatus == "oven" {
		for i := 0; i < len(ovens); i++ {
			if ovens[i].State == 0 {
				ovens[i].State = 1
				return true
			}
		}
	}

	return false

}

func release_aparatus(aparatus string) {

	if aparatus == "stove" {
		for i := 0; i < len(stoves); i++ {
			if stoves[i].State == 1 {
				stoves[i].State = 0
			}
		}
	}

	if aparatus == "oven" {
		for i := 0; i < len(ovens); i++ {
			if ovens[i].State == 1 {
				ovens[i].State = 0
			}
		}
	}

}

func search_dish_to_make(order models.Order, rank int) int {

	dishes := order.Items
	if len(dishes) == 0 {
		return -1 // no more dishes in order
	}

	for dish_idx := 0; dish_idx < len(dishes); dish_idx++ {
		dish := appData.GetDish(dishes[dish_idx] - 1)
		// fmt.Println("Dish:", dish, " |dc", dish.Complexity, " |r", rank)
		if dish.Complexity == rank || dish.Complexity == rank-1 {
			return dish_idx
		}
	}
	return -2 // rank too low or high
}

func search_dishe_by_time(order models.Order, rank int) int {

	smallestTimeDishIdx := -1
	smallestTime := 120

	dishes := order.Items

	for dish_idx := 0; dish_idx < len(dishes); dish_idx++ {
		dish := appData.GetDish(dishes[dish_idx] - 1)
		if dish.Complexity == rank || dish.Complexity == rank-1 {
			if dish.Preparation_time < smallestTime {
				smallestTime = dish.Preparation_time
				smallestTimeDishIdx = dish_idx
			}
		}
	}
	return smallestTimeDishIdx // rank too low or high
}

func prepareDish(dish_id int, cook_id int, order_id int) {
	fmt.Println("working on dish", dish_id)
	dish := appData.GetDish(dish_id - 1)

	time.Sleep(time.Duration(dish.Preparation_time) * time.Second)

	m.Lock()
	activity[cook_id] = activity[cook_id] - 1

	for dr_idx := 0; dr_idx < len(dishesReady); dr_idx++ {
		if dishesReady[dr_idx].Order_id == order_id {
			dishesReady[dr_idx].Cooking_details = append(dishesReady[dr_idx].Cooking_details, models.CookingDetails{dish_id, cook_id})
		}
	}
	release_aparatus(dish.Cooking_aparatus)
	m.Unlock()

	fmt.Println("Cook ", cook_id, " finished dish ", dish_id)
	w.Done()
}

func return_order(response models.KitchenResponse) {

	response.Cooking_time = int(time.Now().Unix()) - response.Pick_up_time

	json_data, err_marshall := json.Marshal(response)
	if err_marshall != nil {
		log.Fatal(err_marshall)
	}

	// resp, err := http.Post("http://"+appData.GetHallAddress()+"/distribution", "application/json",
	resp, err := http.Post("http://"+appData.GetHallAddress()+"/distribution", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Dishes sent to hall. Took %d seconds. Order id: %d. Status: %d.\n", response.Cooking_time, response, resp.StatusCode)
}

func print_resources() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("\n---Current waiting orders:", receivedOrders)
		fmt.Println("---Stoves:", stoves)
		fmt.Println("---Ovens:", ovens)
		fmt.Println("---Current waiting responses:", dishesReady, "\n")
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", call_kitchen).Methods("GET")
	myRouter.HandleFunc("/order", post_order).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+appData.GetKitchenPort(), myRouter))
}

func main() {
	// Initialize stoves
	for i := 0; i < appData.XStoves(); i++ {
		stoves = append(stoves, models.Aparatus{0})
	}

	// Initialize ovens
	for i := 0; i < appData.XOvens(); i++ {
		ovens = append(ovens, models.Aparatus{0})
	}

	// Initialize cooks
	n_cooks := appData.XCooks()
	for i := 0; i < n_cooks; i++ {
		activity = append(activity, 0)
	}

	for i := 0; i < n_cooks; i++ {
		w.Add(1)
		go cook(i)
	}

	w.Add(1)
	go print_resources()

	handleRequests()

	w.Wait()
}

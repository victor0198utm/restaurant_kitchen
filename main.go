package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Cooking_details_type struct {
	Food_id int `json:"food_id"`
	Cook_id int `json:"cook_id"`
}

type Order struct {
	Order_id     int   `json:"order_id"`
	Table_id     int   `json:"table_id"`
	Waiter_id    int   `json:"waiter_id"`
	Items        []int `json:"items"`
	Priority     int   `json:"priority"`
	Max_wait     int   `json:"max_wait"`
	Pick_up_time int   `json:"pick_up_time"`
}

type Kitchen_response struct {
	Order_id        int                    `json:"order_id"`
	Table_id        int                    `json:"table_id"`
	Waiter_id       int                    `json:"waiter_id"`
	Items           []int                  `json:"items"`
	Priority        int                    `json:"priority"`
	Max_wait        int                    `json:"max_wait"`
	Pick_up_time    int                    `json:"pick_up_time"`
	Cooking_time    int                    `json:"cooking_time"`
	Cooking_details []Cooking_details_type `json:"cooking_details"`
}

func call_kitchen(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Kitchen")
	fmt.Fprintf(w, "Welcome to the Kitchen!")
}

func post_order(w http.ResponseWriter, r *http.Request) {
	// get order from request body
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process the order in a separate thread
	fmt.Print(time.Now().Clock())
	fmt.Print(": Started to process the order with id: ")
	fmt.Println(order.Order_id)
	go process_order(order)

}

func process_order(order Order) {
	// processing logic

	// sleep for 3-10 seconds
	preparation_time := rand.Intn(7) + 3
	time.Sleep(time.Duration(preparation_time) * time.Second)

	// finished
	make_dishes(order, preparation_time)
}

func make_dishes(order Order, preparation_time int) {
	items := []int{3, 4, 4, 2}
	cooking_details := []Cooking_details_type{
		{3, 1}, {4, 1}, {4, 2}, {2, 3},
	}
	the_finished_order := Kitchen_response{order.Order_id, 1, 1, items, 3, 45, 1631453140, 65, cooking_details}
	json_data, err_marshall := json.Marshal(the_finished_order)
	if err_marshall != nil {
		log.Fatal(err_marshall)
	}

	resp, err := http.Post("http://hall:8002/distribution", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(time.Now().Clock())
	fmt.Printf(": Dishes sent to hall. Took %d seconds. Order id: %d. Status: %d.\n", preparation_time, order.Order_id, resp.StatusCode)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", call_kitchen).Methods("GET")
	myRouter.HandleFunc("/order", post_order).Methods("POST")
	log.Fatal(http.ListenAndServe(":8001", myRouter))
}

func main() {
	rand.Seed(1583)
	handleRequests()
}

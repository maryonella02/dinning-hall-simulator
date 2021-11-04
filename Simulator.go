package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"
)

type Dishes struct {
	Dishes []Dish `json:"dishes"`
}

type Dish struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	PreparationTime  int    `json:"preparation-time"`
	Complexity       int    `json:"complexity"`
	CookingApparatus string `json:"cooking-apparatus"`
}
type Order struct {
	ID         int   `json:"id"`
	Items      []int `json:"items"`
	TableID    int   `json:"table_id"`
	WaiterID   int   `json:"waiter_id"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max_wait"`
	PickUpTime int64 `json:"pick_up_time"`
}
type Table struct {
	ID    int
	State State
}
type State int

const (
	Free      = 1
	WaitMenu  = 2
	WaitOrder = 3
)

type Waiter struct {
	ID int
}

const NumberOfTables = 4
const NumberOfWaiters = 2

type Orders struct {
	sync.RWMutex
	Orders []Order
}

var OrdersMap = Orders{Orders: make([]Order, 0)}

var ratingSystem = NewRating()

//var orderList *[]Orders
func main() {

	jsonFile, err := os.Open("menu.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened menu.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var dishes Dishes
	json.Unmarshal(byteValue, &dishes)

	/*for i := 0; i < len(dishes.Dishes); i++ {
		fmt.Println("ID: " + strconv.Itoa(dishes.Dishes[i].ID))
		fmt.Println("Name: " + dishes.Dishes[i].Name)
		fmt.Println("PreparationTime: " + strconv.Itoa(dishes.Dishes[i].PreparationTime))
		fmt.Println("Complexity: " + strconv.Itoa(dishes.Dishes[i].Complexity))
		fmt.Println("CookingApparatus: " + dishes.Dishes[i].CookingApparatus)
	}*/

	tables := generateTables(NumberOfTables)
	serveTables(NumberOfWaiters, tables, dishes)
	//if all tables are free, generate more busy tables

	http.HandleFunc("/distribution", HandleRequest)
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func serveTables(nrOfWaiters int, tables []Table, dishes Dishes) {
	for i := 0; i < nrOfWaiters; i++ {
		waiter := &Waiter{ID: i + 1}
		go getOrder(waiter, tables, dishes)
		for len(OrdersMap.Orders) > 0 {
			go giveOrder(waiter, tables)
		}
	}
}

func giveOrder(waiter *Waiter, tables []Table) {
	for i := 0; i < len(tables); i++ {
		if tables[i].State == 3 {
			if IsServed(waiter.ID, tables[i].ID) {
				tables[i].State = 1
			}
		}
	}
}

func IsServed(waiterID int, tableID int) bool {
	for i := 0; i < len(OrdersMap.Orders); i++ {
		if waiterID == OrdersMap.Orders[i].WaiterID && tableID == OrdersMap.Orders[i].TableID {
			DeleteOrderFromList(i)
			rating := CalculateRating(OrdersMap.Orders[i])
			fmt.Printf("%s = %d\n", "Order rating", rating)
			ratingSystem.AddRating(rating)
			fmt.Printf("%s = %f\n", "Rating overall", ratingSystem.ReturnRating())
			return true
		}
	}
	return false
}

func getOrder(waiter *Waiter, tables []Table, dishes Dishes) {
	var rwMutex = &sync.RWMutex{}
	for i := 0; i < len(tables); i++ {
		if tables[i].State == 2 {
			rwMutex.Lock()
			order := createOrder(dishes, tables[i].ID, waiter.ID)
			makeRequest(order)
			tables[i].State = 3
			rwMutex.Unlock()
		}
	}
}
func generateTables(nrOfTables int) []Table {
	tables := make([]Table, nrOfTables)
	for i := 0; i < nrOfTables; i++ {
		state := generateRandomNumber(1, 2)
		table := &Table{ID: i + 1,
			State: State(state)}
		tables[i] = *table
	}
	return tables
}

func createOrder(dishes Dishes, tableId int, waiterId int) []byte {
	order := &Order{ID: generateRandomNumber(1, 100000),
		TableID:    tableId,
		WaiterID:   waiterId,
		Items:      generateItems(),
		Priority:   generatePriority(),
		MaxWait:    getMaxWaitTime(dishes),
		PickUpTime: getUnixTimestamp()}
	b, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return b

}

func makeRequest(b []byte) {
	url := "http://localhost:8082/order"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("skip the error")
	} else {
		defer resp.Body.Close()
		fmt.Println(string(b))
		fmt.Println("Request sent")
	}

}

func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}
	fmt.Println("Request Handled")
	log.Println(order)
	AddOrderToList(order)

}

func AddOrderToList(order Order) {
	OrdersMap.Lock()
	defer OrdersMap.Unlock()
	OrdersMap.Orders = append(OrdersMap.Orders, order)
}

func DeleteOrderFromList(index int) {
	OrdersMap.Lock()
	defer OrdersMap.Unlock()
	OrdersMap.Orders = append(OrdersMap.Orders[:index], OrdersMap.Orders[index+1:]...)
}

func generateRandomNumber(min, max int64) int {

	bg := big.NewInt(max + 1 - min)
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}
	return int(n.Int64() + min)
}

func generateItems() []int {
	nrOfItems := generateRandomNumber(1, 10)
	var items = make([]int, nrOfItems)
	for i := 0; i < nrOfItems; i++ {
		items[i] = generateRandomNumber(1, 10)
	}
	return items
}

func generatePriority() int {
	priority := generateRandomNumber(1, 5)
	return priority
}

func getMaxWaitTime(dishes Dishes) int {
	var max int
	for i := 0; i < len(dishes.Dishes); i++ {
		if dishes.Dishes[i].PreparationTime > max {
			max = dishes.Dishes[i].PreparationTime
		}
	}
	var maxTime = float32(max) * 1.3
	return int(maxTime)
}

func getUnixTimestamp() int64 {
	now := time.Now()
	sec := now.Unix()
	return sec
}

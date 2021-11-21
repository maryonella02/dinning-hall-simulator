package services

import (
	"bytes"
	"crypto/rand"
	"dinning-hall/models"
	"dinning-hall/utils"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

// func generateTables(nrOfTables int) []models.Table {
// 	tables := make([]models.Table, nrOfTables)
// 	for i := 0; i < nrOfTables; i++ {
// 		state := generateRandomNumber(1, 2)
// 		table := &models.Table{ID: i + 1,
// 			State: models.TableState(state)}
// 		tables[i] = *table
// 	}
// 	return tables
// }

var OrdersMap = models.Orders{Orders: make([]models.Order, 0)}

func GenerateInitialTables(nrOfTables int) models.Tables {
	tablesSlice := make([]*models.Table, nrOfTables)
	for i := 0; i < nrOfTables; i++ {
		// state := generateRandomNumber(1, 2)
		tablesSlice[i] = &models.Table{ID: i + 1,
			State: models.TableState(1)}
	}

	allTables := models.Tables{AllTables: tablesSlice}
	fmt.Println("Generetad ", 4, "Tables")
	return allTables
}

func ServeTables(nrOfWaiters int, tables *models.Tables, dishes models.Dishes) {
	for i := 0; i < nrOfWaiters; i++ {
		waiter := &models.Waiter{ID: i + 1}
		go getOrder(waiter, tables, dishes)
	}
}

func getOrder(waiter *models.Waiter, tables *models.Tables, dishes models.Dishes) {
	fmt.Println("Serving as ", waiter.ID, "waiter")
	for {
		for idx, table := range (*tables).AllTables {
			table.Lock()
			if table.State == 2 {
				order := createOrder(dishes, table.ID, waiter.ID)
				fmt.Println("Order created ", order)
				requested := makeRequest(order)
				if requested {
					table.State = 3
				}
			}
			if table.State == 3 {
				fmt.Println("Verify if served")
				if IsServed(waiter.ID, (*tables).AllTables[idx].ID) {
					fmt.Println("Served for ", table.ID)
					table.State = 1
				}
			}
			table.Unlock()
			time.Sleep(1 * time.Second)
		}

	}

}

func IsServed(waiterID int, tableID int) bool {
	OrdersMap.Lock()
	defer OrdersMap.Unlock()
	for i := 0; i < len(OrdersMap.Orders); i++ {
		if waiterID == OrdersMap.Orders[i].WaiterID && tableID == OrdersMap.Orders[i].TableID {
			rating := CalculateRating(OrdersMap.Orders[i])
			DeleteOrderFromList(i)
			fmt.Printf("%s = %d\n", "Order rating", rating)
			ratingSystem.AddRating(rating)
			fmt.Printf("%s = %f\n", "Rating overall", ratingSystem.ReturnRating())
			return true
		}
	}

	return false
}
func AddOrderToList(order models.Order) {
	OrdersMap.Lock()
	defer OrdersMap.Unlock()
	OrdersMap.Orders = append(OrdersMap.Orders, order)
}

func DeleteOrderFromList(index int) {
	OrdersMap.Orders = append(OrdersMap.Orders[:index], OrdersMap.Orders[index+1:]...)
}

func createOrder(dishes models.Dishes, tableId int, waiterId int) []byte {
	order := &models.Order{ID: generateRandomNumber(1, 100000),
		TableID:    tableId,
		WaiterID:   waiterId,
		Items:      generateItems(),
		Priority:   generatePriority(),
		MaxWait:    getMaxWaitTime(dishes),
		PickUpTime: utils.GetUnixTimestamp()}
	b, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return b

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

func getMaxWaitTime(dishes models.Dishes) int {
	var max int
	for i := 0; i < len(dishes.Dishes); i++ {
		if dishes.Dishes[i].PreparationTime > max {
			max = dishes.Dishes[i].PreparationTime
		}
	}
	var maxTime = float32(max) * 1.3
	return int(maxTime)
}

func makeRequest(b []byte) bool {
	url := "http://localhost:8082/order"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("skip the error")
		return false
	} else {
		defer resp.Body.Close()
		fmt.Println(string(b))
		fmt.Println("Request sent")
		return true
	}

}

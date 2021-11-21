package main

import (
	"dinning-hall/models"
	"dinning-hall/services"
	"dinning-hall/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const NumberOfTables = 4
const NumberOfWaiters = 2

func main() {
	dishes := utils.ReadMenus("menu.json")

	tables := services.GenerateInitialTables(NumberOfTables)
	go services.OccupyTables(&tables)

	services.ServeTables(NumberOfWaiters, &tables, dishes)
	//if all tables are free, generate more busy tables

	http.HandleFunc("/distribution", HandleRequest)
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func HandleRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order models.Order
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}
	fmt.Println("Request Handled")
	log.Println(order)
	services.AddOrderToList(order)

}

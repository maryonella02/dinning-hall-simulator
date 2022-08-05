package main

import (
	"dinning-hall/models"
	"dinning-hall/services"
	"dinning-hall/utils"
	"encoding/json"
	"log"
	"net/http"
)

const NumberOfTables = 10
const NumberOfWaiters = 4

func main() {
	dishes := utils.ReadMenus("menu.json")

	tables := services.GenerateInitialTables(NumberOfTables)
	go services.OccupyTables(&tables)

	services.ServeTables(NumberOfWaiters, &tables, dishes)

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
	log.Println("Request Handled")
	log.Println(order)
	services.AddOrderToList(order)

}

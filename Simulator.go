package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max_wait"`
	PickUpTime int64 `json:"pick_up_time"`
}

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
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			worker(dishes)
		}()
	}
	wg.Wait()

	http.HandleFunc("/distribution", test)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
func worker(dishes Dishes) {
	order := createOrder(dishes)
	time.Sleep(time.Second)
	makeRequest(order)

}

func createOrder(dishes Dishes) []byte {
	order := &Order{ID: generateRandomNumber(1, 100000),
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
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(b))
	fmt.Println("Request sent")
}

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil {
		panic(err)
	}
	fmt.Println("Request Handled")
	log.Println(order)
}

func generateRandomNumber(min int, max int) int {
	digit := 0

	rand.Seed(time.Now().UnixNano())
	digit = min + rand.Intn(max-min)
	return digit
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

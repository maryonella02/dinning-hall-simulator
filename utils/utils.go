package utils

import (
	"dinning-hall/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func ReadMenus(filename string) models.Dishes {

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Successfully Opened %s \n", filename)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var dishes models.Dishes

	err = json.Unmarshal(byteValue, &dishes)
	if err != nil {
		log.Panic("Cannot unmarshal json menu")
	}
	return dishes
}

func GetUnixTimestamp() int64 {
	now := time.Now()
	sec := now.Unix()
	return sec
}

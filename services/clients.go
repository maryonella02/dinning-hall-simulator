package services

import (
	"dinning-hall/models"
	"log"
	"time"
)

func OccupyTables(tables *models.Tables) {
	for {
		for _, table := range (*tables).AllTables {
			table.Lock()
			if table.State == 1 {
				table.State = 2
				log.Println("As a client I occupied table nr ", table.ID)
			}
			table.Unlock()
			time.Sleep(1 * time.Second)
		}
	}
}

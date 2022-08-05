package services

import (
	"dinning-hall/models"
	"log"
	"time"
)

type RatingSystem struct {
	values []int
}

const Multiplier = 1.3

var TimeUnit time.Duration
var ratingSystem = NewRating()

func NewRating() *RatingSystem {
	return &RatingSystem{
		values: make([]int, 0),
	}
}

func (s *RatingSystem) AddRating(value int) {
	s.values = append(s.values, value)
}

func (s *RatingSystem) ReturnRating() float32 {
	var result float32 = 0
	for _, value := range s.values {
		result += float32(value)
	}

	return result / float32(len(s.values))
}

func CalculateRating(order models.Order) int {
	totalOrderPreparingTime := float64(time.Now().Unix() - order.PickUpTime)

	maxWait := float64(order.MaxWait) * Multiplier

	log.Printf("Total Order Preparing Time : %f : Max Wait : %f\n", totalOrderPreparingTime, maxWait)

	if totalOrderPreparingTime < maxWait {
		return 5
	}
	if totalOrderPreparingTime < maxWait*1.1 {
		return 4
	}
	if totalOrderPreparingTime < maxWait*1.2 {
		return 3
	}
	if totalOrderPreparingTime < maxWait*1.3 {
		return 2
	}
	if totalOrderPreparingTime < maxWait*1.4 {
		return 1
	}

	return 0
}

package main

import (
	"fmt"
	"time"
)

var TimeUnit time.Duration

type RatingSystem struct {
	values []int
}

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

func CalculateRating(order Order) int {
	var totalOrderPreparingTime = float32(order.PickUpTime - getUnixTimestamp())
	multiplier := float32(TimeUnit) / float32(time.Second)
	maxWait := float32(order.MaxWait) * multiplier

	fmt.Printf("Total Order Preparing Time : %f : Max Wait : %f\n", totalOrderPreparingTime, maxWait)

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

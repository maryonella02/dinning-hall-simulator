package models

import "sync"

const (
	Free      = 1
	WaitMenu  = 2
	WaitOrder = 3
)

type Orders struct {
	sync.RWMutex
	Orders []Order
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
	sync.RWMutex
	ID    int
	State TableState
}

type Tables struct {
	AllTables []*Table
}

type TableState int

type Waiter struct {
	ID int
}

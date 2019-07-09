package services

import (
	"fmt"
	"time"
)

// Cleaner rear func
func Cleaner() {
	time.Sleep(time.Second * 2)
	fmt.Println("cleaner")
}

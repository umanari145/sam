package main

import (
	"fmt"
	"testing"
)

func TestHandler(t *testing.T) {
	area, err := loadAreaFromZip("2740077")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(area)
	}
}

package main

import (
	"fmt"
	"testing"
)

func TestHandler(t *testing.T) {
	area := loadAreaFromZip("27400773")
	fmt.Println(area)
}

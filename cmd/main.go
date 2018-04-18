package main

import (
	"fmt"

	"github.com/radioinmyhead/ipmi"
)

func main() {
	bmc, err := ipmi.GetLocalIPMI()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bmc.GetFanSpeed())
}

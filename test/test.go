package main

import (
	"fmt"

	"github.com/ramin0x53/telegram_alertbot/src/price"
)

func main() {

	i := price.GetPrice("linkbtc")

	if i == "" {
		fmt.Println("Wrong!!")
	}
	fmt.Println(i)
}

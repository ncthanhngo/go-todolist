package common

import "fmt"

func Recovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered-", r)
	}
}

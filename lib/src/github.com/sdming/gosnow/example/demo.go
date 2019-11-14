/*
github.com/twitter/snowflake in golang
*/

package main

import (
	"fmt"
	"github.com/sdming/gosnow"
)

func main() {

	v := gosnow.Default()
	//v := gosnow.NewSnowFlake(100)
	for i := 0; i < 10; i++ {
		id := v.Next()
		fmt.Println(id)
	}
}

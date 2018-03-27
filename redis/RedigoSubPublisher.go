package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {
	c, _ := redis.Dial("tcp", "[ip]:6379")
	//checkError(err)
	defer c.Close()

	for {
		var s string
		fmt.Scanln(&s)
		r, err := c.Do("PUBLISH", "newChatRoom", s)
		if err != nil {
			fmt.Println("Error when publish:", err)
		}
		fmt.Println("Respnse when publish:", r)
	}
}

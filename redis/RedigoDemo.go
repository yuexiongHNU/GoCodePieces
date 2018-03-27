package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"log"
	"time"
)

func main() {
	// New connection
	c, err := redis.Dial("tcp", "[ip]:6379")
	checkError(err)
	defer c.Close()

	// SET KEY
	v, err := c.Do("SET", "name", "moon")
	checkError(err)
	fmt.Println("SET KEY name:", v)

	// GET KEY
	v, err =  redis.String(c.Do("GET", "name"))
	checkError(err)
	fmt.Println("GET KEY name:", v)

	// SET LIST
	dataList := []interface{}{"xiong", "yue", 123141}
	for _, v := range dataList {
		fmt.Println("Value when set to list:", v)
		c.Do("lpush", "testlist", v)
		checkError(err)
		fmt.Println("SET LIST testlist:", v)
	}

	// GET LIST
	values, _ := redis.Values(c.Do("lrange", "testlist", 0, 100))
	fmt.Println("GET LIST VALUES:", values)
	/* GET LIST METHOD 1
	for _, v2 := range values {
		fmt.Println(string(v2.([]byte)))
	}*/
	// GET LIST METHOD 2
	var v2,v3,v4,v5,v6,v7 string
	// var vn []interface{}
	redis.Scan(values, &v2, &v3, &v4, &v5, &v6, &v7)
	fmt.Printf("GET FROM LIST: %s %s %s %s %s %s \n", v2, v3, v4, v5, v6, v7)

	// PIPELINE
	c.Send("SET", "age", 16)
	c.Send("SET", "sex", "female")
	c.Send("GET", "age")
	c.Send("GET", "sex")
	c.Flush()
	r1, err := c.Receive()
	checkError(err)
	r2, err := c.Receive()
	checkError(err)
	r3, err := c.Receive()
	checkError(err)
	r4, err := c.Receive()
	checkError(err)
	fmt.Printf("PIPELINE Received MESSAGE: %s %s %s %s \n", r1, r2, r3, r4)

	// SUBSCRIBE
	go subscribe("newChatRoom")
	go subscribe("newChatRoom")
	go subscribe("newChatRoom")
	go subscribe("newChatRoom")

	for {
		time.Sleep(1 * time.Second)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		return
	}
}

func subscribe(channel string)  {
	// New connection
	c, err := redis.Dial("tcp", "[ip]:6379")
	checkError(err)
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe(channel)
	for {
		time.Sleep(1 * time.Second)
		switch recv := psc.Receive().(type) {
		case redis.Message:
			if recv.Data == nil {
			} else {
				fmt.Println("recv.Data:", recv.Data)
				fmt.Printf("%s: message: %s\n", recv.Channel, recv.Data)
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d \n", recv.Channel, recv.Kind, recv.Count)
		case error:
			fmt.Println(recv)
			return
		}
	}
}

package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

func main() {

	re := redis.Pool{MaxIdle: 10, MaxActive: 10, IdleTimeout: 300 * time.Second, Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "127.0.0.1:6379")
	}}

	conn := re.Get()
	defer conn.Close()

	_,err:=conn.Do("LPUSH","blackiplist","99.0.0.0")
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("success")

	//reply, err := conn.Do("BLPOP", "blackiplist", 100)
	//str,err:=redis.Strings(reply, err)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(str)
}

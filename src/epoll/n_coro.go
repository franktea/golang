package main 

/*
出自：http://my.oschina.net/yunfound/blog/141222?fromerr=CmIFKSjD
其实这种实现方式仍然是每个连接一个coroutine，只不过是限制了coroutine的数量，
根本达不到epoll地的效果
*/

import (
	"fmt"
	"net"
	//"os"
	//"time"
)

const (
	MAX_CONN_NUM = 4
)

//echo server Goroutine
func EchoFunc(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	for {
		ll, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		fmt.Println("recv ll:", ll)
		// send reply
		send_buf := buff[:ll]
		ll, err = conn.Write(send_buf) // 不能够直接用buff，否则会发送整个1024字节的长度
		if err != nil {
			fmt.Println("Error sending: ", err.Error())
			return
		}
		fmt.Println("send ll:", ll)
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		return
	}
	
	defer listener.Close()
	
	fmt.Println("listening...")
	
	var cur_conn_num int = 0
	conn_chan := make(chan net.Conn)
	ch_conn_change := make(chan int)
	
	go func() {
		for conn_change := range ch_conn_change {
			cur_conn_num += conn_change
		}
	}()
	
//	go func() {
//		for _ = range time.Tick(1e8) {
//			fmt.Printf("current conn num: %f\n", cur_conn_num)
//		}
//	}()
	
	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range conn_chan {
				ch_conn_change <- 1
				EchoFunc(conn)
				ch_conn_change <- -1
			}
		}()
	}
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accept:", err.Error())
			return
		}
		conn_chan <- conn
	}
}


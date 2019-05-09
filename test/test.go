package main

import (
	"bytes"
	"cts2/socket"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"time"
)

func main() {
	req := socket.HTTP_REQ{}
	req.CRT_TIME = time.Now().Format("2006-01-02 15:04:05")
	req.METHOD = "GET"
	req.MSG_ID = 10000
	req.PROJ_ID = 10000
	req.URL = "http://www.baidu.com"
	headerItem := []socket.ITEM{}
	headerItem = append(headerItem, socket.ITEM{KEY: "header_key1", VALUE: "header_value1"}, socket.ITEM{KEY: "header_key2", VALUE: "header_value2"})
	req.HEADERS = headerItem
	paramItem := []socket.ITEM{}
	paramItem = append(paramItem, socket.ITEM{KEY: "param_key1", VALUE: "param_value1"}, socket.ITEM{KEY: "param_key2", VALUE: "param_value2"})
	req.PARAMS = paramItem

	reqbytes, err := xml.Marshal(&req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(reqbytes))
	}

	var a string = "9999999999999999"
	buf := bytes.NewBufferString(a)

	fmt.Println(len(buf.Bytes()))
	fmt.Printf("% x", buf.Bytes())
	fmt.Println("-------")
	fmt.Printf("Byte 2-->%c\n", byte(0x02))
	fmt.Printf("Byte 1-->%x\n ", byte(49))
	var one rune = '1'
	i1 := int(one)
	fmt.Println("'1' convert to", i1)
	fmt.Printf("%", 1)

	x := int32(8888)
	//x := "a"
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	byteB := bytes.NewBuffer(bytesBuffer.Bytes())
	var y int32
	binary.Read(byteB, binary.BigEndian, &y)
	fmt.Println("*********")
	fmt.Println(y)

	//runeI, n, e1 := bytesBuffer.ReadRune()

	//	fmt.Printf("%c,%v,%v,%v", string(runeI), runeI, n, e1)
	config := socket.Config{}
	config.Addr = "101.200.138.213"
	config.Port = "9123"
	config.MaxRetry = "5"
	config.Timeout = "60"
	socket.SocketConfig = config
	socket.Start()
	for {

		select {

		case <-time.After(2 * time.Second):
			fmt.Println("HTTP_QUEUE_REQ IN")
			fmt.Println("REQ==>\n", string(reqbytes))
			socket.SendREQ(reqbytes)

		}
	}

}

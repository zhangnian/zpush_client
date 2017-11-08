package main

import (
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"io"
	"log"
	"net"
	msg "zpush_client/message"
	"time"
)

func mockLogin(username, password string) {
	loginReq := msg.LoginReq{
		Username: username,
		Password: password,
	}

	b, err := proto.Marshal(&loginReq)
	if err != nil {
		log.Fatalf("marshal client msg error: %s\n", err.Error())
	}

	buf := make([]byte, 2+4+len(b))

	n := 0
	binary.BigEndian.PutUint16(buf[n:], 1)
	n += 2

	binary.BigEndian.PutUint32(buf[n:], uint32(len(b)))
	n += 4

	copy(buf[n:], b)

	conn, err := net.Dial("tcp", "192.168.123.199:10001")
	if err != nil {
		log.Fatalln("connect to gateway error")
	}
	defer conn.Close()

	r, err := conn.Write(buf)
	if err != nil {
		log.Fatalln("write error")
	}

	log.Printf("send login message to gateway, size: %d\n", r)

	packetHeaderBuf := make([]byte, 6)
	_, err = io.ReadFull(conn, packetHeaderBuf)
	if err != nil {
		if err == io.EOF {
			log.Println("connection closed by peer")
			return
		}

		log.Printf("io.ReadFull error: %s\n", err.Error())
		return
	}

	log.Println(packetHeaderBuf)

	packetBodyLen := binary.BigEndian.Uint32(packetHeaderBuf[2:6])

	packetBodyBuf := make([]byte, packetBodyLen)
	_, err = io.ReadFull(conn, packetBodyBuf)
	if err != nil {
		if err == io.EOF {
			log.Println("connection closed by peer")
			return
		}

		log.Printf("io.ReadFull error: %s\n", err.Error())
		return
	}

	log.Println(packetBodyBuf)

	var loginResp msg.LoginResp
	err = proto.Unmarshal(packetBodyBuf, &loginResp)
	if err != nil{
		log.Println("unmarshal error", err.Error())
		return
	}

	log.Println(loginResp)

	for{
		tmpBuf := make([]byte, 1024)
		conn.Read(tmpBuf)
		log.Println(tmpBuf)
	}
}

func main() {
	for i:=0; i<10; i++{
		go mockLogin("zhangnian", "123456")
	}

	time.Sleep(time.Second * 1000)
}

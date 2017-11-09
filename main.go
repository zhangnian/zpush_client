package main

import (
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"io"
	"log"
	"net"
	"time"
	msg "zpush_client/message"
)

func buildMsg(cmd int, b []byte) []byte{
	buf := make([]byte, 2+4+len(b))

	n := 0
	binary.BigEndian.PutUint16(buf[n:], uint16(cmd))
	n += 2

	binary.BigEndian.PutUint32(buf[n:], uint32(len(b)))
	n += 4

	copy(buf[n:], b)
	return buf
}

func mockLogin(username, password string) {
	loginReq := msg.LoginReq{
		Username: username,
		Password: password,
	}

	b, err := proto.Marshal(&loginReq)
	if err != nil {
		log.Fatalf("marshal client msg error: %s\n", err.Error())
	}

	buf := buildMsg(1, b)

	conn, err := net.Dial("tcp", "localhost:1001")
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
	var loginResp msg.LoginResp
	err = proto.Unmarshal(packetBodyBuf, &loginResp)
	if err != nil {
		log.Println("unmarshal error", err.Error())
		return
	}

	for {
		mockHeartBeat(conn, loginResp.Userid)
		time.Sleep(time.Second * 3)
	}
}

func mockHeartBeat(conn net.Conn, userid int32) {
	hbReq := msg.HBReq{
		Userid: userid,
	}

	b, _ := proto.Marshal(&hbReq)
	buf := buildMsg(9999, b)

	_, err := conn.Write(buf)
	if err != nil {
		log.Fatalln("write error")
	}

	//log.Printf("send heartbeat message to gateway, size: %d\n", r)

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

	hbResp := msg.HBResp{}
	err = proto.Unmarshal(packetBodyBuf, &hbResp)
	if err != nil {
		log.Println("unmarshal error", err.Error())
		return
	}
}

func main() {
	for i := 0; i < 10; i++ {
		go mockLogin("zhangnian", "123456")
	}

	time.Sleep(time.Second * 1000)
}

package main

import (
	"flag"
	"github.com/Allenxuxu/gev/example/protocol_gd/common"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"

	"github.com/Allenxuxu/gev"
)

type example struct{}

func (s *example) OnConnect(c *gev.Connection) {
	log.Println(" OnConnect ： ", c.PeerAddr())
}
func (s *example) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	log.Println("OnMessage：", data)

	msgReq := common.SendMsg{}
	err := proto.Unmarshal(data, &msgReq)
	if err != nil {
		log.Panic(" proto.Unmarshal err:", err)
		return
	}
	log.Println("msgReq:", msgReq)
	log.Println("msgReq.Cmd:", msgReq.Cmd)

	switch msgReq.Cmd {
	case common.CMD_CMD_HEART_BEAT_REQ:
		log.Println("- msgReq.Cmd@CMD_CMD_HEART_BEAT")
	case common.CMD_CMD_UPLOAD_COOKIE_REQ:
		log.Println("- msgReq.Cmd@CMD_CMD_UPLOAD_COOKIE")
		body := common.UploadCookie{}
		err := proto.Unmarshal(msgReq.Body, &body)
		if err != nil {
			log.Panic(" proto.Unmarshal(msgReq.Body, &body) err:", err)
			return
		}
		log.Println("body:", body)
	default:
		log.Println("- msgReq.Cmd@default")
	}

	// 响应；注释则不响应
	//out = data
	//out = []byte("aaaaa123")

	return
}

func (s *example) OnClose(c *gev.Connection) {
	log.Println("OnClose")
}

func main() {
	handler := new(example)
	var port int
	var loops int

	flag.IntVar(&port, "port", 1834, "server port")
	flag.IntVar(&loops, "loops", -1, "num loops")
	flag.Parse()

	s, err := gev.NewServer(handler,
		gev.Network("tcp"),
		gev.Address(":"+strconv.Itoa(port)),
		gev.NumLoops(loops),
		gev.CustomProtocol(&ExampleProtocol{}))
	if err != nil {
		panic(err)
	}

	log.Println("server start")
	s.Start()
}

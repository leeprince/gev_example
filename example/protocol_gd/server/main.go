package main

import (
	"flag"
	"github.com/Allenxuxu/gev/example/protocol_gd/common"
	"gitlab.yewifi.com/golden-cloud/common/gclog"
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

	msgReq := &common.SendMsg{}
	err := proto.Unmarshal(data, msgReq)
	if err != nil {
		log.Panic(" proto.Unmarshal err:", err)
		return
	}
	log.Println("msgReq.Cmd:", msgReq.Cmd)

	log.Println("msgReq:", msgReq)
	gclog.Info("gclog.Info msgReq:", msgReq)

	log.Println("msgReq.Body:", msgReq.Body)
	log.Println("msgReq.Body string:", string(msgReq.Body))
	gclog.Info("gclog.Info msgReq.Body:", msgReq.Body)

	gclog.WithField("msgReq.Body", msgReq.Body).WithField("msgReq.Body string", string(msgReq.Body)).Info("gclog.Info msgReq.Body info")

	var respCmd common.CMD
	switch msgReq.Cmd {
	case common.CMD_CMD_HEART_BEAT_REQ:
		log.Println("- msgReq.Cmd@CMD_CMD_HEART_BEAT")
		respCmd = common.CMD_CMD_HEART_BEAT_RSP

	case common.CMD_CMD_UPLOAD_LOGIN_REQ:
		log.Println("- msgReq.Cmd@CMD_CMD_UPLOAD_LOGIN_REQ")
		respCmd = common.CMD_CMD_UPLOAD_LOGIN_RSP

	case common.CMD_CMD_UPLOAD_COOKIE_REQ:
		log.Println("- msgReq.Cmd@CMD_CMD_UPLOAD_COOKIE")
		body := common.UploadCookie{}
		err := proto.Unmarshal(msgReq.Body, &body)
		if err != nil {
			log.Panic(" proto.Unmarshal(msgReq.Body, &body) err:", err)
			return
		}
		log.Println("body:", &body)
		gclog.Info("gclog.Info body:", &body)
		respCmd = common.CMD_CMD_UPLOAD_COOKIE_RSP

	default:
		log.Println("- msgReq.Cmd@default")
	}

	// 响应；注释则不响应
	msg := msgReq
	msg.Cmd = respCmd
	log.Println("msg:", msg)
	gclog.Info("gclog.Info msg:", msg)
	out, err = proto.Marshal(msg)
	if err != nil {
		log.Panic(" proto.Marshal(msg) err:", err)
	}

	return
}

func (s *example) OnClose(c *gev.Connection) {
	log.Println("OnClose")
}

func init() {
	logger := gclog.Default(
		"tools_test",
		".",
		5,
		"",
	)
	gclog.SetDefaultLogger(logger)
}
func main() {
	handler := new(example)
	var port int
	var loops int

	flag.IntVar(&port, "port", 18000, "server port")
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

	log.Println("server start, port:", port)
	s.Start()
}

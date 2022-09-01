package main

import (
	"encoding/binary"
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
	"github.com/spf13/cast"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Server example
type Server struct {
	mu     sync.RWMutex
	server *gev.Server
}

var clienConnList []connectData

type connectData struct {
	conn *gev.Connection
	data string
}

const exampleHeaderLen = 4

type ExampleProtocol struct{}

func (d *ExampleProtocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if buffer.VirtualLength() > exampleHeaderLen {
		buf := pbytes.GetLen(exampleHeaderLen)
		defer pbytes.Put(buf)
		_, _ = buffer.VirtualRead(buf)
		dataLen := binary.BigEndian.Uint32(buf)

		if buffer.VirtualLength() >= int(dataLen) {
			ret := make([]byte, dataLen)
			_, _ = buffer.VirtualRead(ret)

			buffer.VirtualFlush()
			return nil, ret
		} else {
			buffer.VirtualRevert()
		}
	}
	return nil, nil
}

func (d *ExampleProtocol) Packet(c *gev.Connection, data interface{}) []byte {
	dd := data.([]byte)
	dataLen := len(dd)
	ret := make([]byte, exampleHeaderLen+dataLen)
	binary.BigEndian.PutUint32(ret, uint32(dataLen))
	copy(ret[4:], dd)
	return ret
}

// New server
func New(ip, port string) (*Server, error) {
	var err error
	s := new(Server)
	s.server, err = gev.NewServer(s,
		gev.Address(ip+":"+port),
		gev.CustomProtocol(&ExampleProtocol{}))
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Start server
func (s *Server) Start() {
	s.server.Start()
}

// Stop server
func (s *Server) Stop() {
	log.Println(" defer Stop")

	s.server.Stop()
}

// OnConnect callback
func (s *Server) OnConnect(c *gev.Connection) {
	log.Println(" OnConnect ： ", c.PeerAddr())

	clienConnList = append(clienConnList, connectData{
		conn: c,
		data: c.PeerAddr(),
	})

	go func(c *gev.Connection) {
		for {
			randTime := cast.ToDuration(rand.Uint64()%10) + 1
			log.Println("time.Sleep:", randTime)
			time.Sleep(time.Second * randTime)

			err := c.Send([]byte(c.PeerAddr()))
			if err != nil {
				log.Println("c.Send err：", err.Error())
				continue
			}
			log.Println("c.Send success")
		}
	}(c)
}

// OnMessage callback
func (s *Server) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	log.Println("OnMessage", string(data))
	out = data
	return
}

// OnClose callback
func (s *Server) OnClose(c *gev.Connection) {
	log.Println("OnClose")
}

func main() {
	s, err := New("", "1833")
	if err != nil {
		panic(err)
	}
	defer s.Stop()

	s.Start()
}

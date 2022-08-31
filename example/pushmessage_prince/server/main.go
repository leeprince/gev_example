package main

import (
	"container/list"
	"encoding/binary"
	"github.com/gobwas/pool/pbytes"
	"log"
	"sync"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
)

const clientsKey = "demo_push_message_key"

// Server example
type Server struct {
	conn   *list.List
	mu     sync.RWMutex
	server *gev.Server
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
	s.conn = list.New()
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
	s.server.RunEvery(1*time.Second, s.RunPush)
	s.server.Start()
}

// Stop server
func (s *Server) Stop() {
	s.server.Stop()
}

// RunPush push message
func (s *Server) RunPush() {
	var next *list.Element

	s.mu.RLock()
	defer s.mu.RUnlock()

	for e := s.conn.Front(); e != nil; e = next {
		next = e.Next()

		c := e.Value.(*gev.Connection)
		_ = c.Send([]byte("hello\n"))
	}
}

// OnConnect callback
func (s *Server) OnConnect(c *gev.Connection) {
	log.Println(" OnConnect ： ", c.PeerAddr())

	s.mu.Lock()
	e := s.conn.PushBack(c)
	s.mu.Unlock()
	c.Set(clientsKey, e)
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
	v, ok := c.Get(clientsKey)
	if !ok {
		log.Println("OnClose : get key fail")
		return
	}

	s.mu.Lock()
	s.conn.Remove(v.(*list.Element))
	s.mu.Unlock()
}

func main() {
	s, err := New("", "1833")
	if err != nil {
		panic(err)
	}
	defer s.Stop()

	s.Start()
}

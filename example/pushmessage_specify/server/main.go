package main

import (
	"encoding/binary"
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
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

var clientConnList = make(map[string]connectData)

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

	// 加入全局连接管理器中
	s.mu.Lock()
	clientConnList[c.PeerAddr()] = connectData{
		conn: c,
		data: c.PeerAddr(),
	}
	s.mu.Unlock()

	// 为每个连接启动一个协程，并定时（10秒内随机）发送消息
	// 	- 每次发送消息前检查当前协程的局部变量（socket连接）是否再全局管理连接中，不在则结束当前协程（避免当前协程一直占用资源）
	go func(c *gev.Connection) {
		for {
			randTime := int64(rand.Uint64() % 10)
			log.Println("time.Sleep:", randTime)
			time.Sleep(time.Second * time.Duration(randTime))

			if _, ok := clientConnList[c.PeerAddr()]; !ok {
				log.Println(" clientConnList c.PeerAddr() connenct close")
				return
			}
			err := c.Send([]byte(c.PeerAddr()))
			if err != nil {
				log.Println("c.Send err：", err.Error())
				continue
			}
			log.Println("c.Send success:", clientConnList)
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

	// 客户端主动关闭连接后，服务端应剔除该连接，否则会报错：connection closed
	if _, ok := clientConnList[c.PeerAddr()]; ok {
		s.mu.Lock()
		delete(clientConnList, c.PeerAddr())
		s.mu.Unlock()
	}
}

func main() {
	s, err := New("", "1833")
	if err != nil {
		panic(err)
	}
	defer s.Stop()

	s.Start()
}

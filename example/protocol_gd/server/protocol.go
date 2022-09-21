package main

import (
	"encoding/binary"
	"log"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
)

const exampleHeaderLen = 4

type ExampleProtocol struct{}

func (d *ExampleProtocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if buffer.VirtualLength() > exampleHeaderLen {
		buf := pbytes.GetLen(exampleHeaderLen)
		log.Println("*ExampleProtocol.UnPacket:buf:", buf)

		defer pbytes.Put(buf)
		_, _ = buffer.VirtualRead(buf)
		dataLen := binary.BigEndian.Uint32(buf)

		log.Println("*ExampleProtocol.UnPacket:dataLen:", dataLen)

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

package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/Allenxuxu/gev/example/protocol_gd/client/mockdata"
	"github.com/Allenxuxu/gev/example/protocol_gd/common"
	"gitlab.yewifi.com/golden-cloud/common/gclog"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func Packet(data []byte) []byte {
	buffer := make([]byte, 4+len(data))
	// 将buffer前面四个字节设置为包长度，大端序
	binary.BigEndian.PutUint32(buffer, uint32(len(data)))

	log.Println("Packet len(data):", len(data))
	//log.Println("Packet buffer(copy data before):", buffer)
	copy(buffer[4:], data)
	//log.Println("Packet buffer:", buffer)

	return buffer
}

func UnPacket(c net.Conn) ([]byte, error) {
	var header = make([]byte, 4)

	_, err := io.ReadFull(c, header)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)
	log.Println("UnPacket binary.BigEndian.Uint32 length:", length)

	contentByte := make([]byte, length)
	_, e := io.ReadFull(c, contentByte) //读取内容
	if e != nil {
		return nil, e
	}

	return contentByte, nil
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
	//conn, e := net.Dial("tcp", "10.20.16.49:30106") // gd 开发环境
	//conn, e := net.Dial("tcp", "10.20.16.49:30107") // gd 测试环境

	conn, e := net.Dial("tcp", "127.0.0.1:18000") // docker容器内部IP:端口
	//conn, e := net.Dial("tcp", "10.98.10.61:18000") // 本地宿主机IP:端口

	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	for {
		// --- 向TCP服务端发送请求
		sendMessage(conn)

		// --- 监听TCP服务端响应
		receveMessage(conn)
	}
}

func sendMessage(conn net.Conn) {
	// --- 模拟TCP客户端发送的数据 -------------------------
	log.Println("准备数据>选择发送消息的命令(可能被打印出的消息覆盖命令列表)...")

	// 发送的数据：读取标准输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n\n===============选择发送消息的命令============\n")
	fmt.Printf("0:心跳\n")
	fmt.Printf("1:上报cookie\n")
	fmt.Printf("2:上报税号登录态(客户端掉线重连)\n")
	fmt.Printf("3:响应代理请求\n")
	fmt.Printf("其他:重新输入\n")
	fmt.Printf("请输入命令:")

	cmd, _ := reader.ReadString('\n')
	cmdInt, _ := strconv.Atoi(cmd)
	fmt.Printf("输入的 cmd:%s; T:%T; cmdName:%s \n", cmd, cmd, common.CMD_name[int32(common.CMD(cmdInt))])
	// 去除右边的换行符
	cmd = strings.TrimRight(cmd, "\n")
	fmt.Printf("输入的 cmd 并去除右边的换行符后的 cmd:%s; T:%T \n", cmd, cmd)

	var sendMsg common.SendMsg
	switch cmd {
	case "0":
		//心跳
		sendMsg = mockdata.Pong()
	case "1":
		// 上报cookies
		sendMsg = mockdata.UploadCookie()
	case "2":
		// 上报税号登录态
		sendMsg = mockdata.UploadLogin()
	case "3":
		// 响应代理请求
		sendMsg = mockdata.SendProxyResponse()
	default:
		log.Println("不支持该指令", cmd, &sendMsg)
		return
	}
	log.Println("sendMsg", &sendMsg)

	sendMsgByte, err := proto.Marshal(&sendMsg)
	if err != nil {
		log.Panic("proto.Marshal err", err)
	}
	buffer := Packet(sendMsgByte)
	log.Println("数据打包结束")
	// --- 模拟TCP客户端发送的数据 -end

	// --- 发送数据
	log.Println("发送数据...")
	_, err = conn.Write(buffer)
	if err != nil {
		log.Panic("conn.Write err", err)
	}
	log.Println("发送数据结束")
	// --- 发送数据 -end
}

func receveMessage(conn net.Conn) {
	// --- 监听TCP服务端响应 =============================
	log.Println("接受数据...")
	respByte, err := UnPacket(conn)
	if err != nil {
		log.Panic("UnPacket respByte err:", err)
	}
	msgReq := common.SendMsg{}
	err = proto.Unmarshal(respByte, &msgReq)
	if err != nil {
		log.Panic(" proto.Unmarshal respByte err:", err)
		return
	}
	log.Println("msgReq:", &msgReq)
	log.Println("msgReq.Cmd:", msgReq.Cmd)
	log.Println("msgReq.Cmd CMD_name:", common.CMD_name[int32(msgReq.Cmd)])
	switch msgReq.Cmd {
	case common.CMD_CMD_NONE:
		log.Println("- switch CMD_CMD_NONE")
	case common.CMD_CMD_UPLOAD_LOGIN_RSP:
		log.Println("- switch CMD_CMD_UPLOAD_LOGIN_RSP")
	case common.CMD_CMD_HEART_BEAT_RSP:
		log.Println("- switch CMD_CMD_HEART_BEAT_RSP")
	case common.CMD_CMD_OPEN_URL_REQ:
		log.Println("- switch CMD_CMD_OPEN_URL")
		body := common.OpenUrl{}
		err := proto.Unmarshal(msgReq.Body, &body)
		if err != nil {
			log.Panic(" proto.Unmarshal(msgReq.Body, &body) err:", err)
			return
		}
		log.Println("body:", &body)
		log.Println("body.Url:", body.Url)
	case common.CMD_CMD_PROXY_RSP:
		log.Println("- switch CMD_CMD_PROXY_RSP")
	default:
		log.Println("- 监听的响应命令字暂不处理")
	}
	log.Println("接受数据结束")

	//log.Println("time.sleep...")
	//time.Sleep(time.Second * 1)
	// --- 监听TCP服务端响应 -end =============================
}

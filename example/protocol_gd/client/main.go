package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
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
	conn, e := net.Dial("tcp", "127.0.0.1:18000") // gd 本地

	//conn, e := net.Dial("tcp", "127.0.0.1:1834") // gev 本地

	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	for {
		// --- 模拟TCP客户端发送数据 -------------------------
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
		fmt.Printf("输入的 cmd: %s; T:%T; cmdName: %s \n", cmd, cmd, common.CMD_name[int32(common.CMD(cmdInt))])
		// 去除右边的换行符
		cmd = strings.TrimRight(cmd, "\n")
		fmt.Printf("输入的 cmd 并去除右边的换行符后的 cmd: %s; T:%T \n", cmd, cmd)

		// 配置信息
		token := "v5_FlfITAGQlOeINlVK1euASi7e6eHmFP5U1154764563"
		openEnterpriseId := "YpQHNOTfZ0KWvZZno1A2BOvojyShiQ3+7um3VkAxYUA="
		// 配置信息 -end

		var sendMsg common.SendMsg
		switch cmd {
		case "0":
			//心跳
			sendMsg = common.SendMsg{
				Cmd:   common.CMD_CMD_HEART_BEAT_REQ,
				Token: token,
				LogId: "prince-logID",
				Body:  nil,
			}
			//--- 心跳 -end
		case "1":
			// 上报cookies
			sendMsgBody := common.UploadCookie{
				BusinessId:       10107,
				OpenEnterpriseId: openEnterpriseId,
				Cookies: []*common.Cookies{
					{
						Key:    "SSO_SECURITY_CHECK_TOKEN",
						Value:  "90ce583ed7294d73a8548733e8fdaead",
						Domain: "https://d1.com",
						Expire: 0,
					},
					{
						Key:    "dzfp-ssotoken",
						Value:  "ade3f3f9a79447faa5a6867a38d6cc72",
						Domain: "https://d2.com",
						Expire: 0,
					},
					{
						Key:    "SSOTicket",
						Value:  "f1030b22-2523-4cb5-820f-e408d25f6871",
						Domain: "https://d2.com",
						Expire: 0,
					},
				},
				BrowserRequest: "",
			}
			bodyBytes, err := proto.Marshal(&sendMsgBody)
			if err != nil {
				log.Panic(" proto.Marshal(&sendMsgBody) err", err)
			}
			sendMsg = common.SendMsg{
				Cmd:   common.CMD_CMD_UPLOAD_COOKIE_REQ,
				Token: token,
				LogId: "prince-logID",
				Body:  bodyBytes,
			}
			// --- 上报cookies -end
		case "2":
			// 上报税号登录态
			sendMsgBody := common.UploadLogin{
				BusinessId:       10107,
				OpenEnterpriseId: openEnterpriseId,
				LoginFlag:        1,
			}
			bodyBytes, err := proto.Marshal(&sendMsgBody)
			if err != nil {
				log.Panic(" proto.Marshal(&sendMsgBody) err", err)
			}
			sendMsg = common.SendMsg{
				Cmd:   common.CMD_CMD_UPLOAD_LOGIN_REQ,
				Token: token,
				LogId: "prince-logID",
				Body:  bodyBytes,
			}
			// --- 上报税号登录态 -end
		case "3":
			// 响应代理请求
			// 客户端发送代理请求后，得到的响应内容是json
			rspDataByte := []byte("{\"Response\":{\"RequestId\":\"a991b4a3b415f77f\",\"Error\":null,\"Data\":{\"Message\":null,\"Code\":null,\"Wkfpuuid\":\"8e0eae4d29c6490da098912fecc4982f\",\"Fpkjdto\":null,\"SmkpHqgfxxRequestVO\":null}}}")
			sendMsgBody := common.ProxyRequestRsp{
				Seq:     "prince-test-seq",
				RspData: rspDataByte,
			}
			bodyBytes, err := proto.Marshal(&sendMsgBody)
			if err != nil {
				log.Panic(" proto.Marshal(&sendMsgBody) err", err)
			}
			sendMsg = common.SendMsg{
				Cmd:   common.CMD_CMD_PROXY_RSP,
				Token: token,
				LogId: "prince-logID",
				Body:  bodyBytes,
			}
			// --- 响应代理请求 -end
		default:
			log.Println("不支持该指令")
			continue
		}

		sendMsgByte, err := proto.Marshal(&sendMsg)
		if err != nil {
			log.Panic("proto.Marshal err", err)
		}
		buffer := Packet(sendMsgByte)
		log.Println("准备数据结束")
		// --- 模拟TCP客户端发送数据 -end -------------------------

		// --- 发送数据 ++++++++++++++++++++++++++++++++++
		log.Println("发送数据...")
		_, err = conn.Write(buffer)
		if err != nil {
			log.Panic("conn.Write err", err)
		}
		log.Println("发送数据结束")
		// --- 发送数据 -end ++++++++++++++++++++++++++++++++++

		// --- 监听TCP服务端响应
		receveMessage(conn)
	}
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
	log.Println("msgReq:", msgReq)
	log.Println("msgReq.Cmd:", msgReq.Cmd)
	log.Println("msgReq.Cmd CMD_name:", common.CMD_name[int32(msgReq.Cmd)])
	switch msgReq.Cmd {
	case common.CMD_CMD_NONE:
		log.Println("- switch CMD_CMD_NONE")
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
		log.Println("body:", body)
		log.Println("body.Url:", body.Url)
	case common.CMD_CMD_PROXY_RSP:
		log.Println("- switch CMD_CMD_PROXY_RSP")
	default:
		log.Println("- switch @default")
	}
	log.Println("接受数据结束")

	//log.Println("time.sleep...")
	//time.Sleep(time.Second * 1)
	// --- 监听TCP服务端响应 -end =============================
}

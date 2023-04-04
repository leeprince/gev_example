package mockdata

import (
	"github.com/Allenxuxu/gev/example/protocol_gd/common"
	"github.com/golang/protobuf/proto"
	"log"
)

/**
 * @Author: prince.lee <leeprince@foxmail.com>
 * @Date:   2023/4/4 16:35
 * @Desc:
 */

func SendProxyResponse() common.SendMsg {
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
	return common.SendMsg{
		Cmd:   common.CMD_CMD_PROXY_RSP,
		Token: token,
		LogId: "prince-logID",
		Body:  bodyBytes,
	}
}

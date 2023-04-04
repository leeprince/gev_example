package main

import (
	"github.com/Allenxuxu/gev/example/protocol_gd/common"
	"github.com/golang/protobuf/proto"
	"log"
)

/**
 * @Author: prince.lee <leeprince@foxmail.com>
 * @Date:   2023/4/4 16:32
 * @Desc:
 */

func uploadLogin() common.SendMsg {
	sendMsgBody := common.UploadLogin{
		BusinessId:       10107,
		OpenEnterpriseId: openEnterpriseId,
		LoginFlag:        1,
	}
	bodyBytes, err := proto.Marshal(&sendMsgBody)
	if err != nil {
		log.Panic(" proto.Marshal(&sendMsgBody) err", err)
	}
	return common.SendMsg{
		Cmd:   common.CMD_CMD_UPLOAD_LOGIN_REQ,
		Token: token,
		LogId: "prince-logID",
		Body:  bodyBytes,
	}
}

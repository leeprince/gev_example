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

func uploadCookie() common.SendMsg {
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
	return common.SendMsg{
		Cmd:   common.CMD_CMD_UPLOAD_COOKIE_REQ,
		Token: token,
		LogId: "prince-logID",
		Body:  bodyBytes,
	}
}

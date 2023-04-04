package mockdata

import "github.com/Allenxuxu/gev/example/protocol_gd/common"

/**
 * @Author: prince.lee <leeprince@foxmail.com>
 * @Date:   2023/4/4 16:27
 * @Desc:
 */

func Pong() common.SendMsg {
	return common.SendMsg{
		Cmd:   common.CMD_CMD_HEART_BEAT_REQ,
		Token: token,
		LogId: "prince-logID",
		Body:  nil,
	}
}

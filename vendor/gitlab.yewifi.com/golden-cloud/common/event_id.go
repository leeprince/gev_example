/*
 * @Date: 2021-08-18 10:34:58
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2021-08-18 14:32:05
 */

package common

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	tableUpperCaseLetters                           = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tableLowerCaseLetters                           = "abcdefghijklmnopqrstuvwxyz"
	tableNumber                                     = "0123456789"
	tableUpperCaseLetterAndLowerCaseLetterAndNumber = tableUpperCaseLetters + tableLowerCaseLetters + tableNumber
	defaultEventIDLength                            = 24
	defaultEventIDTable                             = tableUpperCaseLetterAndLowerCaseLetterAndNumber
)

type rander struct {
	r *rand.Rand
}

func newRander(seed int64) *rander {
	r := rand.New(rand.NewSource(seed))
	return &rander{r: r}
}

func (r *rander) RandString(length int, table string) string {
	if length == 0 || table == "" {
		return ""
	}
	n := len(table)
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = table[r.r.Intn(n)]
	}

	return string(bytes)
}

func EventID(logID ...string) string {
	var s string
	if len(logID) != 0 {
		s = logID[0]
	}
	buf := make([]byte, (defaultEventIDLength-len(s))/2)
	_, err := rand.Read(buf)
	if err != nil {
		// 当使用 rand.Read 发生错误，使用补偿函数来生成 event id
		return s + newRander(time.Now().UnixNano()).RandString(defaultEventIDLength-len(s), defaultEventIDTable) // 补偿方法
	}

	s += fmt.Sprintf("%x", buf)
	n := len(s)
	for i := 0; i < defaultEventIDLength-n; i++ {
		s += "0"
	}
	return s
}

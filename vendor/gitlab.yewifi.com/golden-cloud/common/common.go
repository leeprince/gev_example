package common

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.yewifi.com/golden-cloud/common/crypto"
	"gitlab.yewifi.com/golden-cloud/common/tracer"
	"google.golang.org/grpc/metadata"
)

func Uniqid() string {
	now := time.Now()
	return fmt.Sprintf("%08x%08x", now.Unix(), now.UnixNano()%0x100000)
}

func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		t = append(t, 'X')
		i++
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && IsASCIIUpper(s[i+1]) {
			continue
		}
		if IsASCIIDigit(c) {
			t = append(t, c)
			continue
		}

		if IsASCIIUpper(c) {
			c ^= ' '
		}
		t = append(t, c)

		for i+1 < len(s) && IsASCIIUpper(s[i+1]) {
			i++
			t = append(t, '_')
			t = append(t, bytes.ToLower([]byte{s[i]})[0])
		}
	}
	return string(t)
}

func IsASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func IsASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// TimeFunc 打印从调用到结束的耗时
// usage: defer TimeFunc("Hello world")()
func TimeFunc(v ...interface{}) func() {
	start := time.Now()
	return func() {
		fmt.Println(append(v, "|", time.Since(start)))
	}
}

func TokenByCtx(ctx context.Context) (string, error) { //context.Context为interface 无需声明为指针类型
	return ValByCtx(ctx, "access-token")
}

func AppkeyByCtx(ctx context.Context) (string, error) {
	return ValByCtx(ctx, "appkey")
}

func ValByCtx(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMetadataNotFound
	}

	list, ok := md[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	if len(list) <= 0 || list[0] == "" {
		return "", ErrKeyNotFound
	}

	return list[0], nil
}

func LogIdByCtx(ctx context.Context) string {
	tracerId, err := tracer.TraceID(&ctx)
	if err != nil || tracerId.String() == "" {
		// NOTE: fallback
		return Uniqid()
	}
	// NOTE: fill tracerId to 16 byte
	if len(tracerId.String()) == 15 {
		return "0" + tracerId.String()
	}

	return tracerId.String()
}

// CtxCopy returns new context with same metadata
// Modified by aiden.deng 2021-09-09: 加上 span 的拷贝，从而保持 trace id 的传递。
func CtxCopy(ctx context.Context) context.Context {
	span := opentracing.SpanFromContext(ctx)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return opentracing.ContextWithSpan(context.Background(), span)
	}

	return opentracing.ContextWithSpan(metadata.NewIncomingContext(context.Background(), md), span)
}

type SignType string

const (
	SIGN_HMAC SignType = "HMAC-SHA256"
	SIGN_RSA  SignType = "RSA-SHA256"

	HeaderKeyLogId = "log_id"
)

//HMAC-SHA256签名
func SignWithHmacSha256(appkey string, path string, payLoad string, secretKey string) (string, string, string) {
	randNum := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	src := "algorithm=HMAC-SHA256|appkey="
	src += appkey
	src += "|nonce="
	src += randNum
	src += "|timestamp="
	src += timeStamp
	src += "|"
	src += path
	src += "|"
	src += payLoad

	return crypto.HmacSha256(src, secretKey), randNum, timeStamp
}

//RSA-SHA256签名
func SignWithRsaSha256(appkey string, path string, payLoad string, privateKey string) (string, string, string, error) {
	randNum := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	src := "algorithm=RSA-SHA256"
	if appkey != "" {
		//推送时签名不需要appkey
		src += "|appkey="
		src += appkey
	}
	src += "|nonce="
	src += randNum
	src += "|timestamp="
	src += timeStamp
	src += "|"
	src += path
	src += "|"
	src += payLoad

	signature, err := crypto.RsaSignWithSha256([]byte(src), []byte(privateKey))
	if err != nil {
		return "", "", "", err
	}

	return signature, randNum, timeStamp, nil
}

//RSA-SHA256签名
func SignWithRsaSha256_v2(appkey string, path string, payLoad string, privateKey string) ([]byte, string, string, error) {
	randNum := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	src := "algorithm=RSA-SHA256"
	if appkey != "" {
		//推送时签名不需要appkey
		src += "|appkey="
		src += appkey
	}
	src += "|nonce="
	src += randNum
	src += "|timestamp="
	src += timeStamp
	src += "|"
	src += path
	src += "|"
	src += payLoad

	fmt.Println("待加密串:", src)
	signature, err := crypto.RsaSignWithSha256_v2([]byte(src), []byte(privateKey))
	if err != nil {
		return nil, "", "", err
	}

	return signature, randNum, timeStamp, nil
}

func GenAuthorization(siginType SignType, appkey string, path string, payLoad string, secretKey string) (string, error) {

	var algorithm, signature, randNum, timeStamp string
	var err error

	if SIGN_HMAC == siginType {
		algorithm = "HMAC-SHA256"
		signature, randNum, timeStamp = SignWithHmacSha256(appkey, path, payLoad, secretKey)
	} else if SIGN_RSA == siginType {
		algorithm = "RSA-SHA256"
		signature, randNum, timeStamp, err = SignWithRsaSha256(appkey, path, payLoad, secretKey)
		if err != nil {
			return "", err
		}
	}

	autuStr := "algorithm="
	autuStr += algorithm
	if appkey != "" {
		autuStr += ",appkey="
		autuStr += appkey
	}
	autuStr += ",nonce="
	autuStr += randNum
	autuStr += ",timestamp="
	autuStr += timeStamp
	autuStr += ",signature="
	autuStr += base64.StdEncoding.EncodeToString([]byte(signature))

	return autuStr, nil
}

func GenAuthorization_v2(siginType SignType, appkey string, path string, payLoad string, secretKey string) (string, error) {

	var algorithm, signature, randNum, timeStamp string
	var signBytes []byte
	var err error

	if SIGN_HMAC == siginType {
		algorithm = "HMAC-SHA256"
		signature, randNum, timeStamp = SignWithHmacSha256(appkey, path, payLoad, secretKey)
		signBytes = []byte(signature)
	} else if SIGN_RSA == siginType {
		algorithm = "RSA-SHA256"
		signBytes, randNum, timeStamp, err = SignWithRsaSha256_v2(appkey, path, payLoad, secretKey)
		if err != nil {
			return "", err
		}
	}

	autuStr := "algorithm="
	autuStr += algorithm
	if appkey != "" {
		autuStr += ",appkey="
		autuStr += appkey
	}
	autuStr += ",nonce="
	autuStr += randNum
	autuStr += ",timestamp="
	autuStr += timeStamp
	autuStr += ",signature="
	autuStr += base64.StdEncoding.EncodeToString(signBytes)

	return autuStr, nil
}

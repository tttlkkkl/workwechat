package workwechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/sbabiv/xml2map"
)

// ReceiveMessage 事件与消息接收
type ReceiveMessage struct {
	Token          string
	EncodingAESKey string
	aesKey         []byte
}

// ReceiveMessageBody 原始消息内容
type ReceiveMessageBody struct {
	// Message 解密后的xml消息内容
	Message []byte
	// Receiveid 企业应用的回调，表示corpid，第三方事件的回调，表示suiteid
	Receiveid []byte
	// Data xml解码后的消息
	Data Values
}

// xmlMessage 原始的事件消息
type xmlMessage struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
	AgentID    string `xml:"AgentID"`
}

// NewReceiveMessage 消息事件处理
func newReceiveMessage(token, encodingAesKey string) (*ReceiveMessage, error) {
	e := &ReceiveMessage{
		Token:          token,
		EncodingAESKey: encodingAesKey,
	}
	b, err := base64.StdEncoding.DecodeString(encodingAesKey + "=")
	if err != nil {
		return nil, fmt.Errorf("请填写正确的 EncodingAESKey %w", err)
	}
	e.aesKey = b
	return e, nil
}

// Handle 消息接收处理
func (e *ReceiveMessage) Handle(req *http.Request, rep http.ResponseWriter) (*ReceiveMessageBody, error) {
	var err error
	body := &ReceiveMessageBody{
		Data: make(Values),
	}
	// url 响应
	query := req.URL.Query()
	if req.Method == http.MethodGet {
		loger.Info("收到消息接收服务器配置请求", req.URL.RequestURI())
		echoStr, err := url.PathUnescape(query.Get("echostr"))
		if err != nil {
			return body, fmt.Errorf("消息接收服务器响应失败，无法解析的查询：%w", err)
		}
		echoStr, err = e.verifyURL(query.Get("msg_signature"), query.Get("timestamp"), query.Get("nonce"), echoStr)
		if err != nil {
			return body, err
		}
		rep.Write([]byte(echoStr))
		return body, nil
	}
	postData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return body, err
	}
	mp := xmlMessage{}
	err = xml.Unmarshal(postData, &mp)
	if err != nil {
		return body, err
	}
	body.Message, body.Receiveid, err = e.decryptMsg(query.Get("msg_signature"), query.Get("timestamp"), query.Get("nonce"), mp.Encrypt)
	if err != nil {
		return body, err
	}
	xm := xml2map.NewDecoder(bytes.NewReader(body.Message))
	result, err := xm.Decode()
	if err != nil {
		return body, err
	}
	body.Data = Values(result)
	return body, nil
}

// 签名校验
func (e *ReceiveMessage) verifySignature(sMsgSignature, sTimeStamp, sNonce, data string) error {
	sorted := []string{e.Token, sTimeStamp, sNonce, data}
	sort.Strings(sorted)
	sha1 := sha1.New()
	sortedString := []byte(strings.Join(sorted, ""))
	sha1.Write(sortedString)
	sha1String := hex.EncodeToString(sha1.Sum([]byte(nil)))
	if sha1String != sMsgSignature {
		return errors.New("参数签名校验失败")
	}
	return nil
}

// url 校验
func (e *ReceiveMessage) verifyURL(sMsgSignature, sTimeStamp, sNonce, sEchoStr string) (string, error) {
	m, _, err := e.decryptMsg(sMsgSignature, sTimeStamp, sNonce, sEchoStr)

	return string(m), err
}

// 消息解密
func (e *ReceiveMessage) decryptMsg(sMsgSignature, sTimeStamp, sNonce, sPostData string) (message, receiveid []byte, err error) {
	// 校验签名
	err = e.verifySignature(sMsgSignature, sTimeStamp, sNonce, sPostData)
	if err != nil {
		return
	}
	b, err := base64.StdEncoding.DecodeString(sPostData)
	if err != nil {
		return nil, nil, fmt.Errorf("消息解密失败：%w", err)
	}
	rBytes, err := e.aesDecode(b)
	if err != nil {
		return nil, nil, err
	}
	// 消息长度
	buf := bytes.NewBuffer(rBytes[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	if length <= 0 || len(rBytes) <= 0 {
		return nil, nil, errors.New("消息长度错误，请检查配置")
	}
	appIDstart := 20 + length
	if len(rBytes) < int(appIDstart) {
		return nil, nil, errors.New("数据长度错误，请检查配置")
	}
	// 读取消息内容
	message = rBytes[20 : 20+length]
	// 读取 receiveid
	receiveid = rBytes[20+length:]
	return message, receiveid, nil
}

// 消息加密
func (e *ReceiveMessage) encryptMsg(sReplyMsg, sTimeStamp, sNonce string) (string, error) {
	// TODO
	return "", nil
}
func (e *ReceiveMessage) aesDecode(sb []byte) (bt []byte, err error) {
	block, err := aes.NewCipher(e.aesKey)
	if err != nil {
		log.Printf("aes解密失败: %v", err)
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, e.aesKey[:blockSize])
	origData := make([]byte, len(sb))
	blockMode.CryptBlocks(origData, sb)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

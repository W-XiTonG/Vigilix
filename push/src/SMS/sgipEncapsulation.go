package SMS

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"sync"
	"sync/atomic"
	"time"
)

// 定义命令ID常量
const (
	SGIP_BIND               = 0x00000001
	SGIP_UNBIND             = 0x00000002
	SGIP_SUBMIT             = 0x00000003
	SGIP_UNBIND_RESP        = 0x80000002
	SGIP_SUBMIT_RESP        = 0x80000003
	SGIP_DELIVER            = 0x4
	SGIP_DELIVER_RESP       = 0x80000004
	SGIP_REPORT             = 0x5
	SGIP_REPORT_RESP        = 0x80000005
	SGIP_ADDSP              = 0x6
	SGIP_ADDSP_RESP         = 0x80000006
	SGIP_MODIFYSP           = 0x7
	SGIP_MODIFYSP_RESP      = 0x80000007
	SGIP_DELETESP           = 0x8
	SGIP_DELETESP_RESP      = 0x80000008
	SGIP_QUERYROUTE         = 0x9
	SGIP_QUERYROUTE_RESP    = 0x80000009
	SGIP_ADDTELESEG         = 0xa
	SGIP_ADDTELESEG_RESP    = 0x8000000a
	SGIP_MODIFYTELESEG      = 0xb
	SGIP_MODIFYTELESEG_RESP = 0x8000000b
	SGIP_DELETETELESEG      = 0xc
	SGIP_DELETETELESEG_RESP = 0x8000000c
	SGIP_ADDSMG             = 0xd
	SGIP_ADDSMG_RESP        = 0x8000000d
	SGIP_MODIFYSMG          = 0xe
	SGIP_MODIFYSMG_RESP     = 0x0000000e
	SGIP_DELETESMG          = 0xf
	SGIP_DELETESMG_RESP     = 0x8000000f
	SGIP_CHECKUSER          = 0x10
	SGIP_CHECKUSER_RESP     = 0x80000010
	SGIP_USERRPT            = 0x11
	SGIP_USERRPT_RESP       = 0x80000011
	SGIP_TRACE              = 0x1000
	SGIP_TRACE_RESP         = 0x80001000
)

// Encode 对 SGIPBind 请求进行编码
func (req SGIPBind) Encode() ([]byte, error) {
	var buf bytes.Buffer

	// 计算消息总长度
	req.Header.TotalLength = uint32(binary.Size(req.Header) + binary.Size(req.LoginType) +
		binary.Size(req.LoginName) + binary.Size(req.LoginPassword) + binary.Size(req.Reserve))

	// 编码消息头
	if err := binary.Write(&buf, binary.BigEndian, req.Header.TotalLength); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Header.CommandID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Header.SequenceID); err != nil {
		return nil, err
	}

	// 编码消息体
	if err := binary.Write(&buf, binary.BigEndian, req.LoginType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.LoginName); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.LoginPassword); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Reserve); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var (
	sequenceNumber uint32
	mutex          sync.Mutex
)

// 生成 SequenceID
func generateSequenceID(nodeID uint32) [12]byte {
	var sequenceID [12]byte

	// NodeID (4 bytes)
	binary.BigEndian.PutUint32(sequenceID[0:4], nodeID)

	// Timestamp (4 bytes)
	timestamp := uint32(time.Now().Unix())
	binary.BigEndian.PutUint32(sequenceID[4:8], timestamp)

	// SequenceNumber (4 bytes)
	binary.BigEndian.PutUint32(sequenceID[8:12], atomic.AddUint32(&sequenceNumber, 1))

	return sequenceID
}

func (req SGIPSubmit) Encode() ([]byte, error) {
	var buf bytes.Buffer

	// 计算消息总长度
	req.Header.TotalLength = uint32(binary.Size(req.Header) + binary.Size(req.SPNumber) +
		binary.Size(req.ChargeNumber) + binary.Size(req.UserCount) + binary.Size(req.UserNumber) +
		binary.Size(req.CorpID) + binary.Size(req.ServiceType) + binary.Size(req.FeeType) +
		binary.Size(req.FeeValue) + binary.Size(req.GivenValue) + binary.Size(req.AgentFlag) +
		binary.Size(req.MorelatetoMTFlag) + binary.Size(req.Priority) + binary.Size(req.ExpireTime) +
		binary.Size(req.ScheduleTime) + binary.Size(req.ReportFlag) + binary.Size(req.TP_pid) +
		binary.Size(req.TP_udhi) + binary.Size(req.MessageCoding) + binary.Size(req.MessageType) +
		binary.Size(req.MessageLength) + len(req.MessageContent) + binary.Size(req.Reserve))
	// 编码消息头
	if err := binary.Write(&buf, binary.BigEndian, req.Header.TotalLength); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Header.CommandID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Header.SequenceID); err != nil {
		return nil, err
	}
	// 编码消息体
	if err := binary.Write(&buf, binary.BigEndian, req.SPNumber); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.ChargeNumber); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.UserCount); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.UserNumber); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.CorpID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.ServiceType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.FeeType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.FeeValue); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.GivenValue); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.AgentFlag); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.MorelatetoMTFlag); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Priority); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.ExpireTime); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.ScheduleTime); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.ReportFlag); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.TP_pid); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.TP_udhi); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.MessageCoding); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.MessageType); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.MessageLength); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.MessageContent); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Reserve); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// encodeMessage 根据 MessageCoding 编码短信内容
func encodeMessage(message string, messageCoding uint8) ([]byte, error) {
	switch messageCoding {
	case 8: // UTF-16BE
		encoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
		encoded, _, err := transform.Bytes(encoder, []byte(message))
		return encoded, err
	case 15: // GBK
		// 这里省略 GBK 编码的实现
		return []byte(message), nil
	default: // 默认使用 UTF-8
		return []byte(message), nil
	}
}

// getSGIPTime 生成 SGIP 协议的时间字符串（yymmddhhmmsstnnp 格式）
func getSGIPTime(t time.Time) [16]byte {
	var timeStr [16]byte

	// 格式化为 yymmddhhmmsstnnp
	timeFormat := t.Format("060102150405") + "032+"
	copy(timeStr[:], []byte(timeFormat))

	return timeStr
}

func (req SGIPUnbind) Encode() ([]byte, error) {
	var buf bytes.Buffer

	// 计算消息总长度
	req.Header.TotalLength = uint32(binary.Size(req.Header))
	// 编码消息头
	if err := binary.Write(&buf, binary.BigEndian, req.Header.TotalLength); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Header.CommandID); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, req.Header.SequenceID); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

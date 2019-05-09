package socket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	//"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

type MsgHeader struct {
	Header  rune
	Version uint32
	Type    uint8
	DataLen uint32
	Body    string
}

type HTTP_REQ struct {
	MSG_ID   uint32 `xml:"MSG_ID,attr"`
	PROJ_ID  uint32 `xml:"PROJ_ID,attr"`
	URL      string `xml:"URL,attr"`
	METHOD   string `xml:"METHOD,attr"`
	CRT_TIME string `xml:"CRT_TIME,attr"`
	HEADERS  []ITEM `xml:"HEADERS>HEADER"`
	PARAMS   []ITEM `xml:"PARAMS>PARAM"`
}

type ITEM struct {
	KEY   string `xml:"KEY"`
	VALUE string `xml:"VALUE"`
}

type HTTP_RSP struct {
	MSG_ID         uint32 `xml:"MSG_ID,attr"`
	PROJ_ID        uint32 `xml:"PROJ_ID,attr"`
	CRT_TIME       string `xml:"CRT_TIME,attr"`
	MNT_ST         string `xml:"MNT_ST"`
	MNT_ET         string `xml:"MNT_ET"`
	MNT_CD         uint32 `xml:"MNT_CD"`
	MNT_IF         string `xml:"MNT_IF"`
	MNT_URL        string `xml:"MNT_URL"`
	MNT_RSP_HEADER string `xml:"MNT_RSP_HEADER"`
	MNT_RSP_BODY   string `xml:"MNT_RSP_BODY"`
}

const (
	MSG_VERSION        = 1000
	MSG_HEADER         = 0x02
	MSG_CHECK          = "9999999999999999"
	MSG_TYPE_HTTP_NORM = 0x31
	MSG_TYPE_HTTP_FLOW = 0x32
	MSG_TYPE_UNKNOWN   = 0x33
	MSG_TAIL           = 0x03
	MSG_RESERVE        = 1111
)

type TCP struct {
	addr     string
	port     string
	conn     *net.TCPConn
	maxRetry int
	timeout  int
}

type Config struct {
	Addr     string
	Port     string
	MaxRetry string
	Timeout  string
}

var SocketConfig Config

func NewTCP() *TCP {
	tcp := new(TCP)
	tcp.addr = SocketConfig.Addr
	tcp.port = SocketConfig.Port
	tcp.maxRetry, _ = strconv.Atoi(SocketConfig.MaxRetry)
	tcp.timeout, _ = strconv.Atoi(SocketConfig.Timeout)
	tcp.conn = nil
	return tcp
}

func (tcp *TCP) connect() error {
	addr, err := net.ResolveTCPAddr("tcp", tcp.addr+":"+tcp.port)
	if err != nil {
		return err
	}

	var i int = 0
	for {
		conn, connErr := net.DialTCP("tcp", nil, addr)
		//default timeout 1 minute
		conn.SetDeadline(time.Now().Add(time.Duration(tcp.timeout) * time.Second))
		if connErr == nil && conn != nil {
			//set buffer 1M
			conn.SetReadBuffer(1 << 20)
			//set buffer 1M
			conn.SetWriteBuffer(1 << 20)
			tcp.conn = conn
			break
		}

		if i > tcp.maxRetry {
			return connErr
		}
		i += 1

	}

	return nil
}

func (tcp *TCP) Write(conn *net.TCPConn) error {

	return nil
}

//func (tcp *TCP) ReadWrite(rw func(conn *net.TCPConn) error) error {
//	for tcp.conn != nil {
//		fmt.Println("connection [%s-%s] in use", tcp.addr, tcp.port)
//		time.Sleep(1 * time.Second)
//	}
//	//connect TCP
//	connErr := tcp.connect()
//	if connErr != nil {
//		return connErr
//	}

//	defer (func() {
//		closeErr := tcp.close()
//		if closeErr != nil {
//			fmt.Println("close the [%s-%s] connection fail", tcp.addr, tcp.port)
//		}
//	})()

//	return rw(tcp.conn)

//}

func (tcp *TCP) close() error {
	if tcp.conn == nil {
		return nil
	}

	closeErr := tcp.conn.Close()
	if closeErr != nil {
		return closeErr
	}
	//release conn
	tcp.conn = nil
	return nil
}

var HTTP_QUEUE_REQ chan []byte = make(chan []byte, 100)
var HTTP_QUEUE_RSP chan []byte = make(chan []byte, 100)

func SendREQ(req []byte) error {
	HTTP_QUEUE_REQ <- req
	return nil
}

func (tcp *TCP) sendMsg(buf []byte) error {
	sendBytes, err2 := Pack(buf, 1)
	if err2 != nil {
		fmt.Println(err2)
	}
	n, err := tcp.conn.Write(sendBytes)
	if n != len(sendBytes) || err != nil {
		return err
	}
	bytes, err1 := ioutil.ReadAll(tcp.conn)
	if err1 == nil {
		HTTP_QUEUE_RSP <- bytes
	} else {
		return err1
	}
	defer tcp.close()

	return nil
}

//func (tcp *TCP) recvMsg(c *net.Conn) ([]byte, error) {
//	bytes, err := ioutil.ReadAll(tcp.conn)
//	if err != nil {
//		return nil, err
//	}
//	return bytes, nil
//}

func Pack(req []byte, msgType int) ([]byte, error) {
	buff := new(bytes.Buffer)

	//step1: write header
	binary.Write(buff, binary.BigEndian, byte(MSG_HEADER))
	//buff.WriteByte(byte(MSG_HEADER))
	//step2: write version
	//buff.WriteRune(rune(1000))
	binary.Write(buff, binary.BigEndian, uint32(MSG_VERSION))
	//step3: write check
	binary.Write(buff, binary.BigEndian, []byte(MSG_CHECK))
	//buff.WriteString(MSG_CHECK)
	//step4: write type
	if msgType == 1 {
		binary.Write(buff, binary.BigEndian, byte(MSG_TYPE_HTTP_NORM))
		//buff.WriteByte(MSG_TYPE_HTTP_NORM)
	} else if msgType == 2 {
		binary.Write(buff, binary.BigEndian, byte(MSG_TYPE_HTTP_FLOW))
		//buff.WriteByte(MSG_TYPE_HTTP_NORM)
	} else {
		binary.Write(buff, binary.BigEndian, byte(MSG_TYPE_UNKNOWN))
		//buff.WriteByte(MSG_TYPE_UNKNOWN)
	}
	//step5: write datalen
	binary.Write(buff, binary.BigEndian, uint32(len(req)))
	//buff.WriteRune(rune(len(req)))
	//step6: write reserve
	binary.Write(buff, binary.BigEndian, uint32(MSG_RESERVE))
	//step7: write data
	binary.Write(buff, binary.BigEndian, req)
	binary.Write(buff, binary.BigEndian, byte(MSG_TAIL))
	return buff.Bytes(), nil

}

func UnPack(rsp []byte) (string, error) {

	if len(rsp) < 31 {
		return "", errors.New("Bad Message!")
	}

	var err error

	//reader contain rsp bytes
	reader := bytes.NewReader(rsp)
	//step1: read header
	var header byte
	err = binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return "", err
	}
	//step2: read version
	var version uint32
	binary.Read(reader, binary.BigEndian, &version)
	//step4: read check
	var check [16]byte
	binary.Read(reader, binary.BigEndian, &check)
	var checkBytes []byte = check[:]
	fmt.Println("Message check", string(checkBytes))
	//step4: read type
	var msgType byte
	binary.Read(reader, binary.BigEndian, &msgType)
	//step5: read datalen
	var datalen uint32
	binary.Read(reader, binary.BigEndian, &datalen)
	//step6: read data
	var data = make([]byte, datalen)
	binary.Read(reader, binary.BigEndian, &data)
	//step7: read tail
	var tail byte
	binary.Read(reader, binary.BigEndian, &tail)
	return string(data), nil
	//	//-----------------------
	//	//step1 : read header
	//	binary.Read()
	//	header, errHeader := buf.ReadByte()
	//	if errHeader != nil {
	//		return "", errors.New("Read Version error -" + errHeader.Error())
	//	}
	//	if header != byte(0x02) {
	//		return "", errors.New("Bad Header")
	//	}
	//	//step2 : read version
	//	version := buf.Next(4)
	//	fmt.Println("Message.Version->", string(version))
	//	check := buf.Next(16)
	//	//step3: read check
	//	fmt.Println("Message check->", string(check))
	//	//step4: read type
	//	msgType, msgTypeError := buf.ReadByte()
	//	if msgTypeError != nil {
	//		return "", errors.New("Read Type Error - " + msgTypeError.Error())
	//	}
	//	if msgType != byte(MSG_TYPE_HTTP_NORM) || msgType != byte(MSG_TYPE_HTTP_FLOW) {
	//		return "", errors.New("Bad Type")
	//	}
	//	//step5: read datalen
	//	datalenRune, n, datalenError := buf.ReadRune()
	//	if datalenError != nil || n != 4 {
	//		return "", errors.New("Read DataLen error")
	//	}
	//	dataLenInt := int(datalenRune)
	//	if dataLenInt <= 0 {
	//		return "", errors.New("Bad dataLen value:" + strconv.Itoa(dataLenInt))
	//	}
	//	content := buf.Next(dataLenInt)
	//	result := string(content)
	//	fmt.Println("Content==>\n", result)
	//	tail, tailErr := buf.ReadByte()
	//	if tail != byte(MSG_TAIL) || tailErr != nil {
	//		return "", tailErr
	//	}
	return "", nil
}

func Start() {
	fmt.Println("start.....")
	for {
		select {
		case req := <-HTTP_QUEUE_REQ:
			fmt.Println("get HTTP_QUEUE_REQ")
			go func() {
				tcp := NewTCP()
				tcp.connect()
				tcp.sendMsg(req)
			}()
		case rsp := <-HTTP_QUEUE_RSP:
			fmt.Println("get HTTP_QUEUE_RSP")
			go func() {
				xml, err := UnPack(rsp)
				if err == nil {
					fmt.Println(xml)
				}
			}()

		default:
			fmt.Println("nothing")

		}
	}

}

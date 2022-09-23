package utils

import (
	"client/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// Transfer 将这些方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //传输时，使用缓冲
}

func (t *Transfer) ReadPkg() (mes message.Message, err error) {
	//conn.Read在conn没有被关闭的情况下，才会阻塞，如果客户端关闭了conn，就不会阻塞
	_, err = t.Conn.Read(t.Buf[:4])

	if err != nil {
		//err = errors.New("read pkg header err")
		return
	}
	//根据bug[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(t.Buf[0:4])
	//	根据pkgLen读取消息内容
	n, err := t.Conn.Read(t.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body err")
		return
	}
	//	把pkgLen反序列化成=》message.Message
	err = json.Unmarshal(t.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

func (t *Transfer) WritePkg(data []byte) (err error) {
	//发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(t.Buf[0:4], pkgLen)
	//发送长度
	n, err := t.Conn.Write(t.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	//	发送data本身
	n, err = t.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}

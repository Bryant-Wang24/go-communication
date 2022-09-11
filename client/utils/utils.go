package utils

import (
	"client/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func ReadPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据...")
	//conn.Read在conn没有被关闭的情况下，才会阻塞，如果客户端关闭了conn，就不会阻塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header err")
		return
	}
	//根据bug[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])
	//	根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body err")
		return
	}
	//	把pkgLen反序列化成=》message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

func WritePkg(conn net.Conn, data []byte) (err error) {
	//发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	//	发送data本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}

// Package vnc implements a VNC client.
//
// References:
//   [PROTOCOL]: http://tools.ietf.org/html/rfc6143
package vnc

import (
	"encoding/binary"
	"io"
)

// type DataSource struct {
// 	conn        io.Reader
// 	output      io.Writer
// 	passThrough bool
// 	PixelFormat PixelFormat
// }Color

func (d *ClientConn) readBytes(count int) ([]byte, error) {
	buff := make([]byte, count)

	_, err := io.ReadFull(d.conn, buff)
	if err != nil {
		//if err := binary.Read(d.conn, binary.BigEndian, &buff); err != nil {
		return nil, err
	}
	return buff, nil
}

func (d *ClientConn) readUint8() (uint8, error) {
	var myUint uint8
	if err := binary.Read(d.conn, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (d *ClientConn) readUint16() (uint16, error) {
	var myUint uint16
	if err := binary.Read(d.conn, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (d *ClientConn) readUint32() (uint32, error) {
	var myUint uint32
	if err := binary.Read(d.conn, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (d *ClientConn) readCompactLen() (int, error) {
	var err error
	part, err := d.readUint8()
	//byteCount := 1
	len := uint32(part & 0x7F)
	if (part & 0x80) != 0 {
		part, err = d.readUint8()
		//byteCount++
		len |= uint32(part&0x7F) << 7
		if (part & 0x80) != 0 {
			part, err = d.readUint8()
			//byteCount++
			len |= uint32(part&0xFF) << 14
		}
	}

	//   for  i := 0; i < byteCount; i++{
	//     rec.writeByte(portion[i]);
	//   }

	return int(len), err
}

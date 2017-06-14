// Package vnc implements a VNC client.
//
// References:
//   [PROTOCOL]: http://tools.ietf.org/html/rfc6143
package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// type DataSource struct {
// 	conn        io.Reader
// 	output      io.Writer
// 	passThrough bool
// 	PixelFormat PixelFormat
// }Color
type RfbReader struct {
	reader    io.Reader
	saveBytes bool
	savedBuff bytes.Buffer
}

func (r *RfbReader) Read(p []byte) (n int, err error) {
	readLen, err := r.reader.Read(p)
	r.savedBuff.Write(p)
	return readLen, err
}
func (r *RfbReader) SavedBuff() bytes.Buffer {
	return r.savedBuff
}

type RfbReadHelper struct {
	io.Reader
}

func (d *RfbReadHelper) ReadBytes(count int) ([]byte, error) {
	buff := make([]byte, count)

	_, err := io.ReadFull(d.Reader, buff)
	if err != nil {
		//if err := binary.Read(d.conn, binary.BigEndian, &buff); err != nil {
		return nil, err
	}
	return buff, nil
}

func (d *RfbReadHelper) ReadUint8() (uint8, error) {
	var myUint uint8
	if err := binary.Read(d.Reader, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (d *RfbReadHelper) ReadUint16() (uint16, error) {
	var myUint uint16
	if err := binary.Read(d.Reader, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (d *RfbReadHelper) ReadUint32() (uint32, error) {
	var myUint uint32
	if err := binary.Read(d.Reader, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (d *RfbReadHelper) ReadCompactLen() (int, error) {
	var err error
	part, err := d.ReadUint8()
	//byteCount := 1
	len := uint32(part & 0x7F)
	if (part & 0x80) != 0 {
		part, err = d.ReadUint8()
		//byteCount++
		len |= uint32(part&0x7F) << 7
		if (part & 0x80) != 0 {
			part, err = d.ReadUint8()
			//byteCount++
			len |= uint32(part&0xFF) << 14
		}
	}

	//   for  i := 0; i < byteCount; i++{
	//     rec.writeByte(portion[i]);
	//   }

	return int(len), err
}

var TightMinToCompress int = 12

func (r *RfbReadHelper) ReadTightData(dataSize int) ([]byte, error) {
	if int(dataSize) < TightMinToCompress {
		return r.ReadBytes(int(dataSize))
	}
	zlibDataLen, err := r.ReadCompactLen()
	fmt.Printf("compactlen=%d\n", zlibDataLen)
	if err != nil {
		return nil, err
	}
	//byte[] zlibData = new byte[zlibDataLen];
	//rfb.readFully(zlibData);
	return r.ReadBytes(zlibDataLen)
}

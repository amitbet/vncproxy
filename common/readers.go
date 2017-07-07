package common

import (
	"encoding/binary"
	"fmt"
	"io"
)

var TightMinToCompress = 12

const (
	SegmentBytes SegmentType = iota
	SegmentMessageSeparator
	SegmentRectSeparator
	SegmentFullyParsedClientMessage
	SegmentFullyParsedServerMessage
	SegmentServerInitMessage
)

type SegmentType int

type RfbSegment struct {
	Bytes              []byte
	SegmentType        SegmentType
	UpcomingObjectType int
	Message            interface{}
}

type SegmentConsumer interface {
	Consume(*RfbSegment) error
}

type RfbReadHelper struct {
	io.Reader
	Listener SegmentConsumer
}

func (r *RfbReadHelper) ReadDiscrete(p []byte) (int, error) {
	return r.Read(p)
}

func (r *RfbReadHelper) SendRectSeparator(upcomingRectType int) error {
	seg := &RfbSegment{SegmentType: SegmentRectSeparator, UpcomingObjectType: upcomingRectType}
	if r.Listener != nil {
		return nil
	}
	return r.Listener.Consume(seg)

}

func (r *RfbReadHelper) SendMessageSeparator(upcomingMessageType ServerMessageType) error {
	seg := &RfbSegment{SegmentType: SegmentMessageSeparator, UpcomingObjectType: int(upcomingMessageType)}
	if r.Listener == nil {
		return nil
	}
	return r.Listener.Consume(seg)
}

func (r *RfbReadHelper) PublishBytes(p []byte) error {
	seg := &RfbSegment{Bytes: p, SegmentType: SegmentBytes}
	if r.Listener == nil {
		return nil
	}
	return r.Listener.Consume(seg)
}

func (r *RfbReadHelper) Read(p []byte) (n int, err error) {
	readLen, err := r.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	//write the bytes to the Listener for further processing
	seg := &RfbSegment{Bytes: p, SegmentType: SegmentBytes}
	if r.Listener == nil {
		return 0, nil
	}
	r.Listener.Consume(seg)
	if err != nil {
		return 0, err
	}
	return readLen, err
}

func (r *RfbReadHelper) ReadBytes(count int) ([]byte, error) {
	buff := make([]byte, count)

	_, err := io.ReadFull(r, buff)
	if err != nil {
		//if err := binary.Read(d.conn, binary.BigEndian, &buff); err != nil {
		return nil, err
	}
	return buff, nil
}

func (r *RfbReadHelper) ReadUint8() (uint8, error) {
	var myUint uint8
	if err := binary.Read(r, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (r *RfbReadHelper) ReadUint16() (uint16, error) {
	var myUint uint16
	if err := binary.Read(r, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (r *RfbReadHelper) ReadUint32() (uint32, error) {
	var myUint uint32
	if err := binary.Read(r, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}
	//fmt.Printf("myUint=%d", myUint)
	return myUint, nil
}
func (r *RfbReadHelper) ReadCompactLen() (int, error) {
	var err error
	part, err := r.ReadUint8()
	//byteCount := 1
	len := uint32(part & 0x7F)
	if (part & 0x80) != 0 {
		part, err = r.ReadUint8()
		//byteCount++
		len |= uint32(part&0x7F) << 7
		if (part & 0x80) != 0 {
			part, err = r.ReadUint8()
			//byteCount++
			len |= uint32(part&0xFF) << 14
		}
	}

	//   for  i := 0; i < byteCount; i++{
	//     rec.writeByte(portion[i]);
	//   }

	return int(len), err
}

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

package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"vncproxy/logger"
)

var TightMinToCompress = 12

const (
	SegmentBytes SegmentType = iota
	SegmentMessageSeparator
	SegmentRectSeparator
	SegmentFullyParsedClientMessage
	SegmentFullyParsedServerMessage
	SegmentServerInitMessage
	SegmentConnectionClosed
)

type SegmentType int

func (seg SegmentType) String() string {
	switch seg {
	case SegmentBytes:
		return "SegmentBytes"
	case SegmentMessageSeparator:
		return "SegmentMessageSeparator"
	case SegmentRectSeparator:
		return "SegmentRectSeparator"
	case SegmentFullyParsedClientMessage:
		return "SegmentFullyParsedClientMessage"
	case SegmentFullyParsedServerMessage:
		return "SegmentFullyParsedServerMessage"
	case SegmentServerInitMessage:
		return "SegmentServerInitMessage"
	case SegmentConnectionClosed:
		return "SegmentConnectionClosed"
	}

	return ""
}

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
	Listeners  *MultiListener
	savedBytes *bytes.Buffer
}

func NewRfbReadHelper(r io.Reader) *RfbReadHelper {
	return &RfbReadHelper{Reader: r, Listeners: &MultiListener{}}
}

func (r *RfbReadHelper) StartByteCollection() {
	r.savedBytes = &bytes.Buffer{}
}

func (r *RfbReadHelper) EndByteCollection() []byte {
	bts := r.savedBytes.Bytes()
	r.savedBytes = nil
	return bts
}

func (r *RfbReadHelper) ReadDiscrete(p []byte) (int, error) {
	return r.Read(p)
}

func (r *RfbReadHelper) SendRectSeparator(upcomingRectType int) error {
	seg := &RfbSegment{SegmentType: SegmentRectSeparator, UpcomingObjectType: upcomingRectType}
	return r.Listeners.Consume(seg)
}

func (r *RfbReadHelper) SendMessageSeparator(upcomingMessageType ServerMessageType) error {
	seg := &RfbSegment{SegmentType: SegmentMessageSeparator, UpcomingObjectType: int(upcomingMessageType)}
	return r.Listeners.Consume(seg)
}

func (r *RfbReadHelper) PublishBytes(p []byte) error {
	seg := &RfbSegment{Bytes: p, SegmentType: SegmentBytes}
	return r.Listeners.Consume(seg)
}

//var prevlen int

func (r *RfbReadHelper) Read(p []byte) (n int, err error) {
	readLen, err := r.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	//if saving up our bytes, write them into the predefined buffer
	if r.savedBytes != nil {
		_, err := r.savedBytes.Write(p)
		if err != nil {
			logger.Warn("RfbReadHelper.Read: failed to collect bytes in mem buffer:", err)
		}
	}
	/////////
	// modLen := (prevlen % 10000)
	// if len(p) == modLen && modLen != prevlen {
	// 	logger.Warn("RFBReadHelper debug!! plen=", prevlen, " len=", len(p))
	// }
	// prevlen = len(p)
	/////////

	logger.Debugf("RfbReadHelper.Read: publishing bytes, bytes:%v", p)

	//write the bytes to the Listener for further processing
	seg := &RfbSegment{Bytes: p, SegmentType: SegmentBytes}
	err = r.Listeners.Consume(seg)
	if err != nil {
		return 0, err
	}

	return readLen, err
}

func (r *RfbReadHelper) ReadBytes(count int) ([]byte, error) {
	buff := make([]byte, count)

	lengthRead, err := io.ReadFull(r, buff)

	//lengthRead, err := r.Read(buff)
	if lengthRead != count {
		logger.Errorf("RfbReadHelper.ReadBytes unable to read bytes: lengthRead=%d, countExpected=%d", lengthRead, count)
		return nil, errors.New("RfbReadHelper.ReadBytes unable to read bytes")
	}

	//err := binary.Read(r, binary.BigEndian, &buff)

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

	return myUint, nil
}
func (r *RfbReadHelper) ReadUint16() (uint16, error) {
	var myUint uint16
	if err := binary.Read(r, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}

	return myUint, nil
}
func (r *RfbReadHelper) ReadUint32() (uint32, error) {
	var myUint uint32
	if err := binary.Read(r, binary.BigEndian, &myUint); err != nil {
		return 0, err
	}

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
	logger.Debugf("RfbReadHelper.ReadTightData: compactlen=%d", zlibDataLen)
	if err != nil {
		return nil, err
	}

	return r.ReadBytes(zlibDataLen)
}

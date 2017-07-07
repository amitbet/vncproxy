package listeners

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"
	"vncproxy/common"
	"vncproxy/server"
)

type Recorder struct {
	//common.BytesListener
	RBSFileName         string
	writer              *os.File
	logger              common.Logger
	startTime           int
	buffer              bytes.Buffer
	serverInitMessage   *common.ServerInit
	sessionStartWritten bool
	segmentChan         chan *common.RfbSegment
}

func getNowMillisec() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

func NewRecorder(saveFilePath string) *Recorder {
	//delete file if it exists
	if _, err := os.Stat(saveFilePath); err == nil {
		os.Remove(saveFilePath)
	}

	rec := Recorder{RBSFileName: saveFilePath, startTime: getNowMillisec()}
	var err error

	rec.writer, err = os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("unable to open file: %s, error: %v", saveFilePath, err)
		return nil
	}

	//buffer the channel so we don't halt the proxying flow for slow writes when under pressure
	rec.segmentChan = make(chan *common.RfbSegment, 100)
	go func() {
		for {
			data := <-rec.segmentChan
			rec.HandleRfbSegment(data)
		}
	}()

	return &rec
}

// func (rec *Recorder) startSession(desktopName string, fbWidth uint16, fbHeight uint16) error {

// 	err := rec.writeStartSession(desktopName, fbWidth, fbHeight)

// 	if err != nil {
// 		fmt.Printf("Recorder was unable to write StartSession to file error: %v", err)
// 		return nil
// 	}

// 	return nil
// }

const versionMsg_3_3 = "RFB 003.003\n"
const versionMsg_3_7 = "RFB 003.007\n"
const versionMsg_3_8 = "RFB 003.008\n"

// Security types
const (
	SecTypeInvalid = 0
	SecTypeNone    = 1
	SecTypeVncAuth = 2
	SecTypeTight   = 16
)

// func (r *Recorder) writeHeader() error {
// 	_, err := r.writer.WriteString("FBS 001.000\n")
// 	return err
// 	// df.write("FBS 001.000\n".getBytes());
// }

func (r *Recorder) writeStartSession(initMsg *common.ServerInit) error {

	desktopName := string(initMsg.NameText)
	framebufferWidth := initMsg.FBWidth
	framebufferHeight := initMsg.FBHeight

	//write rfb header information (the only part done without the [size|data|timestamp] block wrapper)
	r.buffer.WriteString("FBS 001.000\n")
	r.buffer.WriteTo(r.writer)
	r.buffer.Reset()

	//push the version message into the buffer so it will be written in the first rbs block
	r.buffer.WriteString(versionMsg_3_3)

	//push sec type and fb dimensions
	binary.Write(&r.buffer, binary.BigEndian, int32(SecTypeNone))
	binary.Write(&r.buffer, binary.BigEndian, int16(framebufferWidth))
	binary.Write(&r.buffer, binary.BigEndian, int16(framebufferHeight))

	var fbsServerInitMsg = []byte{32, 24, 0, 1, 0, byte(0xFF), 0, byte(0xFF), 0, byte(0xFF), 16, 8, 0, 0, 0, 0}
	r.buffer.Write(fbsServerInitMsg)

	binary.Write(&r.buffer, binary.BigEndian, uint32(len(desktopName)+1))

	r.buffer.WriteString(desktopName)
	binary.Write(&r.buffer, binary.BigEndian, byte(0)) // add null termination for desktop string

	return nil
}

func (r *Recorder) Consume(data *common.RfbSegment) error {

	select {
	case r.segmentChan <- data:
	default:
		fmt.Println("error: recorder queue is full")
	}

	return nil
}

func (r *Recorder) HandleRfbSegment(data *common.RfbSegment) error {

	switch data.SegmentType {
	case common.SegmentMessageSeparator:
		if !r.sessionStartWritten {
			r.writeStartSession(r.serverInitMessage)
		}

		switch common.ServerMessageType(data.UpcomingObjectType) {
		case common.FramebufferUpdate:
			r.writeToDisk()
		case common.SetColourMapEntries:
		case common.Bell:
		case common.ServerCutText:
		default:
			return errors.New("unknown message type:" + string(data.UpcomingObjectType))
		}

	case common.SegmentRectSeparator:
		r.writeToDisk()
	case common.SegmentBytes:
		_, err := r.buffer.Write(data.Bytes)
		return err
	case common.SegmentServerInitMessage:
		r.serverInitMessage = data.Message.(*common.ServerInit)
	case common.SegmentFullyParsedClientMessage:
		clientMsg := data.Message.(common.ClientMessage)

		switch clientMsg.Type() {
		case common.SetPixelFormatMsgType:
			clientMsg := data.Message.(*server.SetPixelFormat)
			r.serverInitMessage.PixelFormat = clientMsg.PF
		default:
			return errors.New("unknown client message type:" + string(data.UpcomingObjectType))
		}

	default:
		//return errors.New("undefined RfbSegment type")
	}
	return nil
}

func (r *Recorder) writeToDisk() error {
	timeSinceStart := getNowMillisec() - r.startTime
	if r.buffer.Len() == 0 {
		return nil
	}

	//write buff length
	bytesLen := r.buffer.Len()
	binary.Write(r.writer, binary.BigEndian, uint32(bytesLen))
	paddedSize := (bytesLen + 3) & 0x7FFFFFFC
	paddingSize := paddedSize - bytesLen

	fmt.Printf("paddedSize=%d paddingSize=%d bytesLen=%d", paddedSize, paddingSize, bytesLen)
	//write buffer padded to 32bit
	_, err := r.buffer.WriteTo(r.writer)
	padding := make([]byte, paddingSize)
	fmt.Printf("padding=%v ", padding)

	binary.Write(r.writer, binary.BigEndian, padding)

	//write timestamp
	binary.Write(r.writer, binary.BigEndian, uint32(timeSinceStart))
	r.buffer.Reset()
	return err
}

// func (r *Recorder) WriteUInt8(data uint8) error {
// 	buf := make([]byte, 1)
// 	buf[0] = byte(data) // cast int8 to byte
// 	return r.Write(buf)
// }

func (r *Recorder) Close() {
	r.writer.Close()
}

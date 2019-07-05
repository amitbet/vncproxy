package vnc_rec

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/amitbet/vncproxy/common"
	"github.com/amitbet/vncproxy/logger"
	"github.com/amitbet/vncproxy/server"
)

type ClientRecorder struct {
	//common.BytesListener
	RBSFileName string
	writer      *os.File
	//logger              common.Logger
	startTime           int
	buffer              bytes.Buffer
	serverInitMessage   *common.ServerInit
	sessionStartWritten bool
	segmentChan         chan *common.RfbSegment
	maxWriteSize        int
}

func NewClientRecorder(saveFilePath string) (*ClientRecorder, error) {
	//delete file if it exists
	if _, err := os.Stat(saveFilePath); err == nil {
		os.Remove(saveFilePath)
	}

	rec := ClientRecorder{RBSFileName: saveFilePath, startTime: getNowMillisec()}
	var err error

	rec.maxWriteSize = 65535

	rec.writer, err = os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("unable to open file: %s, error: %v", saveFilePath, err)
		return nil, err
	}

	//buffer the channel so we don't halt the proxying flow for slow writes when under pressure
	rec.segmentChan = make(chan *common.RfbSegment, 100)
	go func() {
		for {
			data := <-rec.segmentChan
			rec.HandleRfbSegment(data)
		}
	}()

	return &rec, nil
}

// func (r *ClientRecorder) writeHeader() error {
// 	_, err := r.writer.WriteString("FBS 001.000\n")
// 	return err
// 	// df.write("FBS 001.000\n".getBytes());
// }

func (r *ClientRecorder) writeStartSession(initMsg *common.ServerInit) error {
	r.sessionStartWritten = true
	desktopName := string(initMsg.NameText)
	framebufferWidth := initMsg.FBWidth
	framebufferHeight := initMsg.FBHeight

	//write rfb header information (the only part done without the [size|data|timestamp] block wrapper)
	r.writer.WriteString("FBS 001.000\n")

	//push the version message into the buffer so it will be written in the first rbs block
	r.buffer.WriteString(versionMsg_3_3)

	//push sec type and fb dimensions
	binary.Write(&r.buffer, binary.BigEndian, int32(SecTypeNone))
	binary.Write(&r.buffer, binary.BigEndian, int16(framebufferWidth))
	binary.Write(&r.buffer, binary.BigEndian, int16(framebufferHeight))

	buff := bytes.Buffer{}
	//binary.Write(&buff, binary.BigEndian, initMsg.FBWidth)
	//binary.Write(&buff, binary.BigEndian, initMsg.FBHeight)
	binary.Write(&buff, binary.BigEndian, initMsg.PixelFormat)
	buff.Write([]byte{0, 0, 0}) //padding
	r.buffer.Write(buff.Bytes())
	//logger.Debugf(">>>>>>buffer for initMessage:%v ", buff.Bytes())

	//var fbsServerInitMsg = []byte{32, 24, 0, 1, 0, byte(0xFF), 0, byte(0xFF), 0, byte(0xFF), 16, 8, 0, 0, 0, 0}
	//r.buffer.Write(fbsServerInitMsg)

	binary.Write(&r.buffer, binary.BigEndian, uint32(len(desktopName)))

	r.buffer.WriteString(desktopName)
	//binary.Write(&r.buffer, binary.BigEndian, byte(0)) // add null termination for desktop string

	return nil
}

func (r *ClientRecorder) Consume(data *common.RfbSegment) error {
	//using async writes so if chan buffer overflows, proxy will not be affected
	select {
	case r.segmentChan <- data:
		// default:
		// 	logger.Error("error: ClientRecorder queue is full")
	}

	return nil
}

func (r *ClientRecorder) HandleRfbSegment(data *common.RfbSegment) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered in HandleRfbSegment: ", r)
		}
	}()

	switch data.SegmentType {
	case common.SegmentMessageStart:
		if !r.sessionStartWritten {
			logger.Debugf("ClientRecorder.HandleRfbSegment: writing start session segment: %v", r.serverInitMessage)
			r.writeStartSession(r.serverInitMessage)
		}

		switch common.ServerMessageType(data.UpcomingObjectType) {
		case common.FramebufferUpdate:
			logger.Debugf("ClientRecorder.HandleRfbSegment: saving FramebufferUpdate segment")
			//r.writeToDisk()
		case common.SetColourMapEntries:
		case common.Bell:
		case common.ServerCutText:
		default:
			logger.Warn("ClientRecorder.HandleRfbSegment: unknown message type:" + string(data.UpcomingObjectType))
		}
	case common.SegmentConnectionClosed:
		r.writeToDisk()
	case common.SegmentRectSeparator:
		logger.Debugf("ClientRecorder.HandleRfbSegment: not writing rect")
		//r.writeToDisk()
	case common.SegmentBytes:
		logger.Debug("ClientRecorder.HandleRfbSegment: not writing bytes, len:", len(data.Bytes))
		if r.buffer.Len()+len(data.Bytes) > r.maxWriteSize-4 {
			r.writeToDisk()
		}
		// _, err := r.buffer.Write(data.Bytes)
		// return err
	case common.SegmentServerInitMessage:
		r.serverInitMessage = data.Message.(*common.ServerInit)
	case common.SegmentFullyParsedClientMessage:
		clientMsg := data.Message.(common.ClientMessage)

		switch clientMsg.Type() {
		case common.SetPixelFormatMsgType:
			clientMsg := data.Message.(*server.MsgSetPixelFormat)
			logger.Debugf("ClientRecorder.HandleRfbSegment: client message %v", *clientMsg)
			r.serverInitMessage.PixelFormat = clientMsg.PF
		case common.KeyEventMsgType:
			clientMsg := data.Message.(*server.MsgKeyEvent)
			logger.Debug("Recorder.HandleRfbSegment: writing bytes for KeyEventMsgType, len:", *clientMsg)
			clientMsg.Write(r.writer)
		case common.PointerEventMsgType:
			clientMsg := data.Message.(*server.MsgPointerEvent)
			logger.Debug("Recorder.HandleRfbSegment: writing bytes for PointerEventMsgType, len:", *clientMsg)
			clientMsg.Write(r.writer)
		default:
			//return errors.New("unknown client message type:" + string(data.UpcomingObjectType))
		}

	default:
		//return errors.New("undefined RfbSegment type")
	}
	return nil
}

func (r *ClientRecorder) writeToDisk() error {
	timeSinceStart := getNowMillisec() - r.startTime
	if r.buffer.Len() == 0 {
		return nil
	}

	//write buff length
	bytesLen := r.buffer.Len()
	binary.Write(r.writer, binary.BigEndian, uint32(bytesLen))
	paddedSize := (bytesLen + 3) & 0x7FFFFFFC
	paddingSize := paddedSize - bytesLen

	//logger.Debugf("paddedSize=%d paddingSize=%d bytesLen=%d", paddedSize, paddingSize, bytesLen)
	//write buffer padded to 32bit
	_, err := r.buffer.WriteTo(r.writer)
	padding := make([]byte, paddingSize)
	//logger.Debugf("padding=%v ", padding)

	binary.Write(r.writer, binary.BigEndian, padding)

	//write timestamp
	binary.Write(r.writer, binary.BigEndian, uint32(timeSinceStart))
	r.buffer.Reset()
	return err
}

// func (r *ClientRecorder) WriteUInt8(data uint8) error {
// 	buf := make([]byte, 1)
// 	buf[0] = byte(data) // cast int8 to byte
// 	return r.Write(buf)
// }

func (r *ClientRecorder) Close() {
	r.writer.Close()
}

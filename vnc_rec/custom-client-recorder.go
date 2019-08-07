package vnc_rec

import (
	"bytes"
	"os"

	"github.com/amitbet/vncproxy/common"
	"github.com/amitbet/vncproxy/logger"
	"github.com/amitbet/vncproxy/server"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	pb "github.com/sibeshkar/vncproxy/proto"
)

type CustomClientRecorder struct {
	//common.BytesListener
	RBSFileName string
	writer      *os.File
	//logger              common.Logger
	Rectbuffer          bytes.Buffer
	FramebufferUpdate   *pb.FramebufferUpdate
	Rect                *pb.Rectangle
	startTime           int
	buffer              bytes.Buffer
	serverInitMessage   *common.ServerInit
	sessionStartWritten bool
	segmentChan         chan *common.RfbSegment
	maxWriteSize        int
}

func NewCustomClientRecorder(saveFilePath string) (*CustomClientRecorder, error) {
	//delete file if it exists
	if _, err := os.Stat(saveFilePath); err == nil {
		os.Remove(saveFilePath)
	}

	rec := CustomClientRecorder{RBSFileName: saveFilePath, startTime: getNowMillisec()}
	var err error

	rec.maxWriteSize = 65535

	rec.writer, err = os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("unable to open file: %s, error: %v", saveFilePath, err)
		return nil, err
	}

	// rec.writer, err = os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	logger.Errorf("unable to open file: %s, error: %v", saveFilePath, err)
	// 	return nil, err
	// }

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

func (r *CustomClientRecorder) writeStartSession(initMsg *common.ServerInit) error {
	r.sessionStartWritten = true
	desktopName := string(initMsg.NameText)
	framebufferWidth := initMsg.FBWidth
	framebufferHeight := initMsg.FBHeight
	// //write rfb header information (the only part done without the [size|data|timestamp] block wrapper)
	// r.writer.WriteString("FBS 001.000\n")
	// r.demonstration.Initmsg.

	// 	//push the version message into the buffer so it will be written in the first rbs block
	// 	r.buffer.WriteString(versionMsg_3_3)

	// //push sec type and fb dimensions
	// binary.Write(&r.buffer, binary.BigEndian, int32(SecTypeNone))
	// binary.Write(&r.buffer, binary.BigEndian, int16(framebufferWidth))
	// binary.Write(&r.buffer, binary.BigEndian, int16(framebufferHeight))

	// buff := bytes.Buffer{}
	// //binary.Write(&buff, binary.BigEndian, initMsg.FBWidth)
	// //binary.Write(&buff, binary.BigEndian, initMsg.FBHeight)
	// binary.Write(&buff, binary.BigEndian, initMsg.PixelFormat)
	// buff.Write([]byte{0, 0, 0}) //padding
	// r.buffer.Write(buff.Bytes())
	// //logger.Debugf(">>>>>>buffer for initMessage:%v ", buff.Bytes())

	// //var fbsServerInitMsg = []byte{32, 24, 0, 1, 0, byte(0xFF), 0, byte(0xFF), 0, byte(0xFF), 16, 8, 0, 0, 0, 0}
	// //r.buffer.Write(fbsServerInitMsg)

	// binary.Write(&r.buffer, binary.BigEndian, uint32(len(desktopName)))

	// r.buffer.WriteString(desktopName)

	pixel_format := &pb.PixelFormat{
		BPP:        uint32(initMsg.PixelFormat.BPP),
		Depth:      uint32(initMsg.PixelFormat.Depth),
		BigEndian:  uint32(initMsg.PixelFormat.BigEndian),
		TrueColor:  uint32(initMsg.PixelFormat.TrueColor),
		RedMax:     uint32(initMsg.PixelFormat.RedMax),
		GreenMax:   uint32(initMsg.PixelFormat.GreenMax),
		BlueMax:    uint32(initMsg.PixelFormat.BlueMax),
		RedShift:   uint32(initMsg.PixelFormat.RedShift),
		GreenShift: uint32(initMsg.PixelFormat.GreenShift),
		BlueShift:  uint32(initMsg.PixelFormat.BlueShift),
	}

	initMsgProto := &pb.InitMsg{
		RfbHeader:   "FBS 001.000",
		RfbVersion:  versionMsg_3_3,
		FBHeight:    uint32(framebufferHeight),
		FBWidth:     uint32(framebufferWidth),
		SecType:     uint32(SecTypeNone),
		StartTime:   uint32(r.startTime),
		DesktopName: desktopName,
		PixelFormat: pixel_format,
	}

	pbutil.WriteDelimited(r.writer, initMsgProto)
	//binary.Write(&r.buffer, binary.BigEndian, byte(0)) // add null termination for desktop string

	return nil
}

func (r *CustomClientRecorder) Consume(data *common.RfbSegment) error {
	//using async writes so if chan buffer overflows, proxy will not be affected
	select {
	case r.segmentChan <- data:
		// default:
		// 	logger.Error("error: CustomClientRecorder queue is full")
	}

	return nil
}

func (r *CustomClientRecorder) HandleRfbSegment(data *common.RfbSegment) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered in HandleRfbSegment: ", r)
		}
	}()

	timeSinceStart := uint32(getNowMillisec() - r.startTime)

	switch data.SegmentType {
	case common.SegmentMessageStart:
		if !r.sessionStartWritten {
			logger.Debugf("CustomClientRecorder.HandleRfbSegment: writing start session segment: %v", r.serverInitMessage)
			r.writeStartSession(r.serverInitMessage)
		}

		switch common.ServerMessageType(data.UpcomingObjectType) {
		case common.FramebufferUpdate:
			logger.Debugf("CustomClientRecorder.HandleRfbSegment: saving FramebufferUpdate segment")

		case common.SetColourMapEntries:
		case common.Bell:
		case common.ServerCutText:
		default:
			logger.Warn("CustomClientRecorder.HandleRfbSegment: unknown message type:" + string(data.UpcomingObjectType))
		}
	case common.SegmentConnectionClosed:
		logger.Debugf("CustomClientRecorder.HandleRfbSegment: connection closed")
	case common.SegmentRectSeparator:
		logger.Debugf("CustomClientRecorder.HandleRfbSegment: writing rect")
		//r.Rect.Reset()
		//r.writeToDisk()
	case common.SegmentBytes:
		logger.Debug("CustomClientRecorder.HandleRfbSegment: writing bytes, len:", len(data.Bytes))
		// if r.buffer.Len()+len(data.Bytes) > r.maxWriteSize-4 {
		// 	r.writeToDisk()
		// }
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
			//clientMsg.Write(r.writer)
			keyevent := &pb.KeyEvent{
				Down:      uint32(clientMsg.Down),
				Key:       uint32(clientMsg.Key),
				Timestamp: timeSinceStart,
			}
			logger.Debug("CustomClientRecorder.HandleRfbSegment: writing bytes for KeyEventMsgType, len:", *keyevent)

			pbutil.WriteDelimited(r.writer, &pb.MessageType{Type: uint32(4)})
			pbutil.WriteDelimited(r.writer, keyevent)
			//r.demonstration.Keyevents = append(r.demonstration.Keyevents, keyevent)
		case common.PointerEventMsgType:
			clientMsg := data.Message.(*server.MsgPointerEvent)

			//clientMsg.Write(r.writer)
			pointerevent := &pb.PointerEvent{
				Mask:      uint32(clientMsg.Mask),
				X:         uint32(clientMsg.X),
				Y:         uint32(clientMsg.Y),
				Timestamp: timeSinceStart,
			}
			logger.Debug("CustomClientRecorder.HandleRfbSegment: writing bytes for PointerEventMsgType, len:", *pointerevent)
			//r.demonstration.Pointerevents = append(r.demonstration.Pointerevents, pointerevent)
			pbutil.WriteDelimited(r.writer, &pb.MessageType{Type: uint32(5)})
			pbutil.WriteDelimited(r.writer, pointerevent)
		default:
			//return errors.New("unknown client message type:" + string(data.UpcomingObjectType))
		}

	default:
		//return errors.New("undefined RfbSegment type")
	}
	return nil
}

func (r *CustomClientRecorder) Close() {
	r.writer.Close()
}

package player

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
)

type RfbReader struct {
	reader           io.Reader
	buffer           bytes.Buffer
	currentTimestamp int
	pixelFormat      *common.PixelFormat
	encodings        []common.IEncoding
}

/**************************************************************
** RFB File documentation:
** Sections:
** 0. header:
	* index seek position
	*
** 1. init message
** 2. content
	* frame message:
		* size, timestamp, type, content
** 3. index:
	* each frame message start position, full/incremental, timestamp
	*
***************************************************************/

func (rfb *RfbReader) CurrentTimestamp() int {
	return rfb.currentTimestamp
}

func (rfb *RfbReader) Read(p []byte) (n int, err error) {
	if rfb.buffer.Len() < len(p) {
		seg, err := rfb.ReadSegment()

		if err != nil {
			logger.Error("rfbReader.Read: error reading rfbsegment: ", err)
			return 0, err
		}
		rfb.buffer.Write(seg.bytes)
		rfb.currentTimestamp = int(seg.timestamp)
	}
	return rfb.buffer.Read(p)
}

func (rfb *RfbReader) CurrentPixelFormat() *common.PixelFormat { return rfb.pixelFormat }

//func (rfb *rfbReader) CurrentColorMap() *common.ColorMap       { return &common.ColorMap{} }
func (rfb *RfbReader) Encodings() []common.IEncoding { return rfb.encodings }

func NewRfbReader(rfbFile string) (*RfbReader, error) {

	reader, err := os.OpenFile(rfbFile, os.O_RDONLY, 0644)
	if err != nil {
		logger.Error("NewrfbReader: can't open rfb file: ", rfbFile)
		return nil, err
	}
	return &RfbReader{reader: reader,
		encodings: []common.IEncoding{
			&encodings.CopyRectEncoding{},
			&encodings.ZLibEncoding{},
			&encodings.ZRLEEncoding{},
			&encodings.CoRREEncoding{},
			&encodings.HextileEncoding{},
			&encodings.TightEncoding{},
			&encodings.TightPngEncoding{},
			&encodings.EncCursorPseudo{},
			&encodings.RawEncoding{},
			&encodings.RREEncoding{},
		},
	}, nil

}

func (rfb *RfbReader) ReadStartSession() (*common.ServerInit, error) {

	initMsg := common.ServerInit{}
	reader := rfb.reader

	var framebufferWidth uint16
	var framebufferHeight uint16
	var SecTypeNone uint32
	//read rfb header information (the only part done without the [size|data|timestamp] block wrapper)
	//.("rfb 001.000\n")
	bytes := make([]byte, 12)
	_, err := reader.Read(bytes)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init message - rfb file Version:", err)
		return nil, err
	}

	//read the version message into the buffer so it will be written in the first rbs block
	//RFB 003.008\n
	bytes = make([]byte, 12)
	_, err = rfb.Read(bytes)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - RFB Version: ", err)
		return nil, err
	}

	//push sec type and fb dimensions
	binary.Read(rfb, binary.BigEndian, &SecTypeNone)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - SecType: ", err)
	}

	//read frame buffer width, height
	binary.Read(rfb, binary.BigEndian, &framebufferWidth)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - FBWidth: ", err)
		return nil, err
	}
	initMsg.FBWidth = framebufferWidth

	binary.Read(rfb, binary.BigEndian, &framebufferHeight)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - FBHeight: ", err)
		return nil, err
	}
	initMsg.FBHeight = framebufferHeight

	//read pixel format
	pixelFormat := &common.PixelFormat{}
	binary.Read(rfb, binary.BigEndian, pixelFormat)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - Pixelformat: ", err)
		return nil, err
	}
	initMsg.PixelFormat = *pixelFormat
	//read padding
	bytes = make([]byte, 3)
	rfb.Read(bytes)
	rfb.pixelFormat = pixelFormat

	//read desktop name
	var desknameLen uint32
	binary.Read(rfb, binary.BigEndian, &desknameLen)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - deskname Len: ", err)
		return nil, err
	}
	initMsg.NameLength = desknameLen

	bytes = make([]byte, desknameLen)
	rfb.Read(bytes)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: error reading rbs init - desktopName: ", err)
		return nil, err
	}

	initMsg.NameText = bytes

	return &initMsg, nil
}

func (rfb *RfbReader) ReadSegment() (*FbsSegment, error) {
	reader := rfb.reader
	var bytesLen uint32

	//read length
	err := binary.Read(reader, binary.BigEndian, &bytesLen)
	if err != nil {
		logger.Error("rfbReader.ReadStartSession: read len, error reading rbs file: ", err)
		return nil, err
	}

	paddedSize := (bytesLen + 3) & 0x7FFFFFFC

	//read bytes
	bytes := make([]byte, paddedSize)
	_, err = reader.Read(bytes)
	if err != nil {
		logger.Error("rfbReader.ReadSegment: read bytes, error reading rbs file: ", err)
		return nil, err
	}

	//remove padding
	actualBytes := bytes[:bytesLen]

	//read timestamp
	var timeSinceStart uint32
	binary.Read(reader, binary.BigEndian, &timeSinceStart)
	if err != nil {
		logger.Error("rfbReader.ReadSegment: read timestamp, error reading rbs file: ", err)
		return nil, err
	}

	//timeStamp := time.Unix(timeSinceStart, 0)
	seg := &FbsSegment{bytes: actualBytes, timestamp: timeSinceStart}
	return seg, nil
}

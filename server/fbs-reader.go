package server

import (
	"encoding/binary"
	"io"
	"os"
	"vncproxy/common"
	"vncproxy/logger"
)

type FbsReader struct {
	reader io.Reader
}

func NewFbsReader(fbsFile string) (*FbsReader, error) {

	reader, err := os.OpenFile(fbsFile, os.O_RDONLY, 0644)
	if err != nil {
		logger.Error("NewFbsReader: can't open fbs file: ", fbsFile)
		return nil, err
	}
	return &FbsReader{reader: reader}, nil
}

func (player *FbsReader) ReadStartSession() (*common.ServerInit, error) {

	initMsg := common.ServerInit{}
	reader := player.reader

	var framebufferWidth uint16
	var framebufferHeight uint16
	var SecTypeNone uint32
	//read rfb header information (the only part done without the [size|data|timestamp] block wrapper)
	//.("FBS 001.000\n")
	bytes := make([]byte, 12)
	_, err := reader.Read(bytes)
	if err != nil {
		logger.Error("error reading rbs init message - FBS file Version:", err)
		return nil, err
	}

	//read the version message into the buffer so it will be written in the first rbs block
	//RFB 003.008\n
	bytes = make([]byte, 12)
	_, err = reader.Read(bytes)
	if err != nil {
		logger.Error("error reading rbs init - RFB Version: ", err)
		return nil, err
	}

	//push sec type and fb dimensions
	binary.Read(reader, binary.BigEndian, &SecTypeNone)
	if err != nil {
		logger.Error("error reading rbs init - SecType: ", err)
	}

	//read frame buffer width, height
	binary.Read(reader, binary.BigEndian, &framebufferWidth)
	if err != nil {
		logger.Error("error reading rbs init - FBWidth: ", err)
		return nil, err
	}
	initMsg.FBWidth = framebufferWidth

	binary.Read(reader, binary.BigEndian, &framebufferHeight)
	if err != nil {
		logger.Error("error reading rbs init - FBHeight: ", err)
		return nil, err
	}
	initMsg.FBHeight = framebufferHeight

	//read pixel format
	pixelFormat := &common.PixelFormat{}
	binary.Read(reader, binary.BigEndian, pixelFormat)
	if err != nil {
		logger.Error("error reading rbs init - Pixelformat: ", err)
		return nil, err
	}
	initMsg.PixelFormat = *pixelFormat
	//read padding
	bytes = make([]byte, 3)
	reader.Read(bytes)

	//read desktop name
	var desknameLen uint32
	binary.Read(reader, binary.BigEndian, &desknameLen)
	if err != nil {
		logger.Error("error reading rbs init - deskname Len: ", err)
		return nil, err
	}
	initMsg.NameLength = desknameLen

	bytes = make([]byte, desknameLen)
	reader.Read(bytes)
	if err != nil {
		logger.Error("error reading rbs init - desktopName: ", err)
		return nil, err
	}

	initMsg.NameText = bytes

	return &initMsg, nil
}

func (player *FbsReader) ReadSegment() (*FbsSegment, error) {
	reader := player.reader
	var bytesLen uint32

	//read length
	err := binary.Read(reader, binary.BigEndian, &bytesLen)
	if err != nil {
		logger.Error("error reading rbs file: ", err)
		return nil, err
	}

	paddedSize := (bytesLen + 3) & 0x7FFFFFFC

	//read bytes
	bytes := make([]byte, paddedSize)
	_, err = reader.Read(bytes)
	if err != nil {
		logger.Error("error reading rbs file: ", err)
		return nil, err
	}

	//remove padding
	actualBytes := bytes[:bytesLen]

	//read timestamp
	var timeSinceStart uint32
	binary.Read(reader, binary.BigEndian, &timeSinceStart)
	if err != nil {
		logger.Error("error reading rbs file: ", err)
		return nil, err
	}

	//timeStamp := time.Unix(timeSinceStart, 0)
	return &FbsSegment{actualBytes, timeSinceStart}, nil
}

type FbsSegment struct {
	bytes          []byte
	timeSinceStart uint32
}

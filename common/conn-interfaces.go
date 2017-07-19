package common

import "io"

type IServerConn interface {
	io.ReadWriter
	//IServerConn() io.ReadWriter
	Protocol() string
	CurrentPixelFormat() *PixelFormat
	SetPixelFormat(*PixelFormat) error
	//ColorMap() *ColorMap
	SetColorMap(*ColorMap)
	Encodings() []IEncoding
	SetEncodings([]EncodingType) error
	Width() uint16
	Height() uint16
	SetWidth(uint16)
	SetHeight(uint16)
	DesktopName() string
	SetDesktopName(string)
	//Flush() error
	SetProtoVersion(string)
	// Write([]byte) (int, error)
}

type IClientConn interface {
	CurrentPixelFormat() *PixelFormat
	CurrentColorMap() *ColorMap
	Encodings() []IEncoding
}

package common

import (
	"bytes"
	"encoding/binary"
	"io"
)

// An Encoding implements a method for encoding pixel data that is
// sent by the server to the client.
type Encoding interface {
	// The number that uniquely identifies this encoding type.
	Type() int32

	// Read reads the contents of the encoded pixel data from the reader.
	// This should return a new Encoding implementation that contains
	// the proper data.
	Read(*PixelFormat, *Rectangle, *RfbReadHelper) (Encoding, error)
}

// EncodingType represents a known VNC encoding type.
type EncodingType int32

const (
	EncRaw                           EncodingType = 0
	EncCopyRect                      EncodingType = 1
	EncRRE                           EncodingType = 2
	EncCoRRE                         EncodingType = 4
	EncHextile                       EncodingType = 5
	EncZlib                          EncodingType = 6
	EncTight                         EncodingType = 7
	EncZlibHex                       EncodingType = 8
	EncUltra1                        EncodingType = 9
	EncUltra2                        EncodingType = 10
	EncJPEG                          EncodingType = 21
	EncJRLE                          EncodingType = 22
	EncTRLE                          EncodingType = 15
	EncZRLE                          EncodingType = 16
	EncJPEGQualityLevelPseudo10      EncodingType = -23
	EncJPEGQualityLevelPseudo9       EncodingType = -24
	EncJPEGQualityLevelPseudo8       EncodingType = -25
	EncJPEGQualityLevelPseudo7       EncodingType = -26
	EncJPEGQualityLevelPseudo6       EncodingType = -27
	EncJPEGQualityLevelPseudo5       EncodingType = -28
	EncJPEGQualityLevelPseudo4       EncodingType = -29
	EncJPEGQualityLevelPseudo3       EncodingType = -30
	EncJPEGQualityLevelPseudo2       EncodingType = -31
	EncJPEGQualityLevelPseudo1       EncodingType = -32
	EncColorPseudo                   EncodingType = -239
	EncDesktopSizePseudo             EncodingType = -223
	EncLastRectPseudo                EncodingType = -224
	EncCompressionLevel10            EncodingType = -247
	EncCompressionLevel9             EncodingType = -248
	EncCompressionLevel8             EncodingType = -249
	EncCompressionLevel7             EncodingType = -250
	EncCompressionLevel6             EncodingType = -251
	EncCompressionLevel5             EncodingType = -252
	EncCompressionLevel4             EncodingType = -253
	EncCompressionLevel3             EncodingType = -254
	EncCompressionLevel2             EncodingType = -255
	EncCompressionLevel1             EncodingType = -256
	EncQEMUPointerMotionChangePseudo EncodingType = -257
	EncQEMUExtendedKeyEventPseudo    EncodingType = -258
	EncTightPng                      EncodingType = -260
	EncExtendedDesktopSizePseudo     EncodingType = -308
	EncXvpPseudo                     EncodingType = -309
	EncFencePseudo                   EncodingType = -312
	EncContinuousUpdatesPseudo       EncodingType = -313
	EncClientRedirect                EncodingType = -311
	EncTightPNGBase64                EncodingType = 21 + 0x574d5600
	EncTightDiffComp                 EncodingType = 22 + 0x574d5600
	EncVMWDefineCursor               EncodingType = 100 + 0x574d5600
	EncVMWCursorState                EncodingType = 101 + 0x574d5600
	EncVMWCursorPosition             EncodingType = 102 + 0x574d5600
	EncVMWTypematicInfo              EncodingType = 103 + 0x574d5600
	EncVMWLEDState                   EncodingType = 104 + 0x574d5600
	EncVMWServerPush2                EncodingType = 123 + 0x574d5600
	EncVMWServerCaps                 EncodingType = 122 + 0x574d5600
	EncVMWFrameStamp                 EncodingType = 124 + 0x574d5600
	EncOffscreenCopyRect             EncodingType = 126 + 0x574d5600
)

// PixelFormat describes the way a pixel is formatted for a VNC connection.
//
// See RFC 6143 Section 7.4 for information on each of the fields.
type PixelFormat struct {
	BPP        uint8
	Depth      uint8
	BigEndian  bool
	TrueColor  bool
	RedMax     uint16
	GreenMax   uint16
	BlueMax    uint16
	RedShift   uint8
	GreenShift uint8
	BlueShift  uint8
}

func (format *PixelFormat) WriteTo(w io.Writer) error {
	var buf bytes.Buffer

	// Byte 1
	if err := binary.Write(&buf, binary.BigEndian, format.BPP); err != nil {
		return err
	}

	// Byte 2
	if err := binary.Write(&buf, binary.BigEndian, format.Depth); err != nil {
		return err
	}

	var boolByte byte
	if format.BigEndian {
		boolByte = 1
	} else {
		boolByte = 0
	}

	// Byte 3 (BigEndian)
	if err := binary.Write(&buf, binary.BigEndian, boolByte); err != nil {
		return err
	}

	if format.TrueColor {
		boolByte = 1
	} else {
		boolByte = 0
	}

	// Byte 4 (TrueColor)
	if err := binary.Write(&buf, binary.BigEndian, boolByte); err != nil {
		return err
	}

	// If we have true color enabled then we have to fill in the rest of the
	// structure with the color values.
	if format.TrueColor {
		if err := binary.Write(&buf, binary.BigEndian, format.RedMax); err != nil {
			return err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.GreenMax); err != nil {
			return err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.BlueMax); err != nil {
			return err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.RedShift); err != nil {
			return err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.GreenShift); err != nil {
			return err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.BlueShift); err != nil {
			return err
		}
	}

	w.Write(buf.Bytes()[0:16])
	return nil
}

func NewPixelFormat(bpp uint8) *PixelFormat {
	bigEndian := false
	//	rgbMax := uint16(math.Exp2(float64(bpp))) - 1
	rMax := uint16(255)
	gMax := uint16(255)
	bMax := uint16(255)
	var (
		tc         = true
		rs, gs, bs uint8
		depth      uint8
	)
	switch bpp {
	case 8:
		tc = false
		depth = 8
		rs, gs, bs = 0, 0, 0
	case 16:
		depth = 16
		rs, gs, bs = 0, 4, 8
	case 32:
		depth = 24
		//	rs, gs, bs = 0, 8, 16
		rs, gs, bs = 16, 8, 0
	}

	return &PixelFormat{bpp, depth, bigEndian, tc, rMax, gMax, bMax, rs, gs, bs}
}

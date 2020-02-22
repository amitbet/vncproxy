package common

import (
	"bytes"
	"encoding/binary"
	"io"
)

// An IEncoding implements a method for encoding pixel data that is
// sent by the server to the client.
type IEncoding interface {
	// The number that uniquely identifies this encoding type.
	Type() int32
	//WriteTo ...
	WriteTo(w io.Writer) (n int, err error)
	// Read reads the contents of the encoded pixel data from the reader.
	// This should return a new IEncoding implementation that contains
	// the proper data.
	Read(*PixelFormat, *Rectangle, *RfbReadHelper) (IEncoding, error)
}

// EncodingType represents a known VNC encoding type.
type EncodingType int32

func (enct EncodingType) String() string {
	switch enct {
	case EncRaw:
		return "EncRaw"
	case EncCopyRect:
		return "EncCopyRect"
	case EncRRE:
		return "EncRRE"
	case EncCoRRE:
		return "EncCoRRE"
	case EncHextile:
		return "EncHextile"
	case EncZlib:
		return "EncZlib"
	case EncTight:
		return "EncTight"
	case EncZlibHex:
		return "EncZlibHex"
	case EncUltra1:
		return "EncUltra1"
	case EncUltra2:
		return "EncUltra2"
	case EncJPEG:
		return "EncJPEG"
	case EncJRLE:
		return "EncJRLE"
	case EncTRLE:
		return "EncTRLE"
	case EncZRLE:
		return "EncZRLE"
	case EncJPEGQualityLevelPseudo10:
		return "EncJPEGQualityLevelPseudo10"
	case EncJPEGQualityLevelPseudo9:
		return "EncJPEGQualityLevelPseudo9"
	case EncJPEGQualityLevelPseudo8:
		return "EncJPEGQualityLevelPseudo8"
	case EncJPEGQualityLevelPseudo7:
		return "EncJPEGQualityLevelPseudo7"
	case EncJPEGQualityLevelPseudo6:
		return "EncJPEGQualityLevelPseudo6"
	case EncJPEGQualityLevelPseudo5:
		return "EncJPEGQualityLevelPseudo5"
	case EncJPEGQualityLevelPseudo4:
		return "EncJPEGQualityLevelPseudo4"
	case EncJPEGQualityLevelPseudo3:
		return "EncJPEGQualityLevelPseudo3"
	case EncJPEGQualityLevelPseudo2:
		return "EncJPEGQualityLevelPseudo2"
	case EncJPEGQualityLevelPseudo1:
		return "EncJPEGQualityLevelPseudo1"
	case EncCursorPseudo:
		return "EncCursorPseudo"
	case EncLedStatePseudo:
		return "EncLedStatePseudo"
	case EncDesktopSizePseudo:
		return "EncDesktopSizePseudo"
	case EncLastRectPseudo:
		return "EncLastRectPseudo"
	case EncPointerPosPseudo:
		return "EncPointerPosPseudo"
	case EncCompressionLevel10:
		return "EncCompressionLevel10"
	case EncCompressionLevel9:
		return "EncCompressionLevel9"
	case EncCompressionLevel8:
		return "EncCompressionLevel8"
	case EncCompressionLevel7:
		return "EncCompressionLevel7"
	case EncCompressionLevel6:
		return "EncCompressionLevel6"
	case EncCompressionLevel5:
		return "EncCompressionLevel5"
	case EncCompressionLevel4:
		return "EncCompressionLevel4"
	case EncCompressionLevel3:
		return "EncCompressionLevel3"
	case EncCompressionLevel2:
		return "EncCompressionLevel2"
	case EncCompressionLevel1:
		return "EncCompressionLevel1"
	case EncQEMUPointerMotionChangePseudo:
		return "EncQEMUPointerMotionChangePseudo"
	case EncQEMUExtendedKeyEventPseudo:
		return "EncQEMUExtendedKeyEventPseudo"
	case EncTightPng:
		return "EncTightPng"
	case EncExtendedDesktopSizePseudo:
		return "EncExtendedDesktopSizePseudo"
	case EncXvpPseudo:
		return "EncXvpPseudo"
	case EncFencePseudo:
		return "EncFencePseudo"
	case EncContinuousUpdatesPseudo:
		return "EncContinuousUpdatesPseudo"
	case EncClientRedirect:
		return "EncClientRedirect"
	case EncTightPNGBase64:
		return "EncTightPNGBase64"
	case EncTightDiffComp:
		return "EncTightDiffComp"
	case EncVMWDefineCursor:
		return "EncVMWDefineCursor"
	case EncVMWCursorState:
		return "EncVMWCursorState"
	case EncVMWCursorPosition:
		return "EncVMWCursorPosition"
	case EncVMWTypematicInfo:
		return "EncVMWTypematicInfo"
	case EncVMWLEDState:
		return "EncVMWLEDState"
	case EncVMWServerPush2:
		return "EncVMWServerPush2"
	case EncVMWServerCaps:
		return "EncVMWServerCaps"
	case EncVMWFrameStamp:
		return "EncVMWFrameStamp"
	case EncOffscreenCopyRect:
		return "EncOffscreenCopyRect"
	}
	return ""
}

//EncodingType ...
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
	EncCursorPseudo                  EncodingType = -239
	EncDesktopSizePseudo             EncodingType = -223
	EncLastRectPseudo                EncodingType = -224
	EncPointerPosPseudo              EncodingType = -232
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
	EncLedStatePseudo                EncodingType = -261
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
	BigEndian  uint8
	TrueColor  uint8
	RedMax     uint16
	GreenMax   uint16
	BlueMax    uint16
	RedShift   uint8
	GreenShift uint8
	BlueShift  uint8
}

//WriteTo ...
func (format *PixelFormat) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	// Byte 1
	if err := binary.Write(&buf, binary.BigEndian, format.BPP); err != nil {
		return 0, err
	}

	// Byte 2
	if err := binary.Write(&buf, binary.BigEndian, format.Depth); err != nil {
		return 0, err
	}

	var boolByte byte
	if format.BigEndian == 1 {
		boolByte = 1
	} else {
		boolByte = 0
	}

	// Byte 3 (BigEndian)
	if err := binary.Write(&buf, binary.BigEndian, boolByte); err != nil {
		return 0, err
	}

	if format.TrueColor == 1 {
		boolByte = 1
	} else {
		boolByte = 0
	}

	// Byte 4 (TrueColor)
	if err := binary.Write(&buf, binary.BigEndian, boolByte); err != nil {
		return 0, err
	}

	// If we have true color enabled then we have to fill in the rest of the
	// structure with the color values.
	if format.TrueColor == 1 {
		if err := binary.Write(&buf, binary.BigEndian, format.RedMax); err != nil {
			return 0, err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.GreenMax); err != nil {
			return 0, err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.BlueMax); err != nil {
			return 0, err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.RedShift); err != nil {
			return 0, err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.GreenShift); err != nil {
			return 0, err
		}

		if err := binary.Write(&buf, binary.BigEndian, format.BlueShift); err != nil {
			return 0, err
		}
	}

	w.Write(buf.Bytes()[0:16])
	return 0, nil
}

// NewPixelFormat ...
func NewPixelFormat(bpp uint8) *PixelFormat {
	bigEndian := 0
	//	rgbMax := uint16(math.Exp2(float64(bpp))) - 1
	rMax := uint16(255)
	gMax := uint16(255)
	bMax := uint16(255)
	var (
		tc         = 1
		rs, gs, bs uint8
		depth      uint8
	)
	switch bpp {
	case 8:
		tc = 0
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

	return &PixelFormat{bpp, depth, uint8(bigEndian), uint8(tc), rMax, gMax, bMax, rs, gs, bs}
}

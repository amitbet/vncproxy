package encodings

import (
	"bytes"
	"errors"
	"io"
	"vncproxy/common"
	"vncproxy/logger"
)

var TightMinToCompress int = 12

const (
	TightExplicitFilter = 0x04
	TightFill           = 0x08
	TightJpeg           = 0x09
	TightPNG            = 0x10

	TightFilterCopy     = 0x00
	TightFilterPalette  = 0x01
	TightFilterGradient = 0x02
)

type TightEncoding struct {
	bytes []byte
}

func (*TightEncoding) Type() int32 { return int32(common.EncTight) }

func calcTightBytePerPixel(pf *common.PixelFormat) int {
	bytesPerPixel := int(pf.BPP / 8)

	var bytesPerPixelTight int
	if 24 == pf.Depth && 32 == pf.BPP {
		bytesPerPixelTight = 3
	} else {
		bytesPerPixelTight = bytesPerPixel
	}
	return bytesPerPixelTight
}

func (z *TightEncoding) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(z.bytes)
}

func StoreBytes(bytes *bytes.Buffer, data []byte) {
	_, err := bytes.Write(data)
	if err != nil {
		logger.Error("Error in encoding while saving bytes: ", err)
	}
}

func (t *TightEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.IEncoding, error) {
	bytesPixel := calcTightBytePerPixel(pixelFmt)

	r.StartByteCollection()
	defer func() {
		t.bytes = r.EndByteCollection()
	}()

	compctl, err := r.ReadUint8()

	if err != nil {
		logger.Debugf("error in handling tight encoding: %v\n", err)
		return nil, err
	}
	logger.Debugf("bytesPixel= %d, subencoding= %d\n", bytesPixel, compctl)

	//move it to position (remove zlib flush commands)
	compType := compctl >> 4 & 0x0F

	logger.Debugf("afterSHL:%d\n", compType)
	switch compType {
	case TightFill:
		logger.Debugf("reading fill size=%d\n", bytesPixel)
		//read color
		_, err := r.ReadBytes(int(bytesPixel))
		if err != nil {
			logger.Debugf("error in handling tight encoding: %v\n", err)
			return nil, err
		}

		return t, nil
	case TightJpeg:
		if pixelFmt.BPP == 8 {
			return nil, errors.New("Tight encoding: JPEG is not supported in 8 bpp mode")
		}

		len, err := r.ReadCompactLen()

		if err != nil {
			return nil, err
		}
		logger.Debugf("reading jpeg, size=%d\n", len)
		_, err = r.ReadBytes(len)
		if err != nil {
			return nil, err
		}

		return t, nil
	default:

		if compType > TightJpeg {
			logger.Debug("Compression control byte is incorrect!")
		}

		handleTightFilters(compctl, pixelFmt, rect, r)

		return t, nil
	}
}

func handleTightFilters(subencoding uint8, pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) {

	var FILTER_ID_MASK uint8 = 0x40

	var filterid uint8
	var err error

	if (subencoding & FILTER_ID_MASK) > 0 { // filter byte presence
		filterid, err = r.ReadUint8()

		if err != nil {
			logger.Debugf("error in handling tight encoding, reading filterid: %v\n", err)
			return
		}
		logger.Debugf("read filter: %d\n", filterid)
	}

	bytesPixel := calcTightBytePerPixel(pixelFmt)

	logger.Debugf("filter: %d\n", filterid)

	lengthCurrentbpp := int(bytesPixel) * int(rect.Width) * int(rect.Height)

	switch filterid {
	case TightFilterPalette: //PALETTE_FILTER

		colorCount, err := r.ReadUint8()
		paletteSize := colorCount + 1 // add one more
		logger.Debugf("----PALETTE_FILTER: paletteSize=%d bytesPixel=%d\n", paletteSize, bytesPixel)
		//complete palette
		_, err = r.ReadBytes(int(paletteSize) * bytesPixel)

		var dataLength int
		if paletteSize == 2 {
			dataLength = int(rect.Height) * ((int(rect.Width) + 7) / 8)
		} else {
			dataLength = int(rect.Width * rect.Height)
		}
		_, err = r.ReadTightData(dataLength)
		if err != nil {
			logger.Debugf("error in handling tight encoding, Reading Palette: %v\n", err)
			return
		}

	case TightFilterGradient: //GRADIENT_FILTER
		logger.Debugf("----GRADIENT_FILTER: bytesPixel=%d\n", bytesPixel)
		logger.Debugf("usegrad: %d\n", filterid)
		_, err := r.ReadTightData(lengthCurrentbpp)
		if err != nil {
			logger.Debugf("error in handling tight encoding, Reading GRADIENT_FILTER: %v\n", err)
			return
		}

	case TightFilterCopy: //BASIC_FILTER
		logger.Debugf("----BASIC_FILTER: bytesPixel=%d\n", bytesPixel)
		_, err := r.ReadTightData(lengthCurrentbpp)
		if err != nil {
			logger.Debugf("error in handling tight encoding, Reading BASIC_FILTER: %v\n", err)
			return
		}

	default:
		logger.Debugf("Bad tight filter id: %d\n", filterid)
		return
	}

	return
}

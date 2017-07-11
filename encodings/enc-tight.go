package encodings

import (
	"errors"
	"fmt"
	"vncproxy/common"
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
	//output io.Writer
	//logger common.Logger
}

// func (t *TightEncoding) SetOutput(output io.Writer) {
// 	t.output = output
// }

func (*TightEncoding) Type() int32 { return int32(common.EncTight) }

// func ReadAndRecBytes(conn io.Reader, rec io.Writer, count int) ([]byte, error) {
// 	buf, err := readBytes(conn, count)
// 	rec.Write(buf)
// 	return buf, err
// }
// func ReadAndRecUint8(conn io.Reader, rec io.Writer) (uint8, error) {
// 	myUint, err := readUint8(conn)
// 	buf := make([]byte, 1)
// 	buf[0] = byte(myUint) // cast int8 to byte
// 	rec.Write(buf)
// 	return myUint, err
// }

// func ReadAndRecUint16(conn io.Reader, rec io.Writer) (uint16, error) {
// 	myUint, err := readUint16(conn)
// 	buf := make([]byte, 2)
// 	//buf[0] = byte(myUint) // cast int8 to byte
// 	//var i int16 = 41
// 	//b := make([]byte, 2)
// 	binary.LittleEndian.PutUint16(buf, uint16(myUint))

// 	rec.Write(buf)
// 	return myUint, err
// }

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

func (t *TightEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	bytesPixel := calcTightBytePerPixel(pixelFmt)
	//conn := common.RfbReadHelper{Reader:reader}
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}

	//var subencoding uint8
	compctl, err := r.ReadUint8()
	if err != nil {
		fmt.Printf("error in handling tight encoding: %v\n", err)
		return nil, err
	}
	fmt.Printf("bytesPixel= %d, subencoding= %d\n", bytesPixel, compctl)

	//move it to position (remove zlib flush commands)
	compType := compctl >> 4 & 0x0F

	fmt.Printf("afterSHL:%d\n", compType)
	switch compType {
	case TightFill:
		fmt.Printf("reading fill size=%d\n", bytesPixel)
		//read color
		r.ReadBytes(int(bytesPixel))
		//byt, _ := r.ReadBytes(3)
		//fmt.Printf(">>>>>>>>>TightFillBytes=%v", byt)
		return t, nil
	case TightJpeg:
		if pixelFmt.BPP == 8 {
			return nil, errors.New("Tight encoding: JPEG is not supported in 8 bpp mode")
		}

		len, err := r.ReadCompactLen()
		if err != nil {
			return nil, err
		}
		fmt.Printf("reading jpeg size=%d\n", len)
		r.ReadBytes(len)
		return t, nil
	default:

		if compType > TightJpeg {
			fmt.Println("Compression control byte is incorrect!")
		}

		handleTightFilters(compctl, pixelFmt, rect, r)
		return t, nil
	}
}

func handleTightFilters(subencoding uint8, pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) {
	//conn := common.RfbReadHelper{Reader:reader}
	var FILTER_ID_MASK uint8 = 0x40
	//var STREAM_ID_MASK uint8 = 0x30

	//decoderId := (subencoding & STREAM_ID_MASK) >> 4
	var filterid uint8
	var err error

	if (subencoding & FILTER_ID_MASK) > 0 { // filter byte presence
		filterid, err = r.ReadUint8()
		if err != nil {
			fmt.Printf("error in handling tight encoding, reading filterid: %v\n", err)
			return
		}
		fmt.Printf("read filter: %d\n", filterid)
	}

	//var numColors uint8
	bytesPixel := calcTightBytePerPixel(pixelFmt)

	fmt.Printf("filter: %d\n", filterid)
	// if rfb.rec != null {
	// 	rfb.rec.writeByte(filter_id)
	// }
	lengthCurrentbpp := int(bytesPixel) * int(rect.Width) * int(rect.Height)

	switch filterid {
	case TightFilterPalette: //PALETTE_FILTER

		colorCount, err := r.ReadUint8()
		paletteSize := colorCount + 1 // add one more
		fmt.Printf("----PALETTE_FILTER: paletteSize=%d bytesPixel=%d\n", paletteSize, bytesPixel)
		//complete palette
		r.ReadBytes(int(paletteSize) * bytesPixel)

		var dataLength int
		if paletteSize == 2 {
			dataLength = int(rect.Height) * ((int(rect.Width) + 7) / 8)
		} else {
			dataLength = int(rect.Width * rect.Height)
		}
		_, err = r.ReadTightData(dataLength)
		if err != nil {
			fmt.Printf("error in handling tight encoding, Reading Palette: %v\n", err)
			return
		}
	case TightFilterGradient: //GRADIENT_FILTER
		fmt.Printf("----GRADIENT_FILTER: bytesPixel=%d\n", bytesPixel)
		//useGradient = true
		fmt.Printf("usegrad: %d\n", filterid)
		r.ReadTightData(lengthCurrentbpp)
	case TightFilterCopy: //BASIC_FILTER
		fmt.Printf("----BASIC_FILTER: bytesPixel=%d\n", bytesPixel)
		r.ReadTightData(lengthCurrentbpp)
	default:
		fmt.Printf("Bad tight filter id: %d\n", filterid)
		return
	}

	////////////

	// if numColors == 0 && bytesPixel == 4 {
	// 	rowSize1 *= 3
	// }
	// rowSize := (int(rect.Width)*bitsPixel + 7) / 8
	// dataSize := int(rect.Height) * rowSize

	// dataSize1 := rect.Height * rowSize1
	// fmt.Printf("datasize: %d, origDatasize: %d", dataSize, dataSize1)
	// // Read, optionally uncompress and decode data.
	// if int(dataSize1) < TightMinToCompress {
	// 	// Data size is small - not compressed with zlib.
	// 	if numColors != 0 {
	// 		// Indexed colors.
	// 		//indexedData := make([]byte, dataSize)
	// 		readBytes(conn.c, int(dataSize1))
	// 		//readFully(indexedData);
	// 		// if (rfb.rec != null) {
	// 		//   rfb.rec.write(indexedData);
	// 		// }
	// 		// if (numColors == 2) {
	// 		//   // Two colors.
	// 		//   if (bytesPixel == 1) {
	// 		//     decodeMonoData(x, y, w, h, indexedData, palette8);
	// 		//   } else {
	// 		//     decodeMonoData(x, y, w, h, indexedData, palette24);
	// 		//   }
	// 		// } else {
	// 		//   // 3..255 colors (assuming bytesPixel == 4).
	// 		//   int i = 0;
	// 		//   for (int dy = y; dy < y + h; dy++) {
	// 		//     for (int dx = x; dx < x + w; dx++) {
	// 		//       pixels24[dy * rfb.framebufferWidth + dx] = palette24[indexedData[i++] & 0xFF];
	// 		//     }
	// 		//   }
	// 		// }
	// 	} else if useGradient {
	// 		// "Gradient"-processed data
	// 		//buf := make ( []byte,w * h * 3);
	// 		dataByteCount := int(3) * int(rect.Width) * int(rect.Height)
	// 		readBytes(conn.c, dataByteCount)
	// 		// rfb.readFully(buf);
	// 		// if (rfb.rec != null) {
	// 		//   rfb.rec.write(buf);
	// 		// }
	// 		// decodeGradientData(x, y, w, h, buf);
	// 	} else {
	// 		// Raw truecolor data.
	// 		dataByteCount := int(bytesPixel) * int(rect.Width) * int(rect.Height)
	// 		readBytes(conn.c, dataByteCount)

	// 		// if (bytesPixel == 1) {
	// 		//   for (int dy = y; dy < y + h; dy++) {

	// 		//     rfb.readFully(pixels8, dy * rfb.framebufferWidth + x, w);
	// 		//     if (rfb.rec != null) {
	// 		//       rfb.rec.write(pixels8, dy * rfb.framebufferWidth + x, w);
	// 		//     }
	// 		//   }
	// 		// } else {
	// 		//   byte[] buf = new byte[w * 3];
	// 		//   int i, offset;
	// 		//   for (int dy = y; dy < y + h; dy++) {
	// 		//     rfb.readFully(buf);
	// 		//     if (rfb.rec != null) {
	// 		//       rfb.rec.write(buf);
	// 		//     }
	// 		//     offset = dy * rfb.framebufferWidth + x;
	// 		//     for (i = 0; i < w; i++) {
	// 		//       pixels24[offset + i] = (buf[i * 3] & 0xFF) << 16 | (buf[i * 3 + 1] & 0xFF) << 8 | (buf[i * 3 + 2] & 0xFF);
	// 		//     }
	// 		//   }
	// 		// }
	// 	}
	// } else {
	// 	// Data was compressed with zlib.
	// 	zlibDataLen, err := readCompactLen(conn.c)
	// 	fmt.Printf("compactlen=%d\n", zlibDataLen)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	//byte[] zlibData = new byte[zlibDataLen];
	// 	//rfb.readFully(zlibData);
	// 	readBytes(conn.c, zlibDataLen)

	// 	//   if (rfb.rec != null) {
	// 	//     rfb.rec.write(zlibData);
	// 	//   }
	// 	//   int stream_id = comp_ctl & 0x03;
	// 	//   if (tightInflaters[stream_id] == null) {
	// 	//     tightInflaters[stream_id] = new Inflater();
	// 	//   }
	// 	//   Inflater myInflater = tightInflaters[stream_id];
	// 	//   myInflater.setInput(zlibData);
	// 	//   byte[] buf = new byte[dataSize];
	// 	//   myInflater.inflate(buf);
	// 	//   if (rfb.rec != null && !rfb.recordFromBeginning) {
	// 	//     rfb.recordCompressedData(buf);
	// 	//   }

	// 	//   if (numColors != 0) {
	// 	//     // Indexed colors.
	// 	//     if (numColors == 2) {
	// 	//       // Two colors.
	// 	//       if (bytesPixel == 1) {
	// 	//         decodeMonoData(x, y, w, h, buf, palette8);
	// 	//       } else {
	// 	//         decodeMonoData(x, y, w, h, buf, palette24);
	// 	//       }
	// 	//     } else {
	// 	//       // More than two colors (assuming bytesPixel == 4).
	// 	//       int i = 0;
	// 	//       for (int dy = y; dy < y + h; dy++) {
	// 	//         for (int dx = x; dx < x + w; dx++) {
	// 	//           pixels24[dy * rfb.framebufferWidth + dx] = palette24[buf[i++] & 0xFF];
	// 	//         }
	// 	//       }
	// 	//     }
	// 	//   } else if (useGradient) {
	// 	//     // Compressed "Gradient"-filtered data (assuming bytesPixel == 4).
	// 	//     decodeGradientData(x, y, w, h, buf);
	// 	//   } else {
	// 	//     // Compressed truecolor data.
	// 	//     if (bytesPixel == 1) {
	// 	//       int destOffset = y * rfb.framebufferWidth + x;
	// 	//       for (int dy = 0; dy < h; dy++) {
	// 	//         System.arraycopy(buf, dy * w, pixels8, destOffset, w);
	// 	//         destOffset += rfb.framebufferWidth;
	// 	//       }
	// 	//     } else {
	// 	//       int srcOffset = 0;
	// 	//       int destOffset, i;
	// 	//       for (int dy = 0; dy < h; dy++) {
	// 	//         myInflater.inflate(buf);
	// 	//         destOffset = (y + dy) * rfb.framebufferWidth + x;
	// 	//         for (i = 0; i < w; i++) {
	// 	//           pixels24[destOffset + i] = (buf[srcOffset] & 0xFF) << 16 | (buf[srcOffset + 1] & 0xFF) << 8
	// 	//               | (buf[srcOffset + 2] & 0xFF);
	// 	//           srcOffset += 3;
	// 	//         }
	// 	//       }
	// 	//     }
	// 	//   }
	// }

	return
}

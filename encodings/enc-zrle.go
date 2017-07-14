package encodings

import (
	"bytes"
	"encoding/binary"
	"io"
	"vncproxy/common"
)

type ZRLEEncoding struct {
	//Colors []Color
	bytes []byte
}

func (z *ZRLEEncoding) Type() int32 {
	return 16
}

func (z *ZRLEEncoding) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(z.bytes)
}

func (z *ZRLEEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	//conn := common.RfbReadHelper{Reader: r}
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	//bytesPerPixel := c.PixelFormat.BPP / 8
	bytes := &bytes.Buffer{}
	len, err := r.ReadUint32()
	if err != nil {
		return nil, err
	}

	binary.Write(bytes, binary.BigEndian, len)
	bts, err := r.ReadBytes(int(len))
	if err != nil {
		return nil, err
	}
	StoreBytes(bytes, bts)
	z.bytes = bytes.Bytes()
	return z, nil
}

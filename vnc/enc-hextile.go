package vnc

import "io"

const (
	HextileRaw                 = 1
	HextileBackgroundSpecified = 2
	HextileForegroundSpecified = 4
	HextileAnySubrects         = 8
	HextileSubrectsColoured    = 16
)

type HextileEncoding struct {
	Colors []Color
}

func (z *HextileEncoding) Type() int32 {
	return 5
}
func (z *HextileEncoding) Read(conn *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	bytesPerPixel := int(conn.PixelFormat.BPP) / 8
	//buf := make([]byte, bytesPerPixel)
	for ty := rect.Y; ty < rect.Y+rect.Height; ty += 16 {
		th := 16
		if rect.Y+rect.Height-ty < 16 {
			th = int(rect.Y) + int(rect.Height) - int(ty)
		}

		for tx := rect.X; tx < rect.X+rect.Width; tx += 16 {
			tw := 16
			if rect.X+rect.Width-tx < 16 {
				tw = int(rect.X) + int(rect.Width) - int(tx)
			}

			//handle Hextile Subrect(tx, ty, tw, th):
			subencoding, err := conn.readUint8()
			//fmt.Printf("hextile reader tile: (%d,%d) subenc=%d\n", ty, tx, subencoding)
			if err != nil {
				//fmt.Printf("error in hextile reader: %v\n", err)
				return nil, err
			}

			if (subencoding & HextileRaw) != 0 {
				//ReadRawRect(c, rect, r)
				conn.readBytes(tw * th * bytesPerPixel)
				//fmt.Printf("hextile reader: HextileRaw\n")
				continue
			}
			if (subencoding & HextileBackgroundSpecified) != 0 {
				conn.readBytes(int(bytesPerPixel))
			}
			if (subencoding & HextileForegroundSpecified) != 0 {
				conn.readBytes(int(bytesPerPixel))
			}
			if (subencoding & HextileAnySubrects) == 0 {
				//fmt.Printf("hextile reader: no Subrects\n")
				continue
			}
			//fmt.Printf("hextile reader: handling Subrects\n")
			nSubrects, err := conn.readUint8()
			if err != nil {
				return nil, err
			}
			bufsize := int(nSubrects) * 2
			if (subencoding & HextileSubrectsColoured) != 0 {
				bufsize += int(nSubrects) * int(bytesPerPixel)
			}
			//byte[] buf = new byte[bufsize];
			conn.readBytes(bufsize)
		}
	}

	// len, _ := readUint32(c.c)
	// _, err := readBytes(c.c, int(len))

	// if err != nil {
	// 	return nil, err
	// }
	return z, nil
}

package server

import (
	"encoding/binary"
	"fmt"
	"github.com/amitbet/vncproxy/common"

	"io"
	"github.com/amitbet/vncproxy/logger"
)

//ProtoVersionLength ...
const ProtoVersionLength = 12

//ProtoVersion ...
const (
	ProtoVersionUnknown = ""
	ProtoVersion33      = "RFB 003.003\n"
	ProtoVersion38      = "RFB 003.008\n"
)

//ParseProtoVersion ...
func ParseProtoVersion(pv []byte) (uint, uint, error) {
	var major, minor uint

	if len(pv) < ProtoVersionLength {
		return 0, 0, fmt.Errorf("ProtocolVersion message too short (%v < %v)", len(pv), ProtoVersionLength)
	}

	l, err := fmt.Sscanf(string(pv), "RFB %d.%d\n", &major, &minor)
	if l != 2 {
		return 0, 0, fmt.Errorf("error parsing ProtocolVersion")
	}
	if err != nil {
		return 0, 0, err
	}

	return major, minor, nil
}

//ServerVersionHandler ...
func ServerVersionHandler(cfg *ServerConfig, c *ServerConn) error {
	var version [ProtoVersionLength]byte
	if err := binary.Write(c, binary.BigEndian, []byte(ProtoVersion38)); err != nil {
		return err
	}
	// if err := c.Flush(); err != nil {
	// 	return err
	// }
	if err := binary.Read(c, binary.BigEndian, &version); err != nil {
		return err
	}

	major, minor, err := ParseProtoVersion(version[:])
	if err != nil {
		return err
	}

	pv := ProtoVersionUnknown
	if major == 3 {
		if minor >= 8 {
			pv = ProtoVersion38
		} else if minor >= 3 {
			pv = ProtoVersion33
		}
	}
	if pv == ProtoVersionUnknown {
		return fmt.Errorf("ProtocolVersion handshake failed; unsupported version '%v'", string(version[:]))
	}

	c.SetProtoVersion(pv)
	return nil
}

func ServerSecurityHandler(cfg *ServerConfig, c *ServerConn) error {
	if err := binary.Write(c, binary.BigEndian, uint8(len(cfg.SecurityHandlers))); err != nil {
		return err
	}

	for _, sectype := range cfg.SecurityHandlers {
		if err := binary.Write(c, binary.BigEndian, sectype.Type()); err != nil {
			return err
		}
	}

	// if err := c.Flush(); err != nil {
	// 	return err
	// }

	var secType SecurityType
	if err := binary.Read(c, binary.BigEndian, &secType); err != nil {
		return err
	}

	secTypes := make(map[SecurityType]SecurityHandler)
	for _, sType := range cfg.SecurityHandlers {
		secTypes[sType.Type()] = sType
	}

	sType, ok := secTypes[secType]
	if !ok {
		return fmt.Errorf("server type %d not implemented", secType)
	}

	var authCode uint32
	authErr := sType.Auth(c)
	if authErr != nil {
		authCode = uint32(1)
	}

	if err := binary.Write(c, binary.BigEndian, authCode); err != nil {
		return err
	}
	// if err := c.Flush(); err != nil {
	// 	return err
	// }

	if authErr != nil {
		if err := binary.Write(c, binary.BigEndian, len(authErr.Error())); err != nil {
			return err
		}
		if err := binary.Write(c, binary.BigEndian, []byte(authErr.Error())); err != nil {
			return err
		}
		// if err := c.Flush(); err != nil {
		// 	return err
		// }
		return authErr
	}

	return nil
}

func ServerServerInitHandler(cfg *ServerConfig, c *ServerConn) error {
	srvInit := &common.ServerInit{
		FBWidth:     c.Width(),
		FBHeight:    c.Height(),
		PixelFormat: *c.CurrentPixelFormat(),
		NameLength:  uint32(len(cfg.DesktopName)),
		NameText:    []byte(cfg.DesktopName),
	}
	logger.Debugf("Server.ServerServerInitHandler initMessage: %v", srvInit)
	if err := binary.Write(c, binary.BigEndian, srvInit.FBWidth); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, srvInit.FBHeight); err != nil {
		return err
	}

	if _, err := srvInit.PixelFormat.WriteTo(c); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, srvInit.NameLength); err != nil {
		return err
	}

	if err := binary.Write(c, binary.BigEndian, srvInit.NameText); err != nil {
		return err
	}
	//
	//serverCaps:=[]TightCapability{
	//	TightCapability{uint32(1), [4]byte(StandardVendor), [8]byte("12345678")},
	//}
	//clientCaps:=[]TightCapability{
	//	TightCapability{uint32(1), [4]byte(StandardVendor), [8]byte("12345678")},
	//}
	//encodingCaps:=[]TightCapability{
	//	TightCapability{uint32(1), [4]byte(StandardVendor), [8]byte("12345678")},
	//}
	//
	//tightInit:=TightServerInit{
	//	serverCaps,clientCaps,encodingCaps,
	//}
	//tightInit.WriteTo(c)

	return nil
}

const (
	StandardVendor  = "STDV"
	TridiaVncVendor = "TRDV"
	TightVncVendor  = "TGHT"
)

/*
  void initCapabilities() {
    tunnelCaps    = new CapsContainer();
    authCaps      = new CapsContainer();
    serverMsgCaps = new CapsContainer();
    clientMsgCaps = new CapsContainer();
    encodingCaps  = new CapsContainer();

    // Supported authentication methods
    authCaps.add(AuthNone, StandardVendor, SigAuthNone,
		 "No authentication");
    authCaps.add(AuthVNC, StandardVendor, SigAuthVNC,
		 "Standard VNC password authentication");

    // Supported non-standard server-to-client messages
    // [NONE]

    // Supported non-standard client-to-server messages
    // [NONE]

    // Supported encoding types
    encodingCaps.add(EncodingCopyRect, StandardVendor,
		     SigEncodingCopyRect, "Standard CopyRect encoding");
    encodingCaps.add(EncodingRRE, StandardVendor,
		     SigEncodingRRE, "Standard RRE encoding");
    encodingCaps.add(EncodingCoRRE, StandardVendor,
		     SigEncodingCoRRE, "Standard CoRRE encoding");
    encodingCaps.add(EncodingHextile, StandardVendor,
		     SigEncodingHextile, "Standard Hextile encoding");
    encodingCaps.add(EncodingZRLE, StandardVendor,
		     SigEncodingZRLE, "Standard ZRLE encoding");
    encodingCaps.add(EncodingZlib, TridiaVncVendor,
		     SigEncodingZlib, "Zlib encoding");
    encodingCaps.add(EncodingTight, TightVncVendor,
		     SigEncodingTight, "Tight encoding");

    // Supported pseudo-encoding types
    encodingCaps.add(EncodingCompressLevel0, TightVncVendor,
		     SigEncodingCompressLevel0, "Compression level");
    encodingCaps.add(EncodingQualityLevel0, TightVncVendor,
		     SigEncodingQualityLevel0, "JPEG quality level");
    encodingCaps.add(EncodingXCursor, TightVncVendor,
		     SigEncodingXCursor, "X-style cursor shape update");
    encodingCaps.add(EncodingRichCursor, TightVncVendor,
		     SigEncodingRichCursor, "Rich-color cursor shape update");
    encodingCaps.add(EncodingPointerPos, TightVncVendor,
		     SigEncodingPointerPos, "Pointer position update");
    encodingCaps.add(EncodingLastRect, TightVncVendor,
		     SigEncodingLastRect, "LastRect protocol extension");
    encodingCaps.add(EncodingNewFBSize, TightVncVendor,
		     SigEncodingNewFBSize, "Framebuffer size change");
  }
*/
type TightServerInit struct {
	ServerMessageCaps []TightCapability
	ClientMessageCaps []TightCapability
	EncodingCaps      []TightCapability
}

func (t *TightServerInit) ReadFrom(r io.Reader) error {
	var numSrvCaps uint16
	var numCliCaps uint16
	var numEncCaps uint16
	var padding uint16

	if err := binary.Read(r, binary.BigEndian, &numSrvCaps); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &numCliCaps); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &numEncCaps); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &padding); err != nil {
		return err
	}

	for i := 0; i < int(numSrvCaps); i++ {
		cap := TightCapability{}
		cap.ReadFrom(r)
		t.ServerMessageCaps = append(t.ServerMessageCaps, cap)
	}

	for i := 0; i < int(numCliCaps); i++ {
		cap := TightCapability{}
		cap.ReadFrom(r)
		t.ClientMessageCaps = append(t.ClientMessageCaps, cap)
	}

	for i := 0; i < int(numEncCaps); i++ {
		cap := TightCapability{}
		cap.ReadFrom(r)
		t.EncodingCaps = append(t.EncodingCaps, cap)
	}
	return nil
}

func (t *TightServerInit) WriteTo(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, uint16(len(t.ServerMessageCaps))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(len(t.ClientMessageCaps))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, uint16(len(t.EncodingCaps))); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint16(0)); err != nil {
		return err
	}

	for _, s := range t.ServerMessageCaps {
		s.WriteTo(w)
	}
	for _, s := range t.ClientMessageCaps {
		s.WriteTo(w)
	}
	for _, s := range t.EncodingCaps {
		s.WriteTo(w)
	}
	return nil
}

type TightCapability struct {
	code   uint32
	vendor [4]byte
	name   [8]byte
}

func (t *TightCapability) WriteTo(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, t.code); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, t.vendor); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, t.name); err != nil {
		return err
	}
	return nil
}

func (t *TightCapability) ReadFrom(r io.Reader) error {

	if err := binary.Read(r, binary.BigEndian, &t.code); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &t.vendor); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &t.name); err != nil {
		return err
	}
	return nil
}

func ServerClientInitHandler(cfg *ServerConfig, c *ServerConn) error {
	var shared uint8
	if err := binary.Read(c, binary.BigEndian, &shared); err != nil {
		return err
	}
	/* TODO
	if shared != 1 {
		c.SetShared(false)
	}
	*/
	return nil
}

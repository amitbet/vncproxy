package server

import (
	"bytes"
	"crypto/des"
	"crypto/rand"
	"errors"
	"log"
	"vncproxy/common"
)

type SecurityType uint8

const (
	SecTypeUnknown  = SecurityType(0)
	SecTypeNone     = SecurityType(1)
	SecTypeVNC      = SecurityType(2)
	SecTypeVeNCrypt = SecurityType(19)
)

type SecuritySubType uint32

const (
	SecSubTypeUnknown = SecuritySubType(0)
)

const (
	SecSubTypeVeNCrypt01Unknown   = SecuritySubType(0)
	SecSubTypeVeNCrypt01Plain     = SecuritySubType(19)
	SecSubTypeVeNCrypt01TLSNone   = SecuritySubType(20)
	SecSubTypeVeNCrypt01TLSVNC    = SecuritySubType(21)
	SecSubTypeVeNCrypt01TLSPlain  = SecuritySubType(22)
	SecSubTypeVeNCrypt01X509None  = SecuritySubType(23)
	SecSubTypeVeNCrypt01X509VNC   = SecuritySubType(24)
	SecSubTypeVeNCrypt01X509Plain = SecuritySubType(25)
)

const (
	SecSubTypeVeNCrypt02Unknown   = SecuritySubType(0)
	SecSubTypeVeNCrypt02Plain     = SecuritySubType(256)
	SecSubTypeVeNCrypt02TLSNone   = SecuritySubType(257)
	SecSubTypeVeNCrypt02TLSVNC    = SecuritySubType(258)
	SecSubTypeVeNCrypt02TLSPlain  = SecuritySubType(259)
	SecSubTypeVeNCrypt02X509None  = SecuritySubType(260)
	SecSubTypeVeNCrypt02X509VNC   = SecuritySubType(261)
	SecSubTypeVeNCrypt02X509Plain = SecuritySubType(262)
)

type SecurityHandler interface {
	Type() SecurityType
	SubType() SecuritySubType
	Auth(common.ServerConn) error
}

// type ClientAuthNone struct{}

// func (*ClientAuthNone) Type() SecurityType {
// 	return SecTypeNone
// }

// func (*ClientAuthNone) SubType() SecuritySubType {
// 	return SecSubTypeUnknown
// }

// func (*ClientAuthNone) Auth(conn common.ServerConn) error {
// 	return nil
// }

// ServerAuthNone is the "none" authentication. See 7.2.1.
type ServerAuthNone struct{}

func (*ServerAuthNone) Type() SecurityType {
	return SecTypeNone
}

func (*ServerAuthNone) Auth(c common.ServerConn) error {
	return nil
}

func (*ServerAuthNone) SubType() SecuritySubType {
	return SecSubTypeUnknown
}

// func (*ClientAuthVeNCrypt02Plain) Type() SecurityType {
// 	return SecTypeVeNCrypt
// }

// func (*ClientAuthVeNCrypt02Plain) SubType() SecuritySubType {
// 	return SecSubTypeVeNCrypt02Plain
// }

// // ClientAuthVeNCryptPlain see https://www.berrange.com/~dan/vencrypt.txt
// type ClientAuthVeNCrypt02Plain struct {
// 	Username []byte
// 	Password []byte
// }

// func (auth *ClientAuthVeNCrypt02Plain) Auth(c common.ServerConn) error {
// 	if err := binary.Write(c, binary.BigEndian, []uint8{0, 2}); err != nil {
// 		return err
// 	}
// 	if err := c.Flush(); err != nil {
// 		return err
// 	}
// 	var (
// 		major, minor uint8
// 	)

// 	if err := binary.Read(c, binary.BigEndian, &major); err != nil {
// 		return err
// 	}
// 	if err := binary.Read(c, binary.BigEndian, &minor); err != nil {
// 		return err
// 	}
// 	res := uint8(1)
// 	if major == 0 && minor == 2 {
// 		res = uint8(0)
// 	}
// 	if err := binary.Write(c, binary.BigEndian, res); err != nil {
// 		return err
// 	}
// 	c.Flush()
// 	if err := binary.Write(c, binary.BigEndian, uint8(1)); err != nil {
// 		return err
// 	}
// 	if err := binary.Write(c, binary.BigEndian, auth.SubType()); err != nil {
// 		return err
// 	}
// 	if err := c.Flush(); err != nil {
// 		return err
// 	}
// 	var secType SecuritySubType
// 	if err := binary.Read(c, binary.BigEndian, &secType); err != nil {
// 		return err
// 	}
// 	if secType != auth.SubType() {
// 		binary.Write(c, binary.BigEndian, uint8(1))
// 		c.Flush()
// 		return fmt.Errorf("invalid sectype")
// 	}
// 	if len(auth.Password) == 0 || len(auth.Username) == 0 {
// 		return fmt.Errorf("Security Handshake failed; no username and/or password provided for VeNCryptAuth.")
// 	}
// 	/*
// 		if err := binary.Write(c, binary.BigEndian, uint32(len(auth.Username))); err != nil {
// 			return err
// 		}

// 		if err := binary.Write(c, binary.BigEndian, uint32(len(auth.Password))); err != nil {
// 			return err
// 		}

// 		if err := binary.Write(c, binary.BigEndian, auth.Username); err != nil {
// 			return err
// 		}

// 		if err := binary.Write(c, binary.BigEndian, auth.Password); err != nil {
// 			return err
// 		}
// 	*/
// 	var (
// 		uLength, pLength uint32
// 	)
// 	if err := binary.Read(c, binary.BigEndian, &uLength); err != nil {
// 		return err
// 	}
// 	if err := binary.Read(c, binary.BigEndian, &pLength); err != nil {
// 		return err
// 	}

// 	username := make([]byte, uLength)
// 	password := make([]byte, pLength)
// 	if err := binary.Read(c, binary.BigEndian, &username); err != nil {
// 		return err
// 	}

// 	if err := binary.Read(c, binary.BigEndian, &password); err != nil {
// 		return err
// 	}
// 	if !bytes.Equal(auth.Username, username) || !bytes.Equal(auth.Password, password) {
// 		return fmt.Errorf("invalid username/password")
// 	}
// 	return nil
// }

// ServerAuthVNC is the standard password authentication. See 7.2.2.
type ServerAuthVNC struct {
	Pass string
}

func (*ServerAuthVNC) Type() SecurityType {
	return SecTypeVNC
}

func (*ServerAuthVNC) SubType() SecuritySubType {
	return SecSubTypeUnknown
}

const AUTH_FAIL = "Authentication Failure"

func (auth *ServerAuthVNC) Auth(c common.ServerConn) error {
	buf := make([]byte, 8+len([]byte(AUTH_FAIL)))
	rand.Read(buf[:16]) // Random 16 bytes in buf
	sndsz, err := c.Write(buf[:16])
	if err != nil {
		log.Printf("Error sending challenge to client: %s\n", err.Error())
		return errors.New("Error sending challenge to client:" + err.Error())
	}
	if sndsz != 16 {
		log.Printf("The full 16 byte challenge was not sent!\n")
		return errors.New("The full 16 byte challenge was not sent")
	}
	//c.Flush()
	buf2 := make([]byte, 16)
	_, err = c.Read(buf2)
	if err != nil {
		log.Printf("The authentication result was not read: %s\n", err.Error())
		return errors.New("The authentication result was not read" + err.Error())
	}
	AuthText := auth.Pass
	bk, err := des.NewCipher([]byte(fixDesKey(AuthText)))
	if err != nil {
		log.Printf("Error generating authentication cipher: %s\n", err.Error())
		return errors.New("Error generating authentication cipher")
	}
	buf3 := make([]byte, 16)
	bk.Encrypt(buf3, buf)               //Encrypt first 8 bytes
	bk.Encrypt(buf3[8:], buf[8:])       // Encrypt second 8 bytes
	if bytes.Compare(buf2, buf3) != 0 { // If the result does not decrypt correctly to what we sent then a problem
		SetUint32(buf, 0, 1)
		SetUint32(buf, 4, uint32(len([]byte(AUTH_FAIL))))
		copy(buf[8:], []byte(AUTH_FAIL))
		c.Write(buf)
		//c.Flush()
		return errors.New("Authentication failed")
	}
	return nil
}

// SetUint32 set 4 bytes at pos in buf to the val (in big endian format)
// A test is done to ensure there are 4 bytes available at pos in the buffer
func SetUint32(buf []byte, pos int, val uint32) {
	if pos+4 > len(buf) {
		return
	}
	for i := 0; i < 4; i++ {
		buf[3-i+pos] = byte(val)
		val >>= 8
	}
}

// fixDesKeyByte is used to mirror a byte's bits
// This is not clearly indicated by the document, but is in actual fact used
func fixDesKeyByte(val byte) byte {
	var newval byte = 0
	for i := 0; i < 8; i++ {
		newval <<= 1
		newval += (val & 1)
		val >>= 1
	}
	return newval
}

// fixDesKey will make sure that exactly 8 bytes is used either by truncating or padding with nulls
// The bytes are then bit mirrored and returned
func fixDesKey(key string) []byte {
	tmp := []byte(key)
	buf := make([]byte, 8)
	if len(tmp) <= 8 {
		copy(buf, tmp)
	} else {
		copy(buf, tmp[:8])
	}
	for i := 0; i < 8; i++ {
		buf[i] = fixDesKeyByte(buf[i])
	}
	return buf
}

// // ClientAuthVNC is the standard password authentication. See 7.2.2.
// type ClientAuthVNC struct {
// 	Challenge [16]byte
// 	Password  []byte
// }

// func (*ClientAuthVNC) Type() SecurityType {
// 	return SecTypeVNC
// }
// func (*ClientAuthVNC) SubType() SecuritySubType {
// 	return SecSubTypeUnknown
// }

// func (auth *ClientAuthVNC) Auth(c common.ServerConn) error {
// 	if len(auth.Password) == 0 {
// 		return fmt.Errorf("Security Handshake failed; no password provided for VNCAuth.")
// 	}

// 	if err := binary.Read(c, binary.BigEndian, auth.Challenge); err != nil {
// 		return err
// 	}

// 	auth.encode()

// 	// Send the encrypted challenge back to server
// 	if err := binary.Write(c, binary.BigEndian, auth.Challenge); err != nil {
// 		return err
// 	}

// 	return c.Flush()
// }

// func (auth *ClientAuthVNC) encode() error {
// 	// Copy password string to 8 byte 0-padded slice
// 	key := make([]byte, 8)
// 	copy(key, auth.Password)

// 	// Each byte of the password needs to be reversed. This is a
// 	// non RFC-documented behaviour of VNC clients and servers
// 	for i := range key {
// 		key[i] = (key[i]&0x55)<<1 | (key[i]&0xAA)>>1 // Swap adjacent bits
// 		key[i] = (key[i]&0x33)<<2 | (key[i]&0xCC)>>2 // Swap adjacent pairs
// 		key[i] = (key[i]&0x0F)<<4 | (key[i]&0xF0)>>4 // Swap the 2 halves
// 	}

// 	// Encrypt challenge with key.
// 	cipher, err := des.NewCipher(key)
// 	if err != nil {
// 		return err
// 	}
// 	for i := 0; i < len(auth.Challenge); i += cipher.BlockSize() {
// 		cipher.Encrypt(auth.Challenge[i:i+cipher.BlockSize()], auth.Challenge[i:i+cipher.BlockSize()])
// 	}

// 	return nil
// }

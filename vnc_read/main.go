package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/sibeshkar/vncproxy/logger"
	pb "github.com/sibeshkar/vncproxy/proto"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]

	reader, err := os.OpenFile(fname, os.O_RDWR, 0644)
	if err != nil {
		logger.Errorf("unable to open file: %s, error: %v", fname, err)

	}

	// [START unmarshal_proto]
	// Read the existing address book.
	// in, err := ioutil.ReadFile(fname)
	// if err != nil {
	// 	log.Fatalln("Error reading file:", err)
	// }
	// demonstration := &pb.Demonstration{}
	// if err := proto.Unmarshal(in, demonstration); err != nil {
	// 	log.Fatalln("Failed to parse demonstration file:", err)
	// }

	// pf := &pb.PixelFormat{}
	// pbutil.ReadDelimited(reader, pf)

	// fmt.Println(pf)

	initMsg := &pb.InitMsg{}
	pbutil.ReadDelimited(reader, initMsg)

	fmt.Printf("FBHeight: %v \n", initMsg.GetFBHeight())
	fmt.Printf("FBWidth: %v \n", initMsg.GetFBWidth())
	fmt.Printf("RfbHeader: %v \n", initMsg.GetRfbHeader())
	fmt.Printf("RfbVersion: %v \n", initMsg.GetRfbVersion())
	fmt.Printf("SecType: %v \n", initMsg.GetSecType())
	fmt.Printf("StartTime: %v \n", initMsg.GetStartTime())
	fmt.Printf("DesktopName: %v \n", initMsg.GetDesktopName())
	fmt.Printf("PixelFormat: %v \n", initMsg.GetPixelFormat())

	i := 0

	for {

		// msgType := &pb.MessageType{}
		// pbutil.ReadDelimited(reader, msgType)
		// if msgType.GetType() == uint32(4) {
		// 	keyEvent := &pb.KeyEvent{}
		// 	pbutil.ReadDelimited(reader, keyEvent)
		// 	fmt.Printf("Key event is %v", keyEvent)

		// } else if msgType.GetType() == uint32(5) {
		// 	pointerEvent := &pb.PointerEvent{}
		// 	pbutil.ReadDelimited(reader, pointerEvent)
		// 	fmt.Println("Pointer event is ", pointerEvent)

		// }
		fbupdate := &pb.FramebufferUpdate{}
		pbutil.ReadDelimited(reader, fbupdate)
		writeFbupdate(fbupdate, &i)
		time.Sleep(1 * time.Second)
	}

	//listPeople(os.Stdout, demonstration)

}

func listPeople(w io.Writer, demo *pb.Demonstration) {
	fmt.Printf("FBHeight: %v \n", demo.Initmsg.GetFBHeight())
	fmt.Printf("FBWidth: %v \n", demo.Initmsg.GetFBWidth())
	fmt.Printf("RfbHeader: %v \n", demo.Initmsg.GetRfbHeader())
	fmt.Printf("RfbVersion: %v \n", demo.Initmsg.GetRfbVersion())
	fmt.Printf("SecType: %v \n", demo.Initmsg.GetSecType())
	fmt.Printf("StartTime: %v \n", demo.Initmsg.GetStartTime())
	fmt.Printf("DesktopName: %v \n", demo.Initmsg.GetDesktopName())
	fmt.Printf("PixelFormat: %v \n", demo.Initmsg.GetPixelFormat())

	// for _, p := range demo.Segments {
	// 	writeSegment(w, p)
	// }

	// for _, k := range demo.Keyevents {
	// 	writeKeyEvent(w, k)
	// }

	// for _, p := range demo.Pointerevents {
	// 	writePointerEvent(w, p)
	// }

	// i := 0

	// for _, p := range demo.Fbupdates {
	// 	writeFbupdate(w, p, &i)
	// }

	// fmt.Println(i)

}

func writeSegment(w io.Writer, p *pb.FbsSegment) {
	fmt.Printf("Length: %v Timestamp: %v \n", len(p.GetBytes()), p.GetTimestamp())

}

func writeFbupdate(p *pb.FramebufferUpdate, i *int) {

	*i++
	fmt.Printf("----------FRAMEBUFFERUPDATE NUMBER %v -------------- \n", *i)
	for _, r := range p.Rectangles {
		fmt.Printf("(%d,%d) (width: %d, height: %d), Enc= %d , Bytes: %v \n", r.X, r.Y, r.Width, r.Height, r.Enc, len(r.Bytes))
	}

}

func writeKeyEvent(w io.Writer, p *pb.KeyEvent) {

	fmt.Printf(" KeyEvent : Down : %v, Key: %v) \n", p.GetDown(), p.GetKey())

}

func writePointerEvent(w io.Writer, p *pb.PointerEvent) {

	fmt.Printf(" PointerEvent : X : %v, Y: %v , Mask: %v \n", p.GetX(), p.GetY(), p.GetMask())

}

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/sibeshkar/vncproxy/proto"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]

	// [START unmarshal_proto]
	// Read the existing address book.
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	demonstration := &pb.Demonstration{}
	if err := proto.Unmarshal(in, demonstration); err != nil {
		log.Fatalln("Failed to parse demonstration file:", err)
	}

	listPeople(os.Stdout, demonstration)

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

	i := 0

	for _, p := range demo.Fbupdates {
		writeFbupdate(w, p, &i)
	}

	fmt.Println(i)

}

func writeSegment(w io.Writer, p *pb.FbsSegment) {
	fmt.Printf("Length: %v Timestamp: %v \n", len(p.GetBytes()), p.GetTimestamp())

}

func writeFbupdate(w io.Writer, p *pb.FramebufferUpdate, i *int) {

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

package recorder

import (
	"time"
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/logger"
)

type RfbRequester struct {
	Conn                   *client.ClientConn
	Name                   string
	Width                  uint16
	Height                 uint16
	lastRequestTime        time.Time
	nextFullScreenRefresh  time.Time
	FullScreenRefreshInSec int // refresh interval (creates keyframes) if 0, disables keyframe creation
}

func (p *RfbRequester) Consume(seg *common.RfbSegment) error {

	logger.Debugf("WriteTo.Consume ("+p.Name+"): got segment type=%s", seg.SegmentType)
	switch seg.SegmentType {
	case common.SegmentServerInitMessage:
		serverInitMessage := seg.Message.(*common.ServerInit)
		p.Conn.FrameBufferHeight = serverInitMessage.FBHeight
		p.Conn.FrameBufferWidth = serverInitMessage.FBWidth
		p.Conn.DesktopName = string(serverInitMessage.NameText)
		p.Conn.SetPixelFormat(&serverInitMessage.PixelFormat)
		p.Width = serverInitMessage.FBWidth
		p.Height = serverInitMessage.FBHeight
		p.lastRequestTime = time.Now()
		p.Conn.FramebufferUpdateRequest(false, 0, 0, p.Width, p.Height)
		p.nextFullScreenRefresh = time.Now().Add(time.Duration(p.FullScreenRefreshInSec) * time.Second)

	case common.SegmentMessageStart:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
	case common.SegmentFullyParsedClientMessage:
	case common.SegmentMessageEnd:
		// minTimeBetweenReq := 300 * time.Millisecond
		// timeForNextReq := p.lastRequestTime.Unix() + minTimeBetweenReq.Nanoseconds()/1000
		// if seg.UpcomingObjectType == int(common.FramebufferUpdate) && time.Now().Unix() > timeForNextReq {
		//time.Sleep(300 * time.Millisecond)
		p.lastRequestTime = time.Now()
		incremental := true

		if p.FullScreenRefreshInSec > 0 {
			// if p.nextFullScreenRefresh.IsZero() {
			// 	p.nextFullScreenRefresh = time.Now().Add(time.Duration(p.FullScreenRefreshInSec) * time.Second)
			// }
			if time.Now().Sub(p.nextFullScreenRefresh) <= 0 {
				logger.Warn(">>Creating keyframe")
				p.nextFullScreenRefresh = time.Now().Add(time.Duration(p.FullScreenRefreshInSec) * time.Second)
				incremental = false
			}
		}
		p.Conn.FramebufferUpdateRequest(incremental, 0, 0, p.Width, p.Height)
		//}
	default:
	}
	return nil
}

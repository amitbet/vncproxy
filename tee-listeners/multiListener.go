package listeners

import "vncproxy/common"

type MultiListener struct {
	listeners []common.SegmentConsumer
}

func (m *MultiListener) AddListener(listener common.SegmentConsumer) {
	m.listeners = append(m.listeners, listener)
}

func (m *MultiListener) Consume(seg *common.RfbSegment) error {
	for _, li := range m.listeners {
		//fmt.Println(li)
		err := li.Consume(seg)
		if err != nil {
			return err
		}
	}
	return nil
}

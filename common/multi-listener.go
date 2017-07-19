package common

type MultiListener struct {
	listeners []SegmentConsumer
}

func (m *MultiListener) AddListener(listener SegmentConsumer) {
	m.listeners = append(m.listeners, listener)
}

func (m *MultiListener) Consume(seg *RfbSegment) error {
	for _, li := range m.listeners {

		err := li.Consume(seg)
		if err != nil {
			return err
		}
	}
	return nil
}

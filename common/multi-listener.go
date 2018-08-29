package common

//MultiListener ...
type MultiListener struct {
	listeners []SegmentConsumer
}

//AddListener ...
func (m *MultiListener) AddListener(listener SegmentConsumer) {
	m.listeners = append(m.listeners, listener)
}

//Consume ...
func (m *MultiListener) Consume(seg *RfbSegment) error {
	for _, li := range m.listeners {

		err := li.Consume(seg)
		if err != nil {
			return err
		}
	}
	return nil
}

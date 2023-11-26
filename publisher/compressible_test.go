package publisher

type MockPayload struct {
	err     error
	payload []byte
}

func (c *MockPayload) Bytes() ([]byte, error) {
	return c.payload, c.err
}

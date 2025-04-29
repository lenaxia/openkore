// Package testing provides testing utilities for the send package
package testing

// MockConnection is a mock connection for testing
type MockConnection struct {
	Sent []byte
}

// Send sends data to the mock connection
func (c *MockConnection) Send(data []byte) error {
	c.Sent = data
	return nil
}

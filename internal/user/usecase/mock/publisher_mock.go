package mock

type MockPublisher struct {
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{}
}

func (m *MockPublisher) Publish(subj string, t any) error {
	return nil
}

package func05_test

import (
	func05 "test/func_05"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByID(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func TestGetUserName(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockRepo.On("FindByID", 1).Return("Alice", nil)

	result := func05.GetUserName(mockRepo, 1)
	assert.Equal(t, "Alice", result)

	mockRepo.AssertExpectations(t)
}

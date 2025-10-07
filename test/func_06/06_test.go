package func06_test

import (
	"errors"
	func06 "test/func_06"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ----- Mock -----

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByID(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

// ----- Suite -----

type UserTestSuite struct {
	suite.Suite
	repo *MockUserRepo
}

func (s *UserTestSuite) SetupTest() {
	s.repo = new(MockUserRepo)
}

func (s *UserTestSuite) TestGetUserName_Success() {
	s.repo.On("FindByID", 1).Return("Alice", nil)

	result := func06.GetUsername(s.repo, 1)
	assert.Equal(s.T(), "Alice", result)
	s.repo.AssertExpectations(s.T())
}

func (s *UserTestSuite) TestGetUserName_Fail() {
	s.repo.On("FindByID", 2).Return("", errors.New("not found"))

	result := func06.GetUsername(s.repo, 2)
	assert.Equal(s.T(), "unknown", result)
	s.repo.AssertExpectations(s.T())
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

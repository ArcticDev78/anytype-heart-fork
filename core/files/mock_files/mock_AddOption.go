// Code generated by mockery. DO NOT EDIT.

package mock_files

import (
	files "github.com/anyproto/anytype-heart/core/files"
	mock "github.com/stretchr/testify/mock"
)

// MockAddOption is an autogenerated mock type for the AddOption type
type MockAddOption struct {
	mock.Mock
}

type MockAddOption_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAddOption) EXPECT() *MockAddOption_Expecter {
	return &MockAddOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockAddOption) Execute(_a0 *files.AddOptions) {
	_m.Called(_a0)
}

// MockAddOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockAddOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *files.AddOptions
func (_e *MockAddOption_Expecter) Execute(_a0 interface{}) *MockAddOption_Execute_Call {
	return &MockAddOption_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockAddOption_Execute_Call) Run(run func(_a0 *files.AddOptions)) *MockAddOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*files.AddOptions))
	})
	return _c
}

func (_c *MockAddOption_Execute_Call) Return() *MockAddOption_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockAddOption_Execute_Call) RunAndReturn(run func(*files.AddOptions)) *MockAddOption_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAddOption creates a new instance of MockAddOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAddOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAddOption {
	mock := &MockAddOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

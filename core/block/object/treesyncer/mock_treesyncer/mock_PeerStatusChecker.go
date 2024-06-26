// Code generated by mockery. DO NOT EDIT.

package mock_treesyncer

import (
	app "github.com/anyproto/any-sync/app"
	mock "github.com/stretchr/testify/mock"
)

// MockPeerStatusChecker is an autogenerated mock type for the PeerStatusChecker type
type MockPeerStatusChecker struct {
	mock.Mock
}

type MockPeerStatusChecker_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPeerStatusChecker) EXPECT() *MockPeerStatusChecker_Expecter {
	return &MockPeerStatusChecker_Expecter{mock: &_m.Mock}
}

// Init provides a mock function with given fields: a
func (_m *MockPeerStatusChecker) Init(a *app.App) error {
	ret := _m.Called(a)

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*app.App) error); ok {
		r0 = rf(a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPeerStatusChecker_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type MockPeerStatusChecker_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - a *app.App
func (_e *MockPeerStatusChecker_Expecter) Init(a interface{}) *MockPeerStatusChecker_Init_Call {
	return &MockPeerStatusChecker_Init_Call{Call: _e.mock.On("Init", a)}
}

func (_c *MockPeerStatusChecker_Init_Call) Run(run func(a *app.App)) *MockPeerStatusChecker_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*app.App))
	})
	return _c
}

func (_c *MockPeerStatusChecker_Init_Call) Return(err error) *MockPeerStatusChecker_Init_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockPeerStatusChecker_Init_Call) RunAndReturn(run func(*app.App) error) *MockPeerStatusChecker_Init_Call {
	_c.Call.Return(run)
	return _c
}

// IsPeerOffline provides a mock function with given fields: peerId
func (_m *MockPeerStatusChecker) IsPeerOffline(peerId string) bool {
	ret := _m.Called(peerId)

	if len(ret) == 0 {
		panic("no return value specified for IsPeerOffline")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(peerId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockPeerStatusChecker_IsPeerOffline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsPeerOffline'
type MockPeerStatusChecker_IsPeerOffline_Call struct {
	*mock.Call
}

// IsPeerOffline is a helper method to define mock.On call
//   - peerId string
func (_e *MockPeerStatusChecker_Expecter) IsPeerOffline(peerId interface{}) *MockPeerStatusChecker_IsPeerOffline_Call {
	return &MockPeerStatusChecker_IsPeerOffline_Call{Call: _e.mock.On("IsPeerOffline", peerId)}
}

func (_c *MockPeerStatusChecker_IsPeerOffline_Call) Run(run func(peerId string)) *MockPeerStatusChecker_IsPeerOffline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPeerStatusChecker_IsPeerOffline_Call) Return(_a0 bool) *MockPeerStatusChecker_IsPeerOffline_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPeerStatusChecker_IsPeerOffline_Call) RunAndReturn(run func(string) bool) *MockPeerStatusChecker_IsPeerOffline_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *MockPeerStatusChecker) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockPeerStatusChecker_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockPeerStatusChecker_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockPeerStatusChecker_Expecter) Name() *MockPeerStatusChecker_Name_Call {
	return &MockPeerStatusChecker_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockPeerStatusChecker_Name_Call) Run(run func()) *MockPeerStatusChecker_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPeerStatusChecker_Name_Call) Return(name string) *MockPeerStatusChecker_Name_Call {
	_c.Call.Return(name)
	return _c
}

func (_c *MockPeerStatusChecker_Name_Call) RunAndReturn(run func() string) *MockPeerStatusChecker_Name_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPeerStatusChecker creates a new instance of MockPeerStatusChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPeerStatusChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPeerStatusChecker {
	mock := &MockPeerStatusChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery. DO NOT EDIT.

package mock_idresolver

import (
	app "github.com/anyproto/any-sync/app"

	mock "github.com/stretchr/testify/mock"
)

// MockResolver is an autogenerated mock type for the Resolver type
type MockResolver struct {
	mock.Mock
}

type MockResolver_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResolver) EXPECT() *MockResolver_Expecter {
	return &MockResolver_Expecter{mock: &_m.Mock}
}

// Init provides a mock function with given fields: a
func (_m *MockResolver) Init(a *app.App) error {
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

// MockResolver_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type MockResolver_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - a *app.App
func (_e *MockResolver_Expecter) Init(a interface{}) *MockResolver_Init_Call {
	return &MockResolver_Init_Call{Call: _e.mock.On("Init", a)}
}

func (_c *MockResolver_Init_Call) Run(run func(a *app.App)) *MockResolver_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*app.App))
	})
	return _c
}

func (_c *MockResolver_Init_Call) Return(err error) *MockResolver_Init_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockResolver_Init_Call) RunAndReturn(run func(*app.App) error) *MockResolver_Init_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *MockResolver) Name() string {
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

// MockResolver_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockResolver_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockResolver_Expecter) Name() *MockResolver_Name_Call {
	return &MockResolver_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockResolver_Name_Call) Run(run func()) *MockResolver_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockResolver_Name_Call) Return(name string) *MockResolver_Name_Call {
	_c.Call.Return(name)
	return _c
}

func (_c *MockResolver_Name_Call) RunAndReturn(run func() string) *MockResolver_Name_Call {
	_c.Call.Return(run)
	return _c
}

// ResolveSpaceID provides a mock function with given fields: objectID
func (_m *MockResolver) ResolveSpaceID(objectID string) (string, error) {
	ret := _m.Called(objectID)

	if len(ret) == 0 {
		panic("no return value specified for ResolveSpaceID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(objectID)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(objectID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockResolver_ResolveSpaceID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResolveSpaceID'
type MockResolver_ResolveSpaceID_Call struct {
	*mock.Call
}

// ResolveSpaceID is a helper method to define mock.On call
//   - objectID string
func (_e *MockResolver_Expecter) ResolveSpaceID(objectID interface{}) *MockResolver_ResolveSpaceID_Call {
	return &MockResolver_ResolveSpaceID_Call{Call: _e.mock.On("ResolveSpaceID", objectID)}
}

func (_c *MockResolver_ResolveSpaceID_Call) Run(run func(objectID string)) *MockResolver_ResolveSpaceID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockResolver_ResolveSpaceID_Call) Return(_a0 string, _a1 error) *MockResolver_ResolveSpaceID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockResolver_ResolveSpaceID_Call) RunAndReturn(run func(string) (string, error)) *MockResolver_ResolveSpaceID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockResolver creates a new instance of MockResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockResolver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockResolver {
	mock := &MockResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
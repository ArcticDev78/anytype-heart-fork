// Code generated by mockery. DO NOT EDIT.

package mock_linkpreview

import (
	context "context"

	app "github.com/anyproto/any-sync/app"

	mock "github.com/stretchr/testify/mock"

	model "github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

// MockLinkPreview is an autogenerated mock type for the LinkPreview type
type MockLinkPreview struct {
	mock.Mock
}

type MockLinkPreview_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLinkPreview) EXPECT() *MockLinkPreview_Expecter {
	return &MockLinkPreview_Expecter{mock: &_m.Mock}
}

// Fetch provides a mock function with given fields: ctx, url
func (_m *MockLinkPreview) Fetch(ctx context.Context, url string) (model.LinkPreview, []byte, bool, error) {
	ret := _m.Called(ctx, url)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 model.LinkPreview
	var r1 []byte
	var r2 bool
	var r3 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.LinkPreview, []byte, bool, error)); ok {
		return rf(ctx, url)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.LinkPreview); ok {
		r0 = rf(ctx, url)
	} else {
		r0 = ret.Get(0).(model.LinkPreview)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) []byte); ok {
		r1 = rf(ctx, url)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) bool); ok {
		r2 = rf(ctx, url)
	} else {
		r2 = ret.Get(2).(bool)
	}

	if rf, ok := ret.Get(3).(func(context.Context, string) error); ok {
		r3 = rf(ctx, url)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// MockLinkPreview_Fetch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fetch'
type MockLinkPreview_Fetch_Call struct {
	*mock.Call
}

// Fetch is a helper method to define mock.On call
//   - ctx context.Context
//   - url string
func (_e *MockLinkPreview_Expecter) Fetch(ctx interface{}, url interface{}) *MockLinkPreview_Fetch_Call {
	return &MockLinkPreview_Fetch_Call{Call: _e.mock.On("Fetch", ctx, url)}
}

func (_c *MockLinkPreview_Fetch_Call) Run(run func(ctx context.Context, url string)) *MockLinkPreview_Fetch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockLinkPreview_Fetch_Call) Return(linkPreview model.LinkPreview, responseBody []byte, isFile bool, err error) *MockLinkPreview_Fetch_Call {
	_c.Call.Return(linkPreview, responseBody, isFile, err)
	return _c
}

func (_c *MockLinkPreview_Fetch_Call) RunAndReturn(run func(context.Context, string) (model.LinkPreview, []byte, bool, error)) *MockLinkPreview_Fetch_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function with given fields: a
func (_m *MockLinkPreview) Init(a *app.App) error {
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

// MockLinkPreview_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type MockLinkPreview_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - a *app.App
func (_e *MockLinkPreview_Expecter) Init(a interface{}) *MockLinkPreview_Init_Call {
	return &MockLinkPreview_Init_Call{Call: _e.mock.On("Init", a)}
}

func (_c *MockLinkPreview_Init_Call) Run(run func(a *app.App)) *MockLinkPreview_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*app.App))
	})
	return _c
}

func (_c *MockLinkPreview_Init_Call) Return(err error) *MockLinkPreview_Init_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockLinkPreview_Init_Call) RunAndReturn(run func(*app.App) error) *MockLinkPreview_Init_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *MockLinkPreview) Name() string {
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

// MockLinkPreview_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockLinkPreview_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockLinkPreview_Expecter) Name() *MockLinkPreview_Name_Call {
	return &MockLinkPreview_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockLinkPreview_Name_Call) Run(run func()) *MockLinkPreview_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLinkPreview_Name_Call) Return(name string) *MockLinkPreview_Name_Call {
	_c.Call.Return(name)
	return _c
}

func (_c *MockLinkPreview_Name_Call) RunAndReturn(run func() string) *MockLinkPreview_Name_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLinkPreview creates a new instance of MockLinkPreview. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLinkPreview(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLinkPreview {
	mock := &MockLinkPreview{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

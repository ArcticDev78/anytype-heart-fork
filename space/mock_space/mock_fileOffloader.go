// Code generated by mockery v2.38.0. DO NOT EDIT.

package mock_space

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockfileOffloader is an autogenerated mock type for the fileOffloader type
type MockfileOffloader struct {
	mock.Mock
}

type MockfileOffloader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockfileOffloader) EXPECT() *MockfileOffloader_Expecter {
	return &MockfileOffloader_Expecter{mock: &_m.Mock}
}

// FileSpaceOffload provides a mock function with given fields: ctx, spaceId, includeNotPinned
func (_m *MockfileOffloader) FileSpaceOffload(ctx context.Context, spaceId string, includeNotPinned bool) (int, uint64, error) {
	ret := _m.Called(ctx, spaceId, includeNotPinned)

	if len(ret) == 0 {
		panic("no return value specified for FileSpaceOffload")
	}

	var r0 int
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (int, uint64, error)); ok {
		return rf(ctx, spaceId, includeNotPinned)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) int); ok {
		r0 = rf(ctx, spaceId, includeNotPinned)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) uint64); ok {
		r1 = rf(ctx, spaceId, includeNotPinned)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, bool) error); ok {
		r2 = rf(ctx, spaceId, includeNotPinned)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockfileOffloader_FileSpaceOffload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FileSpaceOffload'
type MockfileOffloader_FileSpaceOffload_Call struct {
	*mock.Call
}

// FileSpaceOffload is a helper method to define mock.On call
//   - ctx context.Context
//   - spaceId string
//   - includeNotPinned bool
func (_e *MockfileOffloader_Expecter) FileSpaceOffload(ctx interface{}, spaceId interface{}, includeNotPinned interface{}) *MockfileOffloader_FileSpaceOffload_Call {
	return &MockfileOffloader_FileSpaceOffload_Call{Call: _e.mock.On("FileSpaceOffload", ctx, spaceId, includeNotPinned)}
}

func (_c *MockfileOffloader_FileSpaceOffload_Call) Run(run func(ctx context.Context, spaceId string, includeNotPinned bool)) *MockfileOffloader_FileSpaceOffload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *MockfileOffloader_FileSpaceOffload_Call) Return(filesOffloaded int, totalSize uint64, err error) *MockfileOffloader_FileSpaceOffload_Call {
	_c.Call.Return(filesOffloaded, totalSize, err)
	return _c
}

func (_c *MockfileOffloader_FileSpaceOffload_Call) RunAndReturn(run func(context.Context, string, bool) (int, uint64, error)) *MockfileOffloader_FileSpaceOffload_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockfileOffloader creates a new instance of MockfileOffloader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockfileOffloader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockfileOffloader {
	mock := &MockfileOffloader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

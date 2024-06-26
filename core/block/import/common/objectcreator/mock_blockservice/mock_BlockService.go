// Code generated by mockery. DO NOT EDIT.

package mock_blockservice

import (
	context "context"

	domain "github.com/anyproto/anytype-heart/core/domain"
	mock "github.com/stretchr/testify/mock"

	pb "github.com/anyproto/anytype-heart/pb"

	smartblock "github.com/anyproto/anytype-heart/core/block/editor/smartblock"
)

// MockBlockService is an autogenerated mock type for the BlockService type
type MockBlockService struct {
	mock.Mock
}

type MockBlockService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBlockService) EXPECT() *MockBlockService_Expecter {
	return &MockBlockService_Expecter{mock: &_m.Mock}
}

// DeleteObject provides a mock function with given fields: objectId
func (_m *MockBlockService) DeleteObject(objectId string) error {
	ret := _m.Called(objectId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteObject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(objectId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBlockService_DeleteObject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteObject'
type MockBlockService_DeleteObject_Call struct {
	*mock.Call
}

// DeleteObject is a helper method to define mock.On call
//   - objectId string
func (_e *MockBlockService_Expecter) DeleteObject(objectId interface{}) *MockBlockService_DeleteObject_Call {
	return &MockBlockService_DeleteObject_Call{Call: _e.mock.On("DeleteObject", objectId)}
}

func (_c *MockBlockService_DeleteObject_Call) Run(run func(objectId string)) *MockBlockService_DeleteObject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockBlockService_DeleteObject_Call) Return(err error) *MockBlockService_DeleteObject_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockBlockService_DeleteObject_Call) RunAndReturn(run func(string) error) *MockBlockService_DeleteObject_Call {
	_c.Call.Return(run)
	return _c
}

// GetObject provides a mock function with given fields: ctx, objectID
func (_m *MockBlockService) GetObject(ctx context.Context, objectID string) (smartblock.SmartBlock, error) {
	ret := _m.Called(ctx, objectID)

	if len(ret) == 0 {
		panic("no return value specified for GetObject")
	}

	var r0 smartblock.SmartBlock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (smartblock.SmartBlock, error)); ok {
		return rf(ctx, objectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) smartblock.SmartBlock); ok {
		r0 = rf(ctx, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(smartblock.SmartBlock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBlockService_GetObject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetObject'
type MockBlockService_GetObject_Call struct {
	*mock.Call
}

// GetObject is a helper method to define mock.On call
//   - ctx context.Context
//   - objectID string
func (_e *MockBlockService_Expecter) GetObject(ctx interface{}, objectID interface{}) *MockBlockService_GetObject_Call {
	return &MockBlockService_GetObject_Call{Call: _e.mock.On("GetObject", ctx, objectID)}
}

func (_c *MockBlockService_GetObject_Call) Run(run func(ctx context.Context, objectID string)) *MockBlockService_GetObject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockBlockService_GetObject_Call) Return(sb smartblock.SmartBlock, err error) *MockBlockService_GetObject_Call {
	_c.Call.Return(sb, err)
	return _c
}

func (_c *MockBlockService_GetObject_Call) RunAndReturn(run func(context.Context, string) (smartblock.SmartBlock, error)) *MockBlockService_GetObject_Call {
	_c.Call.Return(run)
	return _c
}

// GetObjectByFullID provides a mock function with given fields: ctx, id
func (_m *MockBlockService) GetObjectByFullID(ctx context.Context, id domain.FullID) (smartblock.SmartBlock, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetObjectByFullID")
	}

	var r0 smartblock.SmartBlock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.FullID) (smartblock.SmartBlock, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.FullID) smartblock.SmartBlock); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(smartblock.SmartBlock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.FullID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBlockService_GetObjectByFullID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetObjectByFullID'
type MockBlockService_GetObjectByFullID_Call struct {
	*mock.Call
}

// GetObjectByFullID is a helper method to define mock.On call
//   - ctx context.Context
//   - id domain.FullID
func (_e *MockBlockService_Expecter) GetObjectByFullID(ctx interface{}, id interface{}) *MockBlockService_GetObjectByFullID_Call {
	return &MockBlockService_GetObjectByFullID_Call{Call: _e.mock.On("GetObjectByFullID", ctx, id)}
}

func (_c *MockBlockService_GetObjectByFullID_Call) Run(run func(ctx context.Context, id domain.FullID)) *MockBlockService_GetObjectByFullID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.FullID))
	})
	return _c
}

func (_c *MockBlockService_GetObjectByFullID_Call) Return(sb smartblock.SmartBlock, err error) *MockBlockService_GetObjectByFullID_Call {
	_c.Call.Return(sb, err)
	return _c
}

func (_c *MockBlockService_GetObjectByFullID_Call) RunAndReturn(run func(context.Context, domain.FullID) (smartblock.SmartBlock, error)) *MockBlockService_GetObjectByFullID_Call {
	_c.Call.Return(run)
	return _c
}

// SetPageIsArchived provides a mock function with given fields: req
func (_m *MockBlockService) SetPageIsArchived(req pb.RpcObjectSetIsArchivedRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for SetPageIsArchived")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(pb.RpcObjectSetIsArchivedRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBlockService_SetPageIsArchived_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetPageIsArchived'
type MockBlockService_SetPageIsArchived_Call struct {
	*mock.Call
}

// SetPageIsArchived is a helper method to define mock.On call
//   - req pb.RpcObjectSetIsArchivedRequest
func (_e *MockBlockService_Expecter) SetPageIsArchived(req interface{}) *MockBlockService_SetPageIsArchived_Call {
	return &MockBlockService_SetPageIsArchived_Call{Call: _e.mock.On("SetPageIsArchived", req)}
}

func (_c *MockBlockService_SetPageIsArchived_Call) Run(run func(req pb.RpcObjectSetIsArchivedRequest)) *MockBlockService_SetPageIsArchived_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(pb.RpcObjectSetIsArchivedRequest))
	})
	return _c
}

func (_c *MockBlockService_SetPageIsArchived_Call) Return(err error) *MockBlockService_SetPageIsArchived_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockBlockService_SetPageIsArchived_Call) RunAndReturn(run func(pb.RpcObjectSetIsArchivedRequest) error) *MockBlockService_SetPageIsArchived_Call {
	_c.Call.Return(run)
	return _c
}

// SetPageIsFavorite provides a mock function with given fields: req
func (_m *MockBlockService) SetPageIsFavorite(req pb.RpcObjectSetIsFavoriteRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for SetPageIsFavorite")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(pb.RpcObjectSetIsFavoriteRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBlockService_SetPageIsFavorite_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetPageIsFavorite'
type MockBlockService_SetPageIsFavorite_Call struct {
	*mock.Call
}

// SetPageIsFavorite is a helper method to define mock.On call
//   - req pb.RpcObjectSetIsFavoriteRequest
func (_e *MockBlockService_Expecter) SetPageIsFavorite(req interface{}) *MockBlockService_SetPageIsFavorite_Call {
	return &MockBlockService_SetPageIsFavorite_Call{Call: _e.mock.On("SetPageIsFavorite", req)}
}

func (_c *MockBlockService_SetPageIsFavorite_Call) Run(run func(req pb.RpcObjectSetIsFavoriteRequest)) *MockBlockService_SetPageIsFavorite_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(pb.RpcObjectSetIsFavoriteRequest))
	})
	return _c
}

func (_c *MockBlockService_SetPageIsFavorite_Call) Return(err error) *MockBlockService_SetPageIsFavorite_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockBlockService_SetPageIsFavorite_Call) RunAndReturn(run func(pb.RpcObjectSetIsFavoriteRequest) error) *MockBlockService_SetPageIsFavorite_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockBlockService creates a new instance of MockBlockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBlockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBlockService {
	mock := &MockBlockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

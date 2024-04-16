// Code generated by mockery. DO NOT EDIT.

package mock_participantwatcher

import (
	context "context"

	app "github.com/anyproto/any-sync/app"
	clientspace "github.com/anyproto/anytype-heart/space/clientspace"

	list "github.com/anyproto/any-sync/commonspace/object/acl/list"

	mock "github.com/stretchr/testify/mock"
)

// MockParticipantWatcher is an autogenerated mock type for the ParticipantWatcher type
type MockParticipantWatcher struct {
	mock.Mock
}

type MockParticipantWatcher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockParticipantWatcher) EXPECT() *MockParticipantWatcher_Expecter {
	return &MockParticipantWatcher_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields: ctx
func (_m *MockParticipantWatcher) Close(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockParticipantWatcher_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockParticipantWatcher_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockParticipantWatcher_Expecter) Close(ctx interface{}) *MockParticipantWatcher_Close_Call {
	return &MockParticipantWatcher_Close_Call{Call: _e.mock.On("Close", ctx)}
}

func (_c *MockParticipantWatcher_Close_Call) Run(run func(ctx context.Context)) *MockParticipantWatcher_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockParticipantWatcher_Close_Call) Return(err error) *MockParticipantWatcher_Close_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockParticipantWatcher_Close_Call) RunAndReturn(run func(context.Context) error) *MockParticipantWatcher_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function with given fields: a
func (_m *MockParticipantWatcher) Init(a *app.App) error {
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

// MockParticipantWatcher_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type MockParticipantWatcher_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - a *app.App
func (_e *MockParticipantWatcher_Expecter) Init(a interface{}) *MockParticipantWatcher_Init_Call {
	return &MockParticipantWatcher_Init_Call{Call: _e.mock.On("Init", a)}
}

func (_c *MockParticipantWatcher_Init_Call) Run(run func(a *app.App)) *MockParticipantWatcher_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*app.App))
	})
	return _c
}

func (_c *MockParticipantWatcher_Init_Call) Return(err error) *MockParticipantWatcher_Init_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockParticipantWatcher_Init_Call) RunAndReturn(run func(*app.App) error) *MockParticipantWatcher_Init_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *MockParticipantWatcher) Name() string {
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

// MockParticipantWatcher_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockParticipantWatcher_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockParticipantWatcher_Expecter) Name() *MockParticipantWatcher_Name_Call {
	return &MockParticipantWatcher_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockParticipantWatcher_Name_Call) Run(run func()) *MockParticipantWatcher_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockParticipantWatcher_Name_Call) Return(name string) *MockParticipantWatcher_Name_Call {
	_c.Call.Return(name)
	return _c
}

func (_c *MockParticipantWatcher_Name_Call) RunAndReturn(run func() string) *MockParticipantWatcher_Name_Call {
	_c.Call.Return(run)
	return _c
}

// Run provides a mock function with given fields: ctx
func (_m *MockParticipantWatcher) Run(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockParticipantWatcher_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockParticipantWatcher_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockParticipantWatcher_Expecter) Run(ctx interface{}) *MockParticipantWatcher_Run_Call {
	return &MockParticipantWatcher_Run_Call{Call: _e.mock.On("Run", ctx)}
}

func (_c *MockParticipantWatcher_Run_Call) Run(run func(ctx context.Context)) *MockParticipantWatcher_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockParticipantWatcher_Run_Call) Return(err error) *MockParticipantWatcher_Run_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockParticipantWatcher_Run_Call) RunAndReturn(run func(context.Context) error) *MockParticipantWatcher_Run_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAccountParticipantFromProfile provides a mock function with given fields: ctx, space
func (_m *MockParticipantWatcher) UpdateAccountParticipantFromProfile(ctx context.Context, space clientspace.Space) error {
	ret := _m.Called(ctx, space)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAccountParticipantFromProfile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, clientspace.Space) error); ok {
		r0 = rf(ctx, space)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAccountParticipantFromProfile'
type MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call struct {
	*mock.Call
}

// UpdateAccountParticipantFromProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - space clientspace.Space
func (_e *MockParticipantWatcher_Expecter) UpdateAccountParticipantFromProfile(ctx interface{}, space interface{}) *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call {
	return &MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call{Call: _e.mock.On("UpdateAccountParticipantFromProfile", ctx, space)}
}

func (_c *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call) Run(run func(ctx context.Context, space clientspace.Space)) *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(clientspace.Space))
	})
	return _c
}

func (_c *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call) Return(_a0 error) *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call) RunAndReturn(run func(context.Context, clientspace.Space) error) *MockParticipantWatcher_UpdateAccountParticipantFromProfile_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateParticipantFromAclState provides a mock function with given fields: ctx, space, accState
func (_m *MockParticipantWatcher) UpdateParticipantFromAclState(ctx context.Context, space clientspace.Space, accState list.AccountState) error {
	ret := _m.Called(ctx, space, accState)

	if len(ret) == 0 {
		panic("no return value specified for UpdateParticipantFromAclState")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, clientspace.Space, list.AccountState) error); ok {
		r0 = rf(ctx, space, accState)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockParticipantWatcher_UpdateParticipantFromAclState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateParticipantFromAclState'
type MockParticipantWatcher_UpdateParticipantFromAclState_Call struct {
	*mock.Call
}

// UpdateParticipantFromAclState is a helper method to define mock.On call
//   - ctx context.Context
//   - space clientspace.Space
//   - accState list.AccountState
func (_e *MockParticipantWatcher_Expecter) UpdateParticipantFromAclState(ctx interface{}, space interface{}, accState interface{}) *MockParticipantWatcher_UpdateParticipantFromAclState_Call {
	return &MockParticipantWatcher_UpdateParticipantFromAclState_Call{Call: _e.mock.On("UpdateParticipantFromAclState", ctx, space, accState)}
}

func (_c *MockParticipantWatcher_UpdateParticipantFromAclState_Call) Run(run func(ctx context.Context, space clientspace.Space, accState list.AccountState)) *MockParticipantWatcher_UpdateParticipantFromAclState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(clientspace.Space), args[2].(list.AccountState))
	})
	return _c
}

func (_c *MockParticipantWatcher_UpdateParticipantFromAclState_Call) Return(_a0 error) *MockParticipantWatcher_UpdateParticipantFromAclState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockParticipantWatcher_UpdateParticipantFromAclState_Call) RunAndReturn(run func(context.Context, clientspace.Space, list.AccountState) error) *MockParticipantWatcher_UpdateParticipantFromAclState_Call {
	_c.Call.Return(run)
	return _c
}

// WatchParticipant provides a mock function with given fields: ctx, space, accState
func (_m *MockParticipantWatcher) WatchParticipant(ctx context.Context, space clientspace.Space, accState list.AccountState) error {
	ret := _m.Called(ctx, space, accState)

	if len(ret) == 0 {
		panic("no return value specified for WatchParticipant")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, clientspace.Space, list.AccountState) error); ok {
		r0 = rf(ctx, space, accState)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockParticipantWatcher_WatchParticipant_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WatchParticipant'
type MockParticipantWatcher_WatchParticipant_Call struct {
	*mock.Call
}

// WatchParticipant is a helper method to define mock.On call
//   - ctx context.Context
//   - space clientspace.Space
//   - accState list.AccountState
func (_e *MockParticipantWatcher_Expecter) WatchParticipant(ctx interface{}, space interface{}, accState interface{}) *MockParticipantWatcher_WatchParticipant_Call {
	return &MockParticipantWatcher_WatchParticipant_Call{Call: _e.mock.On("WatchParticipant", ctx, space, accState)}
}

func (_c *MockParticipantWatcher_WatchParticipant_Call) Run(run func(ctx context.Context, space clientspace.Space, accState list.AccountState)) *MockParticipantWatcher_WatchParticipant_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(clientspace.Space), args[2].(list.AccountState))
	})
	return _c
}

func (_c *MockParticipantWatcher_WatchParticipant_Call) Return(_a0 error) *MockParticipantWatcher_WatchParticipant_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockParticipantWatcher_WatchParticipant_Call) RunAndReturn(run func(context.Context, clientspace.Space, list.AccountState) error) *MockParticipantWatcher_WatchParticipant_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockParticipantWatcher creates a new instance of MockParticipantWatcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockParticipantWatcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockParticipantWatcher {
	mock := &MockParticipantWatcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
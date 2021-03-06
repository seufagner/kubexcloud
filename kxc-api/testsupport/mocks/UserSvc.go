// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	requests "github.com/didil/kubexcloud/kxc-api/requests"
	mock "github.com/stretchr/testify/mock"

	responses "github.com/didil/kubexcloud/kxc-api/responses"
)

// UserSvc is an autogenerated mock type for the UserSvc type
type UserSvc struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, reqData
func (_m *UserSvc) Create(ctx context.Context, reqData *requests.CreateUser) error {
	ret := _m.Called(ctx, reqData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *requests.CreateUser) error); ok {
		r0 = rf(ctx, reqData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HasRole provides a mock function with given fields: ctx, userName, role
func (_m *UserSvc) HasRole(ctx context.Context, userName string, role string) (bool, error) {
	ret := _m.Called(ctx, userName, role)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, userName, role)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userName, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx
func (_m *UserSvc) List(ctx context.Context) (*responses.ListUser, error) {
	ret := _m.Called(ctx)

	var r0 *responses.ListUser
	if rf, ok := ret.Get(0).(func(context.Context) *responses.ListUser); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*responses.ListUser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, userName, password
func (_m *UserSvc) Login(ctx context.Context, userName string, password string) (string, error) {
	ret := _m.Called(ctx, userName, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, userName, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userName, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

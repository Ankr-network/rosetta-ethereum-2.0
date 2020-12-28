// Code generated by mockery v1.0.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/coinbase/rosetta-sdk-go/types"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Block provides a mock function with given fields: _a0, _a1
func (_m *Client) Block(_a0 context.Context, _a1 *types.PartialBlockIdentifier) (*types.Block, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.Block
	if rf, ok := ret.Get(0).(func(context.Context, *types.PartialBlockIdentifier) *types.Block); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.PartialBlockIdentifier) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Status provides a mock function with given fields: _a0
func (_m *Client) Status(_a0 context.Context) (*types.BlockIdentifier, int64, *types.SyncStatus, []*types.Peer, error) {
	ret := _m.Called(_a0)

	var r0 *types.BlockIdentifier
	if rf, ok := ret.Get(0).(func(context.Context) *types.BlockIdentifier); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BlockIdentifier)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context) int64); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 *types.SyncStatus
	if rf, ok := ret.Get(2).(func(context.Context) *types.SyncStatus); ok {
		r2 = rf(_a0)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*types.SyncStatus)
		}
	}

	var r3 []*types.Peer
	if rf, ok := ret.Get(3).(func(context.Context) []*types.Peer); ok {
		r3 = rf(_a0)
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).([]*types.Peer)
		}
	}

	var r4 error
	if rf, ok := ret.Get(4).(func(context.Context) error); ok {
		r4 = rf(_a0)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}

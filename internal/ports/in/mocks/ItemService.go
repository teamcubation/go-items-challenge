// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	item "github.com/teamcubation/go-items-challenge/internal/domain/item"

	mock "github.com/stretchr/testify/mock"
)

// ItemService is an autogenerated mock type for the ItemService type
type ItemService struct {
	mock.Mock
}

// CreateItem provides a mock function with given fields: ctx, itm
func (_m *ItemService) CreateItem(ctx context.Context, itm *item.Item) (*item.Item, error) {
	ret := _m.Called(ctx, itm)

	if len(ret) == 0 {
		panic("no return value specified for CreateItem")
	}

	var r0 *item.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *item.Item) (*item.Item, error)); ok {
		return rf(ctx, itm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *item.Item) *item.Item); ok {
		r0 = rf(ctx, itm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*item.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *item.Item) error); ok {
		r1 = rf(ctx, itm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteItem provides a mock function with given fields: ctx, id
func (_m *ItemService) DeleteItem(ctx context.Context, id int) (*item.Item, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteItem")
	}

	var r0 *item.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*item.Item, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *item.Item); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*item.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItemByID provides a mock function with given fields: ctx, id
func (_m *ItemService) GetItemByID(ctx context.Context, id int) (*item.Item, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetItemByID")
	}

	var r0 *item.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*item.Item, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *item.Item); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*item.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ItemExistsByCode provides a mock function with given fields: ctx, code
func (_m *ItemService) ItemExistsByCode(ctx context.Context, code string) bool {
	ret := _m.Called(ctx, code)

	if len(ret) == 0 {
		panic("no return value specified for ItemExistsByCode")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ListItems provides a mock function with given fields: ctx, status, limit, page
func (_m *ItemService) ListItems(ctx context.Context, status string, limit int, page int) ([]*item.Item, int, error) {
	ret := _m.Called(ctx, status, limit, page)

	if len(ret) == 0 {
		panic("no return value specified for ListItems")
	}

	var r0 []*item.Item
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]*item.Item, int, error)); ok {
		return rf(ctx, status, limit, page)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*item.Item); ok {
		r0 = rf(ctx, status, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*item.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int); ok {
		r1 = rf(ctx, status, limit, page)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, status, limit, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateItem provides a mock function with given fields: ctx, itm
func (_m *ItemService) UpdateItem(ctx context.Context, itm *item.Item) (*item.Item, error) {
	ret := _m.Called(ctx, itm)

	if len(ret) == 0 {
		panic("no return value specified for UpdateItem")
	}

	var r0 *item.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *item.Item) (*item.Item, error)); ok {
		return rf(ctx, itm)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *item.Item) *item.Item); ok {
		r0 = rf(ctx, itm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*item.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *item.Item) error); ok {
		r1 = rf(ctx, itm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewItemService creates a new instance of ItemService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewItemService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ItemService {
	mock := &ItemService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

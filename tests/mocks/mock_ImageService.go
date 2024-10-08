// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	readme "github.com/liveoaklabs/readme-api-go-client/readme"
	mock "github.com/stretchr/testify/mock"
)

// MockImageService is an autogenerated mock type for the ImageService type
type MockImageService struct {
	mock.Mock
}

type MockImageService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockImageService) EXPECT() *MockImageService_Expecter {
	return &MockImageService_Expecter{mock: &_m.Mock}
}

// Upload provides a mock function with given fields: source, filename
func (_m *MockImageService) Upload(source []byte, filename ...string) (readme.Image, *readme.APIResponse, error) {
	_va := make([]interface{}, len(filename))
	for _i := range filename {
		_va[_i] = filename[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, source)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 readme.Image
	var r1 *readme.APIResponse
	var r2 error
	if rf, ok := ret.Get(0).(func([]byte, ...string) (readme.Image, *readme.APIResponse, error)); ok {
		return rf(source, filename...)
	}
	if rf, ok := ret.Get(0).(func([]byte, ...string) readme.Image); ok {
		r0 = rf(source, filename...)
	} else {
		r0 = ret.Get(0).(readme.Image)
	}

	if rf, ok := ret.Get(1).(func([]byte, ...string) *readme.APIResponse); ok {
		r1 = rf(source, filename...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*readme.APIResponse)
		}
	}

	if rf, ok := ret.Get(2).(func([]byte, ...string) error); ok {
		r2 = rf(source, filename...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockImageService_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type MockImageService_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - source []byte
//   - filename ...string
func (_e *MockImageService_Expecter) Upload(source interface{}, filename ...interface{}) *MockImageService_Upload_Call {
	return &MockImageService_Upload_Call{Call: _e.mock.On("Upload",
		append([]interface{}{source}, filename...)...)}
}

func (_c *MockImageService_Upload_Call) Run(run func(source []byte, filename ...string)) *MockImageService_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].([]byte), variadicArgs...)
	})
	return _c
}

func (_c *MockImageService_Upload_Call) Return(_a0 readme.Image, _a1 *readme.APIResponse, _a2 error) *MockImageService_Upload_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockImageService_Upload_Call) RunAndReturn(run func([]byte, ...string) (readme.Image, *readme.APIResponse, error)) *MockImageService_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockImageService creates a new instance of MockImageService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockImageService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockImageService {
	mock := &MockImageService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

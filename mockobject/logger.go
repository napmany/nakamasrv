// Code generated by mockery v2.42.1. DO NOT EDIT.

package mockobject

import (
	runtime "github.com/heroiclabs/nakama-common/runtime"
	mock "github.com/stretchr/testify/mock"
)

// LoggerMock is an autogenerated mock type for the Logger type
type LoggerMock struct {
	mock.Mock
}

// Debug provides a mock function with given fields: format, v
func (_m *LoggerMock) Debug(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: format, v
func (_m *LoggerMock) Error(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Fields provides a mock function with given fields:
func (_m *LoggerMock) Fields() map[string]interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Fields")
	}

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// Info provides a mock function with given fields: format, v
func (_m *LoggerMock) Info(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: format, v
func (_m *LoggerMock) Warn(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// WithField provides a mock function with given fields: key, v
func (_m *LoggerMock) WithField(key string, v interface{}) runtime.Logger {
	ret := _m.Called(key, v)

	if len(ret) == 0 {
		panic("no return value specified for WithField")
	}

	var r0 runtime.Logger
	if rf, ok := ret.Get(0).(func(string, interface{}) runtime.Logger); ok {
		r0 = rf(key, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(runtime.Logger)
		}
	}

	return r0
}

// WithFields provides a mock function with given fields: fields
func (_m *LoggerMock) WithFields(fields map[string]interface{}) runtime.Logger {
	ret := _m.Called(fields)

	if len(ret) == 0 {
		panic("no return value specified for WithFields")
	}

	var r0 runtime.Logger
	if rf, ok := ret.Get(0).(func(map[string]interface{}) runtime.Logger); ok {
		r0 = rf(fields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(runtime.Logger)
		}
	}

	return r0
}

// NewLoggerMock creates a new instance of LoggerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoggerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *LoggerMock {
	mock := &LoggerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

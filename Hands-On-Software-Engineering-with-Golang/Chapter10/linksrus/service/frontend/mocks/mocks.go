// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ImagineDevOps/Hands-On-Software-Engineering-with-Golang/Chapter10/linksrus/service/frontend (interfaces: GraphAPI,IndexAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	graph "github.com/ImagineDevOps/Hands-On-Software-Engineering-with-Golang/Chapter06/linkgraph/graph"
	index "github.com/ImagineDevOps/Hands-On-Software-Engineering-with-Golang/Chapter06/textindexer/index"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGraphAPI is a mock of GraphAPI interface
type MockGraphAPI struct {
	ctrl     *gomock.Controller
	recorder *MockGraphAPIMockRecorder
}

// MockGraphAPIMockRecorder is the mock recorder for MockGraphAPI
type MockGraphAPIMockRecorder struct {
	mock *MockGraphAPI
}

// NewMockGraphAPI creates a new mock instance
func NewMockGraphAPI(ctrl *gomock.Controller) *MockGraphAPI {
	mock := &MockGraphAPI{ctrl: ctrl}
	mock.recorder = &MockGraphAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGraphAPI) EXPECT() *MockGraphAPIMockRecorder {
	return m.recorder
}

// UpsertLink mocks base method
func (m *MockGraphAPI) UpsertLink(arg0 *graph.Link) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertLink", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertLink indicates an expected call of UpsertLink
func (mr *MockGraphAPIMockRecorder) UpsertLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertLink", reflect.TypeOf((*MockGraphAPI)(nil).UpsertLink), arg0)
}

// MockIndexAPI is a mock of IndexAPI interface
type MockIndexAPI struct {
	ctrl     *gomock.Controller
	recorder *MockIndexAPIMockRecorder
}

// MockIndexAPIMockRecorder is the mock recorder for MockIndexAPI
type MockIndexAPIMockRecorder struct {
	mock *MockIndexAPI
}

// NewMockIndexAPI creates a new mock instance
func NewMockIndexAPI(ctrl *gomock.Controller) *MockIndexAPI {
	mock := &MockIndexAPI{ctrl: ctrl}
	mock.recorder = &MockIndexAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIndexAPI) EXPECT() *MockIndexAPIMockRecorder {
	return m.recorder
}

// Search mocks base method
func (m *MockIndexAPI) Search(arg0 index.Query) (index.Iterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0)
	ret0, _ := ret[0].(index.Iterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockIndexAPIMockRecorder) Search(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexAPI)(nil).Search), arg0)
}

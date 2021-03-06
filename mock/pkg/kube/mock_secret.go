// Code generated by MockGen. DO NOT EDIT.
// Source: ../../pkg/translator/kube/secret.go

// Package mock_kube is a generated GoMock package.
package mock_kube

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSecretClient is a mock of SecretClient interface
type MockSecretClient struct {
	ctrl     *gomock.Controller
	recorder *MockSecretClientMockRecorder
}

// MockSecretClientMockRecorder is the mock recorder for MockSecretClient
type MockSecretClientMockRecorder struct {
	mock *MockSecretClient
}

// NewMockSecretClient creates a new mock instance
func NewMockSecretClient(ctrl *gomock.Controller) *MockSecretClient {
	mock := &MockSecretClient{ctrl: ctrl}
	mock.recorder = &MockSecretClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecretClient) EXPECT() *MockSecretClientMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockSecretClient) Delete(namespace, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", namespace, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSecretClientMockRecorder) Delete(namespace, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSecretClient)(nil).Delete), namespace, name)
}

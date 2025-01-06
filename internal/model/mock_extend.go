package model

import (
	context "context"

	gomock "github.com/golang/mock/gomock"
)

type ExtendMockModel struct {
	*MockModelInterface
}

func NewExtendedMockModelInterface(ctrl *gomock.Controller) *ExtendMockModel {
	mock := NewMockModelInterface(ctrl)
	return &ExtendMockModel{mock}
}

func (e *ExtendMockModel) RunTransaction(ctx context.Context, f func(model ModelInterface) error) error {
	return f(e)
}

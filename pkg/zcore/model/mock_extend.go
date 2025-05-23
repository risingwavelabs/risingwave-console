package model

import (
	context "context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/mock/gomock"
)

type ExtendMockModel struct {
	*MockModelInterface
}

func NewMockModelInterfaceWithTransaction(ctrl *gomock.Controller) *ExtendMockModel {
	mock := NewMockModelInterface(ctrl)
	return &ExtendMockModel{mock}
}

func (e *ExtendMockModel) RunTransaction(ctx context.Context, f func(model ModelInterface) error) error {
	return f(e)
}

func (e *ExtendMockModel) RunTransactionWithTx(ctx context.Context, f func(tx pgx.Tx, model ModelInterface) error) error {
	return f(nil, e)
}

func (e *ExtendMockModel) SpawnWithTx(tx pgx.Tx) ModelInterface {
	return e
}

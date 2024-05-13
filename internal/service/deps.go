//go:generate mockgen -source=$GOFILE -destination=mock_test.go -package=$GOPACKAGE
package service

import (
	"context"

	"github.com/mirhijinam/pg-project/internal/model"
	"github.com/mirhijinam/pg-project/internal/repository"
)

type Repo interface {
	Create(context.Context, *model.Command) error
	SetError(context.Context, int, string) error
	SetSuccess(context.Context, int) error
	GetList(context.Context) ([]model.Command, error)
	GetOne(context.Context, int) (model.Command, error)
	Writer(context.Context, int) repository.WriterFunc
}

type Pool interface {
	activate()
	Go(func())
}

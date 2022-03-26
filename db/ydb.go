package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	ydb "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	yc "github.com/ydb-platform/ydb-go-yc"

	"github.com/pav5000/easy-yc/auth"
)

type Service struct {
	conn ydb.Connection
	txc  *table.TransactionControl
}

func New(ctx context.Context, endpoint, path string) (_ *Service, err error) {
	service := &Service{}

	service.conn, err = ydb.New(
		ctx,
		ydb.WithEndpoint(endpoint),
		ydb.WithDatabase(path),
		ydb.WithSecure(true),
		yc.WithInternalCA(),
		//yc.WithServiceAccountKeyFileCredentials("iam.txt"),
		ydb.WithCredentials(auth.GetYdbCredentials()),
		ydb.WithDialTimeout(3*time.Second),
	)
	if err != nil {
		return nil, errors.Wrap(err, "ydb.New")
	}

	service.txc = table.TxControl(
		table.BeginTx(table.WithSerializableReadWrite()),
		table.CommitTx(),
	)

	return service, nil
}

func (s *Service) DefaultTXC() *table.TransactionControl {
	return s.txc
}

func (s *Service) Execute(
	ctx context.Context,
	query string, params *table.QueryParameters,
	dataFunc func(result.Result) error,
) error {
	return s.conn.Table().Do(
		ctx,
		func(ctx context.Context, session table.Session) (err error) {
			_, res, err := session.Execute(ctx, s.DefaultTXC(), query, params)
			if err != nil {
				return err
			}
			defer res.Close()

			return dataFunc(res)
		},
	)
}

func (s *Service) Do(ctx context.Context, op table.Operation, opts ...table.Option) {
	s.conn.Table().Do(ctx, op, opts...)
}

func (s *Service) DoTx(ctx context.Context, op table.TxOperation, opts ...table.Option) {
	s.conn.Table().DoTx(ctx, op, opts...)
}

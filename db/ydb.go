package db

import (
	"context"
	"time"

	"github.com/pav5000/easy-yc/auth"
	"github.com/pkg/errors"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
)

type Service struct {
	conn ydb.Connection
	txc  *table.TransactionControl
}

func New(ctx context.Context, endpoint, path string) (*Service, error) {
	service := &Service{}

	token, err := auth.GetIAMToken(ctx)
	if err != nil {
		return nil, err
	}

	service.conn, err = ydb.New(
		ctx,
		ydb.WithConnectionString(endpoint),
		ydb.WithDialTimeout(3*time.Second),
		ydb.WithCertificatesFromPem([]byte(rootPEM)),
		ydb.WithSessionPoolIdleThreshold(time.Second*5),
		ydb.WithSessionPoolKeepAliveMinSize(-1),
		ydb.WithDiscoveryInterval(5*time.Second),
		ydb.WithAccessTokenCredentials(token),
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
	userFunc func(result.Result) error,
	query string, params *table.QueryParameters,
) error {
	return s.conn.Table().Do(
		ctx,
		func(ctx context.Context, session table.Session) (err error) {
			_, res, err := session.Execute(ctx, s.DefaultTXC(), query, params)
			if err != nil {
				return err
			}
			defer res.Close()

			return userFunc(res)
		},
	)
}

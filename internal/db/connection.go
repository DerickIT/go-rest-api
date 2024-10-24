package db

import (
	"context"
	"errors"
	"time"

	"github.com/derickit/go-rest-api/internal/logger"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidConnURL      = errors.New("failed to connect to DB:as the connection string is invalid")
	ErrConnectionEstablish = errors.New("failed to establish connection to DB")
	ErrClientInit          = errors.New("failed to initialize db client")
	ErrConnectionLeak      = errors.New("unable to disconnect from DB, connection leak")
	ErrPingDB              = errors.New("failed to ping DB")
)

const (
	DefConnectionTimeOut = 10 * time.Second
	DefDatabase          = "ecommerce"
)

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type MongoManager interface {
	Database() MongoDatabase
	Ping() error
	Disconnect() error
}

type ConnectionOpts struct {
	ConnectionTimeout time.Duration
	PrintQueries      bool
	Database          string
}

type ConnectionManager struct {
	connectionURL string
	client        *mongo.Client
	database      *mongo.Database
	credentials   *MongoDBCredentials
	logger        *logger.AppLogger
}

func FillConnectionOpts(opts *ConnectionOpts) *ConnectionOpts {
	if opts == nil {
		return &ConnectionOpts{
			PrintQueries:      false,
			ConnectionTimeout: DefConnectionTimeOut,
			Database:          DefDatabase,
		}
	}
	if opts.ConnectionTimeout == 0 {
		opts.ConnectionTimeout = DefConnectionTimeOut
	}
	if opts.Database == "" {
		opts.Database = DefDatabase
	}
	return opts
}

func NewMongoManager(mc *MongoDBCredentials, opts *ConnectionOpts, lgr *logger.AppLogger) (*ConnectionManager, error) {
	connURL := MongoConnectionURL(mc)
	lgr.Info().Str("connURL", MaskedMongoConnectionURL(mc)).Msg("connecting to db")
	if len(connURL) == 0 {
		return nil, ErrInvalidConnURL
	}
	connMgr := &ConnectionManager{
		credentials:   mc,
		logger:        lgr,
		connectionURL: connURL,
	}
	connOpts := FillConnectionOpts(opts)
	var err error
	var c *mongo.Client
	if c, err = connMgr.NewClient(connOpts); err == nil {
		db := c.Database(connOpts.Database)
		connMgr.database = db
		connMgr.client = c
		if pErr := connMgr.Ping(); pErr != nil {
			return nil, ErrConnectionEstablish
		}
		return connMgr, nil
	}
	return nil, err
}

func (c *ConnectionManager) NewClient(connOpts *ConnectionOpts) (*mongo.Client, error) {
	var cmdMonitor *event.CommandMonitor
	if connOpts.PrintQueries {
		cmdMonitor = &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				c.logger.Info().Str("dbQuery", evt.Command.String()).Send()
			},
		}
	}
	clientOptions := options.Client().ApplyURI(c.connectionURL).SetMonitor(cmdMonitor)
	ctx, cancel := context.WithTimeout(context.Background(), connOpts.ConnectionTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to connect to db")
		return nil, ErrClientInit
	}
	return client, nil
}

func (c *ConnectionManager) Database() MongoDatabase {
	return c.database
}

func (c *ConnectionManager) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefConnectionTimeOut)
	defer cancel()
	if err := c.client.Ping(ctx, nil); err != nil {
		c.logger.Error().Err(err).Msg("failed to ping db")
		return ErrPingDB
	}
	return nil
}

func (c *ConnectionManager) Disconnect() error {
	if err := c.client.Disconnect(context.Background()); err != nil {
		c.logger.Error().Err(err).Msg("failed to disconnect from db")
		return ErrConnectionLeak
	}
	c.logger.Info().Msg("disconnected from db")
	return nil
}

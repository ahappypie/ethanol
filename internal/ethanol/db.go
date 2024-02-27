package ethanol

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	dbsql "github.com/databricks/databricks-sql-go"
	"github.com/databricks/databricks-sql-go/auth"
	"github.com/databricks/databricks-sql-go/auth/oauth/m2m"
	"github.com/databricks/databricks-sql-go/auth/oauth/u2m"
	dbsqllog "github.com/databricks/databricks-sql-go/logger"
	"log"
	"os"
	"time"
)

type ClusterArgs struct {
	Host     string
	HttpPath string
}

type M2MArgs struct {
	ClientId     string
	ClientSecret string
}

type InternalTableArgs struct {
	Catalog string
	Schema  string
	Table   string
}

type Client struct {
	db              *sql.DB
	internalCatalog string
	internalSchema  string
	internalTable   string
}

func NewDBSQLClientWithM2M(clusterArgs *ClusterArgs, internalTableArgs *InternalTableArgs, authArgs *M2MArgs) *Client {
	// use this package to set up logging. By default logging level is `warn`. If you want to disable logging, use `disabled`
	if err := dbsqllog.SetLogLevel("debug"); err != nil {
		log.Fatal(err)
	}
	// sets the logging output. By default it will use os.Stderr. If running in terminal, it will use ConsoleWriter to make it pretty
	dbsqllog.SetLogOutput(os.Stdout)

	authenticator := m2m.NewAuthenticator(
		authArgs.ClientId,
		authArgs.ClientSecret,
		clusterArgs.Host,
	)

	connector, err := newConnector(clusterArgs, authenticator)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	// Opening a driver typically will not attempt to connect to the database.
	db := sql.OpenDB(connector)
	// make sure to close it later
	defer db.Close()

	// Pinging should require logging in
	if err := db.Ping(); err != nil {
		fmt.Println(err)
	}

	return &Client{db: db, internalCatalog: internalTableArgs.Catalog, internalSchema: internalTableArgs.Schema, internalTable: internalTableArgs.Table}
}

func NewDBSQLClientWithU2M(clusterArgs *ClusterArgs, internalTableArgs *InternalTableArgs) *Client {
	// use this package to set up logging. By default logging level is `warn`. If you want to disable logging, use `disabled`
	if err := dbsqllog.SetLogLevel("debug"); err != nil {
		log.Fatal(err)
	}
	// sets the logging output. By default it will use os.Stderr. If running in terminal, it will use ConsoleWriter to make it pretty
	dbsqllog.SetLogOutput(os.Stdout)

	authenticator, err := u2m.NewAuthenticator(
		clusterArgs.Host,
		1*time.Minute,
	)
	if err != nil {
		log.Fatal(err)
	}

	connector, err := newConnector(clusterArgs, authenticator)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	// Opening a driver typically will not attempt to connect to the database.
	db := sql.OpenDB(connector)
	// make sure to close it later
	//defer db.Close()

	// Pinging should require logging in
	startupContext, startupCancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer startupCancel()
	if err := db.PingContext(startupContext); err != nil {
		log.Fatal(err)
	}

	return &Client{db: db, internalCatalog: internalTableArgs.Catalog, internalSchema: internalTableArgs.Schema, internalTable: internalTableArgs.Table}
}

func newConnector(clusterArgs *ClusterArgs, authenticator auth.Authenticator) (driver.Connector, error) {
	return dbsql.NewConnector(
		dbsql.WithServerHostname(clusterArgs.Host),
		dbsql.WithHTTPPath(clusterArgs.HttpPath),
		dbsql.WithPort(443),
		dbsql.WithAuthenticator(authenticator),
		//optional configuration
		dbsql.WithSessionParams(map[string]string{"timezone": "America/Los_Angeles", "ansi_mode": "true"}),
		dbsql.WithUserAgentEntry("ethanol"),
		dbsql.WithInitialNamespace("main", "default"),
		dbsql.WithTimeout(time.Minute), // defaults to no timeout. Global timeout. Any query will be canceled if taking more than this time.
		dbsql.WithMaxRows(10),          // defaults to 10000
	)
}

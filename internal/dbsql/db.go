package dbsql

import (
	"database/sql"
	dbsql "github.com/databricks/databricks-sql-go"
	"github.com/databricks/databricks-sql-go/auth/oauth/m2m"
	dbsqllog "github.com/databricks/databricks-sql-go/logger"
	"log"
	"os"
	"time"
)

func NewDB(client_id string, client_secret string, host string, http_path string) *sql.DB {
	// use this package to set up logging. By default logging level is `warn`. If you want to disable logging, use `disabled`
	if err := dbsqllog.SetLogLevel("debug"); err != nil {
		log.Fatal(err)
	}
	// sets the logging output. By default it will use os.Stderr. If running in terminal, it will use ConsoleWriter to make it pretty
	dbsqllog.SetLogOutput(os.Stdout)

	//authenticator := m2m.NewAuthenticator(
	//	os.Getenv("DATABRICKS_CLIENT_ID"),
	//	os.Getenv("DATABRICKS_CLIENT_SECRET"),
	//	os.Getenv("DATABRICKS_SERVER_HOSTNAME"),
	//)
	authenticator := m2m.NewAuthenticator(
		client_id,
		client_secret,
		host,
	)

	connector, err := dbsql.NewConnector(
		//dbsql.WithServerHostname(os.Getenv("DATABRICKS_HOST")),
		//dbsql.WithHTTPPath(os.Getenv("DATABRICKS_HTTP_PATH")),
		dbsql.WithServerHostname(host),
		dbsql.WithHTTPPath(http_path),
		dbsql.WithPort(443),
		dbsql.WithAuthenticator(authenticator),
		//optional configuration
		dbsql.WithSessionParams(map[string]string{"timezone": "America/Los_Angeles", "ansi_mode": "true"}),
		dbsql.WithUserAgentEntry("ethanol"),
		dbsql.WithInitialNamespace("main", "default"),
		dbsql.WithTimeout(time.Minute), // defaults to no timeout. Global timeout. Any query will be canceled if taking more than this time.
		dbsql.WithMaxRows(10),          // defaults to 10000
	)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)

	}
	// Opening a driver typically will not attempt to connect to the database.
	db := sql.OpenDB(connector)
	// make sure to close it later
	defer db.Close()

	return db
}

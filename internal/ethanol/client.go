package ethanol

import (
	"context"
	"fmt"
	"github.com/databricks/databricks-sql-go/driverctx"
	"strings"
)

func (c Client) trackingTable() string {
	return strings.Join([]string{c.internalCatalog, c.internalSchema, c.internalTable}, ".")
}

func (c Client) GetLastMigration() (string, error) {
	var lastVersion string
	err := c.db.QueryRow(fmt.Sprintf(`SELECT version FROM %s ORDER BY run_on DESC LIMIT 1;`, c.trackingTable())).Scan(&lastVersion)
	return lastVersion, err
}

func (c Client) ExecuteMigration(migration Migration) error {
	executionCtx := driverctx.NewContextWithCorrelationId(context.Background(), migration.Version+"_"+migration.Name)
	//multiple statements not supported by golang
	//tx not implemented
	for _, stmt := range migration.ParseStatements() {
		_, err := c.db.ExecContext(executionCtx, stmt)
		if err != nil {
			return err
		}
	}
	//TODO sql injection
	//TODO also check RowsAffected?
	_, err := c.db.ExecContext(executionCtx, fmt.Sprintf(`INSERT INTO %s VALUES (%s, current_timestamp())`, c.trackingTable(), migration.Version))
	if err != nil {
		return err
	}
	return nil
}

func (c Client) Close() error {
	return c.db.Close()
}

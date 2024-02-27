package ethanol

import (
	"fmt"
	"strings"
)

func upSql(catalog string, schema string, table string) string {
	createCatalog := fmt.Sprintf(`CREATE CATALOG %[1]s;
ALTER CATALOG %[1]s SET OWNER TO $owner_principal;
USE CATALOG %[1]s;
GRANT USE CATALOG, CREATE SCHEMA ON %[1]s TO $automation_principal;
`, catalog)

	var createSchema string
	if schema != "default" {
		createSchema = fmt.Sprintf(`CREATE SCHEMA %[1]s;
ALTER SCHEMA %[1]s SET OWNER TO owner_principal;
USE SCHEMA %[1]s;
GRANT USE SCHEMA, CREATE FUNCTION, CREATE MATERIALIZED VIEW, CREATE MODEL, CREATE TABLE, CREATE VOLUME ON %[1]s TO $automation_principal;
`, schema)
	} else {
		createSchema = fmt.Sprintf(`USE SCHEMA %[1]s;
GRANT USE SCHEMA, CREATE FUNCTION, CREATE MATERIALIZED VIEW, CREATE MODEL, CREATE TABLE, CREATE VOLUME ON %[1]s TO $automation_principal;
`, schema)
	}
	//NOTE: default column values only available on Delta Lake 3.1.0+, falling back to current_timestamp() in insert
	createTable := fmt.Sprintf(`CREATE TABLE %[1]s (version STRING, run_on TIMESTAMP) USING DELTA TBLPROPERTIES (delta.enableChangeDataFeed = true);
ALTER TABLE %[1]s SET OWNER TO $owner_principal; --could also be $automation_principal, depending on security posture
GRANT SELECT, MODIFY ON %[1]s TO $automation_principal;
INSERT INTO %[1]s VALUES ("00000000000000", current_timestamp());
`, table)

	return strings.Join([]string{createCatalog, createSchema, createTable}, "\n")
}

func downSql(catalog string, schema string, table string) string {
	use := fmt.Sprintf("USE CATALOG %s;\nUSE SCHEMA %s;", catalog, schema)
	dropTable := fmt.Sprintf("DROP TABLE %s;", table)
	revokeSchema := fmt.Sprintf("REVOKE USE SCHEMA, CREATE FUNCTION, CREATE MATERIALIZED VIEW, CREATE MODEL, CREATE TABLE, CREATE VOLUME ON %s FROM $automation_principal;", schema)
	revokeCatalog := fmt.Sprintf("REVOKE USE CATALOG, CREATE SCHEMA ON %s FROM $automation_principal;", catalog)

	return strings.Join([]string{use, dropTable, revokeSchema, revokeCatalog}, "\n")
}

func GenerateSetup(directory string, catalog string, schema string, table string) {
	writeMigration(directory, "setup", "00000000000000",
		upSql(catalog, schema, table), downSql(catalog, schema, table),
	)
}

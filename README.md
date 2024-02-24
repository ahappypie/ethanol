# Ethanol - run migrations on your Unity Catalog
***Heavily inspired*** by the excellent [diesel_migrations](https://docs.rs/diesel_migrations/latest/diesel_migrations/) library, 
for the Rust ORM [diesel.rs](https://diesel.rs), thus the name, `ethanol`.

### Design
Read from a specific directory structure and apply migrations to Unity Catalog:
```
migrations/
    YYYYMMddHHmmss-$name/
        up.sql
        down.sql
```
Changes will be tracked in an internal table at `$CATALOG.$SCHEMA.ethanol_migrations` (`$SCHEMA=default`, by default).

This poses a challenge for bootstrapping `ethanol` since the Catalog does not exist or may not have proper permissions. 
Therefore, users must run a bootstrap query themselves prior to setting up automation:
```sql
--00000000000000_bootstrap/up.sql
CREATE CATALOG my_catalog;
ALTER CATALOG my_catalog SET OWNER TO owner_principal;
USE CATALOG my_catalog;
GRANT USE CATALOG, CREATE SCHEMA ON my_catalog TO `automation_principal`;
    
--if using a schema other than `default`
CREATE SCHEMA my_schema;
ALTER SCHEMA my_schema SET OWNER TO owner_principal;
USE SCHEMA my_schema;
-- USE SCHEMA default;
GRANT USE SCHEMA, CREATE FUNCTION, CREATE MATERIALIZED VIEW, CREATE MODEL, CREATE TABLE, CREATE VOLUME ON my_schema TO `automation_principal`;

CREATE TABLE ethanol_migrations (version STRING, run_on TIMESTAMP DEFAULT current_timestamp())
    USING DELTA;
ALTER TABLE ethanol_migrations SET OWNER TO owner_principal; --could also be `automation_principal`, depending on security posture


GRANT SELECT, MODIFY ON ethanol_migrations TO `automation_principal`;

INSERT INTO ethanol_migrations VALUES ("00000000000000");
```
In the future maybe a command similar to `ethanol setup` can help create these queries.

Similarly, different securables may have different lifecycles - i.e. Catalogs & Schemas must exist before Tables, Views, Volumes, etc.

In a given set of changes, you must ensure the versions on the migrations follow lifecycle rules - this tool is naive, and will not 
generate a graph of operations for you.

### Commands & Arguments
* `ethanol setup` - generates (but does *not* run) bootstrap as discussed above. You ***should*** review these queries and adjust as necessary.
* `ethanol migration generate $name` - creates a version (timestamp) + name directory with up and down SQL files
* `ethanol migration run` - runs all pending migrations, determined by the tracking table
* `ethanol migration revert` - runs `down.sql` for the last migration, determined by the tracking table
* `ethanol migration redo` - reverts, then runs. Useful shortcut for development and testing.

The following arguments or enviroment variables are required:
* `DIRECTORY, --directory, -d` - the relative path to your `/migrations` directory
* `INTERNAL_CATALOG, --internal-catalog, -i` - catalog where the tracking table was created
  * If the schema and/or table name were customized,
  
    `INTERNAL_SCHEMA, --internal-schema` and `INTERNAL_TABLE, --internal-table`

The following optional arguments may be passed via flag or environment variable to the `migration` subcommand:
* `CATALOG, --use-catalog, -c` - adds `USE CATALOG` to your query, useful for testing
* `SCHEMA, --use-schema, -s` - adds `USE SCHEMA` to your query, useful for testing
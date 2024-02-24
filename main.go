package main

import (
	"github.com/ahappypie/ethanol/cmd/ethanol"
)

func main() {
	//do cli things
	ethanol.Execute()
	//before setting up a cluster connection, check the existence of the migrations path and its structure

	// Pinging should require logging in
	//if err := db.Ping(); err != nil {
	//	fmt.Println(err)
	//}

	//ogCtx := dbsqlctx.NewContextWithCorrelationId(context.Background(), "createdrop-example")
	//
	//// create a table with some data. This has no context timeout, it will follow the timeout of one minute set for the connection.
	//if _, err := db.ExecContext(ogCtx, `CREATE TABLE IF NOT EXISTS diamonds USING CSV LOCATION '/databricks-datasets/Rdatasets/data-001/csv/ggplot2/diamonds.csv' options (header = true, inferSchema = true)`); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if _, err := db.ExecContext(ogCtx, `DROP TABLE diamonds `); err != nil {
	//	log.Fatal(err)
	//}
}

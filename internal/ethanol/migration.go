package ethanol

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

func frontmatter(version string, name string, direction string, ts time.Time) string {
	return fmt.Sprintf(`--%s_%s/%s.sql
--created_by: ethanol
--version: 0.0.1
--created_at: %s
`, version, name, direction, ts,
	)
}

func GenerateMigration(directory string, name string) {
	writeMigration(directory, name, "", "", "")
}

func writeMigration(directory string, name string, version string, up string, down string) {
	//check for existence of migrations directory
	if _, err := os.Stat("./" + directory); err != nil {
		if os.IsNotExist(err) {
			// file does not exists
			log.Fatal(err)
		} else {
			// file exists
			//continue
		}
	}
	ts := time.Now()
	if version == "" {
		//create version string
		//YYYYMMddHHmmss
		version = ts.Format("20060102150405")
	}
	//create directory for this migration
	dir := "./" + directory + "/" + version + "_" + name
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		//couldn't create the directory for some reason
		log.Fatal(err)
	}
	//get frontmatter
	upsql := frontmatter(version, name, "up", ts) + "\n" + up
	downsql := frontmatter(version, name, "down", ts) + "\n" + down
	//write up/down sql file content
	if err := os.WriteFile(dir+"/"+"up.sql", []byte(upsql), os.ModePerm); err != nil {
		//couldn't write up.sql for some reason
		log.Fatal(err)
	}
	if err := os.WriteFile(dir+"/"+"down.sql", []byte(downsql), os.ModePerm); err != nil {
		//couldn't write down.sql for some reason
		log.Fatal(err)
	}
}

type MigrationDirection int8

const (
	UpDirection MigrationDirection = iota
	DownDirection
)

type Migration struct {
	Version   string
	Name      string
	Sql       string
	Direction MigrationDirection
}

func (m Migration) ParseStatements() []string {
	//this driver and golang in general do not support multiline sql statements
	//we're going to break up each statement (and potentially link them with a correlation_id)
	//ignore comments, i.e. lines starting with --
	var statements []string
	for _, stmt := range strings.Split(m.Sql, "\n") {
		if len(stmt) > 0 && !strings.HasPrefix(stmt, "--") {
			statements = append(statements, stmt)
		}
	}
	return statements
}

func buildClient(catalog string, schema string, table string, clientId string, clientSecret string, host string, httpPath string) *Client {
	var client *Client
	if clientId != "" && clientSecret != "" {
		client = NewDBSQLClientWithM2M(
			&ClusterArgs{Host: host, HttpPath: httpPath},
			&InternalTableArgs{Catalog: catalog, Schema: schema, Table: table},
			&M2MArgs{ClientId: clientId, ClientSecret: clientSecret},
		)
	} else {
		fmt.Println("clientId and clientSecret not given, attempting U2M auth...")
		client = NewDBSQLClientWithU2M(
			&ClusterArgs{Host: host, HttpPath: httpPath},
			&InternalTableArgs{Catalog: catalog, Schema: schema, Table: table},
		)
	}
	return client
}

func getMigrationSetFromDisk(directory string, lastVersion string, direction MigrationDirection) []Migration {
	dirs, err := os.ReadDir("./" + directory)
	if err != nil {
		log.Fatal(err)
	}
	var migrations []Migration
	for _, dir := range dirs {
		if dir.IsDir() {
			version, name := parseVersion(dir.Name())
			if direction == UpDirection && version > lastVersion {
				//read up.sql
				sql, err := os.ReadFile("./" + directory + "/" + dir.Name() + "/up.sql")
				if err != nil {
					log.Fatal(err)
				}
				//add up.sql to migration set
				migrations = append(migrations, Migration{Version: version, Name: name, Sql: string(sql), Direction: UpDirection})
			} else if direction == DownDirection && version == lastVersion {
				sql, err := os.ReadFile("./" + directory + "/" + dir.Name() + "/down.sql")
				if err != nil {
					log.Fatal(err)
				}
				//add down.sql to migration set
				migrations = append(migrations, Migration{Version: version, Name: name, Sql: string(sql), Direction: DownDirection})
			}
		}
	}
	return migrations
}

func RunMigration(directory string, catalog string, schema string, table string, clientId string, clientSecret string, host string, httpPath string) {
	//build client
	client := buildClient(catalog, schema, table, clientId, clientSecret, host, httpPath)
	//get latest migration from tracking table
	lastTrackedVersion, err := client.GetLastMigration()
	if err != nil {
		log.Fatal(err)
	}
	lastVersion, _ := parseVersion(lastTrackedVersion)
	//get pending migration set from disk
	migrations := getMigrationSetFromDisk(directory, lastVersion, UpDirection)
	//run pending migrations in order
	slices.SortFunc(migrations, func(a, b Migration) int {
		return cmp.Compare(a.Version, b.Version)
	})
	for _, m := range migrations {
		err := client.ExecuteMigration(m)
		if err != nil {
			log.Fatal(err)
		}
	}
	client.Close()
}

func parseVersion(migration string) (string, string) {
	s := strings.Split(migration, "_")
	version := s[0]
	name := ""
	if len(s) > 1 {
		name = s[1]
	}
	return version, name
}

func RevertMigration(directory string, catalog string, schema string, table string, clientId string, clientSecret string, host string, httpPath string) {
	client := buildClient(catalog, schema, table, clientId, clientSecret, host, httpPath)
	//get latest migration from tracking table
	lastTrackedVersion, err := client.GetLastMigration()
	if err != nil {
		log.Fatal(err)
	}
	//get down.sql from disk
	migrations := getMigrationSetFromDisk(directory, lastTrackedVersion, DownDirection)
	if len(migrations) > 0 {
		err = client.ExecuteMigration(migrations[0])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("no migrations found for version %s", lastTrackedVersion)
	}
	client.Close()
}

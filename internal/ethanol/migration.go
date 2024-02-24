package ethanol

import (
	"fmt"
	"os"
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
			panic(err)
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
		panic(err)
	}
	//get frontmatter
	upsql := frontmatter(version, name, "up", ts) + "\n" + up
	downsql := frontmatter(version, name, "down", ts) + "\n" + down
	//write up/down sql file content
	if err := os.WriteFile(dir+"/"+"up.sql", []byte(upsql), os.ModePerm); err != nil {
		//couldn't write up.sql for some reason
		panic(err)
	}
	if err := os.WriteFile(dir+"/"+"down.sql", []byte(downsql), os.ModePerm); err != nil {
		//couldn't write down.sql for some reason
		panic(err)
	}
}

func RunMigration(directory string) {

}

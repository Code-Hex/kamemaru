package util

import "os"

const configToml = `[database]
dbname = "kamemaru"
host = "localhost"
port = 5432
sslmode = "disable"
user = ""
pass = ""
`

func CreateConfig() error {
	if !IsExist("config.toml") {
		if err := Create("config.toml", configToml); err != nil {
			return err
		}
	}

	return nil
}

func Create(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	f.WriteString(content)
	f.Close()
	return nil
}

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

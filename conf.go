/*config read to verify normal user*/
package main

import (
	"github.com/arstercz/goconfig"
)

func get_config(conf string) (c *goconfig.ConfigFile, err error) {
	c, err = goconfig.ReadConfigFile(conf)
	if err != nil {
		return c, err
	}
	return c, nil
}

func get_backend_dsn(c *goconfig.ConfigFile) (dsn string, err error) {
	dsn, err = c.GetString("backend", "dsn")
	if err != nil {
		return dsn, err
	}
	return dsn, nil
}

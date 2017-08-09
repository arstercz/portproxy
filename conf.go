/*config read to verify normal user*/
package main

import (
	"github.com/chenzhe07/goconfig"
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

func get_read_user(c *goconfig.ConfigFile) (user string, pass string, err error) {
	user, err = c.GetString("onlineread", "user")
	if err != nil {
		return "", "", err
	}
	pass, err = c.GetString("onlineread", "pass")
	if err != nil {
		return user, "", err
	}
	return user, pass, nil
}

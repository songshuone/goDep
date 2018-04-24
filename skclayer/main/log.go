package main

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"fmt"
)

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = Conf.LogPath
	config["level"] = convertLogLevel(Conf.LogLevel)
	configstr, err := json.Marshal(config)
	if err != nil {
		err = fmt.Errorf("json Marshal failed err:%v", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configstr))
	return
}

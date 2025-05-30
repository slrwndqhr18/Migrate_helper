package handlers

import (
	"fmt"
	"os"
	"time"
)

type LOG_CTRL struct {
	path       string
	LogStorage []string
	ErrCnt     int
}

func Init_logger(_logDirPath string) LOG_CTRL {
	currentDate := time.Now().Format("2006-01-02")
	return LOG_CTRL{
		path:       _logDirPath + "/log_" + currentDate + ".log",
		LogStorage: []string{},
		ErrCnt:     0,
	}
}

func (_LOGGER *LOG_CTRL) Add_log(_level string, _mesg ...string) {
	switch len(_mesg) {
	case 1:
		_LOGGER.LogStorage = append(_LOGGER.LogStorage, fmt.Sprintf("[%s] %s\n", _level, _mesg[0]))
	default:
		_LOGGER.LogStorage = append(_LOGGER.LogStorage, fmt.Sprintf("[%s] %s\n", _level, _mesg[0]))
		for i := 1; i < len(_mesg); i++ {
			_LOGGER.LogStorage = append(_LOGGER.LogStorage, fmt.Sprintf("\tㄴ %s\n", _mesg[i]))
		}
	}
	if _level == "ERROR" {
		_LOGGER.ErrCnt += 1
	}
}

func (_LOGGER LOG_CTRL) Is_ok() bool {
	return _LOGGER.ErrCnt == 0
}

func (_LOGGER LOG_CTRL) Write_to_file() {
	if len(_LOGGER.LogStorage) == 0 {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	logString := fmt.Sprintf("TRANSACTION >> %s\n", now)
	for _, line := range _LOGGER.LogStorage {
		logString += line
	}
	now = time.Now().Format("2006-01-02 15:04:05")
	logString += fmt.Sprintf("TRANSACTION END >> %s\n\n", now)
	f, err := os.OpenFile(_LOGGER.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		if err != nil {
			fmt.Println(logString)
			fmt.Println("[ERROR] Failed to write log at file\n", err.Error())
		}
	}()
	if err == nil {
		defer f.Close()
		_, err = f.WriteString(logString)
	}
	fmt.Println("[INFO] Messages are stored in log file:", _LOGGER.path)
}

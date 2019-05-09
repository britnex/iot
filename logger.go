package logging

import (
	"encoding/json"
	"fmt"
	"time"
)

const (

	LDebug int = 0
	LInfo int = 1
	LWarning int =2
	LError int = 3
	LFatal int = 4

)

type Logger struct {
	Level int
}

type LogEntry struct {
	Level      int
	Interfaces map[string]interface{}
	Logger *Logger
}

func (self *LogEntry) Pair(name string, value interface{}) *LogEntry {
	self.Interfaces[name] = value
	return self
}

func (self *LogEntry) String(name string, value string) *LogEntry {
	self.Interfaces[name] = value
	return self
}

func (self *LogEntry) Number(name string, value int) *LogEntry {
	self.Interfaces[name] = value
	return self
}

func (self *LogEntry) Boolean(name string, value bool) *LogEntry {
	self.Interfaces[name] = value
	return self
}

func (self *LogEntry) Time(name string, value time.Time) *LogEntry {
	self.Interfaces[name] = value
	return self
}

func (self *LogEntry) Interface(name string, value interface{}) *LogEntry {
	self.Interfaces[name] = value
	return self
}

func (self *LogEntry) Msg(format string, v ...interface{}) {
	
	if self.Level < self.Logger.Level {
		return
	}

	jsonb, _ := json.Marshal(self.Interfaces)
	jsons := "\t" + string(jsonb)
	if len(jsons) == 3 {
		jsons = ""
	}
	
	level:=""

	switch self.Level {
	    case 0: level="DEBG"
	    case 1: level="INFO"
	    case 2: level="WARN"
	    case 3: level="ERRO"
	    case 4: level="FATA"
	}

	fmt.Println("[" + level + "] " + time.Now().Format(time.RFC3339) + " " + fmt.Sprintf(format, v...) + jsons)
	return
}


func (self *Logger) Debug() *LogEntry {
	return &LogEntry{Level: 0, Logger:self, Interfaces: make(map[string]interface{})}
}

func (self *Logger) Info() *LogEntry {
	return &LogEntry{Level: 1, Logger:self, Interfaces: make(map[string]interface{})}
}

func (self *Logger) Warning() *LogEntry {
	return &LogEntry{Level: 2, Logger:self, Interfaces: make(map[string]interface{})}
}

func (self *Logger) Error() *LogEntry {
	return &LogEntry{Level: 3, Logger:self, Interfaces: make(map[string]interface{})}
}

func (self *Logger) Fatal() *LogEntry {
	return &LogEntry{Level: 4, Logger:self, Interfaces: make(map[string]interface{})}
}

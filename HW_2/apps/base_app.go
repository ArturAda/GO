package apps

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	defaultHostname = "localhost"
	defaultVersion  = "v1.0.0"
)

type logStructure struct {
	output io.Writer
	prefix string
	flags  int
}

type loggFunc func(*logStructure)

func SetOutput(newOut io.Writer) loggFunc {
	return func(logg *logStructure) {
		if newOut != nil {
			logg.output = newOut
		}
	}
}

func SetPrefix(newPrefix string) loggFunc {
	return func(logg *logStructure) {
		logg.prefix = newPrefix
	}
}

func SetFlags(newFlags int) loggFunc {
	return func(logg *logStructure) {
		logg.flags = newFlags
	}
}

type BaseApp struct {
	logg       *log.Logger
	name       string
	version    string
	lastConfig logStructure
}

func NewBaseApp(newName, newVersion string, functions ...loggFunc) *BaseApp {
	if newName == "" {
		newName = defaultHostname
	}
	if newVersion == "" {
		newVersion = defaultVersion
	}
	defaultLogg := logStructure{
		output: os.Stdout,
		prefix: fmt.Sprintf("[%s]", newName),
		flags:  log.LstdFlags | log.Lmicroseconds | log.Llongfile,
	}
	for _, function := range functions {
		if function != nil {
			function(&defaultLogg)
		}
	}
	logger := log.New(defaultLogg.output, defaultLogg.prefix, defaultLogg.flags)
	return &BaseApp{name: newName, version: newVersion, logg: logger, lastConfig: defaultLogg}
}

func (app *BaseApp) GetName() string { return app.name }

func (app *BaseApp) GetVersion() string { return app.version }

func (app *BaseApp) GetLog() *log.Logger { return app.logg }

func (app *BaseApp) UpdateLogger(functions ...loggFunc) {
	newConfig := app.lastConfig
	for _, function := range functions {
		if function != nil {
			function(&newConfig)
		}
	}
	if newConfig.output == nil {
		newConfig.output = os.Stdout
	}
	if newConfig.prefix == "" {
		newConfig.prefix = fmt.Sprintf("[%s]", app.name)
	}
	if newConfig.flags == 0 {
		newConfig.flags = log.LstdFlags | log.Lmicroseconds | log.Llongfile
	}
	app.logg.SetOutput(newConfig.output)
	app.logg.SetPrefix(newConfig.prefix)
	app.logg.SetFlags(newConfig.flags)
	app.lastConfig = newConfig
}

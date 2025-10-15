package apps

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	randv2 "math/rand/v2"
	"net/http"
	"time"
)

type errorJSON struct {
	Message string `json:"message"`
}

func getNewErrorJSON(message string) errorJSON {
	return errorJSON{Message: message}
}

func writeJSON(output http.ResponseWriter, status int, value any, logger *log.Logger) {
	output.Header().Set("Content-Type", "application/json; charset=utf-8")
	bytes_, err := json.Marshal(value)
	if err != nil {
		if logger != nil {
			logger.Printf("encode error: %v", err)
		}
		http.Error(output, `{"message":"internal error"}`, http.StatusInternalServerError)
		return
	}
	output.WriteHeader(status)
	_, err = output.Write(bytes_)
	if err != nil && logger != nil {
		logger.Printf("write error: %v", err)
	}
	if logger != nil {
		var pretty bytes.Buffer
		err = json.Indent(&pretty, bytes_, "", "  ")
		if err == nil {
			logger.Printf("responded %d:\n%s", status, pretty.String())
		} else {
			logger.Printf("responded %d %s", status, string(bytes_))
		}
	}
}

type HTTPApp struct {
	*BaseApp
	port           string
	closingAppTime time.Duration
	server         *http.Server
	mux            *http.ServeMux
}

type httpFunc func(*HTTPApp)

func SetName(newName string) httpFunc {
	return func(app *HTTPApp) {
		if newName != "" {
			app.name = newName
			app.UpdateLogger(SetPrefix("[" + newName + "]"))
		}
	}
}

func SetVersion(newVersion string) httpFunc {
	return func(app *HTTPApp) {
		if newVersion != "" {
			app.version = newVersion
		}
	}
}

func SetLoggOutput(newOutput io.Writer) httpFunc {
	return func(app *HTTPApp) {
		app.UpdateLogger(SetOutput(newOutput))
	}
}

func SetLoggPrefix(newPrefix string) httpFunc {
	return func(app *HTTPApp) {
		app.UpdateLogger(SetPrefix(newPrefix))
	}
}

func SetLoggFlags(newFlags int) httpFunc {
	return func(app *HTTPApp) {
		app.UpdateLogger(SetFlags(newFlags))
	}
}

func SetPort(newPort string) httpFunc {
	return func(app *HTTPApp) {
		if newPort != "" {
			app.port = newPort
		}
	}
}

func SetClosingAppTime(newClosingAppTime time.Duration) httpFunc {
	return func(app *HTTPApp) {
		if newClosingAppTime > 0 {
			app.closingAppTime = newClosingAppTime
		}
	}
}

func NewHTTPApp(functions ...httpFunc) *HTTPApp {
	newApp := &HTTPApp{
		BaseApp:        NewBaseApp("localhost", "1.0.0"),
		port:           ":8080",
		closingAppTime: 20 * time.Second,
		mux:            http.NewServeMux(),
	}
	for _, function := range functions {
		if function != nil {
			function(newApp)
		}
	}
	newApp.mux.HandleFunc("GET /version", newApp.getVersionHandler)
	newApp.mux.HandleFunc("POST /decode", newApp.getDecodeHandler)
	newApp.mux.HandleFunc("GET /hard-op", newApp.getHardOpHandler)
	newApp.server = &http.Server{
		Addr:         newApp.port,
		Handler:      newApp.logMiddleware(newApp.mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 25 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return newApp
}

func (app *HTTPApp) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(output http.ResponseWriter, request *http.Request) {
		currentTime := time.Now()
		next.ServeHTTP(output, request)
		app.GetLog().Printf("%s %s %s", request.Method, request.URL.Path, time.Since(currentTime))
	})
}

func (app *HTTPApp) Start() error {
	app.GetLog().Printf("HTTP listening on %s (version %s)", app.port, app.GetVersion())
	go func() {
		err := app.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.GetLog().Printf("ListenAndServe error: %v", err)
		}
	}()
	return nil
}

func (app *HTTPApp) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, app.closingAppTime)
	defer cancel()
	app.GetLog().Printf("graceful shutdown: server %q...", app.GetName())
	err := app.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("http shutdown failed for %q: %w", app.GetName(), err)
	}
	app.GetLog().Printf("server %q stopped gracefully", app.GetName())
	return nil
}

func (app *HTTPApp) getVersionHandler(output http.ResponseWriter, _ *http.Request) {
	output.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := output.Write([]byte(app.version + "\n"))
	if err != nil {
		app.GetLog().Printf("write error: %v", err)
	}
}

func (app *HTTPApp) getDecodeHandler(output http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(output, request.Body, 1<<20)
	output.Header().Set("Content-Type", "application/json; charset=utf-8")
	var inputJSON struct {
		InputString string `json:"inputString"`
	}
	err := json.NewDecoder(request.Body).Decode(&inputJSON)
	if err != nil {
		writeJSON(output, http.StatusBadRequest, getNewErrorJSON("Wrong JSON"), app.GetLog())
		return
	}
	if inputJSON.InputString == "" {
		writeJSON(output, http.StatusBadRequest, getNewErrorJSON("InputString is empty"), app.GetLog())
		return
	}
	decodeInputString, err := base64.StdEncoding.DecodeString(inputJSON.InputString)
	if err != nil {
		writeJSON(output, http.StatusBadRequest, getNewErrorJSON("InputString is invalid"), app.GetLog())
		return
	}
	type decodeJSON struct {
		OutputString string `json:"outputString"`
	}
	writeJSON(output, http.StatusOK, decodeJSON{OutputString: string(decodeInputString)}, app.GetLog())
}

func (app *HTTPApp) getHardOpHandler(output http.ResponseWriter, request *http.Request) {
	globalTime := 10 + randv2.IntN(11)
	nanoTime := time.NewTimer(time.Duration(globalTime) * time.Second)
	defer nanoTime.Stop()
	select {
	case <-request.Context().Done():
		return
	case <-nanoTime.C:
	}
	if randv2.IntN(2) == 0 {
		writeJSON(output, http.StatusInternalServerError, getNewErrorJSON("Hard operation timeout"), app.GetLog())
		return
	}
	type correctJSON struct {
		Status   string `json:"status"`
		WorkTime int    `json:"work_time"`
	}
	writeJSON(output, http.StatusOK, correctJSON{Status: "All done!", WorkTime: globalTime}, app.GetLog())
}

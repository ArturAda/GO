package apps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

type HTTPUser struct {
	UserName   string
	httpServer *HTTPApp
	enterTime  time.Time
	enter      bool
}

type httpUserFunc func(*HTTPUser)

func (user *HTTPUser) enterNow() {
	user.enterTime = time.Now()
}

func SetUserName(name string) httpUserFunc {
	return func(httpUser *HTTPUser) {
		httpUser.UserName = name
	}
}

func SetHTTPServer(httpServer *HTTPApp) httpUserFunc {
	return func(httpUser *HTTPUser) {
		httpUser.httpServer = httpServer
	}
}

func SetEnterFlag(enterFlag bool) httpUserFunc {
	return func(httpUser *HTTPUser) {
		if enterFlag {
			httpUser.enterNow()
		}
		httpUser.enter = enterFlag
	}
}

func NewClient(functions ...httpUserFunc) *HTTPUser {
	newClient := &HTTPUser{UserName: "", httpServer: nil, enter: false}
	for _, function := range functions {
		if function != nil {
			function(newClient)
		}
	}
	return newClient
}

func (user *HTTPUser) SetNewValues(functions ...httpUserFunc) {
	for _, function := range functions {
		if function != nil {
			function(user)
		}
	}
}

func (user *HTTPUser) ChangeEnter(newEnterFlag bool) {
	user.enter = newEnterFlag
}

func (user *HTTPUser) ChangeName(newName string) {
	user.UserName = newName
}

func (user *HTTPUser) GetName() string {
	return user.UserName
}

func (user *HTTPUser) GetHttpServer() *HTTPApp {
	return user.httpServer
}

func (user *HTTPUser) GetEnterFlag() bool {
	return user.enter
}

func (user *HTTPUser) TimeOnSite() time.Duration {
	if user.enterTime.IsZero() {
		return 0
	}
	return time.Since(user.enterTime)
}

func (user *HTTPUser) check(ctx context.Context, method, path string, body []byte, contentType string) (int, []byte, error) {
	if user.httpServer == nil || user.httpServer.mux == nil {
		return 0, nil, fmt.Errorf("The client didn't select a server\n")
	}
	request, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	recorder := httptest.NewRecorder()
	user.httpServer.mux.ServeHTTP(recorder, request)
	res := recorder.Result()
	defer res.Body.Close()
	slc, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, slc, nil
}

type allConfig struct {
	contentType   string
	inputString   string
	hardOpTimeOut time.Duration
}

type configFunc func(*allConfig)

func SetContentType(newContentType string) configFunc {
	return func(config *allConfig) {
		config.contentType = newContentType
	}
}

func SetInputString(newInputString string) configFunc {
	return func(config *allConfig) {
		config.inputString = newInputString
	}
}

func SetHardOpTimeOut(newHardOpTimeOut time.Duration) configFunc {
	return func(config *allConfig) {
		if newHardOpTimeOut > 0 {
			config.hardOpTimeOut = newHardOpTimeOut
		}
	}
}

func (user *HTTPUser) checkAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, _, err := user.check(ctx, http.MethodGet, "/version", nil, "")
	if err != nil {
		return fmt.Errorf("server not active: %w", err)
	}
	return nil
}

func (user *HTTPUser) checkVersion() error {
	status, body, err := user.check(context.Background(), http.MethodGet, "/version", nil, "")
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("version unexpected status: %d", status)
	}
	fmt.Print(string(body))
	return nil
}

func (user *HTTPUser) checkDecode(cfg *allConfig) error {
	type inputStruct struct {
		InputString string `json:"inputString"`
	}
	type outputStruct struct {
		OutputString string `json:"outputString"`
	}
	someInput := inputStruct{InputString: cfg.inputString}
	byteInput, _ := json.Marshal(someInput)
	status, body, err := user.check(context.Background(), http.MethodPost, "/decode", byteInput, cfg.contentType)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("decode failed: %s", string(body))
	}
	var output outputStruct
	err = json.Unmarshal(body, &output)
	if err != nil {
		return fmt.Errorf("decode parse: %w", err)
	}
	fmt.Println(output.OutputString)
	return nil
}

func (user *HTTPUser) checkHardOp(cfg *allConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.hardOpTimeOut)
	defer cancel()
	status, body, err := user.check(ctx, http.MethodGet, "/hard-op", nil, "")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("false")
			return nil
		}
		return fmt.Errorf("hard-op: %w", err)
	}
	fmt.Printf("true, %d\n", status)
	if status == http.StatusOK {
		var ok struct {
			Status   string `json:"status"`
			WorkTime int    `json:"work_time"`
		}
		err = json.Unmarshal(body, &ok)
		if err == nil {
			fmt.Printf("%s (work_time=%ds)\n", ok.Status, ok.WorkTime)
		}
	}
	return nil
}

func (user *HTTPUser) RunAll(opts ...configFunc) error {
	if user.httpServer == nil {
		return fmt.Errorf("The client didn't select a server\n")
	}
	cfg := allConfig{
		contentType:   "application/json; charset=utf-8",
		inputString:   "yyy",
		hardOpTimeOut: 15 * time.Second,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if err := user.checkAvailability(); err != nil {
		return err
	}
	if err := user.checkVersion(); err != nil {
		return err
	}
	if err := user.checkDecode(&cfg); err != nil {
		return err
	}
	if err := user.checkHardOp(&cfg); err != nil {
		return err
	}
	return nil
}

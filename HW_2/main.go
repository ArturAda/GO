package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/ArturAda/GO/HW_2/apps"
)

func main() {
	s1 := apps.NewHTTPApp(
		apps.SetName("srv-1"),
		apps.SetVersion("v1.0.0"),
		apps.SetPort(":8081"),
	)
	s2 := apps.NewHTTPApp(
		apps.SetName("srv-2"),
		apps.SetVersion("v1.1.0"),
		apps.SetPort(":8082"),
		apps.SetClosingAppTime(10*time.Second),
	)
	customLog := log.New(os.Stdout, "[custom] ", log.LstdFlags|log.Lshortfile)
	s3 := apps.NewHTTPApp(
		apps.SetName("srv-3"),
		apps.SetVersion("v2.0.0"),
		apps.SetPort(":8083"),
		apps.SetLoggOutput(os.Stdout),
		apps.SetLoggPrefix("[srv-3] "),
		apps.SetLoggFlags(log.LstdFlags|log.Lmicroseconds),
	)
	s3.UpdateLogger(apps.SetOutput(customLog.Writer()), apps.SetPrefix("[srv-3-custom]"))
	s4 := apps.NewHTTPApp(
		apps.SetName("srv-4"),
		apps.SetVersion("v3.0.0"),
		apps.SetPort(":8084"),
	)
	s5 := apps.NewHTTPApp(
		apps.SetName("srv-5"),
		apps.SetVersion("v1.2.3"),
		apps.SetPort(":8085"),
	)
	servers := []*apps.HTTPApp{s1, s2, s3, s4, s5}
	for _, s := range servers {
		err := s.Start()
		if err != nil {
			log.Fatalf("start %s failed: %v", s.GetName(), err)
		}
	}
	time.Sleep(150 * time.Millisecond)
	user := apps.NewClient(
		apps.SetUserName("Alice"),
		apps.SetHTTPServer(s1),
		apps.SetEnterFlag(true),
	)
	log.Printf("client %q entered server %q; timeOnSite=%s",
		user.GetName(), user.GetHttpServer().GetName(), user.TimeOnSite().String())
	b64HelloWorld := base64.StdEncoding.EncodeToString([]byte("Hello, World!"))
	err := user.RunAll(apps.SetInputString(b64HelloWorld),
		apps.SetContentType("application/json; charset=utf-8"),
		apps.SetHardOpTimeOut(15*time.Second))
	if err != nil {
		log.Printf("RunAll on %s failed: %v", user.GetHttpServer().GetName(), err)
	}
	b64A := base64.StdEncoding.EncodeToString([]byte("AAAAAAAAAAAAAAAAAAAAAAAAaaa"))
	err = user.RunAll(
		apps.SetInputString(b64A))
	if err != nil {
		log.Printf("RunAll on %s failed: %v", user.GetHttpServer().GetName(), err)
	}
	err = user.RunAll(
		apps.SetInputString(b64HelloWorld),
		apps.SetHardOpTimeOut(3*time.Second))
	if err != nil {
		log.Printf("RunAll on %s failed: %v", user.GetHttpServer().GetName(), err)
	}
	user.SetNewValues(apps.SetHTTPServer(s3))
	log.Printf("client %q switched to server %q", user.GetName(), user.GetHttpServer().GetName())
	err = user.RunAll(
		apps.SetInputString(b64A),
		apps.SetHardOpTimeOut(15*time.Second))
	if err != nil {
		log.Printf("RunAll on %s failed: %v", user.GetHttpServer().GetName(), err)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	for _, s := range servers {
		err := s.Stop(shutdownCtx)
		if err != nil {
			log.Printf("stop %s failed: %v", s.GetName(), err)
		}
	}
	log.Println("all servers stopped gracefully")
}

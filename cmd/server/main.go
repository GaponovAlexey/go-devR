package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"r/api/pkg/config"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	prdLoger, _ := zap.NewProduction()
	defer prdLoger.Sync() // flushes buffer, if any
	logger := prdLoger.Sugar()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalw("failed111111 to parse config", "err", err)
	}
	fmt.Printf("cfg = %v\n", cfg)

	//GORM
	db, err := gorm.Open(sql.Open(cfg))
	if err != nil {
		panic("failed to connect database")
	}
	//create repository

	// create service

	//create router
	r := mux.NewRouter()
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	//create http server
	srv := &http.Server{
		Handler: r,
		Addr:    cfg.HttpAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalw("faile to connect to db", "err", err, "http-addr", cfg.HttpAddr)
		}
	}()

	// waiting fo shutdown

	shutdown := make(chan struct{})
	// go func() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGTERM)
	<-signals
	// close(shutdown)
	// }()
	context.WithTimeout(context.TODO(), 15*time.Second)
	<-shutdown
}

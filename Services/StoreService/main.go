package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	InitPostgresDB()
	InitMetrics()
	r := InitRouter()
	runApiServer(ctx, r)
	CloseConn()
}

func runApiServer(ctx context.Context, r *gin.Engine) {
	server := &http.Server{
		Handler: r.Handler(),
		Addr: ":8080",
	}

	go func() {
		fmt.Printf("Starting listening on: %s\n", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTP server error: %v\n", err)
			return
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP shutdown error: %v", err)
		if err := server.Close(); err != nil {
			fmt.Printf("HTTP close error: %v\n", err)
		}
	} else {
		fmt.Println("Graceful shutdown complete.")
	}
}

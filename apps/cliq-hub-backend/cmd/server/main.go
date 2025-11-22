package main

import (
    "log"
    "net/http"
    "os"

    "cliq-hub-backend/internal/config"
    "cliq-hub-backend/internal/http/router"
    "cliq-hub-backend/internal/llm"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config error: %v", err)
    }

    client, err := llm.NewClient(cfg)
    if err != nil {
        log.Fatalf("llm client init error: %v", err)
    }

    debugMode := os.Getenv("GIN_MODE") != "release"
    r := router.New(client, debugMode)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    srv := &http.Server{Addr: ":" + port, Handler: r}
    log.Printf("server listening on :%s", port)
    if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
        log.Fatalf("server error: %v", err)
    }
}
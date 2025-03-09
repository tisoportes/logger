// demo.go
// Author: jm & zillion

package main

import (
    "github.com/tisoportes/logger"
)

func main() {
    // Initialize the logger
    err := logger.InitLogger(logger.LevelDebug, true, "logs/demo.log")
    if err != nil {
        panic(err)
    }
    defer logger.CloseLogger() // Ensure the logger is closed when done

    // Log messages at different levels
    logger.Debug("This is a debug message")
    logger.Info("This is an info message")
    logger.Warning("This is a warning message")
    logger.Error("This is an error message")
    logger.Fatal("This is a fatal message")
}

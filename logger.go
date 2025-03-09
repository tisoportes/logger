// File: logger.go
// Description:
// Package logger provides a simple and flexible logging utility for Go applications.
// It supports multiple log levels (debug, info, warning, error, fatal) and can log
// to both the console and a specified log file. The logger can be initialized with
// different configurations, including log level and output options. It also includes
// functionality for log rotation and capturing caller information for better debugging.
//
// Author: Juan Mamani & Zillion 
// Release Date: 2023-03-08

package logger

import (
        "fmt"
        "io"
        "log"
        "os"
        "path/filepath"
        "runtime"
        "strings"
        "time"
)

// Log levels
const (
        LevelDebug = iota
        LevelInfo
        LevelWarning
        LevelError
        LevelFatal
)

var (
        // Loggers for different levels
        debugLogger   *log.Logger
        infoLogger    *log.Logger
        warningLogger *log.Logger
        errorLogger   *log.Logger
        fatalLogger   *log.Logger

        // Current log level
        currentLevel = LevelInfo

        // Log file
        logFile *os.File
)

// InitLogger initializes the logging system
func InitLogger(level int, logToFile bool, logFileName string) error {
        currentLevel = level

        // Set up log format: timestamp, file:line, message
        flags := log.Ldate | log.Ltime | log.Lshortfile

        // Set up output writer(s)
        var writers []io.Writer
        writers = append(writers, os.Stdout) // Always log to stdout

        // If logging to file is enabled, set up the file writer
        if logToFile && logFileName != "" {
                // Create logs directory if it doesn't exist
                logsDir := filepath.Dir(logFileName)
                if _, err := os.Stat(logsDir); os.IsNotExist(err) {
                        if err := os.MkdirAll(logsDir, 0755); err != nil {
                                return fmt.Errorf("failed to create logs directory: %v", err)
                        }
                }

                // Open log file with append mode, create if doesn't exist
                var err error
                logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                if err != nil {
                        return fmt.Errorf("failed to open log file: %v", err)
                }

                writers = append(writers, logFile)
        }

        // Create a multiwriter if we have multiple outputs
        var output io.Writer
        if len(writers) == 1 {
                output = writers[0]
        } else {
                output = io.MultiWriter(writers...)
        }

        // Initialize loggers with appropriate prefixes
        debugLogger = log.New(output, "[DEBUG] ", flags)
        infoLogger = log.New(output, "[INFO] ", flags)
        warningLogger = log.New(output, "[WARN] ", flags)
        errorLogger = log.New(output, "[ERROR] ", flags)
        fatalLogger = log.New(output, "[FATAL] ", flags)

        // Use the default logger for general messages
        log.SetOutput(output)
        log.SetFlags(flags)
        log.SetPrefix("[LOG] ")

        return nil
}

// CloseLogger closes any open resources (like log files)
func CloseLogger() {
        if logFile != nil {
                logFile.Close()
        }
}

// GetLogger returns the appropriate logger based on the level
func getLogger(level int) *log.Logger {
        switch level {
        case LevelDebug:
                return debugLogger
        case LevelInfo:
                return infoLogger
        case LevelWarning:
                return warningLogger
        case LevelError:
                return errorLogger
        case LevelFatal:
                return fatalLogger
        default:
                return infoLogger
        }
}

// logWithCallerInfo logs a message with the caller info (file, line, function)
func logWithCallerInfo(level int, format string, v ...interface{}) {
        if level < currentLevel {
                return
        }

        logger := getLogger(level)

        // Get caller information
        _, file, line, ok := runtime.Caller(2)
        if ok {
                if format == "" {
                        msg := fmt.Sprint(v...)
                        logger.Printf("%s:%d: %s", filepath.Base(file), line, msg)
                } else {
                        msg := fmt.Sprintf(format, v...)
                        logger.Printf("%s:%d: %s", filepath.Base(file), line, msg)
                }
        } else {
                if format == "" {
                        logger.Print(v...)
                } else {
                        logger.Printf(format, v...)
                }
        }
}

// Debug logs a debug message
func Debug(v ...interface{}) {
        logWithCallerInfo(LevelDebug, "", v...)
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
        logWithCallerInfo(LevelDebug, format, v...)
}

// Info logs an info message
func Info(v ...interface{}) {
        logWithCallerInfo(LevelInfo, "", v...)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
        logWithCallerInfo(LevelInfo, format, v...)
}

// Warning logs a warning message
func Warning(v ...interface{}) {
        logWithCallerInfo(LevelWarning, "", v...)
}

// Warningf logs a formatted warning message
func Warningf(format string, v ...interface{}) {
        logWithCallerInfo(LevelWarning, format, v...)
}

// Error logs an error message
func Error(v ...interface{}) {
        logWithCallerInfo(LevelError, "", v...)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
        logWithCallerInfo(LevelError, format, v...)
}

// Fatal logs a fatal message and exits the program
func Fatal(v ...interface{}) {
        logWithCallerInfo(LevelFatal, "", v...)
        os.Exit(1)
}

// Fatalf logs a formatted fatal message and exits the program
func Fatalf(format string, v ...interface{}) {
        logWithCallerInfo(LevelFatal, format, v...)
        os.Exit(1)
}

// RotateLogFile rotates the log file (creates a new one with timestamp)
func RotateLogFile() error {
        if logFile == nil {
                return nil // No log file to rotate
        }

        // Close current log file
        logFile.Close()

        // Get the path and base filename
        dir, filename := filepath.Split(logFile.Name())
        ext := filepath.Ext(filename)
        baseFilename := strings.TrimSuffix(filename, ext)

        // Create a new filename with timestamp
        timestamp := time.Now().Format("20060102-150405")
        newFilename := fmt.Sprintf("%s-%s%s", baseFilename, timestamp, ext)
        newPath := filepath.Join(dir, newFilename)

        // Rename the old file
        err := os.Rename(logFile.Name(), newPath)
        if err != nil {
                return fmt.Errorf("failed to rename log file: %v", err)
        }

        // Open a new log file
        logFile, err = os.OpenFile(logFile.Name(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
                return fmt.Errorf("failed to open new log file: %v", err)
        }

        // Update the writers for all loggers
        writers := []io.Writer{os.Stdout, logFile}
        output := io.MultiWriter(writers...)

        debugLogger.SetOutput(output)
        infoLogger.SetOutput(output)
        warningLogger.SetOutput(output)
        errorLogger.SetOutput(output)
        fatalLogger.SetOutput(output)
        log.SetOutput(output)

        Info("Log file rotated to", newPath)
        return nil
}

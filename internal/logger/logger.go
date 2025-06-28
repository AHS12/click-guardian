package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Logger handles application logging with GUI display
type Logger struct {
	logChannel      chan string
	logHistory      []string
	logText         *widget.RichText
	scrollContainer *container.Scroll
	maxEntries      int
	once            sync.Once
}

// NewLogger creates a new logger instance
func NewLogger(logText *widget.RichText, scrollContainer *container.Scroll, maxEntries int) *Logger {
	if maxEntries <= 0 {
		maxEntries = 100
	}

	return &Logger{
		logChannel:      make(chan string, 100),
		logHistory:      make([]string, 0, maxEntries),
		logText:         logText,
		scrollContainer: scrollContainer,
		maxEntries:      maxEntries,
	}
}

// GetChannel returns the log channel for sending messages
func (l *Logger) GetChannel() chan string {
	return l.logChannel
}

// Start begins processing log messages
func (l *Logger) Start() {
	go func() {
		for msg := range l.logChannel {
			l.addLogEntry(msg)
		}
	}()
}

// Stop closes the log channel
func (l *Logger) Stop() {
	l.once.Do(func() {
		close(l.logChannel)
	})
}

// Clear clears the log display
func (l *Logger) Clear() {
	l.logHistory = l.logHistory[:0]
	fyne.Do(func() {
		l.logText.ParseMarkdown("**Log cleared**\n\n")
	})
}

// Log sends a message to the log
func (l *Logger) Log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	select {
	case l.logChannel <- msg:
	default:
		// Channel is full, drop the message
	}
}

func (l *Logger) addLogEntry(msg string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("**[%s]** %s", timestamp, msg)

	l.logHistory = append(l.logHistory, logEntry)

	// Keep only the last maxEntries to prevent memory issues
	if len(l.logHistory) > l.maxEntries {
		l.logHistory = l.logHistory[1:]
	}

	// Update the log display on the main UI thread
	allLogs := strings.Join(l.logHistory, "\n\n")
	fyne.Do(func() {
		l.logText.ParseMarkdown(allLogs)
		if l.scrollContainer != nil {
			l.scrollContainer.ScrollToBottom()
		}
	})
}

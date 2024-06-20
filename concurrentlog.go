package concurrentlog

import (
	"fmt"
	"log"
	"os"
	"io"
)

// Define Logger structure
type Logger struct {
	
	logFile		*os.File // Define location of log file
	logger		*log.Logger // Define the defaults log package
	msgChan		chan string // Message channel for goroutine to send messages into
	done		chan struct{} // Used for stoping the chennel
}


// The NewLogger function creates a new logger instance with a specified file path and buffer size for
// logging messages.
func NewLogger(filePath string, buffersize int) (*Logger, error) {
	// Open the logging location.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Add multiple writer to return stdout.
	multiWriter := io.MultiWriter(os.Stdout, file)
	
	// Logger Construction
	logger := &Logger{
		logFile: 	file,
		logger: 	log.New(multiWriter, "", log.LstdFlags),
		msgChan: 	make(chan string, buffersize),
		done: 		make(chan struct{}),
	}
	// Initialize the logger goroutine
	go logger.run()

	return logger, nil
}



// The `run` method of the `Logger` struct is a goroutine function that continuously listens for
// messages on the `msgChan` channel. It uses a select statement to either consume messages from the
// channel and log them using the internal logger or break the loop if a signal is received on the
// `done` channel.
func (l *Logger) run() {
	// infinity loop
	for {
		// Select case, if consume normal data from buffer channel -> logging.
		// If received done signal -> break the loop
		select {
		case msg := <- l.msgChan:
			l.logger.Println(msg)
		case <- l.done:
			return
		}
	}
}


// The `Log` method of the `Logger` struct is used to log messages with a specified level and content.
func (l *Logger) Log(level string, msg string) {
	payload := fmt.Sprintf("%s %s", level, msg)
	l.msgChan <- payload
}


// The `Close` method of the `Logger` struct is responsible for closing the logger instance. Here's
// what it does:
func (l *Logger) Close() error {
	close(l.done)
	close(l.msgChan)
	return l.logFile.Close()
}
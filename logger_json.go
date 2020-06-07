// (C)2020 Tuomo Kuure
// JSON logger middleware for gin-gonic

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

type LogItems struct {
	ISOTime time.Time
	UnixTime int64

	IP string

	Method string
	Path string
	Query string
	Protocol string

	Host string

	Response int
	ResponseSize int

	RequestProcessingTime int64
	LogProcessingTime int64
}

var FormatJSON = func(log LogItems) string {
	logline, _ := json.Marshal(log)
	return fmt.Sprintf("%s\n" , logline)
}

func Logger_JSON(filename string, w_stdout bool) gin.HandlerFunc {

	// TODO: When is the right time to close the log file?
	logfile, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {	fmt.Println(err) }

	// Write both to the file and stdout if desired
	// TODO: Pretty-print to stdout option?
	var out = io.Writer(logfile)
	if w_stdout { out = io.MultiWriter(logfile, os.Stdout) }

	return func(c *gin.Context) {

		start := time.Now()
		c.Next() // Request is processed here
		stop := time.Now().UnixNano()

		log := LogItems{
			ISOTime: start,
			UnixTime: start.UnixNano(),

			// Context methods
			IP: c.ClientIP(),

			// Context struct
			Method: c.Request.Method,
			Path: c.Request.URL.Path,
			Query: c.Request.URL.RawQuery,
			Protocol: c.Request.Proto,

			Host: c.Request.Host,

			Response: c.Writer.Status(),
			ResponseSize: c.Writer.Size(),
		}

		log.RequestProcessingTime = stop - log.UnixTime

		// Measure the time it took to process the log
		// TODO: this should be the very last operation
		log.LogProcessingTime = time.Now().UnixNano() - stop

		fmt.Fprint(out, FormatJSON(log))
	}
}

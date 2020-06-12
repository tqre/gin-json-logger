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
	ISOTime        time.Time
	UnixTime       int64
	IP             string
	Method         string
	Path           string
	Query          string
	User           string
	Protocol       string
	ContentType    string
	ContentLength  int64
	Host           string
	ResponseStatus int
	ResponseSize   int
	TLSData        TLSData

	RequestProcessingTime int64
	LogProcessingTime     int64
}

type TLSData struct {
	TLSVersion    uint16
	TLSCipherUsed uint16
}

var FormatJSON = func(log LogItems) string {
	logline, _ := json.Marshal(log)
	return fmt.Sprintf("%s\n", logline)
}

func Logger_JSON(filename string, w_stdout bool) gin.HandlerFunc {

	// TODO: When is the right time to close the log file?
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}

	// Write both to the file and stdout if desired
	// TODO: Pretty-print to stdout option?
	var out = io.Writer(logfile)
	if w_stdout {
		out = io.MultiWriter(logfile, os.Stdout)
	}

	return func(c *gin.Context) {

		// All time values are in nanoseconds
		start := time.Now()
		c.Next() // Request is processed here
		stop := time.Now().UnixNano()

		log := LogItems{
			ISOTime:        start,
			UnixTime:       start.UnixNano(),
			IP:             c.ClientIP(),
			Method:         c.Request.Method,
			User:           c.Request.URL.User.Username(),
			Path:           c.Request.URL.EscapedPath(),
			Query:          c.Request.URL.RawQuery,
			Protocol:       c.Request.Proto,
			ContentType:    c.ContentType(),
			ContentLength:  c.Request.ContentLength,
			Host:           c.Request.Host,
			ResponseStatus: c.Writer.Status(),
			ResponseSize:   c.Writer.Size(),
		}

		if c.Request.TLS != nil {
			// https://golang.org/pkg/crypto/tls/#pkg-constants
			log.TLSData = TLSData{
				TLSVersion:    c.Request.TLS.Version,
				TLSCipherUsed: c.Request.TLS.CipherSuite,
			}
		}

		log.RequestProcessingTime = stop - log.UnixTime

		// Measure the time it took to process the log
		// TODO: this should be the very last operation
		log.LogProcessingTime = time.Now().UnixNano() - stop

		fmt.Fprint(out, FormatJSON(log))
	}
}

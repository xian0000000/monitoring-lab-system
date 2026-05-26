package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"lab-monitor-backend/scheduler"

	"github.com/gin-gonic/gin"
)

func ScanStream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ch := scheduler.Hub.Subscribe()
	defer scheduler.Hub.Unsubscribe(ch)

	notify := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-notify:
			return false
		case result, ok := <-ch:
			if !ok {
				return false
			}
			data, _ := json.Marshal(result)
			fmt.Fprintf(w, "data: %s\n\n", data)
			return true
		}
	})
}

func TriggerScan(c *gin.Context) {
	go scheduler.RunScanOnce()
	c.JSON(200, gin.H{"message": "Scan dimulai, pantau via GET /api/scan/stream"})
}

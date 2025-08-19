package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func StartProxyServer(wg *sync.WaitGroup) {
	r := gin.Default()
	r.Any("/proxy/*path", func(c *gin.Context) {
		targetURL := c.Param("path")
		if targetURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No target URL"})
			return
		}

		fullURL := "https://" + strings.TrimLeft(targetURL, "/")
		req, err := http.NewRequest(c.Request.Method, fullURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("This here requeset:  %v\n", req)

		for k, v := range c.Request.Header {
			req.Header[k] = v
		}
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //ignore tls errors
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			c.Writer.Header()[k] = v
		}

		c.Status(resp.StatusCode)
		htmlBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error reading body": err.Error()})
		}
		io.Copy(c.Writer, bytes.NewReader(htmlBody))
	})
	wg.Add(1)

	go func() {
		err := r.Run(":8080")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}

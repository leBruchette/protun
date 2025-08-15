package proxy_server

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
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
		fullURL := "https://" + targetURL
		req, err := http.NewRequest(c.Request.Method, fullURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("This here requeset:  %v", req)

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
		io.Copy(c.Writer, resp.Body)
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

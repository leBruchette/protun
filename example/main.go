package main

import (
	"fmt"
	"github.com/lebruchette/protun/client"
	"io"
	"time"
)

func main() {
	protunClient := client.NewProtunClient()

	for i := 0; i < 10; i++ {
		request := &client.RequestConfig{
			Method: "GET",
			URL:    "https://<some-website>",
			Headers: map[string]string{
				"content-type": "text/html",
			},
			Body: "",
		}

		resp, err := protunClient.MakeRequest(*request)
		if err != nil {
			fmt.Println(err)
		}

		htmlBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(htmlBody))

		time.Sleep(1 * time.Second)
	}

}

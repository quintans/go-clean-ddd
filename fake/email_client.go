package fake

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type EmailClient struct{}

func NewEmailClient() EmailClient {
	return EmailClient{}
}

func (e EmailClient) Send(destination string, url string) {
	fmt.Println("===> faking send email to", destination, "\nurl: ", url)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("===> faking user confirmation")
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("ERROR while calling confirmation:", err)
		}
		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("ERROR while reading body:", err)
			}
			fmt.Println("ERROR: response not OK\n", string(body))
		}
	}()
}

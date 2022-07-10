package fake

import (
	"fmt"
)

type EmailClient struct{}

func NewEmailClient() EmailClient {
	return EmailClient{}
}

func (e EmailClient) Send(destination string, body string) {
	fmt.Println("===> faking send email to", destination, "\nbody: ", body)
	// TODO uncomment
	// go func() {
	// 	time.Sleep(time.Second)
	// 	fmt.Println("===> faking user confirmation")
	// 	resp, err := http.Get(body)
	// 	if err != nil {
	// 		fmt.Println("ERROR while calling confirmation:", err)
	// 	}
	// 	if resp.StatusCode != http.StatusOK {
	// 		defer resp.Body.Close()
	// 		body, err := io.ReadAll(resp.Body)
	// 		if err != nil {
	// 			fmt.Println("ERROR while reading body:", err)
	// 		}
	// 		fmt.Println("ERROR: response not OK\n", string(body))
	// 	}
	// }()
}

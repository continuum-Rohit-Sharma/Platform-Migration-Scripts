package main

import "github.com/ContinuumLLC/platform-common-lib/src/messaging"
import "fmt"
import "time"

func main() {
	conf := messaging.Config{
		Address: []string{"localhost:9092", "localhost:9093"},
		Topics:  []string{"Jain"},
		GroupID: "Jain",
	}
	service := messaging.NewService(conf)

	header := messaging.Header{}
	header.Set("Type", "Type1")
	header.Set("Type", "Type2")

	env := messaging.Envelope{
		Header:  header,
		Topic:   "Jain",
		Message: "message",
	}
	err := service.Publish(&env)
	fmt.Println(1, "  ", err)
	go listen(service)

	time.Sleep(5 * time.Second)

	err = service.Publish(&env)
	fmt.Println(4, "  ", err)
	err = service.Publish(&env)
	fmt.Println(5, "  ", err)

	time.Sleep(10 * time.Second)
}

func listen(service messaging.Service) {
	err := service.Listen(func(m *messaging.Message) {
		fmt.Println(2, "  ", m)
	})
	fmt.Println(3, "  ", err)
}

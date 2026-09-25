// App binary starts api.
package main

import (
	"fmt"
	"log"

	"server-alarm/api/config"
)

func main() {
	fmt.Println("SERVER ALARM")

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Println(cfg)
}

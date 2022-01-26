package main

import (
	"fmt"
	"time"

	"github.com/deestarks/infiniti/config"
)

func main() {
	// Load the environment variables
	config.LoadEnv("../.env")

	fmt.Println("Started at:", time.Now().Format("2022-01-02 15:04:05"))


	fmt.Println("Ended at:", time.Now().Format("2022-01-02 15:04:05"))
}

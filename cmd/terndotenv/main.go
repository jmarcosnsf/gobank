package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil{
		panic( err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations", "./internal/store/pgstore/migrations",
		"--config", "./internal/store/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("command execution failed: ", err)
		fmt.Println("output: ", string(output))
		os.Exit(1)
	}

	fmt.Println("command executed sucessfully ", string(output))
}
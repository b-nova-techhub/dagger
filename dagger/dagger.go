package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	src := client.Host().Directory(".")

	// create the container with the latest golang image
	golang := client.Container().From("golang:latest")

	// set the src dir in the container to the host path
	golang = golang.WithDirectory("/src", src).WithWorkdir("/src")

	// define the application build command
	path := "./"
	golang = golang.WithExec([]string{"go", "build", "-o", path})

	// get reference to executable file in container
	outputFileName := "dagger-techup"
	outputFile := golang.File(outputFileName)

	// write executable file from container to the host build/ directory in the current project
	outputDir := "./build/" + outputFileName
	_, err = outputFile.Export(ctx, outputDir)
	if err != nil {
		return err
	}

	return nil
}

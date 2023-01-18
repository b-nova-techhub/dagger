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

	golang := client.Container().From("golang:latest")
	golang = golang.WithEnvVariable("GOPRIVATE", "github.com/b-nova")
	golang = golang.WithDirectory("/src", src).WithWorkdir("/src")
	githubUser, _ := client.Host().EnvVariable("GOLANG_GITHUB_USERNAME").Value(ctx)
	githubToken, _ := client.Host().EnvVariable("GOLANG_GITHUB_ACCESS_TOKEN").Value(ctx)
	// define the application build command
	path := "./"
	golang = golang.WithExec([]string{"sed", "-i", "s|replace|// replace|", "go.mod"})
	golang = golang.WithExec([]string{"git", "config", "--global", "url.https://" + githubUser + ":" + githubToken + "@github.com.insteadOf", "https://github.com"})
	golang = golang.WithExec([]string{"go", "build", "-o", path})
	golang = golang.WithExec([]string{"sed", "-i", "s|// replace|replace|", "go.mod"})

	// get reference to build output directory in container
	output := golang.Directory(path)

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, path)
	if err != nil {
		return err
	}

	return nil
}

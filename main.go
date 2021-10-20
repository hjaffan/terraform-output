package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"github.com/sethvargo/go-githubactions"
	"log"
)

var TFToken string
var WorkSpace string
var Org string

func main()  {

	TFToken = githubactions.GetInput("tf-api-token")
	if TFToken == "" {
		githubactions.Fatalf("missing input 'TFToken'")
	}

	WorkSpace = githubactions.GetInput("tf-workspace")
	if WorkSpace == "" {
		githubactions.Fatalf("missing input 'WorkSpace'")
	}

	Org = githubactions.GetInput("tf-org")
	if Org == "" {
		githubactions.Fatalf("missing input 'Org'")
	}

	variableName := githubactions.GetInput("output-variable")
	if variableName == "" {
		githubactions.Fatalf("missing input 'output-variable'")
	}


	config := &tfe.Config{
		Token: TFToken,
	}

	// Create a context
	ctx := context.Background()

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Get more dynamic in what is getting called
	options := tfe.ListOptions{
		PageNumber: 1,
		PageSize:   4,
	}

	workspace, err := client.Workspaces.Read(ctx, Org, WorkSpace)
	if err != nil {
		log.Print("Got an error at workspace ID")
		log.Fatal(err)
	}

	stateVersion, err := client.StateVersions.Current(ctx, workspace.ID)
	if err != nil {
		log.Print("Got an error at state version ID")
		log.Fatal(err)
	}
	stateOutputs, err := client.StateVersions.Outputs(ctx, stateVersion.ID, tfe.StateVersionOutputsListOptions{ListOptions: options})
	if err != nil {
		log.Print("Got an error at state version outputs")
		log.Fatal(err)
	}

	for _, output := range stateOutputs {
		if output.Name == variableName {

			rawValue := output.Value
			value := fmt.Sprintf("%s", rawValue)
			fmt.Printf("::set-output name=terraform_variable::%s", value)
		}
	}

}

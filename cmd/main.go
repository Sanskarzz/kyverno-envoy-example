package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sanskarzz/kyverno-envoy-example/kyvernoauth"
	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
	"github.com/kyverno/kyverno-json/pkg/policy"
	"github.com/spf13/cobra"
)

// ServeConfig holds configuration for serving policies via a server.
// The flags will ve parsed and used to configure the server with docker image or go binary.
type ServeConfig struct {
	policies []string
	address  string
}

func main() {
	serveConfig := &ServeConfig{}
	// Cmd is the command that starts the HTTP server.
	// It uses the provided ServeConfig to determine the policies to load
	// and the address to listen on.
	var Cmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the HTTP server",
		Long:  `Starts the HTTP server with the provided Policy and Address .`,
		Run: func(cmd *cobra.Command, args []string) {
			startServer(serveConfig.address, serveConfig.policies)
		},
	}
	Cmd.Flags().StringSliceVar(&serveConfig.policies, "policy", nil, "path to policy")
	Cmd.Flags().StringVar(&serveConfig.address, "address", "localhost:9002", "address to serve on")

	var rootCmd = &cobra.Command{Use: "kyverno"}
	rootCmd.AddCommand(Cmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// startServer starts an HTTP server on the provided address
// that handles requests by calling the provided authHandler
// function on each request. It logs the address it's listening on
// and any errors starting the server.
func startServer(address string, policies []string) {
	http.HandleFunc("/", authHandler(policies))
	log.Printf("Starting server on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func authHandler(policies []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkReq kyvernoauth.CheckRequest
		var checkRes kyvernoauth.CheckResponse
		// Marshal checkReq into JSON bytes to send in the request body.
		reqBytes, err := json.Marshal(checkReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Unmarshal the JSON request bytes into the payload interface.
		// This allows the payload to be of any JSON-compatible type.
		var payload interface{}
		if err := json.Unmarshal([]byte(reqBytes), &payload); err != nil {
			panic(err)
		}
		// Load parses and initializes policy resources.
		policies, err := policy.Load(policies...)
		if err != nil {
			return
		}
		// enginerequest creates a Request struct containing the resources to check
		// and the policies to check them against. This request will be passed to
		// the policy engine for evaluation.
		enginerequest := jsonengine.Request{
			Resources: []interface{}{payload},
			Policies:  policies,
		}
		// Runs the policy engine to evaluate the request against the loaded policies.
		// Iterates through the results, setting the response status code based on the result status.
		// 200 if the policy passes, 403 if it fails.
		engine := jsonengine.New()
		results := engine.Run(context.Background(), enginerequest)
		for _, result := range results {
			if result.Result == jsonengine.StatusPass {
				checkRes.Status.Code = 200
			} else {
				checkRes.Status.Code = 403
			}
			checkRes.HTTPResponse.Status = checkRes.Status.Code
		}
	}
}

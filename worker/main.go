package main

import (
        "log"
        "os"
        "time"

        "go.temporal.io/sdk/client"
        "go.temporal.io/sdk/worker"

        "event_sender/app"
)

func main() {
        temporal := get_env("TEMPORAL_GRPC_ENDPOINT", "0.0.0.0:7233")

        var c client.Client
        var err error
        for {
                c, err = client.NewClient(client.Options{HostPort: temporal })
                if err == nil {
                        break
                }
                log.Printf("Unable to create Temporal client (%s). Sleeping", err)
                time.Sleep(10 * time.Second)
        }
        defer c.Close()

        // This worker hosts both Worker and Activity functions
        w := worker.New(c, app.EventSenderTaskQueue, worker.Options{
        })

        w.RegisterWorkflow(app.SingleSendEvent)
        w.RegisterActivity(app.SendHTTPTime)

        log.Printf("Waiting for interrupt")
        // Start listening to the Task Queue
        err = w.Run(worker.InterruptCh())
        if err != nil {
                log.Fatalln("unable to start Worker", err)
        }

        log.Printf("Closing worker")
}

func get_env(key string, def_value string) string {
        result := os.Getenv(key)
        if len(result) == 0 {
                return def_value
        }
        return result
}

package app

import (
	"time"
	"log"

	"go.temporal.io/sdk/workflow"
        "go.temporal.io/sdk/temporal"
)

func SingleSendEvent(ctx workflow.Context, events []EventDetails) error {

	// processed events count
	var processed uint64
        var errors uint64

        // RetryPolicy specifies how to automatically handle retries if an Activity fails.
        retrypolicy := &temporal.RetryPolicy{
                InitialInterval:    time.Second,
                BackoffCoefficient: 2.0,
                MaximumInterval:    time.Minute,
                MaximumAttempts:    5,
        }

        // options := workflow.ActivityOptions{
        //         // Timeout options specify when to automatically timeout Actvitivy functions.
        //         StartToCloseTimeout:    15*time.Minute,
        //         ScheduleToStartTimeout: 15*time.Minute,
        //         // Optionally provide a customized RetryPolicy.
        //         // Temporal retries failures by default, this is just an example.
        //         RetryPolicy: retrypolicy,
        // }
        // ctx = workflow.WithActivityOptions(ctx, options)

        options := workflow.LocalActivityOptions{
                // Timeout options specify when to automatically timeout Actvitivy functions.
                StartToCloseTimeout:    15*time.Minute,
                // Optionally provide a customized RetryPolicy.
                // Temporal retries failures by default, this is just an example.
                RetryPolicy: retrypolicy,
        }
        ctx = workflow.WithLocalActivityOptions(ctx, options)

	for i := 0; i < len(events); i++ {
		event := events[i]
                //err := workflow.ExecuteActivity(ctx, SendHTTPTime, event).Get(ctx, nil)
                err := workflow.ExecuteLocalActivity(ctx, SendHTTPTime, event).Get(ctx, nil)
                if err != nil {
                        log.Printf("Activity processing failed: %s", err)
                        errors++
                        continue
                }
		processed++
	}

	signalCh := workflow.GetSignalChannel(ctx, EventSenderSignalName)

        // read one event from queue, if any
        var evt EventDetails
        if ok := signalCh.ReceiveAsync(&evt); ok {
                err := workflow.ExecuteLocalActivity(ctx, SendHTTPTime, evt).Get(ctx, nil)
                if err != nil {
                        log.Printf("Activity processing failed: %s", err)
                        errors++
                }
        }

	// read signal queue and restart workflow as new if any signals found
	var unhandled_events []EventDetails
	for {
		var event EventDetails
		if ok := signalCh.ReceiveAsync(&event); !ok {
			break
		}
		unhandled_events = append(unhandled_events, event)
	}

        if len(unhandled_events) > 0 {
                return workflow.NewContinueAsNewError(ctx, SingleSendEvent, unhandled_events)
        }

        return nil
}

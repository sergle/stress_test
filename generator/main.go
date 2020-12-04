package main

import (
	"context"
	"log"
        "time"
        "os"
        "strconv"
        "math/rand"
        "fmt"
        "os/signal"
        "syscall"
        "errors"
        "sync"
        "sync/atomic"

	"go.temporal.io/sdk/client"
        "go.uber.org/ratelimit"

	"stress_test/app"
)

const WORKFLOW_TMPL = "event-sender-%06d"

type config struct {
        wf_type         string
        wf_min          int
        wf_max          int
        rps             int
        concurrency     int
}

func parse_args() (config, error) {

        var c config

        if len(os.Args) < 4 {
            return c, errors.New("Not enough arguments")
        }

        //wf_type := os.Args[1]

        min, err := strconv.Atoi(os.Args[1])
        if err != nil {
            //log.Fatalln("unable to parse min workflow id: ", err)
            return c, errors.New( fmt.Sprintf("unable to parse min_wf: %s", err) )
        }

        max, err := strconv.Atoi(os.Args[2])
        if err != nil {
            //log.Fatalln("unable to parse max workflow id: ", err)
            return c, errors.New( fmt.Sprintf("unable to parse max_wf: %s", err) )
        }

        rps, err := strconv.Atoi(os.Args[3])
        if err != nil {
            //log.Fatalln("unable to parse rps: ", err)
            return c, errors.New( fmt.Sprintf("unable to parse rps: %s", err) )
        }

        n := 1
        if len(os.Args) > 4 {
                n, err = strconv.Atoi(os.Args[4])
                if err != nil {
                        //log.Fatalln("unable to parse rps: ", err)
                        return c, errors.New( fmt.Sprintf("unable to parse n-threads: %s", err) )
                }
        }

        c = config{
                wf_type:        app.EventSenderWorkflowType,
                wf_min:         min,
                wf_max:         max,
                rps:            rps,
                concurrency:    n,
        }
        return c, nil
}

func main() {

        cfg, err := parse_args()
        if err != nil {
            log.Printf("%s", err)
            log.Printf("Use %s <WF range:Start> <WF range:Stop> <rps> <n-threads>\n", os.Args[0])
            return
        }

	// Create the client object just once per process
        temporal := os.Getenv("TEMPORAL_GRPC_ENDPOINT")
        if len(temporal) == 0 {
            temporal = "0.0.0.0:7233"
        }

	c, err := client.NewClient(client.Options{HostPort: temporal})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

        finished := false

        os_chan := make(chan os.Signal)
        signal.Notify(os_chan, os.Interrupt, syscall.SIGTERM)
        go func() {
                <-os_chan
                log.Println("\n Ctrl+C pressed in Terminal")
                finished = true
        }()


        options := client.StartWorkflowOptions{
                TaskQueue: app.EventSenderTaskQueue,
                WorkflowTaskTimeout: 10*time.Minute,
        }

        var seq uint64
        var cnt uint64

        go (func() {
            var prev_c uint64
            for {
                time.Sleep(1 * time.Minute)
                local_cnt := atomic.LoadUint64(&cnt)
                log.Printf("Send RPS: %.2f (%d -> %d)\n", float64(local_cnt - prev_c) / 60.0, prev_c, local_cnt)
                prev_c = local_cnt
            }
        })()

        rl := ratelimit.New( cfg.rps )
        log.Printf("Generating load with RPS %d\n", cfg.rps)

        partition := (cfg.wf_max - cfg.wf_min) / cfg.concurrency

        var wg sync.WaitGroup
        p_start := cfg.wf_min
        p_stop := cfg.wf_min + partition
        for p := 0; p < cfg.concurrency; p++ {
                wg.Add(1)

                go func(pp, pp_start, pp_stop int) {
                        log.Printf("Started partition %d [%d - %d]", pp, pp_start, pp_stop)

                        for {
                                if finished {
                                        break
                                }

                                now := time.Now()
                                local_seq := atomic.LoadUint64(&seq)

                                event := app.EventDetails{
                                        TypeName:   "Account/New",
                                        EventID:    int64(local_seq),
                                        UniqueID:   strconv.FormatInt( int64(local_seq + 1), 10 ),
                                        Created:    now,
                                        Effective:  now,
                                        //CreatedNano: now.UnixNano(),
                                        Variables:  map[string]interface{} {
                                                        "i_env":        rand.Intn(30),
                                                        "i_account":    rand.Intn(1000000),
                                                    },
                                }

                                for i := pp_start; i < pp_stop; i++ {
                                        rl.Take()

                                        if finished {
                                                break
                                        }

                                        // measure latency on receiver
                                        event.CreatedNano = time.Now().UnixNano()

                                        options.ID = fmt.Sprintf(WORKFLOW_TMPL, i)

                                        _, err := c.SignalWithStartWorkflow(context.Background(), options.ID, app.EventSenderSignalName, event, options, cfg.wf_type)
                                        if err != nil {
                                                log.Fatalln("error sending signal to %s workflow", options.ID, err)
                                        }

                                        atomic.AddUint64(&cnt, 1)
                                        atomic.AddUint64(&seq, 1)
                                }
                        }
                        log.Printf("Stop partition %d", pp)
                        wg.Done()

                }(p, p_start, p_stop)

                p_start = p_stop
                p_stop = p_stop + partition
        }
        wg.Wait()

        log.Printf("Stop at %d seq\n", seq)
}

package main

import (
	"log"
        "net/http"
        "os"
        "io/ioutil"
        "time"
        "strconv"
        "sync/atomic"

        "github.com/jamiealquiza/tachymeter"
)

func main() {

        samples, err := strconv.Atoi( get_env("SAMPLES", "1000") )
        if err != nil {
                log.Printf("env SAMPLES must be a number")
                return
        }

        var req_counter uint64

        t := tachymeter.New(&tachymeter.Config{Size: samples})

        go (func() {
            var prev_c uint64
            for {
                time.Sleep(1 * time.Minute)
                local_cnt := atomic.LoadUint64(&req_counter)
                log.Printf("RPS: %.2f (%d -> %d)\n", float64(local_cnt - prev_c) / 60.0, prev_c, local_cnt)

                log.Println(t.Calc())
                if local_cnt == prev_c {
                        // reset old records when no more incoming received
                        t.Reset()
                }


                prev_c = local_cnt
            }
        })()

        http.HandleFunc("/srv", func(w http.ResponseWriter, r *http.Request) {
                atomic.AddUint64(&req_counter, 1)

                if r.ContentLength > 0 {
                    _, err := ioutil.ReadAll(r.Body)
                    if err != nil {
                        log.Printf("Error reading request body: %s\n", err)
                        return
                    }
                }
                r.Body.Close()

                queue_time, err := strconv.ParseInt( r.Header.Get("X-Queue-Time"), 10, 64)
                if err != nil {
                        log.Printf("invalid X-Queue-Time value: %s", err)
                        return
                }

                delta := time.Now().UnixNano() - queue_time
                if delta >= 0 {
                        // skip negative values - timesync issue
                        t.AddTime( time.Duration( time.Now().UnixNano() - queue_time) )
                }
                //log.Printf("Queue latency: %.5f s\n",  float64(delta)/1000000000 )
        })

        http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
                log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
                for k, v := range r.Header {
                    for _, vv := range v {
                        log.Printf("  %s: %s\n", k, vv)
                    }
                }
                log.Print("\n")
                if r.ContentLength > 0 {
                    defer r.Body.Close()
                    body, err := ioutil.ReadAll(r.Body)
                    if err != nil {
                        log.Printf("Error reading request body: %s\n", err)
                        return
                    }
                    log.Print(string(body))
                    log.Print("\n")
                }

                now := time.Now().UnixNano()
                queue_time, err := strconv.ParseInt( r.Header.Get("X-Queue-Time"), 10, 64)
                if err != nil {
                        log.Printf("invalid X-Queue-Time value: %s", err)
                        return
                }

                req_time, err := strconv.ParseInt( r.Header.Get("X-Request-Time"), 10, 64)
                if err != nil {
                        log.Printf("invalid X-Request-Time value: %s", err)
                        return
                }

                log.Printf("Now:     %d\n", now)
                log.Printf("Queue:   %d\n", queue_time)
                log.Printf("Request: %d\n", req_time)

                if queue_time > now {
                        return
                }

                log.Printf("Queue latency:   %.5f s\n",  float64( now - queue_time )/1000000000 )
                log.Printf("Request latency: %.5f s\n", float64( now - req_time    )/1000000000 )

                t.AddTime( time.Duration(now - queue_time) )

                return
        })

    //listen := get_env("SERVER_PORT", "10100")
    listen := get_env("LISTEN", "localhost:10100")

    log.Printf("Listening on %s. Stats sampling size: %d\n", listen, samples)
    log.Fatal(http.ListenAndServe(listen, nil))
}

func get_env(key string, def_value string) string {
    result := os.Getenv(key)
    if len(result) == 0 {
        return def_value
    }
    return result
}

// func abs(n int64) int64 {
//         y := n >> 63
//         return (n ^ y) - y
// }


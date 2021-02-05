package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tylertreat/bench"
	"github.com/tylertreat/bench/requester"
)

func main() {
	var (
		system      = flag.String("s", "", "[kafka, nats]")
		rate        = flag.Uint64("r", 1400, "requests per second")
		size        = flag.Int("sz", 200, "message size")
		duration    = flag.Duration("d", 30*time.Second, "benchmark runtime")
		connections = flag.Uint64("c", 1, "connections")
		url         = flag.String("url", "", "broker url")
		topic       = flag.String("topic", "benchmark", "kafka topic")
	)
	flag.Parse()

	var factory bench.RequesterFactory

	switch *system {
	case "nats":
		factory = &requester.NATSStreamingRequesterFactory{
			URL:         *url,
			PayloadSize: *size,
			Subject:     "foo",
			ClientID:    "benchmark",
		}
	case "kafka":
		factory = &requester.KafkaRequesterFactory{
			URLs:        []string{*url},
			PayloadSize: *size,
			Topic:       *topic,
		}
	default:
		fmt.Printf("Unknown system '%s'\n", *system)
		os.Exit(1)
	}
	run(factory, *rate, *connections, *duration, fmt.Sprintf("%s_%d_%d.txt", *system, *rate, *size))
}

func run(factory bench.RequesterFactory, rate, conns uint64, duration time.Duration,
	output string) {

	benchmark := bench.NewBenchmark(factory, rate, conns, duration, 0)
	summary, err := benchmark.Run()
	if err != nil {
		panic(err)
	}
	if err := summary.GenerateLatencyDistribution(nil, output); err != nil {
		panic(err)
	}
	fmt.Println(summary)
}

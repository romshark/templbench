<a href="https://github.com/romshark/templbench/actions?query=workflow%3ACI">
    <img src="https://github.com/romshark/templbench/workflows/CI/badge.svg" alt="GitHub Actions: CI">
</a>
<a href="https://goreportcard.com/report/github.com/romshark/templbench">
    <img src="https://goreportcard.com/badge/github.com/romshark/templbench" alt="GoReportCard">
</a>

# github.com/a-h/templ benchmark

```
templbench % ./run-std.sh
2024/09/23 20:23:14 Building server
2024/09/23 20:23:14 Running server with env: []
2024/09/23 20:23:14 Server PID: 1719
2024/09/23 20:23:15 Server OK
2024/09/23 20:23:15 Running benchmark with args: [attack -rate=0 -duration=5s -max-workers=12]
2024/09/23 20:23:20 Generating report...
Requests      [total, rate, throughput]  286849, 57369.90, 57367.99
Duration      [total, attack, wait]      5.000157875s, 4.999991375s, 166.5µs
Latencies     [mean, 50, 95, 99, max]    171.065µs, 147.672µs, 336.016µs, 622.256µs, 4.748667ms
Bytes In      [total, mean]              116747543, 407.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:286849
Error Set:
```

This is a synthetic benchmark using
[github.com/tsenart/vegeta](https://github.com/tsenart/vegeta)
to test the performance of a server serving
[github.com/a-h/templ](https://github.com/a-h/templ) templates.

To run the benchmark, simply execute this in the repository root:

```sh
go run . -veg rate=0 -veg duration=5s -veg max-workers=12
```

For help, run:

```sh
go run . --help
```

If you want to run the server process by yourself
then you may execute the benchmark without building and running the server:

```sh
# Assuming your server already runs at http://127.0.0.1:9091
# and the endpoint you want to benchmark is GET http://127.0.0.1:9091/helloworld.
go run . -veg duration=5s -run "" -scheme http -host 127.0.0.1:9091 -method GET -path "/helloworld"
```

## Adding Server Implementations

Add your server implementation under `./cmd/<name>`, then execute:

```sh
go run . -run ./cmd/<name>
```

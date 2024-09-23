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

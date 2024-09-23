# Templ benchmark

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

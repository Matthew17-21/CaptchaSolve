# CaptchaSolve

Program to help streamline captcha solving by routing requests to third-party services like 2Captcha!

## TODO

- [x] Initialize a solver for each site in the config
- [ ] Store captcha tokens until requested
- [ ] Scheduler to delete expired tokens?
- [ ] Each time a captcha is requested, and no tokens are available, start a goroutine
- [x] Global solver vs instances
- [ ] If a API key in a given site is out of funds, do not keep requesting from the site

## Benchmarks

### Slice vs Channel

To store captcha tokens efficiently, we require a data structure with FIFO (First-In-First-Out) capabilities, especially since tokens have expirations. Two commonly used options in Go are slices and channels. Both approaches offer distinct advantages and trade-offs, as demonstrated in the implementations provided in the [slice_v_chan_test.go file](/tests/slice_v_chan_test.go).

**Why consider channels for this**?

Channels simplify concurrent programming in Go. As noted on [Stack Overflow](https://stackoverflow.com/questions/28809094/simple-concurrent-queue):
> "The concurrency of their access is handled automatically by the Go runtime. All writes into the channel are interleaved so as to be a sequential stream. All the reads are also interleaved to extract values sequentially in the same order they were enqueued."

**Benchmark Results**:

```
goos: darwin
goarch: arm64
pkg: github.com/Matthew17-21/CaptchaSolve/tests
BenchmarkSliceQueue/Unbounded-14                56   19667937 ns/op 41678204 B/op       38 allocs/op
BenchmarkSliceQueue/WithCapacity-14             68   17547552 ns/op   117704 B/op        0 allocs/op
BenchmarkChannelQueue/Buffered-14               12  100943580 ns/op  4008129 B/op   208806 allocs/op
BenchmarkFifoQueue/WithCapacity-14              24   57983495 ns/op   841964 B/op    31764 allocs/op
PASS
ok   github.com/Matthew17-21/CaptchaSolve/tests 5.282s
```

- For high-performance scenarios with controlled memory, a pre-allocated slice is the most efficient option.

- If ease of concurrent access is a priority, channels provide a cleaner and safer abstraction at the expense of performance.

# Cloudflare Systems Assessment - HTTP Tool
## Gurleen Singh <gs585@drexel.edu>

### Compile + Run

```
$ go build
```
Tested on `go version go1.15.3 darwin/amd64`

```
$ ./cf-client --help
Usage of ./cf-client:
  -print-body
    	print the response body
  -profile int
    	runs profiler n times (default -1)
  -url string
    	the url to request from (default "https://cloudflare.com/")
```

### Results
* [My General Engineering Assessment page](https://linktree.gurleen.workers.dev)
```
$ ./cf-client --url "https://linktree.gurleen.workers.dev" --profile 5
Requests: 5
Fastest time: 75 ms
Slowest time: 269 ms
Mean time: 124 ms
Median time: 99 ms
Size max: 1846 bytes
Size min: 1813 bytes
```
These things are _insanely_ fast. I can't wait to learn more about Workers.

* [Cloudflare](https://cloudflare.com/)
```
$ ./cf-client --url "https://www.cloudflare.com" --profile 5
Requests: 5
Fastest time: 116 ms
Slowest time: 320 ms
Mean time: 162 ms
Median time: 126 ms
Size max: 99512 bytes
Size min: 99491 bytes
```

* [YouTube](https://youtube.com/)
```
$ ./cf-client --url "https://www.youtube.com" --profile 5
Requests: 5
Fastest time: 104 ms
Slowest time: 309 ms
Mean time: 160 ms
Median time: 131 ms
Size max: 436394 bytes
Size min: 428140 bytes
```
package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strings"
	"strconv"
	"regexp"
	"time"
	"sort"
	"flag"
	"crypto/tls"
)

type URL struct {
	protocol string
	hostname string
	port string
	path string
}

type Response struct {
	status int
	headers map[string]string
	body string
	time int64
}

func main() {

	var url string
	var printBody bool
	flag.StringVar(&url, "url", "https://cloudflare.com/", "the url to request from")
	flag.BoolVar(&printBody, "print-body", false, "print the response body")
	profile := flag.Int("profile", -1, "runs profiler n times")

	flag.Parse()

	if *profile != -1 {
		Profile(url, *profile)
	} else {
		response := GetRequest(url)
		fmt.Println("Status:", response.status)
		if printBody {
			fmt.Println("Body: ", response.body)
		}
		fmt.Println("Took", response.time, "ms")
	}

}

func ParseUrl(url string) URL {
	r := regexp.MustCompile(`(?P<protocol>https?:\/\/)?(?P<hostname>[^:^\/]*)(?P<port>:\\d*)?(?P<path>.*)?`)
	match := r.FindStringSubmatch(url)
	parsedUrl := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			parsedUrl[name] = match[i]
		}
	}

	if parsedUrl["port"] == "" {
		parsedUrl["port"] = "443"
	}
	if parsedUrl["path"] == "" {
		parsedUrl["path"] = "/"
	}

	urlStruct := URL{
		protocol: parsedUrl["protocol"],
		hostname: parsedUrl["hostname"],
		port: parsedUrl["port"],
		path: parsedUrl["path"],
	}
	return urlStruct
}

func GetRequest(urlString string) Response {
	url := ParseUrl(urlString)

	start := time.Now()
	conn, err := tls.Dial("tcp", url.hostname + ":" + url.port, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	request := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", url.path, url.hostname)
	fmt.Fprintf(conn, request)

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	responseMap := make(map[string]string)

	scanner.Scan()
	status := strings.Split(scanner.Text(), " ")
	code, _ := strconv.Atoi(status[1])
	if code >= 400 {
		return Response{
			status: code,
			headers: responseMap,
			body: "",
			time: time.Since(start).Milliseconds(),
		}
	}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		split := strings.Split(line, ": ")
		responseMap[split[0]] = split[1]
	}

	duration := time.Since(start)

	length, _ := strconv.Atoi(responseMap["Content-Length"])
	body := ""
	lines := 0
	scanner.Scan()
	if length > 0 && length > len(scanner.Text()) {
		for scanner.Scan() {
			body += scanner.Text()
			if len(body) >= length-lines-16 {
				break
			}
			lines++
		}
	} else if responseMap["Transfer-Encoding"] == "chunked" {
		for scanner.Scan() {
			if scanner.Text() == "\r\n" {
				break
			}
			body += scanner.Text()
		}
	} else {
		body = scanner.Text()[0:length]
	}

	response := Response{
		status: code,
		headers: responseMap,
		body: body,
		time: duration.Milliseconds(),
	}
	return response
}

func MaxMin(list []int) (max int, min int) {
	max, min = list[0], list[0]
	for i := 0; i < len(list); i++ {
		if list[i] < min {
			min = list[i]
		}
		if list[i] > max {
			max = list[i]
		}
	}
	return max, min
}


func Profile(urlString string, n int) {
	times := make([]int, n)
	codes := make([]int, n)
	sizes := make([]int, n)
	errors := make([]int, 0)
	var timeTotal int = 0

	for i := 0; i < n; i++ {
		response := GetRequest(urlString)
		
		times[i] = int(response.time)
		timeTotal += int(response.time)

		codes[i] = response.status
		sizes[i] = len(response.body)
	}
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })

	max, min := MaxMin(times)
	smax, smin := MaxMin(sizes)

	fmt.Println("Requests:", n)
	fmt.Println("Fastest time:", min, "ms")
	fmt.Println("Slowest time:", max, "ms")
	fmt.Println("Mean time:", timeTotal / n, "ms")
	fmt.Println("Median time:", times[n/2], "ms")

	for j := 0; j < n; j++ {
		if codes[j] != 200 {
			errors = append(errors, codes[j])
		}
	}
	if len(errors) > 0 {
		fmt.Println("Got the following error codes:", errors)
	}

	fmt.Println("Size max:", smax, "bytes")
	fmt.Println("Size min:", smin, "bytes")
}




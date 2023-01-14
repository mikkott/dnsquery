package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const threads int = 4
const server string = "192.168.37.128:53"
const domain string = "testi.hosti"
const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type dnsRunner struct {
	client *dns.Client
	server string
	start  int
	end    int
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func durationLogger(start time.Time, durationLog *[]time.Duration) {
	elapsed := time.Since(start)
	*durationLog = append(*durationLog, elapsed)
}

func threadRunner(dr dnsRunner, wordList []string, domain string, threadId int, durationLog *[]time.Duration) {
	m := new(dns.Msg)
	for i := dr.start; i < dr.end; i++ {
		defer durationLogger(time.Now(), durationLog)
		m.SetQuestion(wordList[i]+"."+domain+".", dns.TypeA)
		in, rtt, err := dr.client.Exchange(m, dr.server)
		fmt.Println("=================== THREAD", threadId, "======================")
		fmt.Println(in.Answer, rtt, err)
	}
}

func main() {
	wordList, _ := readLines("words_alpha.txt")
	var d [threads]dnsRunner
	var wg sync.WaitGroup
	var durationLog *[]time.Duration

	for i := 0; i < threads; i++ {
		d[i] = dnsRunner{
			client: new(dns.Client),
			server: server,
			start:  len(wordList) / threads * i,
			end:    len(wordList) / threads * (i + 1),
		}
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			threadRunner(d[i], wordList, domain, i, durationLog)
		}()
	}

	wg.Wait()
}

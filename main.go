package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/miekg/dns"
)

const threads int = 4
const server string = "192.168.37.128:53"
const domain string = "testi.hosti"

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

func threadRunner(s dnsRunner, wordList []string, domain string, threadId int) {
	m := new(dns.Msg)
	for i := s.start; i < s.end; i++ {
		m.SetQuestion(wordList[i]+"."+domain+".", dns.TypeA)
		in, rtt, err := s.client.Exchange(m, s.server)
		fmt.Println("=================== THREAD", threadId, "======================")
		fmt.Println(in.Answer, rtt, err)
	}
}

func main() {
	wordList, _ := readLines("words_alpha.txt")

	var s [threads]dnsRunner

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		s[i] = dnsRunner{
			client: new(dns.Client),
			server: server,
			start:  len(wordList) / 4 * i,
			end:    len(wordList) / 4 * (i + 1),
		}
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			threadRunner(s[i], wordList, domain, i)
		}()
	}
	wg.Wait()
}

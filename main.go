package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	kvStore = os.Getenv("KVSTORE")
	kvCred  = os.Getenv("KVCRED")
)

func main() {
	var (
		interval = flag.Duration("interval", time.Minute*5, "interval to check details and send to kvstore")
		myname   = flag.String("myname", "", "the name or hostname of machine, used as key")
	)
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	myName := *myname
	if myName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("missing -myname and unable to generate os hostname")
			os.Exit(1)
		}
		myName = hostname
	}
	if kvStore == "" || kvCred == "" {
		myIPs, err := myIPs()
		if err != nil {
			fmt.Printf("unable to get IP Address: %v\n", err)
			os.Exit(1)
		}
		for _, myIP := range myIPs {
			fmt.Printf("IP: %s\n", myIP)
		}
		return
	}

	fmt.Printf("launching mimi agent. Will check and send IP every %s\n", *interval)
	senderDaemon(myName, *interval)
}
func senderDaemon(myName string, interval time.Duration) {
	lastSent := ""
	for {
		func() {
			ips, err := myIPs()
			if err != nil {
				log.Printf("unable to get my ip: %v", err)
				return
			}
			sending := strings.Join(ips, ",")
			if sending == lastSent {
				return
			}
			if err = sendIPs(myName, ips); err != nil {
				log.Printf("unable to send my ip: %v", err)
				return
			}
			lastSent = sending
			log.Printf("my IP updated onto kvstore: %s: %s", myName, sending)
		}()
		time.Sleep(interval)
	}
}
func myIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.New("unable to find network interface addresses")
	}
	var ips []string
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	if len(ips) == 0 {
		return nil, errors.New("no IPv4 addresses found")
	}
	return ips, nil
}

func sendIPs(myName string, myIPs []string) error {
	u, err := url.Parse(kvStore)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("cred", kvCred)
	q.Add("k", myName)
	q.Add("v", strings.Join(myIPs, ","))
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("PUT", u.String(), nil)
	if err != nil {
		return err
	}
	cl := &http.Client{Timeout: time.Second * 10}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s: %s", resp.StatusCode, http.StatusText(resp.StatusCode), dat)
	}
	return nil
}

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
		interval = flag.Duration("interval", time.Minute*5, "interval to check whoami and send ip")
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
		myIP, err := myIPString()
		if err != nil {
			fmt.Printf("unable to get IP Address: %v", err)
			os.Exit(1)
		}
		defer fmt.Println("whoami: set KVSTORE and KVCRED env values to send to a kvstore")
		if strings.Contains(myIP, ",") {
			fmt.Printf("IPs: %s\n", myIP)
			return
		}
		fmt.Printf("IP: %s\n", myIP)
		return
	}

	fmt.Printf("launching whoami sender. Will check and send IP every %s\n", *interval)
	senderDaemon(myName, *interval)
}
func senderDaemon(myName string, interval time.Duration) {
	lastSent := ""
	for {
		func() {
			myIP, err := myIPString()
			if err != nil {
				log.Printf("unable to get my ip: %v", err)
				return
			}
			if myIP == lastSent {
				return
			}
			if err = sendIP(myName, myIP); err != nil {
				log.Printf("unable to send my ip: %v", err)
				return
			}
			lastSent = myIP
			log.Printf("my IP updated onto kvstore: %s: %s", myName, myIP)
		}()
		time.Sleep(interval)
	}
}
func myIPString() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New("unable to find network interface addresses")
	}
	var theIP string
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() != nil {
			ipString := ipnet.IP.String()
			if theIP != "" {
				theIP = fmt.Sprintf("%s,%s", theIP, ipString)
			}
			theIP = ipString
			break
		}
	}
	if theIP == "" {
		return "", errors.New("unable to find ip address")
	}
	return theIP, nil

}

func sendIP(myName, myIP string) error {
	u, err := url.Parse(kvStore)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("cred", kvCred)
	q.Add("k", myName)
	q.Add("v", myIP)
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

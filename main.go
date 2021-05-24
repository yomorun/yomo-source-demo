package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/client"
)

type noiseData struct {
	Noise float32 `y3:"0x11" json:"noise"` // Noise value
	Time  int64   `y3:"0x12" json:"time"`  // Timestamp (ms)
	From  string  `y3:"0x13" json:"from"`  // Source IP
}

// the address of yomo-zipper.
var zipperAddr = os.Getenv("YOMO_ZIPPER_ENDPOINT")

func main() {
	if zipperAddr == "" {
		zipperAddr = "localhost:9000"
	}
	err := emit(zipperAddr)
	if err != nil {
		log.Printf("❌ Emit the data to yomo-zipper %s failure with err: %v", zipperAddr, err)
	}
}

// emit data to yomo-zipper.
// yomo-source (your data) ---> yomo-zipper ---> yomo-flow (stream processing) ---> yomo-sink (to db or web page)
func emit(addr string) error {
	// connect to yomo-zipper.
	urls := strings.Split(addr, ":")
	if len(urls) != 2 {
		return fmt.Errorf(`❌ The format of url "%s" is incorrect, it should be "host:port", f.e. localhost:9000`, addr)
	}
	host := urls[0]
	port, _ := strconv.Atoi(urls[1])
	cli, err := client.NewSource("yomo-source").Connect(host, port)
	if err != nil {
		return err
	}
	log.Printf("✅ Connected to yomo-zipper %s", addr)
	defer cli.Close()

	// generate mock data and send it to yomo-zipper in every 100 ms.
	// you can change the following codes to fit your business.
	generateAndSendData(cli)

	return nil
}

var codec = y3.NewCodec(0x10)

func generateAndSendData(stream io.Writer) {
	ip, _ := getIP()

	for {
		// generate random data.
		data := noiseData{
			Noise: rand.New(rand.NewSource(time.Now().UnixNano())).Float32() * 200,
			Time:  time.Now().UnixNano() / int64(time.Millisecond),
			From:  ip,
		}

		// Encode data via the high performance yomo-codec.
		// See https://github.com/yomorun/yomo-codec-golang for more information.
		sendingBuf, _ := codec.Marshal(data)

		// send data via QUIC stream.
		_, err := stream.Write(sendingBuf)
		if err != nil {
			log.Printf("❌ Emit %v to yomo-zipper failure with err: %v", data, err)
		} else {
			log.Printf("✅ Emit %v to yomo-zipper", data)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// getIP returns the public IP address.
// https://gist.github.com/ankanch/8c8ec5aaf374039504946e7e2b2cdf7f
func getIP() (string, error) {
	// we are using a pulib IP API, we're using ipify here, below are some others https://www.ipify.org, http://myexternalip.com, http://api.ident.me, http://whatismyipaddress.com/api
	url := "https://api.ipify.org?format=text"
	log.Print("Getting IP address from ipify ...")
	resp, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(ip), nil
}

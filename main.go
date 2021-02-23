package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/quic"
)

type noiseData struct {
	Noise float32 `y3:"0x11"` // Noise value
	Time  int64   `y3:"0x12"` // Timestamp (ms)
	From  string  `y3:"0x13"` // Source IP
}

// the address of yomo-zipper.
var zipperAddr = os.Getenv("YOMO_ZIPPER_ENDPOINT")

func main() {
	if zipperAddr == "" {
		zipperAddr = "localhost:9999"
	}
	err := emit(zipperAddr)
	if err != nil {
		log.Printf("❌ Emit the data to yomo-zipper %s failure with err: %v", zipperAddr, err)
	}
}

// emit data to yomo-zipper.
// yomo-source (your data) ---> yomo-zipper [yomo-flow (stream processing) ---> yomo-sink (to db or web page)]
func emit(addr string) error {
	// connect to yomo-zipper via QUIC.
	client, err := quic.NewClient(addr)
	if err != nil {
		return err
	}
	log.Printf("✅ Connected to yomo-zipper %s", addr)

	// create a stream
	stream, err := client.CreateStream(context.Background())
	if err != nil {
		return err
	}

	// generate mock data and send it to yomo-zipper in every 100 ms.
	// you can change the following codes to fit your business.
	generateAndSendData(stream)

	return nil
}

var codec = y3.NewCodec(0x10)

func generateAndSendData(stream quic.Stream) {
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

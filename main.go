package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/yomorun/yomo-codec-golang/pkg/codes"
	"github.com/yomorun/yomo/pkg/quic"
)

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
	} else {
		log.Printf("✅ Connected to yomo-zipper %s", addr)
	}

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

var protoCodec = codes.NewProtoCodec(0x10)

func generateAndSendData(stream quic.Stream) {
	for {
		// generate random data.
		randData := rand.New(rand.NewSource(time.Now().UnixNano())).Float32() * 200

		// Encode data via the high performance yomo-codec.
		// See https://github.com/yomorun/yomo-codec-golang for more information.
		sendingBuf, _ := protoCodec.Marshal(randData)

		// send data via QUIC stream.
		_, err := stream.Write(sendingBuf)
		if err != nil {
			log.Printf("❌ Emit %f to yomo-zipper failure with err: %v", randData, err)
		} else {
			log.Printf("✅ Emit %f to yomo-zipper", randData)
		}

		time.Sleep(100)
	}
}

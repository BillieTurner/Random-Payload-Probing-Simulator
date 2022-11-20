package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
)

func main() {
	addr := flag.String("addr", "", "server address")
	alpn := flag.String("alpn", "", "application-layer protocol negotiation")
	logfile := flag.String("o", "", "log file")
	maxlen := flag.Int("maxlen", 1, "max payload size")
	quicMode := flag.Bool("quic", false, "use quic protocol")
	flag.Parse()

	if len(*addr) == 0 {
		log.Fatalln("missing valid server address")
	}
	// if len(*alpn) == 0 {
	// 	log.Fatalln("missing valid alpn")
	// }

	if len(*logfile) > 0 {
		f, err := os.OpenFile(
			*logfile,
			os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			0666,
		)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.Println("size, duration(ms), error msg, msg")

	tlsConf := tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{*alpn},
	}

	size := 100
	repeat := 10
	fmt.Println(currTime(), "start")
	for i := 1068; i <= *maxlen; i += size {
		var wg sync.WaitGroup
		max := i + size
		if max > *maxlen {
			max = *maxlen
		}
		wg.Add((max - i) * repeat)
		for j := i; j < max; j++ {
			for k := 0; k < repeat; k++ {
				msg := make([]byte, j)
				rand.Read(msg)
				if *quicMode {
					go proberQUIC(*addr, msg, &tlsConf, &wg)
				} else {
					go proberTCP(*addr, msg, &wg)
				}
				time.Sleep(time.Duration(10) * time.Second)
			}
		}
		time.Sleep(time.Duration(60) * time.Second)
		wg.Wait()
		fmt.Printf("%s %d/%d\n", currTime(), max, *maxlen)
	}
	fmt.Println(currTime(), "done")
}

func currTime() string {
	return time.Now().Format("15:04:05.000000")
}

func proberQUIC(
	addr string,
	msg []byte,
	tlsConf *tls.Config,
	wg *sync.WaitGroup,
) {
	var err error
	start := time.Now().UnixMilli()
	var retMsg []byte
	defer func() {
		dur := time.Now().UnixMilli() - start
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		log.Printf(`%d, %d, %s, %s`, len(msg), dur, errMsg, string(retMsg))
		wg.Done()
	}()

	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return
	}
	_, err = stream.Write(msg)
	if err != nil {
		fmt.Println(1)
		return
	}

	buf := make([]byte, 1024)
	size, err := stream.Read(buf)
	if err != nil {
		return
	}
	retMsg = buf[:size]
}

func proberTCP(addr string, msg []byte, wg *sync.WaitGroup) {
	var err error
	var retMsg []byte
	start := time.Now().UnixMilli()
	defer func() {
		dur := time.Now().UnixMilli() - start
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		log.Printf(`%d, %d, %s, %s`, len(msg), dur, errMsg, string(retMsg))
		wg.Done()
	}()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println(1)
		return
	}

	buf := make([]byte, 1024)
	size, err := conn.Read(buf)
	if err != nil {
		return
	}
	retMsg = buf[:size]
}

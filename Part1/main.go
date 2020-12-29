package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
)

type binds []string

func (b binds) String() string {
	return strings.Join(b, ",")
}

func (b *binds) Set(v string) error {
	*b = strings.Split(v, ",")
	return nil
}

// Size is needed by the /demo/upload handler to determine the size of the uploaded file
type Size interface {
	Size() int64
}

// See https://en.wikipedia.org/wiki/Lehmer_random_number_generator
func generatePRData(l int) []byte {
	res := make([]byte, l)
	seed := uint64(1)
	for i := 0; i < l; i++ {
		seed = seed * 48271 % 2147483647
		res[i] = byte(seed)
	}
	return res
}

func UploadFile(response http.ResponseWriter, request *http.Request) {

	request.ParseMultipartForm(10 << 20)
	file, handler, err := request.FormFile("myfile")

	if err != nil {
		utils.DefaultLogger.Infof("Error receiving upload: %#v", err)
		return
	}
	defer file.Close()
	fmt.Println("FILE INFO")
	fmt.Println("File Name: ", handler.Filename)
	fmt.Println("File Size: ", handler.Size)
	fmt.Println("File Type: ", handler.Header.Get("Content-Type"))

	//Uploads File Folder
	tempFile, err := ioutil.TempFile("uploads", "upload-*.pdf")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Println("File Transfer Complete!")

}

func main() {
	// defer profile.Start().Stop()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// runtime.SetBlockProfileRate(1)

	verbose := flag.Bool("v", false, "verbose")
	bs := binds{}
	flag.Var(&bs, "bind", "bind to")
	//www := flag.String("www", "", "www data")
	tcp := flag.Bool("tcp", false, "also listen on TCP")
	enableQlog := flag.Bool("qlog", false, "output a qlog (in the same directory)")
	flag.Parse()

	logger := utils.DefaultLogger

	if *verbose {
		logger.SetLogLevel(utils.LogLevelDebug)
	} else {
		logger.SetLogLevel(utils.LogLevelInfo)
	}
	logger.SetLogTimeFormat("")

	if len(bs) == 0 {
		bs = binds{"localhost:6121"}
	}

	//handler := setupHandler(*www)
	handler := mux.NewRouter()

	handler.HandleFunc("/api/upload", UploadFile).Methods("POST")

	err1 := http.ListenAndServe(":5000", handler)
	if err1 != nil {
		fmt.Println(err1)
	}

	quicConf := &quic.Config{}
	if *enableQlog {
		quicConf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
			filename := fmt.Sprintf("server_%x.qlog", connID)
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Creating qlog file %s.\n", filename)
			return utils.NewBufferedWriteCloser(bufio.NewWriter(f), f)
		})
	}

	var wg sync.WaitGroup
	wg.Add(len(bs))
	for _, b := range bs {
		bCap := b
		go func() {
			var err error
			if *tcp {
				certFile, keyFile := testdata.GetCertificatePaths()
				err = http3.ListenAndServe(bCap, certFile, keyFile, handler)
			} else {
				server := http3.Server{
					Server:     &http.Server{Handler: handler, Addr: bCap},
					QuicConfig: quicConf,
				}
				err = server.ListenAndServeTLS(testdata.GetCertificatePaths())
			}
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

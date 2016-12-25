package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"strings"
	"io"
	"github.com/dustin/go-humanize"
	"gopkg.in/cheggaaa/pb.v1"
	"sync"
)

func main() {
	//version := "3.1.0"
	//windowsVersion := strings.Replace(version, ".", "_", -1)

	files := []string{
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-community-%s-unix.tar.gz", version),
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-community-%s-windows.zip", version),
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-enterprise-%s-unix.tar.gz", version),
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-enterprise-%s-windows.zip", version),
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-community_macos_%s.dmg", windowsVersion),
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-community_windows-x64_%s.exe", windowsVersion),
		//fmt.Sprintf("http://dist.neo4j.org/neo4j-community_windows_%s.exe", windowsVersion),
		"https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/3.1.0.3/apoc-3.1.0.3-all.jar",
	}

	//files := []string{"https://github.com/neo4j-contrib/neo4j-apoc-procedures/releases/download/3.1.0.3/apoc-3.1.0.3-all.jar"}

	done := make(chan string, 0)

	pool, err := pb.StartPool()
	handleError(err)

	wg := new(sync.WaitGroup)

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			downloadFile(file, pool, wg)
			done <- file
		}(file)
	}

	wg.Wait()
	pool.Stop()

	//for i := 0; i < len(files); i++ {
	//	select {
	//	case file := <-done:
	//		fmt.Println(file, " downloaded")
	//	}
	//}
}
func downloadFile(path string, pool *pb.Pool, wg *sync.WaitGroup) {
	resp, err := http.Get(path)
	handleError(err)

	bytesToDownload := resp.ContentLength
	fmt.Println(fmt.Sprintf("%s: %s", path, humanize.IBytes(uint64(bytesToDownload))))

	bar := pb.New(int(bytesToDownload)).Prefix(path)
	pool.Add(bar)

	body := resp.Body
	defer body.Close()

	parts := strings.Split(path, "/")

	file, _ := os.Create(parts[len(parts) - 1])
	defer file.Close()

	copyBuffer(file, body, bar)
	bar.Finish()
	wg.Done()
}
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func copyBuffer(dst io.Writer, src io.Reader, bar *pb.ProgressBar) (written int64, err error) {
	buf := make([]byte, 32 * 1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				bar.Add(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

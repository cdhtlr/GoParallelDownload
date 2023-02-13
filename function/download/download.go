package download

import (
	"errors"
	. "fmt"
	. "io"
	"net/http"
	"os"
	"runtime/debug"
	. "strconv"
	"strings"
	. "time"
)

var (
	downloadCompleted = 0
	url               = "https://dl-cdn.alpinelinux.org/alpine/v3.17/releases/x86_64/alpine-minirootfs-3.17.1-x86_64.tar.gz"
	//url               = "https://kartolo.sby.datautama.net.id/ubuntu-cd/14.04/ubuntu-14.04.6-server-i386.template"
	//url = "https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000"
)

func isAcceptRangeSupported() (bool, int) {
	req, _ := http.NewRequest("HEAD", url, nil)
	client := &http.Client{
		Timeout: 5 * Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		Println(err)
		return false, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		Println("Cannot continue download, status code " + Itoa(resp.StatusCode))
		return false, 0
	}

	acceptRanges := strings.ToLower(resp.Header.Get("Accept-Ranges"))
	if acceptRanges == "" || acceptRanges == "none" {
		return false, int(resp.ContentLength)
	}

	return true, int(resp.ContentLength)
}

func downloadPart(start int, end int, done chan bool) {
	download(Itoa(int(start)) + "-" + Itoa(int(end)))
	done <- true
}

func download(opts ...string) {
	req, _ := http.NewRequest("GET", url, nil)
	fileName := "Downloaded_Single_File"
	if len(opts) > 0 {
		req.Header.Add("Range", "bytes="+opts[0])
		fileName = opts[0] + ".part"
	}
	client := &http.Client{
		Timeout: 30 * Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		Println(err)
		downloadCompleted -= 1
	} else {
		downloadCompleted += 1
	}

	defer resp.Body.Close()

	if resp.Header.Get("Content-Range") != "" {
		Println(resp.Header.Get("Content-Range"))
	}

	if resp.ContentLength > 0 {
		Println("Content Length: " + resp.Header.Get("Content-Length"))
	}

	file, err := os.Create(fileName)
	if err != nil {
		Println(err)
	}
	defer file.Close()

	Copy(file, resp.Body)
}

func Run() (bool, error) {
	defer debug.FreeOSMemory()
	downloadCompleted = 0

	concurentConn := 4

	acceptRangeSupported, fileSize := isAcceptRangeSupported()

	begin := Now()

	if fileSize > 0 {
		if acceptRangeSupported {
			Println("Accept-Ranges supported")

			partSize := fileSize / concurentConn
			done := make(chan bool, concurentConn)

			for i := 0; i < concurentConn; i++ {
				start := i * partSize
				end := (i+1)*partSize - 1
				if i == concurentConn-1 {
					end = fileSize - 1
				}
				go downloadPart(start, end, done)
			}
			for i := 0; i < concurentConn; i++ {
				<-done
			}
		} else {
			concurentConn = 1
			Println("Accept-Ranges not supported")
			download()
		}
	}

	if (Since(begin)) >= 30*Second {
		return false, errors.New("Download timed out!")
	}

	if downloadCompleted != concurentConn {
		return false, errors.New("Download error :(")
	}

	return true, nil
}

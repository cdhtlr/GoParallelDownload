# GoParallelDownload
Golang basic file downloader using Accept-Ranges for parallel downloads to download a file in efficient way with the help of concurrency.

If a URL supports http header - `Accept-Ranges`, it will be divided into several parts and download it concurrently. Otherwise, the files will still be downloaded but not in parallel

![](https://raw.githubusercontent.com/cdhtlr/GoParallelDownload/main/Run.png)

## QuickStart

```go
go get -u github.com/cdhtlr/GoParallelDownload
```

# how does it work?
* Gets head of http response
* Initializes go routines for partial downloads if range downloads are supported
* Run them simultaneously until all is Done

## Compile command

    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -extldflags=-static" -a -o downloader.exe .

## Disclaimer

This program comes with no warranty. You must use this program at your own risk.
This program slightly copies the source code from <a href="https://github.com/raviraa/speedtest">Raviraa Speedtest</a>

### Note

- Using a large number of connections to a single URL can lead to DOS attacks.

## Todo

* Performs memory efficiency without debug.FreeOSMemory()

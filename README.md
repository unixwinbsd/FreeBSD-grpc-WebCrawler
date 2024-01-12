# FreeBSD GRPC Googlebot Crwaler

### Setup GRPC Web Crawler

```bash
root@ns7:~ # cd /usr/local/etc
root@ns7:/usr/local/etc # git clone git@github.com:unixwinbsd/FreeBSD-grpc-WebCrawler.git

root@ns7:/usr/local/etc # cd FreeBSD-grpc-WebCrawler
root@ns7:/usr/local/etc/FreeBSD-grpc-WebCrawler # cd crawl
root@ns7:/usr/local/etc/FreeBSD-grpc-WebCrawler/crawl # go build -o bin/client crawl.go
```

### Run

Start the server:

```bash
root@ns7:/usr/local/etc/FreeBSD-grpc-WebCrawler/server # cd bin
root@ns7:/usr/local/etc/FreeBSD-grpc-WebCrawler/server/bin # ./server
2024/01/08 10:12:04 starting gRPC server...
```


# Reference, Please visit:
https://www.unixwinbsd.site

https://www.unixwinbsd.site/2024/01/setup-grpc-google-web-crawling-with.html

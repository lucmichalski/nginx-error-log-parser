package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var entries = []string{
	"2020/01/02 03:04:05 [error] 7#8: *851624 open() \"/srv/core/client/dist/client/favicon.ico\" failed (2: No such file or directory), client: 1.1.1.1, server: _, request: \"GET /favicon.ico HTTP/1.1\", host: \"example.com\", referrer: \"https://abc.example.com/\"",
	"2020/01/02 03:04:05 [error] 7#8: *851624 FastCGI sent in stderr: \"Primary script unknown\" while reading response header from upstream, client: 1.1.1.1, server: example.com, request: \"GET /login.php HTTP/1.1\", upstream: \"fastcgi://unix:/var/run/fpm.sock:\", host: \"example.com\"",
	"2020/01/02 03:04:05 [error] 7#8: *851624 FastCGI sent in stderr: \"PHP message: PHP Warning:  file_exists(): open_basedir restriction in effect. File(/home/public_html/www/wp-content/themes/dynamic.css) is not within the allowed path(s): (/home/public_html:/usr/share/pear:/usr/share/php:/tmp:/usr/local/lib/php) in /home/public_html/public_html/wp-content/themes/includes/functions.php on line 238\", client: 1.1.1.1, server: example.com, request: \"GET /login.php HTTP/1.1\", upstream: \"fastcgi://unix:/var/run/fpm.sock:\", host: \"example.com\"",
	"2020/01/02 03:04:05 [error] 7#8: *851624 access forbidden by rule, client: 1.1.1.1, server: example.tld, request: \"GET /.git/config HTTP/1.1\", host: \"example.tld\"",
	"2020/01/02 03:04:05 [info] 7#8: *851624 this is exception might be happened",
	"2020/01/02 03:04:05 [info] 7#8: *851624 SSL_do_handshake() failed (SSL: error:1417D102:SSL routines:tls_process_client_hello:unsupported protocol) while SSL handshaking, client: 127.0.0.1, server: 0.0.0.0:443",
	"2020/01/02 03:04:05 [info] 7#8: *851624 client 127.0.0.1 closed keepalive connection",
	"2020/01/02 03:04:05 [info] 7#8: *851624 client closed connection while waiting for request, client: 127.0.0.1, server: 4.2.2.4:80",
	"2020/01/02 03:04:05 [info] 7#8: *851624 NAXSI_EXLOG: ip=192.168.1.1&server=sub.example.tld&uri=%2Findex.php&id=1013&zone=ARGS&var_name=sid&content=147%27%5B0%5D, client: 127.0.0.1, server: 4.2.2.4, request: \"GET /path/to/file.ext HTTP/2.0\", upstream: \"http://192.168.1.1:80/path/to/file.ext\", host: \"example.tld\", referrer: \"https://www.example.com/page.html\"",
	"2020/01/02 03:04:05 [info] 7#8: *851624 NAXSI_FMT: ip=192.168.1.1&server=sub.example.tld&uri=/pass_req.php&learning=0&vers=0.56&total_processed=1024&total_blocked=128&block=1&cscore0=$SQL&score0=4&cscore1=$XSS&score1=8&zone0=ARGS&id0=1013&var_name0=sid, client: 127.0.0.1, server: 4.2.2.4, request: \"GET /path/to/file.ext HTTP/2.0\", upstream: \"http://192.168.1.1:80/path/to/file.ext\", host: \"example.tld\", referrer: \"https://www.example.com/page.html\"",
	"2020/01/02 03:04:05 [info] 7#8: *851624 client canceled stream 1 while sending to client, client: 127.0.0.1, server: 4.2.2.4, request: \"GET /path/to/file.ext HTTP/2.0\", upstream: \"http://192.168.1.1:80/path/to/file.ext\", host: \"example.tld\", referrer: \"https://www.example.com/page.html\"",
}

func TestParser(t *testing.T) {
	for _, entry := range entries {
		ngxParser, e := Parser(entry)
		if e != nil {
			t.Error(e)
		}
		json, _ := ParserJSON(ngxParser)
		fmt.Println(string(json))
	}
}

func TestPacket000(t *testing.T) {
	message := `2020/01/02 03:04:05 [info] 7#8: *851624 `
	message = message + `this is exception might be happened and
	could be multi line`

	_, e := Parser(message)
	if e != nil {
		t.Error(e)
	}

}

func TestPacket001(t *testing.T) {
	_, e := Parser("INVALID PACKET")
	if e == nil {
		t.Error("Invalid packet must throw an error")
	}
}

func TestPacket002(t *testing.T) {
	_, e := parserTime("INVALID DATETIME")
	if e == nil {
		t.Error("Invalid datetime format must throw an error")
	}
}

func BenchmarkParsingRandomOne(b *testing.B) {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(entries)
	for i := 0; i < b.N; i++ {
		Parser(entries[n])
	}
}

func BenchmarkParsingRandomPick(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Seed(time.Now().Unix())
		n := rand.Int() % len(entries)
		Parser(entries[n])
	}
}

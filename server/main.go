package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func statusLog(w http.ResponseWriter, r *http.Request) {
	log.Println("=======================statusLog()========================")
	log.Println("INFO *http.Request:",
		"\n===============================Request================================",
		"\nMethod                     :", r.Method,
		"\nURL                        :", r.URL,
		"\nProto                      :", r.Proto,
		"\nProtoMajor                 :", r.ProtoMajor,
		"\nProtoMinor                 :", r.ProtoMinor,
		"\nHeader[\"Accept\"]           :", r.Header.Get("Accept"),
		"\nHeader[\"Accept-Encoding\"]  :", r.Header.Get("Accept-Encoding"),
		"\nHeader[\"Accept-Language\"]  :", r.Header.Get("Accept-Language"),
		"\nHeader[\"Cache-Control\"]    :", r.Header.Get("Cache-Control"),
		"\nHeader[\"Connection\"]       :", r.Header.Get("Connection"),
		"\nHeader[\"Pragma\"]           :", r.Header.Get("Pragma"),
		"\nHeader[\"Referer\"]          :", r.Header.Get("Referer"),
		"\nHeader[\"User-Agent\"]       :", r.Header.Get("User-Agent"),
		"\nBody                       :", r.Body,
		"\nContentLength              :", r.ContentLength,
		"\nTransferEncoding           :", r.TransferEncoding,
		"\nClose                      :", r.Close,
		"\nHost                       :", r.Host,
		"\nForm                       :", r.Form,
		"\nPostForm                   :", r.PostForm,
		"\nMultipartForm              :", r.MultipartForm,
		"\nTrailer                    :", r.Trailer,
		"\nRemoteAddr                 :", r.RemoteAddr,
		"\nRequestURI                 :", r.RequestURI,
		"\nTLS                        :", r.TLS,
		"\nCancel                     :", r.Cancel,
		"\nResponse                   :", r.Response,
		"\n======================================================================",
	)
}

type configPage struct {
	source string
	Enable string
	Md5    string
	Sha1   string
	Sha256 string
}

func newConfigPage(postForm url.Values) *configPage {
	if postForm.Encode() != "" {
		_, config.AutoCheck.Enable = postForm["enable"]
		_, config.AutoCheck.Md5 = postForm["md5"]
		_, config.AutoCheck.Sha1 = postForm["sha1"]
		_, config.AutoCheck.Sha256 = postForm["sha256"]
		if _, b := postForm["save"]; b {
			if err := config.Save(); err != nil {
				log.Fatalln("ERROR config.Save():", err)
			} // REMINDER: Not the purpose of the function
		}
	}

	p := new(configPage)
	p.source = "./src/html/config.html"
	const checked = "checked"
	if config.AutoCheck.Enable {
		p.Enable = checked
	}
	if config.AutoCheck.Md5 {
		p.Md5 = checked
	}
	if config.AutoCheck.Sha1 {
		p.Sha1 = checked
	}
	if config.AutoCheck.Sha256 {
		p.Sha256 = checked
	}
	return p
}

func (p *configPage) handler(w http.ResponseWriter, r *http.Request) {
	log.Println("==================(*configPage)handler()==================")
	log.Println("INFO  p:", p)
	tpl, err := template.ParseFiles(p.source)
	if err != nil {
		log.Fatalln("ERROR template.ParseFiles():", err)
	}
	if err := tpl.Execute(w, *p); err != nil {
		log.Fatalln("ERROR tpl.Execute():", err)
	}
}

type resultPage struct {
	source string
	Path   string
	Md5    string
	Sha1   string
	Sha256 string
}

func newResultPage(postForm url.Values) *resultPage {
	p := new(resultPage)
	p.source = "./src/html/result.html"

	t := new(Targets)
	if err := t.Load(); err != nil {
		log.Fatalln("ERROR t.Load():", err)
	}
	p.Path = t.paths[0]
	if postForm.Encode() != "" {
		if _, b := postForm["md5"]; b {
			v := md5.Sum(t.data[0])
			p.Md5 = hex.EncodeToString(v[:])
		}
		if _, b := postForm["sha1"]; b {
			v := sha1.Sum(t.data[0])
			p.Sha1 = hex.EncodeToString(v[:])
		}
		if _, b := postForm["sha256"]; b {
			v := sha256.Sum256(t.data[0])
			p.Sha256 = hex.EncodeToString(v[:])
		}
	} else {
		if config.AutoCheck.Md5 {
			v := md5.Sum(t.data[0])
			p.Md5 = hex.EncodeToString(v[:])
		}
		if config.AutoCheck.Sha1 {
			v := sha1.Sum(t.data[0])
			p.Sha1 = hex.EncodeToString(v[:])
		}
		if config.AutoCheck.Sha256 {
			v := sha256.Sum256(t.data[0])
			p.Sha256 = hex.EncodeToString(v[:])
		}
	}
	return p
}

func (p *resultPage) handler(w http.ResponseWriter, r *http.Request) {
	log.Println("==================(*resultPage)handler()==================")
	log.Println("INFO  p:", p)
	tpl, err := template.ParseFiles(p.source)
	if err != nil {
		log.Fatalln("ERROR template.ParseFiles():", err)
	}
	if err := tpl.Execute(w, *p); err != nil {
		log.Fatalln("ERROR tpl.Execute():", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("========================handler()=========================")
	log.Println("INFO  r.URL.Path:", r.URL.Path)
	r.ParseForm()
	statusLog(w, r)

	switch r.URL.Path {
	case "/":
		if config.AutoCheck.Enable {
			http.Redirect(w, r, r.URL.Path+"result/", 303)
		} else {
			http.Redirect(w, r, r.URL.Path+"config/", 303)
		}
	case "/config/":
		page := newConfigPage(r.PostForm)
		page.handler(w, r)
	case "/result/":
		page := newResultPage(r.PostForm)
		page.handler(w, r)
	case "/favicon.ico":
		log.Println("INFO r.URL.Path:", r.URL.Path)
	default:
		log.Fatalln("ERROR r.URL.Path:", r.URL.Path)
	}
}

var config = new(Config)

func main() {
	isLogging := flag.Bool("log", false, "write log for bool")
	path := flag.String("path", "", "file path")
	port := flag.Uint("port", 8080, "port number")
	flag.Parse()
	if *isLogging {
		dir := getConfigDirPath()
		logFp, err := os.OpenFile(dir+"/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic("cannot open debug.log: " + err.Error())
		}
		defer logFp.Close()
		log.SetOutput(io.MultiWriter(logFp, os.Stdout))
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(log.Ltime | log.Lshortfile)
	}
	log.Println("==========================main()==========================")
	log.Println("INFO *path:", *path)
	log.Println("INFO *port:", *port)

	if err := config.Load(); err != nil {
		log.Fatalln("ERROR Config.Load():", err.Error())
	}
	log.Println("INFO config:", config)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./src/css"))))
	http.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("./src/font"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./src/js"))))
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+strconv.Itoa(int(*port)), nil); err != nil {
		log.Fatalln("ERROR http.ListenAndServe():", err.Error())
	}
}

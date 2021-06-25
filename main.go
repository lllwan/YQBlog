package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func main () {
	c, err := ReadYamlConfig("config.yaml")
	if err != nil{
		panic(err)
	}
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", c.GetRepos)
	r.GET("/DocList/:repo", c.DocList)
	r.GET("/Doc/:repo/:slug", c.Doc)
	r.GET("/yuque/*path", c.CDNProxy)
	if c.Manage.AutoSSL {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
			fmt.Fprintf(w, "Hello, HTTPS world!")
		})
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(".cache"),
			HostPolicy: autocert.HostWhitelist(c.Manage.Domain),
		}
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", c.Manage.HttpPort), m.HTTPHandler(nil)))
		}()

		server := &http.Server{
			Addr: fmt.Sprintf(":%s", c.Manage.HttpsPort),
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
				NextProtos:     []string{http2.NextProtoTLS, "http/1.1"},
				MinVersion:     tls.VersionTLS12,
			},
			Handler: r,
			MaxHeaderBytes: 32 << 20,
		}
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(r.Run(fmt.Sprintf(":%s", c.Manage.HttpPort)))
	}

}

package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func HostWhitelist() autocert.HostPolicy {
	return func(_ context.Context, host string) error {
		if host == "" {
			return nil
		}
		return errors.New("非法请求！")
	}
}

func Redirect(w http.ResponseWriter, r *http.Request){
	log.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.Scheme)
	http.Redirect(w, r, fmt.Sprintf("https://%s%s", r.Host, r.URL.Path), http.StatusFound)
}

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
		go http.ListenAndServe(fmt.Sprintf(":%s", c.Manage.HttpPort), http.HandlerFunc(Redirect))
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(".cache"),
			HostPolicy: HostWhitelist(),
		}
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

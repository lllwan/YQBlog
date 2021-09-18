package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"os"
)

func verify(c *Config) {
	if len(c.YuQue.Repos) == 0 {
		log.Fatal("请至少指定一个repo!")
	}
	if _, err := c.ListRepoDoc(fmt.Sprintf("%s/%s", c.YuQue.User, c.YuQue.Repos[0].Repo)); err != nil {
		log.Fatal("读取doc列表失败:", err)
	}
}

func main() {
	c, err := ReadYamlConfig("config.yaml")
	if err != nil {
		log.Fatal("配置文件解析失败: ", err.Error())
	}
	verify(c)
	if c.Manage.Theme == "" {
		c.Manage.Theme = "default"
	}
	if _, err := os.Stat(fmt.Sprintf("themes/%s/index.html", c.Manage.Theme)); err != nil {
		log.Fatal("未找到主题！")
	}
	r := gin.Default()
	r.LoadHTMLGlob(fmt.Sprintf("themes/%s/*", c.Manage.Theme))
	r.GET("/", c.GetRepos)
	r.GET("/DocList/:repo", c.DocList)
	r.GET("/Doc/:repo/:slug", c.Doc)
	r.GET("/yuque/*path", c.CDNProxy)
	if c.Manage.AutoSSL {
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
			Handler:        r,
			MaxHeaderBytes: 32 << 20,
		}
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(r.Run(fmt.Sprintf(":%s", c.Manage.HttpPort)))
	}
}

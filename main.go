package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/wujiyu115/yuqueg"
	"github.com/yanyiwu/gojieba"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

var (
	Store map[string]*Set
	Seg   *gojieba.Jieba
	Cache = make(map[string]*DocDesc)
	cli   *Config
)

func SetDocIndex(detail yuqueg.DocDetail, namespace string) *DocDesc {
	Cache[detail.Data.Slug] = cli.GenerateCache(detail, namespace)
	seg := Seg.Cut(TrimHtml(detail.Data.BodyHTML), true)
	for _, v := range seg {
		if utf8.RuneCountInString(v) < 2 || regexp.MustCompile(`^\s|。|，|;|；|,|-|、|:|：|\.|\?|？|\(|\)|《|》|"|'$`).MatchString(v) {
			continue
		}
		word := ToLower(v)
		_, exist := Store[word]
		if exist {
			Store[word].Set(detail.Data.Slug)
		} else {
			set := NewSet()
			set.Set(detail.Data.Slug)
			Store[word] = set
		}
	}
	if err := PersistenceCache(); err != nil {
		log.Println("缓存持久化失败:", err.Error())
	}
	return Cache[detail.Data.Slug]
}

func PersistenceCache() error {
	b, _ := json.Marshal(&Persistence{
		Cache:      Cache,
		IndexStore: Store,
	})
	return ioutil.WriteFile("data.cache", b, 0777)
}

func createIndex(repo string) error {
	namespace := fmt.Sprintf("%s/%s", cli.YuQue.User, repo)
	docs, _ := cli.ListRepoDoc(namespace)
	for _, doc := range docs.Data {
		detail, err := cli.Client().Doc.Get(namespace, doc.Slug, &yuqueg.DocGet{Raw: 1})
		if err != nil {
			return err
		}
		SetDocIndex(detail, namespace)
	}
	return PersistenceCache()
}

func init() {
	cli = client()
	Store = make(map[string]*Set)
	Seg = gojieba.NewJieba()
	cacheFile := "data.cache"
	_, err := os.Stat(cacheFile)
Index:
	if err != nil {
		for _, v := range cli.YuQue.Repos {
			err := createIndex(v.Repo)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		var p Persistence
		b, _ := ioutil.ReadFile(cacheFile)
		err := json.Unmarshal(b, &p)
		if err != nil {
			if err := os.Remove(cacheFile); err != nil {
				log.Fatal("缓存文件：", cacheFile, "已损坏且重置失败:", err.Error())
			}
			goto Index
		}
		Cache = p.Cache
		Store = p.IndexStore
	}
}

func client() *Config {
	c, err := ReadYamlConfig("config.yaml")
	if err != nil {
		log.Fatal("配置文件解析失败: ", err.Error())
	}
	if len(c.YuQue.Repos) == 0 {
		log.Fatal("请至少指定一个repo!")
	}
	if _, err := c.ListRepoDoc(fmt.Sprintf("%s/%s", c.YuQue.User, c.YuQue.Repos[0].Repo)); err != nil {
		log.Fatal("读取doc列表失败:", err)
	}
	return c
}

func main() {
	if cli.Manage.Theme == "" {
		cli.Manage.Theme = "default"
	}
	if _, err := os.Stat(fmt.Sprintf("themes/%s/templates/index.html", cli.Manage.Theme)); err != nil {
		log.Fatal("未找到主题！")
	}
	r := gin.Default()
	r.LoadHTMLGlob(fmt.Sprintf("themes/%s/templates/*", cli.Manage.Theme))
	static := fmt.Sprintf("themes/%s/static", cli.Manage.Theme)
	if _, ok := os.Stat(static); ok == nil {
		r.Static("/static", static)
	}
	r.GET("/", cli.GetRepos)
	r.GET("/DocList/:repo", cli.DocList)
	r.GET("/Doc/:repo/:slug", cli.Doc)
	r.GET("/yuque/*path", cli.CDNProxy)
	r.GET("/search", cli.SearchDoc)
	if cli.Manage.AutoSSL {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(".cache"),
			HostPolicy: autocert.HostWhitelist(cli.Manage.Domain),
		}
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cli.Manage.HttpPort), m.HTTPHandler(nil)))
		}()

		server := &http.Server{
			Addr: fmt.Sprintf(":%s", cli.Manage.HttpsPort),
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
		log.Fatal(r.Run(fmt.Sprintf(":%s", cli.Manage.HttpPort)))
	}
}

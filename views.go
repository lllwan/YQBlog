package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wujiyu115/yuqueg"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (y Config) ServeHTTP(uri string, w http.ResponseWriter, r *http.Request) {
	// 反向代理
	remote, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	d := proxy.Director
	proxy.Director = func(r *http.Request) {
		r.Header.Set("Referer", "")
		r.Host = remote.Host
		d(r)
	}
	proxy.ServeHTTP(w, r)
}

func (y Config) CDNProxy(g *gin.Context) {
	y.ServeHTTP("https://cdn.nlark.com", g.Writer, g.Request)
}

func (y Config) GetRepos(g *gin.Context) {
	g.HTML(200, "index.html", gin.H{
		"repos": y.YuQue.Repos,
		"blog":  y.Blog,
	})
}

func (y Config) DocList(g *gin.Context) {
	repo := g.Param("repo")
	detail, err := y.ListRepoDoc(fmt.Sprintf("%s/%s", y.YuQue.User, repo))
	if err != nil {
		g.JSON(403, err.Error())
		return
	}
	var docs []yuqueg.DocBookDetail
	for _, v := range detail.Data {
		docs = append(docs, v)
	}
	detail.Data = docs
	for _, v := range y.YuQue.Repos {
		if v.Repo == repo {
			g.HTML(200, "list.html", gin.H{
				"repos": y.YuQue.Repos,
				"docs":  detail,
				"repo":  repo,
				"name":  v.Name,
				"blog":  y.Blog,
			})
			return
		}
	}
	g.JSON(403, "不存在的知识库!")
}

func (y Config) Doc(g *gin.Context) {
	repo := g.Param("repo")
	slug := g.Param("slug")
	detail, err := y.GetDoc(fmt.Sprintf("%s/%s", y.YuQue.User, repo), slug)
	if err != nil {
		g.JSON(403, err.Error())
		return
	}
	html, err := y.GetDocHTMLUseProxy(detail, g.Request.Host)
	if err != nil {
		g.JSON(403, err.Error())
		return
	}
	g.HTML(200, "doc.html", gin.H{
		"repos":  y.YuQue.Repos,
		"blog":   y.Blog,
		"doc":    template.HTML(html),
		"detail": detail,
		"index":  g.Request.Host,
		"vssue":  y.Blog.Vssue,
	})
}

func (y Config) SearchDoc(g *gin.Context) {
	var (
		result []*DocDesc
	)
	key := g.Query("key")
	if key == "" {
		g.JSON(403, gin.H{
			"success": false,
			"msg":     "搜索词为空",
			"data":    result,
		})
		return
	}
	docs := Search(key)
	for _, v := range docs {
		var item = &DocDesc{
			Description: Cache[v].Description,
			Name:        Cache[v].Name,
			UpdatedAt:   Cache[v].UpdatedAt,
			CreatedAt:   Cache[v].CreatedAt,
			Url:         fmt.Sprintf("//%s/%s/%s", g.Request.Host, strings.Split(Cache[v].Namespace, "/")[1], Cache[v].Slug),
		}
		result = append(result, item)
	}
	g.JSON(200, gin.H{
		"success": true,
		"msg":     "OK",
		"data":    result,
	})
}

package main

import (
	"fmt"
	"github.com/wujiyu115/yuqueg"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		if err := yaml.NewDecoder(f).Decode(conf); err != nil {
			return conf, err
		}
	}
	return conf, nil
}

func (y Config) Client() *yuqueg.Service {
	return yuqueg.NewService(y.YuQue.Token)
}

func (y Config) ListRepo(user string, data map[string]string) (yuqueg.UserRepos, error) {
	return y.Client().Repo.List(user, "", data)
}

func (y Config) ListRepoDoc(namespace string) (yuqueg.BookDetail, error) {
	return y.Client().Doc.List(namespace)
}

func (y Config) GetDoc(namespace, slug string) (yuqueg.DocDetail, error) {
	var doc yuqueg.DocDetail
	docs, err := y.ListRepoDoc(namespace)
	if err != nil {
		return yuqueg.DocDetail{}, err
	}
	for _, v := range docs.Data {
		data, ok := Cache[slug]
		if ok && v.Slug == slug && v.UpdatedAt == data.Data.UpdatedAt {
			return data, nil
		}
	}
	doc, err = y.Client().Doc.Get(namespace, slug, &yuqueg.DocGet{Raw: 1})
	if err != nil {
		return yuqueg.DocDetail{}, err
	}
	Cache[slug] = doc
	return doc, nil
}

func (y Config) GetDocHTML(detail yuqueg.DocDetail) (string, error) {
	html := strings.Replace(detail.Data.BodyHTML, "<!doctype html>", "", -1)
	return html, nil
}

func (y Config) GetDocHTMLUseProxy(detail yuqueg.DocDetail, host string) (string, error) {
	html, err := y.GetDocHTML(detail)
	if err != nil {
		return "", err
	}
	// 通过替换html中的cdn链接进行反向代理避免跨域问题。
	result := strings.Replace(html, "https://cdn.nlark.com/", fmt.Sprintf("http://%s/", host), -1)
	return result, nil
}

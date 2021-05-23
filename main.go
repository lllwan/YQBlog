package main

import (
	"github.com/gin-gonic/gin"
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
	r.Run(":80")
}

package main

type Repo struct {
	Name string `yaml:"name"`
	Repo string `yaml:"repo"`
}

type Vssue struct {
	Owner        string `yaml:"owner"`
	Repo         string `yaml:"repo"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type Blog struct {
	Title       string `yaml:"title"`
	Avatar      string `yaml:"avatar"`
	Subtitle    string `yaml:"subtitle"`
	Keywords    string `yaml:"keywords"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Vssue       Vssue  `yaml:"vssue"`
	Link  []Links `yaml:"link"`
}

type Links struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

type YuQue struct {
	Api   string  `yaml:"api"`
	Token string  `yaml:"token"`
	User  string  `yaml:"user"`
	Repos []Repo  `yaml:"repos"`
}

type Manage struct {
	AutoSSL   bool   `yaml:"autoSSL"`
	HttpPort  string `yaml:"httpPort"`
	HttpsPort string `yaml:"httpsPort"`
	Domain    string `yaml:"domain"`
	Theme     string `yaml:"theme"`
}

type Config struct {
	YuQue  YuQue  `yaml:"yuque"`
	Blog   Blog   `yaml:"blog"`
	Manage Manage `yaml:"manage"`
}

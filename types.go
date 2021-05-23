package main

type Repo struct {
	Name string					`yaml:"name"`
	Repo string					`yaml:"repo"`
}

type Links struct {
	Name string					`yaml:"name"`
	Url string					`yaml:"url"`
}

type Blog struct {
	Title 			string		`yaml:"title"`
	Subtitle		string		`yaml:"subtitle"`
	Keywords 		string		`yaml:"keywords"`
	Description 	string		`yaml:"description"`
	Author 			string		`yaml:"author"`
	Links 			[]Links		`yaml:"links"`
}

type YuQue struct {
	Api 	string		`yaml:"api"`
	Token 	string		`yaml:"token"`
	User 	string		`yaml:"user"`
	Repos 	[]Repo		`yaml:"repos"`
}

type Config struct {
	YuQue		YuQue		`yaml:"yuque"`
	Blog		Blog		`yaml:"blog"`
}
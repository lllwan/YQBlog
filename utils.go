package main

import (
	"regexp"
	"strings"
)

type Empty struct{}

type Set struct {
	Store map[string]Empty
}

func NewSet() *Set {
	return &Set{Store: make(map[string]Empty)}
}

func (s Set) List() []string {
	var result []string
	for k := range s.Store {
		result = append(result, k)
	}
	return result
}

func (s Set) Set(key string) {
	s.Store[key] = Empty{}
}

func (s Set) Del(key string) {
	if s.Exist(key) {
		delete(s.Store, key)
	}
}

func (s Set) Exist(key string) bool {
	_, exist := s.Store[key]
	return exist
}

func ToLower(text string) string {
	match, _ := regexp.MatchString(`^[A-Za-z]+$`, text)
	if match {
		return strings.ToLower(text)
	}
	return text
}

func TrimHtml(src string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func Search(text string) []string {
	var result []string
	for _, v := range Seg.Cut(text, true) {
		var key = ToLower(v)
		doc, exist := Store[key]
		if !exist {
			continue
		}
		result = append(result, doc.List()...)
	}
	return result
}

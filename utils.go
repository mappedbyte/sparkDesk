package main

import (
	browser "github.com/EDDYCJY/fake-useragent"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func readLocalYamlConfig() (Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("无法读取文件,请检查文件格式: %v", err)
		return Config{}, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("解析 YAML 失败: %v", err)
		return Config{}, err
	}
	return config, nil
}

func createChatHeader() http.Header {
	header := make(http.Header)
	//header.Set("Accept", "text/event-stream")
	//header.Set("Connection", "keep-alive")
	//header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//header.Set("Accept-Encoding", "gzip, deflate, br")
	header.Set("Cookie", config.Cookie)
	header.Set("Origin", "https://xinghuo.xfyun.cn")
	//header.Set("Host", "xinghuo.xfyun.cn")
	header.Set("Referer", "https://xinghuo.xfyun.cn/desk")
	//header.Set("sec-ch-ua", `"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`)
	//header.Set("sec-ch-ua-mobile", "?0")
	//header.Set("sec-ch-ua-platform", "Windows")
	//header.Set("Sec-Fetch-Dest", "empty")
	//header.Set("Sec-Fetch-Mode", "cors")
	//header.Set("Sec-Fetch-Site", "same-origin")
	header.Set("User-Agent", browser.Chrome())
	return header
}

func createRequestHeader() http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json, text/plain, */*")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	header.Set("Accept-Encoding", "gzip, deflate, br")
	header.Set("X-Requested-With", "XMLHttpRequest")
	header.Set("Content-Type", "application/json")
	header.Set("Cookie", config.Cookie)
	header.Set("Origin", "https://xinghuo.xfyun.cn")
	header.Set("Referer", "https://xinghuo.xfyun.cn/desk")
	header.Set("sec-ch-ua", `"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`)
	header.Set("sec-ch-ua-mobile", "?0")
	header.Set("sec-ch-ua-platform", "Windows")
	header.Set("Sec-Fetch-Dest", "empty")
	header.Set("Sec-Fetch-Mode", "cors")
	header.Set("Sec-Fetch-Site", "same-origin")
	header.Set("User-Agent", browser.Chrome())
	return header
}

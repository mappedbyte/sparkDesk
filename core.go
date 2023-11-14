package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var config Config

const (
	NewChat = "SparkDesk AI" // 请替换为你的 Chat 名称
)

func init() {

	conf, err := readLocalYamlConfig()
	if err != nil {
		os.Exit(-1)
	}
	config = conf
}

func NewSparkWeb() *SparkWeb {
	return &SparkWeb{
		Config: config,
		ChatId: "",
		HttpClient: http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *SparkWeb) generateChatID() string {
	if s.ChatId == "" {
		//https://xinghuo.xfyun.cn/iflygpt/u/chat-list/v1/create-chat-list
		url := "https://xinghuo.xfyun.cn/iflygpt/u/chat-list/v1/create-chat-list"
		payload := []byte(`{}`)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			log.Fatalf("generateChatID  http.NewRequest error")
		}
		request.Header = createRequestHeader()
		response, err := s.HttpClient.Do(request)
		if err != nil {
			fmt.Println("Error:", err)
			return "0"
		}
		defer response.Body.Close()

		var responseData map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
			fmt.Println("Error decoding response:", err)
			return "0"
		}

		if code, ok := responseData["code"].(float64); ok && code == 0 {
			fmt.Println()
			chatListID := responseData["data"].(map[string]interface{})["id"].(float64)
			//return string(chatListID)
			return strconv.FormatFloat(chatListID, 'f', -1, 64)
		} else {
			return "0"
		}
	} else {
		return s.ChatId
	}
}

func (s *SparkWeb) setName(chatName string) {
	url := "https://xinghuo.xfyun.cn/iflygpt/u/chat-list/v1/rename-chat-list"
	chatListName := chatName[:]
	payload := map[string]string{
		"chatListId":   s.ChatId,
		"chatListName": chatListName,
	}
	payloadBytes, _ := json.Marshal(payload)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	request.Header = createRequestHeader()
	response, err := s.HttpClient.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	defer response.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println("Error decoding response:", err)
		os.Exit(-1)
	}

	if code, ok := responseData["code"].(float64); ok && code != 0 {
		fmt.Println("\nERROR: Failed to initialize session name. Please reset Cookie, fd, and GtToken")
		os.Exit(-1)
	}
}

func (s *SparkWeb) createChat(chatName string) {
	s.ChatId = s.generateChatID()

	s.setName(chatName)
}

func (s *SparkWeb) getChatSID() string {
	response := s.getChatHistory()
	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return ""
	}

	response.Body.Close()

	historyList := data["data"].([]interface{})[0].(map[string]interface{})["historyList"].([]interface{})
	if len(historyList) == 0 {
		return ""
	}
	lastHistory := historyList[len(historyList)-1].(map[string]interface{})
	return lastHistory["sid"].(string)
}

func (s *SparkWeb) getChatHistory() *http.Response {
	url := fmt.Sprintf("https://xinghuo.xfyun.cn/iflygpt/u/chat_history/all/%s", s.ChatId)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("getChatHistory  http.NewRequest  error ", err)
	}
	request.Header = createRequestHeader()
	response, err := s.HttpClient.Do(request)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return response
}

func (s *SparkWeb) getResp(question string) *http.Response {
	sid := s.getChatSID()
	requestUrl := "https://xinghuo.xfyun.cn/iflygpt-chat/u/chat_message/chat"
	// 准备 Form-Data 数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 添加字段
	_ = writer.WriteField("fd", s.Config.Fd)
	_ = writer.WriteField("chatId", s.ChatId)
	_ = writer.WriteField("text", question)
	_ = writer.WriteField("clientType", "1")
	_ = writer.WriteField("GtToken", s.Config.GtToken)
	if len(sid) != 0 {
		_ = writer.WriteField("sid", sid)
	}
	err := writer.Close()
	request, err := http.NewRequest("POST", requestUrl, body)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	request.Header = createChatHeader()
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := s.HttpClient.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	//defer response.Body.Close()

	return response
}

func (s *SparkWeb) chat(question string) string {
	s.createChat(NewChat)
	response := s.getResp(question)
	defer response.Body.Close()
	var responseText string
	dataBuf, _ := ioutil.ReadAll(response.Body)
	text := string(dataBuf)
	for _, line := range strings.Split(text, "\n") {
		if len(line) != 0 {
			encodedData := line[len("data:"):]
			missingPadding := len(encodedData) % 4
			if missingPadding != 0 {
				encodedData += strings.Repeat("=", missingPadding)
			}
			decodedData, err := base64.StdEncoding.DecodeString(encodedData)
			if err != nil {
				continue
			}
			if string(decodedData) != "zw" {
				responseText += string(decodedData)
			}
		}
	}
	return responseText
}

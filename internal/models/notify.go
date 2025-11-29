package models

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "crypto/tls"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "html"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "net/smtp"
    "regexp"
    "strings"
    "time"

    "github.com/dghubble/oauth1"
    "gorm.io/gorm"
)

type NotifyConfig struct {
	gorm.Model
	WebhookEnabled           bool   `json:"webhookEnabled"`
	WebhookURL               string `gorm:"type:varchar(255)" json:"webhookURL"`
	TelegramEnabled          bool   `json:"telegramEnabled"`
	TelegramToken            string `gorm:"type:varchar(255)" json:"telegramToken"`
	TelegramChatID           string `gorm:"type:varchar(100)" json:"telegramChatID"`
	WeworkEnabled            bool   `json:"weworkEnabled"`
	WeworkKey                string `gorm:"type:varchar(255)" json:"weworkKey"`
	FeishuEnabled            bool   `json:"feishuEnabled"`
	FeishuWebhook            string `gorm:"type:varchar(255)" json:"feishuWebhook"`
	FeishuSecret             string `gorm:"type:varchar(255)" json:"feishuSecret"`
	TwitterEnabled           bool   `json:"twitterEnabled"`
	TwitterApiKey            string `gorm:"type:varchar(255)" json:"twitterApiKey"`
	TwitterApiSecret         string `gorm:"type:varchar(255)" json:"twitterApiSecret"`
	TwitterAccessToken       string `gorm:"type:varchar(255)" json:"twitterAccessToken"`
	TwitterAccessTokenSecret string `gorm:"type:varchar(255)" json:"twitterAccessTokenSecret"`
	CustomHttpEnabled        bool   `json:"customHttpEnabled"`
	CustomHttpUrl            string `gorm:"type:varchar(255)" json:"customHttpUrl"`
	CustomHttpMethod         string `gorm:"type:varchar(20)" json:"customHttpMethod"`
	CustomHttpHeaders        string `gorm:"type:text" json:"customHttpHeaders"`
	CustomHttpBody           string `gorm:"type:text" json:"customHttpBody"`
}

func GetNotifyConfig() *NotifyConfig {
	db := GetDB()
	var config NotifyConfig
	if err := db.First(&config).Error; err != nil {
		config = NotifyConfig{}
	}
	return &config
}

// 保存配置
func SaveNotifyConfig(config NotifyConfig) error {
	db := GetDB()
	if db == nil {
		return errors.New("数据库连接未初始化")
	}
	if err := validateNotifyConfig(config); err != nil {
		return err
	}
	var existingConfig NotifyConfig
	result := db.First(&existingConfig)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新配置时，如果配置了对应渠道，则自动启用
			if config.WebhookURL != "" {
				config.WebhookEnabled = true
			}
			if config.TelegramToken != "" && config.TelegramChatID != "" {
				config.TelegramEnabled = true
			}
			if config.WeworkKey != "" {
				config.WeworkEnabled = true
			}
			if config.FeishuWebhook != "" {
				config.FeishuEnabled = true
			}
			// 新增：自动启用 Twitter
			if config.TwitterApiKey != "" && config.TwitterApiSecret != "" &&
				config.TwitterAccessToken != "" && config.TwitterAccessTokenSecret != "" {
				config.TwitterEnabled = true
			}
			// 新增：自动启用自定义HTTP
			if config.CustomHttpUrl != "" {
				config.CustomHttpEnabled = true
			}
			return db.Create(&config).Error
		}
		return result.Error
	}

	// 更新现有配置时，如果配置了对应渠道，则自动启用
	if config.WebhookURL != "" {
		config.WebhookEnabled = true
	}
	if config.TelegramToken != "" && config.TelegramChatID != "" {
		config.TelegramEnabled = true
	}
	if config.WeworkKey != "" {
		config.WeworkEnabled = true
	}
	if config.FeishuWebhook != "" {
		config.FeishuEnabled = true
	}
	// 新增：自动启用 Twitter
	// 修改自动启用逻辑
	if config.TwitterApiKey != "" && config.TwitterApiSecret != "" &&
		config.TwitterAccessToken != "" && config.TwitterAccessTokenSecret != "" {
		config.TwitterEnabled = true
	} else {
		// 明确设置为false
		config.TwitterEnabled = false
		config.TwitterApiKey = ""
		config.TwitterApiSecret = ""
		config.TwitterAccessToken = ""
		config.TwitterAccessTokenSecret = ""
	}

	if config.CustomHttpUrl != "" {
		config.CustomHttpEnabled = true
	} else {
		// 明确设置为false
		config.CustomHttpEnabled = false
		config.CustomHttpUrl = ""
		config.CustomHttpMethod = "POST"
		config.CustomHttpHeaders = ""
		config.CustomHttpBody = ""
	}

	// 更新所有字段
	existingConfig.WebhookEnabled = config.WebhookEnabled
	existingConfig.WebhookURL = config.WebhookURL
	existingConfig.TelegramEnabled = config.TelegramEnabled
	existingConfig.TelegramToken = config.TelegramToken
	existingConfig.TelegramChatID = config.TelegramChatID
	existingConfig.WeworkEnabled = config.WeworkEnabled
	existingConfig.WeworkKey = config.WeworkKey
	existingConfig.FeishuEnabled = config.FeishuEnabled
	existingConfig.FeishuWebhook = config.FeishuWebhook
	existingConfig.FeishuSecret = config.FeishuSecret
	// 新增字段同步
	existingConfig.TwitterEnabled = config.TwitterEnabled
	existingConfig.TwitterApiKey = config.TwitterApiKey
	existingConfig.TwitterApiSecret = config.TwitterApiSecret
	existingConfig.TwitterAccessToken = config.TwitterAccessToken
	existingConfig.TwitterAccessTokenSecret = config.TwitterAccessTokenSecret
	existingConfig.CustomHttpEnabled = config.CustomHttpEnabled
	existingConfig.CustomHttpUrl = config.CustomHttpUrl
	existingConfig.CustomHttpMethod = config.CustomHttpMethod
	existingConfig.CustomHttpHeaders = config.CustomHttpHeaders
	existingConfig.CustomHttpBody = config.CustomHttpBody

	return db.Save(&existingConfig).Error
}

// 验证配置
func validateNotifyConfig(config NotifyConfig) error {
	if config.WebhookEnabled {
		if !strings.HasPrefix(config.WebhookURL, "http") {
			return errors.New("Webhook URL必须以http/https开头")
		}
	}
	if config.TelegramEnabled {
		if config.TelegramToken == "" {
			return errors.New("Telegram Token不能为空")
		}
		if config.TelegramChatID == "" {
			return errors.New("Telegram Chat ID不能为空")
		}
	}
	if config.WeworkEnabled && config.WeworkKey == "" {
		return errors.New("企业微信Key不能为空")
	}
	if config.FeishuEnabled {
		if config.FeishuWebhook == "" {
			return errors.New("飞书Webhook不能为空")
		}
	}
	// 新增 Twitter 校验
	if config.TwitterEnabled {
		if config.TwitterApiKey == "" ||
			config.TwitterApiSecret == "" ||
			config.TwitterAccessToken == "" ||
			config.TwitterAccessTokenSecret == "" {
			return errors.New("Twitter配置不完整")
		}
	}
	// 新增自定义HTTP校验
	if config.CustomHttpEnabled {
		if config.CustomHttpUrl == "" {
			return errors.New("自定义HTTP URL不能为空")
		}
	}
	return nil
}

// 发送Webhook消息
func SendWebhook(message string) error {
	config := GetNotifyConfig()
	if !config.WebhookEnabled || config.WebhookURL == "" {
		log.Printf("Webhook未启用或URL为空")
		return nil
	}

	payload := map[string]interface{}{
		"text":     message,
		"markdown": true, // 启用markdown支持
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("webhook请求错误: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("webhook响应解析失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook请求失败: %d, %v", resp.StatusCode, response)
	}

	return nil
}

// 发送Telegram消息
// SendTelegram 发送 Telegram 消息
func SendTelegram(content string, images []string) error {
	config := GetNotifyConfig()
	if config.TelegramToken == "" || config.TelegramChatID == "" {
		return nil
	}

	// 处理 Markdown 内容
	messageText := content

	// 如果有图片，添加到消息末尾
	if len(images) > 0 {
		messageText += "\n\n"
		for _, img := range images {
			// 使用 HTML 格式的图片标签
			messageText += fmt.Sprintf("<a href=\"%s\">&#8205;</a>", img)
		}
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
	payload := map[string]interface{}{
		"chat_id":                  config.TelegramChatID,
		"text":                     messageText,
		"parse_mode":               "HTML",
		"disable_web_page_preview": false, // 启用网页预览以显示图片
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("telegram请求错误: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("telegram响应解析失败: %v", err)
	}

	if !response["ok"].(bool) {
		return fmt.Errorf("telegram发送失败: %v", response["description"])
	}

	return nil
}

// SendWework 发送企业微信消息
func SendWework(content string, images []string) error {
	config := GetNotifyConfig()
	if config.WeworkKey == "" {
		return nil
	}

	if !strings.HasPrefix(config.WeworkKey, "https://") {
		config.WeworkKey = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + config.WeworkKey
	}

	// 企业微信的 markdown 格式
	messageText := content
	if len(images) > 0 {
		messageText += "\n"
		for _, img := range images {
			messageText += fmt.Sprintf("![image](%s)\n", img)
		}
	}

	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": messageText,
		},
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(config.WeworkKey, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("企业微信请求错误: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("企业微信响应解析失败: %v", err)
	}

	if response["errcode"].(float64) != 0 {
		return fmt.Errorf("企业微信发送失败: %v", response["errmsg"])
	}

	return nil
}

// 发送飞书消息
func SendFeishu(message string) error {
	config := GetNotifyConfig()
	if config.FeishuWebhook == "" {
		return nil
	}

	// 飞书使用富文本格式
	payload := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"config": map[string]interface{}{
				"wide_screen_mode": true,
			},
			"elements": []map[string]interface{}{
				{
					"tag":     "markdown",
					"content": message,
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.FeishuWebhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("飞书请求创建失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if config.FeishuSecret != "" {
		timestamp := time.Now().Unix()
		sign := genFeishuSign(timestamp, config.FeishuSecret)
		req.Header.Set("X-Lark-Signature", sign)
		req.Header.Set("X-Lark-Request-Timestamp", fmt.Sprintf("%d", timestamp))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("飞书请求错误: %v", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("飞书响应解析失败: %v", err)
	}

	if code, ok := response["code"].(float64); !ok || code != 0 {
		return fmt.Errorf("飞书请求失败: %v", response["msg"])
	}

	return nil
}

// 生成飞书签名
func genFeishuSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 测试推送函数 - 保持测试推送的特殊格式
// TestNotify 函数修改
func TestNotify(notifyType string) error {
	testMessage := fmt.Sprintf("Echo Noise 推送测试\n\n"+
		"类型: %s\n"+
		"时间: %s\n"+
		"这是一条测试消息，用于验证推送配置是否正确",
		notifyType,
		time.Now().Format("2006-01-02 15:04:05"))

	// 测试消息不包含图片
	var emptyImages []string

	var err error
	switch notifyType {
	case "webhook":
		err = SendWebhook(testMessage)
	case "telegram":
		err = SendTelegram(testMessage, emptyImages)
	case "wework":
		err = SendWework(testMessage, emptyImages)
	case "feishu":
		err = SendFeishu(testMessage)
	case "twitter":
		err = SendTwitter(testMessage)
	case "customHttp":
		err = SendCustomHttp(testMessage)
	default:
		return fmt.Errorf("不支持的推送类型: %s", notifyType)
	}

	return err
}

// SendNotify 函数修改
func SendNotify(content string, images []string, config NotifyConfig) error {
	var errors []string

	// 保持原始 Markdown 格式
	if config.WebhookEnabled && config.WebhookURL != "" {
		if err := SendWebhook(content); err != nil {
			errors = append(errors, fmt.Sprintf("webhook: %v", err))
		}
	}

	if config.TelegramEnabled && config.TelegramToken != "" && config.TelegramChatID != "" {
		if err := SendTelegram(content, images); err != nil {
			errors = append(errors, fmt.Sprintf("telegram: %v", err))
		}
	}

	if config.WeworkEnabled && config.WeworkKey != "" {
		if err := SendWework(content, images); err != nil {
			errors = append(errors, fmt.Sprintf("wework: %v", err))
		}
	}

	if config.FeishuEnabled && config.FeishuWebhook != "" {
		// 飞书需要特殊处理，确保图片正确显示
		fullContent := content
		if len(images) > 0 {
			fullContent += "\n\n"
			for _, img := range images {
				fullContent += fmt.Sprintf("![image](%s)\n", img)
			}
		}
		if err := SendFeishu(fullContent); err != nil {
			errors = append(errors, fmt.Sprintf("feishu: %v", err))
		}
	}
	// Twitter 推送
	if config.TwitterEnabled && config.TwitterApiKey != "" && config.TwitterApiSecret != "" &&
		config.TwitterAccessToken != "" && config.TwitterAccessTokenSecret != "" {
		tweet := content
		if len([]rune(tweet)) > 280 {
			tweet = string([]rune(tweet)[:280]) + "...(内容截断)"
		}
		if err := SendTwitter(tweet); err != nil {
			errors = append(errors, fmt.Sprintf("twitter: %v", err))
		}
	}

	// 自定义 HTTP 推送
	if config.CustomHttpEnabled && config.CustomHttpUrl != "" {
		if err := SendCustomHttp(content); err != nil {
			errors = append(errors, fmt.Sprintf("customHttp: %v", err))
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("推送失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// SendTelegramWithFormat 发送支持 HTML 格式的 Telegram 消息
func SendTelegramWithFormat(content string, images []string, parseHTML bool) error {
	config := GetNotifyConfig() // 使用 config 替代 notifyConfig
	if !config.TelegramEnabled {
		return nil
	}

	// 处理 Markdown 链接格式 [text](url)
	content = markdownLinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := markdownLinkRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			return fmt.Sprintf("<a href=\"%s\">%s</a>", parts[2], html.EscapeString(parts[1]))
		}
		return match
	})

	// 处理 Markdown 图片格式 ![alt](url)
	content = markdownImageRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := markdownImageRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			return fmt.Sprintf("<a href=\"%s\">&#8205;</a>", parts[2])
		}
		return match
	})

	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
	payload := map[string]interface{}{
		"chat_id":    config.TelegramChatID,
		"text":       content,
		"parse_mode": "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 如果有额外的图片，单独发送
	for _, img := range images {
		if img == "" {
			continue
		}

		photoUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", config.TelegramToken)
		photoPayload := map[string]interface{}{
			"chat_id": config.TelegramChatID,
			"photo":   img,
		}

		jsonData, err = json.Marshal(photoPayload)
		if err != nil {
			return err
		}

		resp, err = http.Post(photoUrl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	return nil
}

var (
	markdownLinkRegex  = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	markdownImageRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

// SendTelegramPhotoWithCaption 发送图片和文本作为一条消息
func SendTelegramPhotoWithCaption(photoURL string, caption string) error {
	config := GetNotifyConfig()
	if config == nil || !config.TelegramEnabled {
		return fmt.Errorf("Telegram 推送未启用")
	}

	// 构建请求URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", config.TelegramToken)

	// 检查URL是否为本地URL
	isLocalURL := strings.HasPrefix(photoURL, "http://localhost") ||
		strings.HasPrefix(photoURL, "https://localhost") ||
		strings.Contains(photoURL, "127.0.0.1") ||
		(!strings.HasPrefix(photoURL, "http://") && !strings.HasPrefix(photoURL, "https://"))

	var resp *http.Response
	var err error

	if isLocalURL {
		// 对于本地URL，先下载图片然后作为文件上传
		log.Printf("检测到本地图片URL: %s，尝试下载后上传", photoURL)

		// 创建一个multipart表单
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 添加chat_id字段
		_ = writer.WriteField("chat_id", config.TelegramChatID)

		// 添加caption字段
		_ = writer.WriteField("caption", caption)

		// 添加parse_mode字段
		_ = writer.WriteField("parse_mode", "Markdown")

		// 尝试下载图片
		var imgResp *http.Response
		if strings.HasPrefix(photoURL, "http") {
			imgResp, err = http.Get(photoURL)
		} else {
			// 如果是相对路径，尝试从本地文件系统读取
			return fmt.Errorf("无法处理相对路径图片: %s", photoURL)
		}

		if err != nil {
			return fmt.Errorf("下载图片失败: %v", err)
		}
		defer imgResp.Body.Close()

		// 创建photo部分
		part, err := writer.CreateFormFile("photo", "image.jpg")
		if err != nil {
			return fmt.Errorf("创建表单文件失败: %v", err)
		}

		// 复制图片数据
		_, err = io.Copy(part, imgResp.Body)
		if err != nil {
			return fmt.Errorf("复制图片数据失败: %v", err)
		}

		// 完成multipart表单
		err = writer.Close()
		if err != nil {
			return fmt.Errorf("关闭表单失败: %v", err)
		}

		// 发送请求
		req, err := http.NewRequest("POST", apiURL, body)
		if err != nil {
			return fmt.Errorf("创建请求失败: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err = client.Do(req)
	} else {
		// 对于公网可访问的URL，直接使用URL
		// 构建请求体
		data := map[string]interface{}{
			"chat_id":    config.TelegramChatID,
			"photo":      photoURL,
			"caption":    caption,
			"parse_mode": "Markdown",
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		// 发送请求
		resp, err = http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API 错误: %s", string(body))
	}

	return nil
}

// SendTelegramVideoWithCaption 发送视频和文本作为一条消息
func SendTelegramVideoWithCaption(videoURL string, caption string) error {
	config := GetNotifyConfig()
	if config == nil || !config.TelegramEnabled {
		return fmt.Errorf("Telegram 推送未启用")
	}

	// 构建请求URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendVideo", config.TelegramToken)

	// 检查URL是否为本地URL
	isLocalURL := strings.HasPrefix(videoURL, "http://localhost") ||
		strings.HasPrefix(videoURL, "https://localhost") ||
		strings.Contains(videoURL, "127.0.0.1") ||
		(!strings.HasPrefix(videoURL, "http://") && !strings.HasPrefix(videoURL, "https://"))

	// 如果是本地URL，改为发送消息并附带链接
	if isLocalURL {
		log.Printf("检测到本地视频URL: %s，改为发送消息并附带链接", videoURL)

		// 构建消息内容
		messageText := caption + "\n\n视频链接: " + videoURL

		// 使用sendMessage API
		messageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
		messageData := map[string]interface{}{
			"chat_id":    config.TelegramChatID,
			"text":       messageText,
			"parse_mode": "Markdown",
		}

		jsonData, err := json.Marshal(messageData)
		if err != nil {
			return err
		}

		resp, err := http.Post(messageURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// 检查响应
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("Telegram API 错误: %s", string(body))
		}

		return nil
	}

	// 对于公网可访问的URL，直接使用URL
	// 构建请求体
	data := map[string]interface{}{
		"chat_id":    config.TelegramChatID,
		"video":      videoURL,
		"caption":    caption,
		"parse_mode": "Markdown",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API 错误: %s", string(body))
	}

	return nil
}

func SendTelegramMessage(message string) error {
	config := GetNotifyConfig()
	if config == nil || !config.TelegramEnabled {
		return fmt.Errorf("Telegram 推送未启用")
	}

	// 构建请求URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)

	// 构建请求体
	data := map[string]interface{}{
		"chat_id":    config.TelegramChatID,
		"text":       message,
		"parse_mode": "Markdown",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API 错误: %s", string(body))
	}

	return nil
}

// 发送Twitter消息
func SendTwitter(message string) error {
	config := GetNotifyConfig()
	if !config.TwitterEnabled {
		return nil
	}
	if config.TwitterApiKey == "" || config.TwitterApiSecret == "" ||
		config.TwitterAccessToken == "" || config.TwitterAccessTokenSecret == "" {
		return errors.New("Twitter配置不完整")
	}

	// 先尝试 v1.1 API
	err := sendTwitterV1(message)
	if err == nil {
		return nil // v1.1 成功则直接返回
	}

	log.Printf("Twitter v1.1 API 失败，尝试 v2 API: %v", err)
	return sendTwitterV2(message)
}

// Twitter v1.1 API 实现
func sendTwitterV1(message string) error {
	config := GetNotifyConfig()

	oauthConfig := oauth1.NewConfig(config.TwitterApiKey, config.TwitterApiSecret)
	token := oauth1.NewToken(config.TwitterAccessToken, config.TwitterAccessTokenSecret)
	httpClient := oauth1.NewClient(oauth1.NoContext, oauthConfig, token)

	apiUrl := "https://api.twitter.com/1.1/statuses/update.json"
	params := map[string]string{"status": message}

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求错误: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Twitter v1.1错误[%d]: %s", resp.StatusCode, string(body))
	}

	log.Printf("Twitter v1.1推送成功: %s", string(body))
	return nil
}

// Twitter v2 API 实现
func sendTwitterV2(message string) error {
	config := GetNotifyConfig()

	oauthConfig := oauth1.NewConfig(config.TwitterApiKey, config.TwitterApiSecret)
	token := oauth1.NewToken(config.TwitterAccessToken, config.TwitterAccessTokenSecret)
	httpClient := oauth1.NewClient(oauth1.NoContext, oauthConfig, token)

	apiUrl := "https://api.twitter.com/2/tweets"
	payload := map[string]interface{}{
		"text": message,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求错误: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Twitter v2错误[%d]: %s", resp.StatusCode, string(body))
	}

	log.Printf("Twitter v2推送成功: %s", string(body))
	return nil
}

// 发送自定义HTTP消息
func SendCustomHttp(message string) error {
	config := GetNotifyConfig()
	if !config.CustomHttpEnabled || config.CustomHttpUrl == "" {
		return nil
	}
	method := strings.ToUpper(config.CustomHttpMethod)
	if method == "" {
		method = "POST"
	}
	// 处理headers
	headers := map[string]string{}
	if config.CustomHttpHeaders != "" {
		_ = json.Unmarshal([]byte(config.CustomHttpHeaders), &headers)
	}
	// 处理body模板
	bodyStr := config.CustomHttpBody
	if bodyStr == "" {
		bodyStr = `{"content":"` + message + `"}`
	} else {
		bodyStr = strings.ReplaceAll(bodyStr, "{{content}}", message)
	}
	req, err := http.NewRequest(method, config.CustomHttpUrl, bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("自定义HTTP发送失败: %v", resp.Status)
	}
	return nil
}

func UpdateNotifyConfig(db *gorm.DB, config *NotifyConfig) error {
    // 只存在一条记录，无需主键ID
    return db.Model(&NotifyConfig{}).Updates(config).Error
}

func SendEmail(to string, subject string, body string) error {
    db := GetDB()
    var cfg SiteConfig
    if err := db.Table("site_configs").First(&cfg).Error; err != nil {
        return err
    }
    if !cfg.SmtpEnabled {
        return fmt.Errorf("邮件未启用")
    }
    host := strings.TrimSpace(cfg.SmtpHost)
    port := cfg.SmtpPort
    user := strings.TrimSpace(cfg.SmtpUser)
    pass := cfg.SmtpPass
    from := strings.TrimSpace(cfg.SmtpFrom)
    enc := strings.ToLower(strings.TrimSpace(cfg.SmtpEncryption))
    if from == "" {
        from = user
    }
    addr := fmt.Sprintf("%s:%d", host, port)
    headers := make(map[string]string)
    headers["From"] = from
    headers["To"] = to
    headers["Subject"] = subject
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/plain; charset=utf-8"
    var msg strings.Builder
    for k, v := range headers {
        msg.WriteString(k)
        msg.WriteString(": ")
        msg.WriteString(v)
        msg.WriteString("\r\n")
    }
    msg.WriteString("\r\n")
    msg.WriteString(body)
    auth := smtp.PlainAuth("", user, pass, host)
    if enc == "ssl" {
        tlsConfig := &tls.Config{ServerName: host}
        conn, err := tls.Dial("tcp", addr, tlsConfig)
        if err != nil {
            return err
        }
        c, err := smtp.NewClient(conn, host)
        if err != nil {
            return err
        }
        if err := c.Auth(auth); err != nil {
            c.Close()
            return err
        }
        if err := c.Mail(from); err != nil {
            c.Close()
            return err
        }
        if err := c.Rcpt(to); err != nil {
            c.Close()
            return err
        }
        wc, err := c.Data()
        if err != nil {
            c.Close()
            return err
        }
        if _, err := wc.Write([]byte(msg.String())); err != nil {
            wc.Close()
            c.Close()
            return err
        }
        wc.Close()
        c.Quit()
        return nil
    }
    c, err := smtp.Dial(addr)
    if err != nil {
        return err
    }
    if enc == "tls" || cfg.SmtpTLS {
        tlsConfig := &tls.Config{ServerName: host}
        if err := c.StartTLS(tlsConfig); err != nil {
            c.Close()
            return err
        }
    }
    if err := c.Auth(auth); err != nil {
        c.Close()
        return err
    }
    if err := c.Mail(from); err != nil {
        c.Close()
        return err
    }
    if err := c.Rcpt(to); err != nil {
        c.Close()
        return err
    }
    wc, err := c.Data()
    if err != nil {
        c.Close()
        return err
    }
    if _, err := wc.Write([]byte(msg.String())); err != nil {
        wc.Close()
        c.Close()
        return err
    }
    wc.Close()
    c.Quit()
    return nil
}

func SendEmailWithFrom(to string, subject string, body string, fromName string) error {
    db := GetDB()
    var cfg SiteConfig
    if err := db.Table("site_configs").First(&cfg).Error; err != nil {
        return err
    }
    if !cfg.SmtpEnabled {
        return fmt.Errorf("邮件未启用")
    }
    host := strings.TrimSpace(cfg.SmtpHost)
    port := cfg.SmtpPort
    user := strings.TrimSpace(cfg.SmtpUser)
    pass := cfg.SmtpPass
    from := strings.TrimSpace(cfg.SmtpFrom)
    enc := strings.ToLower(strings.TrimSpace(cfg.SmtpEncryption))
    if from == "" { from = user }
    addr := fmt.Sprintf("%s:%d", host, port)
    headers := make(map[string]string)
    if strings.TrimSpace(fromName) != "" {
        headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
    } else {
        headers["From"] = from
    }
    headers["To"] = to
    headers["Subject"] = subject
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/plain; charset=utf-8"
    var msg strings.Builder
    for k, v := range headers {
        msg.WriteString(k)
        msg.WriteString(": ")
        msg.WriteString(v)
        msg.WriteString("\r\n")
    }
    msg.WriteString("\r\n")
    msg.WriteString(body)
    auth := smtp.PlainAuth("", user, pass, host)
    if enc == "ssl" {
        tlsConfig := &tls.Config{ServerName: host}
        conn, err := tls.Dial("tcp", addr, tlsConfig)
        if err != nil { return err }
        c, err := smtp.NewClient(conn, host)
        if err != nil { return err }
        if err := c.Auth(auth); err != nil { c.Close(); return err }
        if err := c.Mail(from); err != nil { c.Close(); return err }
        if err := c.Rcpt(to); err != nil { c.Close(); return err }
        wc, err := c.Data(); if err != nil { c.Close(); return err }
        if _, err := wc.Write([]byte(msg.String())); err != nil { wc.Close(); c.Close(); return err }
        wc.Close(); c.Quit(); return nil
    }
    c, err := smtp.Dial(addr); if err != nil { return err }
    if enc == "tls" || cfg.SmtpTLS { tlsConfig := &tls.Config{ServerName: host}; if err := c.StartTLS(tlsConfig); err != nil { c.Close(); return err } }
    if err := c.Auth(auth); err != nil { c.Close(); return err }
    if err := c.Mail(from); err != nil { c.Close(); return err }
    if err := c.Rcpt(to); err != nil { c.Close(); return err }
    wc, err := c.Data(); if err != nil { c.Close(); return err }
    if _, err := wc.Write([]byte(msg.String())); err != nil { wc.Close(); c.Close(); return err }
    wc.Close(); c.Quit(); return nil
}

func SendTestEmail(to string) error {
    subject := "Ech0-Noise 邮件发送测试"
    body := "这是一封测试邮件，用于验证SMTP配置是否正确。"
    return SendEmail(to, subject, body)
}
func SendEmailHTML(to string, subject string, htmlBody string) error {
    db := GetDB()
    var cfg SiteConfig
    if err := db.Table("site_configs").First(&cfg).Error; err != nil {
        return err
    }
    if !cfg.SmtpEnabled {
        return fmt.Errorf("邮件未启用")
    }
    host := strings.TrimSpace(cfg.SmtpHost)
    port := cfg.SmtpPort
    user := strings.TrimSpace(cfg.SmtpUser)
    pass := cfg.SmtpPass
    from := strings.TrimSpace(cfg.SmtpFrom)
    enc := strings.ToLower(strings.TrimSpace(cfg.SmtpEncryption))
    if from == "" {
        from = user
    }
    addr := fmt.Sprintf("%s:%d", host, port)
    headers := map[string]string{
        "From":         from,
        "To":           to,
        "Subject":      subject,
        "MIME-Version": "1.0",
        "Content-Type": "text/html; charset=utf-8",
    }
    var msg strings.Builder
    for k, v := range headers {
        msg.WriteString(k)
        msg.WriteString(": ")
        msg.WriteString(v)
        msg.WriteString("\r\n")
    }
    msg.WriteString("\r\n")
    msg.WriteString(htmlBody)
    auth := smtp.PlainAuth("", user, pass, host)
    if enc == "ssl" {
        tlsConfig := &tls.Config{ServerName: host}
        conn, err := tls.Dial("tcp", addr, tlsConfig)
        if err != nil { return err }
        c, err := smtp.NewClient(conn, host)
        if err != nil { return err }
        if err := c.Auth(auth); err != nil { c.Close(); return err }
        if err := c.Mail(from); err != nil { c.Close(); return err }
        if err := c.Rcpt(to); err != nil { c.Close(); return err }
        wc, err := c.Data(); if err != nil { c.Close(); return err }
        if _, err := wc.Write([]byte(msg.String())); err != nil { wc.Close(); c.Close(); return err }
        wc.Close(); c.Quit(); return nil
    }
    c, err := smtp.Dial(addr); if err != nil { return err }
    if enc == "tls" || cfg.SmtpTLS { tlsConfig := &tls.Config{ServerName: host}; if err := c.StartTLS(tlsConfig); err != nil { c.Close(); return err } }
    if err := c.Auth(auth); err != nil { c.Close(); return err }
    if err := c.Mail(from); err != nil { c.Close(); return err }
    if err := c.Rcpt(to); err != nil { c.Close(); return err }
    wc, err := c.Data(); if err != nil { c.Close(); return err }
    if _, err := wc.Write([]byte(msg.String())); err != nil { wc.Close(); c.Close(); return err }
    wc.Close(); c.Quit(); return nil
}

func SendEmailHTMLWithFrom(to string, subject string, htmlBody string, fromName string) error {
    db := GetDB()
    var cfg SiteConfig
    if err := db.Table("site_configs").First(&cfg).Error; err != nil {
        return err
    }
    if !cfg.SmtpEnabled {
        return fmt.Errorf("邮件未启用")
    }
    host := strings.TrimSpace(cfg.SmtpHost)
    port := cfg.SmtpPort
    user := strings.TrimSpace(cfg.SmtpUser)
    pass := cfg.SmtpPass
    from := strings.TrimSpace(cfg.SmtpFrom)
    enc := strings.ToLower(strings.TrimSpace(cfg.SmtpEncryption))
    if from == "" { from = user }
    addr := fmt.Sprintf("%s:%d", host, port)
    headers := map[string]string{
        "To":           to,
        "Subject":      subject,
        "MIME-Version": "1.0",
        "Content-Type": "text/html; charset=utf-8",
    }
    if strings.TrimSpace(fromName) != "" {
        headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
    } else {
        headers["From"] = from
    }
    var msg strings.Builder
    for k, v := range headers {
        msg.WriteString(k)
        msg.WriteString(": ")
        msg.WriteString(v)
        msg.WriteString("\r\n")
    }
    msg.WriteString("\r\n")
    msg.WriteString(htmlBody)
    auth := smtp.PlainAuth("", user, pass, host)
    if enc == "ssl" {
        tlsConfig := &tls.Config{ServerName: host}
        conn, err := tls.Dial("tcp", addr, tlsConfig)
        if err != nil { return err }
        c, err := smtp.NewClient(conn, host)
        if err != nil { return err }
        if err := c.Auth(auth); err != nil { c.Close(); return err }
        if err := c.Mail(from); err != nil { c.Close(); return err }
        if err := c.Rcpt(to); err != nil { c.Close(); return err }
        wc, err := c.Data(); if err != nil { c.Close(); return err }
        if _, err := wc.Write([]byte(msg.String())); err != nil { wc.Close(); c.Close(); return err }
        wc.Close(); c.Quit(); return nil
    }
    c, err := smtp.Dial(addr); if err != nil { return err }
    if enc == "tls" || cfg.SmtpTLS { tlsConfig := &tls.Config{ServerName: host}; if err := c.StartTLS(tlsConfig); err != nil { c.Close(); return err } }
    if err := c.Auth(auth); err != nil { c.Close(); return err }
    if err := c.Mail(from); err != nil { c.Close(); return err }
    if err := c.Rcpt(to); err != nil { c.Close(); return err }
    wc, err := c.Data(); if err != nil { c.Close(); return err }
    if _, err := wc.Write([]byte(msg.String())); err != nil { wc.Close(); c.Close(); return err }
    wc.Close(); c.Quit(); return nil
}

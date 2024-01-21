package sender

import (
	"bytes"
	"encoding/json"
	"golangcode/alertmanager-webhook/model"
	"golangcode/alertmanager-webhook/transformer"
	"io/ioutil"
	"log"
	"net/http"
)

func SendToQywechat(notification model.Notification, qywechatKey string, redisServer, redisPort, redisPassword string) {
	message, err := transformer.TransformToMarkdown(notification, redisServer, redisPort, redisPassword)
	if err != nil {
		log.Println(err)
		return
	}
	if qywechatKey != "" {
		var qywechatRobotURL string
		qywechatRobotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + qywechatKey

		// 如果有告警信息才发送
		if (message.QywechatMessage.MarkdownFiring.Markdown != model.Markdown{}) {
			dataFiring, err := json.Marshal(message.QywechatMessage.MarkdownFiring)
			if err != nil {
				log.Println(err)
				return
			}

			reqFiring, err := http.NewRequest(
				"POST",
				qywechatRobotURL,
				bytes.NewBuffer(dataFiring))

			if err != nil {
				log.Println(err)
				return
			}
			reqFiring.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			respFiring, err := client.Do(reqFiring)
			if err != nil {
				log.Println(err)
				return
			}
			defer respFiring.Body.Close()

			firingbody, err := ioutil.ReadAll(respFiring.Body)
			if err != nil {
				log.Println(err)
			}
			firingresponse := make(map[string]interface{})
			err = json.Unmarshal(firingbody, &firingresponse)
			if err != nil {
				log.Println(err)
			}
			if int(firingresponse["errcode"].(float64)) != 0 {
				log.Println("send alert message to qywechat error: ", firingresponse)
			}
		}

		// 如果有恢复信息才发送
		if (message.QywechatMessage.MarkdownResolved.Markdown != model.Markdown{}) {
			dataResolved, err := json.Marshal(message.QywechatMessage.MarkdownResolved)
			if err != nil {
				log.Println(err)
				return
			}

			reqResolved, err := http.NewRequest(
				"POST",
				qywechatRobotURL,
				bytes.NewBuffer(dataResolved))

			if err != nil {
				log.Println(err)
				return
			}

			reqResolved.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			respResolved, err := client.Do(reqResolved)
			if err != nil {
				log.Println(err)
				return
			}
			defer respResolved.Body.Close()

			resolvedbody, err := ioutil.ReadAll(respResolved.Body)
			if err != nil {
				log.Println(err)
			}
			resolvedresponse := make(map[string]interface{})
			err = json.Unmarshal(resolvedbody, &resolvedresponse)
			if err != nil {
				log.Println(err)
			}
			if int(resolvedresponse["errcode"].(float64)) != 0 {
				log.Println("send resolved message to qywechat error: ", resolvedresponse)
			}
		}
	} else {
		log.Println("qywechat key doesn't exist")
	}
}

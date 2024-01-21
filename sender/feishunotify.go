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

func SendToFeishu(notification model.Notification, feishuKey string, redisServer, redisPort, redisPassword string) {
	message, err := transformer.TransformToMarkdown(notification, redisServer, redisPort, redisPassword)
	if err != nil {
		log.Println(err)
		return
	}
	if feishuKey != "" {
		var feishuRobotURL string
		feishuRobotURL = "https://open.feishu.cn/open-apis/bot/v2/hook/" + feishuKey

		// 如果有告警信息才发送
		if message.FeishuMessage.CardFiring.Card.Elements[0].Text.Content != "" {
			dataFiring, err := json.Marshal(message.FeishuMessage.CardFiring)
			if err != nil {
				log.Println(err)
				return
			}

			reqFiring, err := http.NewRequest(
				"POST",
				feishuRobotURL,
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
			if int(firingresponse["code"].(float64)) != 0 {
				log.Println("send alert message to feishu error: ", firingresponse)
			}
		}

		// 如果有恢复信息才发送
		if message.FeishuMessage.CardResolved.Card.Elements[0].Text.Content != "" {
			dataResolved, err := json.Marshal(message.FeishuMessage.CardResolved)
			if err != nil {
				log.Println(err)
				return
			}

			reqResolved, err := http.NewRequest(
				"POST",
				feishuRobotURL,
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
			if int(resolvedresponse["code"].(float64)) != 0 {
				log.Println("send resolved message to feishu error: ", resolvedresponse)
			}
		}
	} else {
		log.Println("feishu key doesn't exist")
	}
}

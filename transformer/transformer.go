package transformer

import (
	"bytes"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"golangcode/alertmanager-webhook/model"
	"log"
	"reflect"
	"text/template"
	"time"
)

func TransformToMarkdown(notification model.Notification, redisServer, redisPort, redisPassword string) (message *model.Message, err error) {
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisServer, redisPort))
	if err != nil {
		log.Println("connect redis failed: ", err)
		return
	}
	defer c.Close()

	if redisPassword != "" {
		_, err = c.Do("AUTH", redisPassword)
		if err != nil {
			log.Println("Redis auth error: ", err)
			return
		}
	}

	var notificationFiring model.Notification
	var notificationResolved model.Notification

	var cstZone = time.FixedZone("CST", 8*3600)

	var bufferFiring bytes.Buffer
	var bufferResolved bytes.Buffer

	for _, alert := range notification.Alerts {
		status := alert.Status
		if status == "firing" {
			notificationFiring.Version = notification.Version
			notificationFiring.GroupKey = notification.GroupKey
			notificationFiring.Status = "firing"
			notificationFiring.Receiver = notification.Receiver
			notificationFiring.GroupLabels = notification.GroupLabels
			notificationFiring.CommonLabels = notification.CommonLabels
			notificationFiring.ExternalURL = notification.ExternalURL
			notificationFiring.Alerts = append(notificationFiring.Alerts, alert)
		} else if status == "resolved" {
			notificationResolved.Version = notification.Version
			notificationResolved.GroupKey = notification.GroupKey
			notificationResolved.Status = "resolved"
			notificationResolved.Receiver = notification.Receiver
			notificationResolved.GroupLabels = notification.GroupLabels
			notificationResolved.CommonLabels = notification.CommonLabels
			notificationResolved.ExternalURL = notification.ExternalURL
			notificationResolved.Alerts = append(notificationResolved.Alerts, alert)
		}
	}

	if !reflect.DeepEqual(notificationFiring, model.Notification{}) {
		//bufferFiring.WriteString(fmt.Sprintf("# <font color=\"red\">触发告警</font>\n"))
		for _, alert := range notificationFiring.Alerts {
			//annotations := alert.Annotations
			alert.StartTime = alert.StartsAt.In(cstZone).Format("2006-01-02 15:04:05")
			fingerprint := alert.Fingerprint
			_, err = c.Do("HSet", fingerprint, "startTime", alert.StartTime)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = c.Do("Hincrby", fingerprint, "count", 1)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = c.Do("HSet", fingerprint, "message", alert)
			if err != nil {
				log.Println(err)
				return
			}

			count, err := redis.Int(c.Do("HGet", fingerprint, "count"))
			if err != nil {
				log.Println("get alert count error: ", err)
				//return
			}
			alert.Count = count
			tmpl, err := template.ParseFiles("./template/alert.tmpl")
			err = tmpl.Execute(&bufferFiring, alert)
			if err != nil {
				log.Println("get firing message error:\n", err)
			}
		}
	}
	if !reflect.DeepEqual(notificationResolved, model.Notification{}) {
		//bufferResolved.WriteString(fmt.Sprintf("# <font color=\"green\">告警恢复</font>\n"))
		for _, alert := range notificationResolved.Alerts {
			//annotations := alert.Annotations
			fingerprint := alert.Fingerprint
			alert.StartTime, err = redis.String(c.Do("HGet", fingerprint, "startTime"))
			if err != nil {
				log.Println("get alert startTime error: ", err)
				//return
			}
			alert.EndTime = alert.EndsAt.In(cstZone).Format("2006-01-02 15:04:05")
			tmpl, err := template.ParseFiles("./template/alert.tmpl")
			err = tmpl.Execute(&bufferResolved, alert)
			if err != nil {
				log.Println("get resolved message error:\n", err)
			}
			_, err = c.Do("Del", fingerprint)
			if err != nil {
				log.Println("delete key error: ", err)
			}
		}
	}

	// 转换为企业微信可以识别的格式
	var markdownFiring, markdownResolved *model.QyWeChatMarkdown
	var title string
	title = "# <font color=\"red\">触发告警</font>\n"
	if bufferFiring.String() != "" {
		markdownFiring = model.NewQyWeChatMarkdown(title + bufferFiring.String())
	} else {
		markdownFiring = model.NewQyWeChatMarkdown("")
	}
	title = "# <font color=\"green\">告警恢复</font>\n"
	if bufferResolved.String() != "" {
		markdownResolved = model.NewQyWeChatMarkdown(title + bufferResolved.String())
	} else {
		markdownResolved = model.NewQyWeChatMarkdown("")
	}

	// 转换为飞书可以识别的格式
	var cardFiring, cardResolved *model.CardMsg
	cardFiring = model.NewCardMsg("red", "触发告警")
	if bufferFiring.String() != "" {
		cardFiring.AddElement(bufferFiring.String())
	} else {
		cardFiring.AddElement("")
	}

	cardResolved = model.NewCardMsg("green", "告警恢复")
	if bufferResolved.String() != "" {
		cardResolved.AddElement(bufferResolved.String())
	} else {
		cardResolved.AddElement("")
	}

	// 转换为钉钉可以识别的格式
	var dmarkdownFiring, dmarkdownResolved *model.DingDingMarkdown
	title = "# <font color=\"#FF0000\">触发告警</font>\n\n"
	if bufferFiring.String() != "" {
		dmarkdownFiring = model.NewDingDingMarkdown("触发告警", title+bufferFiring.String())
	} else {
		dmarkdownFiring = model.NewDingDingMarkdown("", "")
	}

	title = "# <font color=\"#008000\">告警恢复</font>\n\n"
	if bufferResolved.String() != "" {
		dmarkdownResolved = model.NewDingDingMarkdown("告警恢复", title+bufferResolved.String())
	} else {
		dmarkdownResolved = model.NewDingDingMarkdown("", "")
	}

	// 将企业微信、飞书、钉钉消息进行封装
	message = model.NewMessage(markdownFiring, markdownResolved, cardFiring, cardResolved, dmarkdownFiring, dmarkdownResolved)

	return
}

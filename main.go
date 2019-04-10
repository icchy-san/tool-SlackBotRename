package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

// Env ... Variable for environment loading
var Env envConfig

func init() {
	loadEnvironment()
}

func main() {
	var fp *os.File
	var err error
	var faildChannels []string
	var randomTime int
	rand.Seed(time.Now().UnixNano())

	client := slack.New(Env.AccessToken)

	if len(os.Args) < 2 {
		os.Exit(1)
	} else {
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		randomTime = rand.Intn(2)
		delay := time.Duration((randomTime + 1)) * time.Second
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		channelID, newChannelName := record[0], record[1]
		channelInfo, _ := client.GetConversationInfo(channelID, false)

		channel, err := client.RenameChannel(channelID, newChannelName)
		if err != nil {
			faildChannels = append(faildChannels, "#"+channelID)
			continue
		}

		if err != nil {
			fmt.Printf("Error: %s", err)
			continue
		}

		log.Printf("CID: %#v, previousName: %#v, requestName: %#v, NewCName: %#v\n", channelID, channelInfo.Name ,newChannelName, channel.Name)

		// 数秒くらい待ってあげましょうよ、という気持ちの現れ
		time.Sleep(delay)
	}

	attachmentParams := setResultMessageParameters()

	faildChannelsText := strings.Join(faildChannels, ", ")
	textOpt := slack.MsgOptionText(faildChannelsText, false)
	client.PostMessage(Env.AdminChannelID, textOpt, attachmentParams)
}

// setResultMessageParameters ... 実行結果ごのメッセージ設定関数
func setResultMessageParameters() slack.MsgOption {
	attachment := slack.Attachment{
		Text: "プロセスが終了しました。チャネル名が表示された場合は、表示さているチャンネルをリネームすることができませんでした。",
	}
	return slack.MsgOptionAttachments(attachment)
}



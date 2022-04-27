package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	// vcsession         *discordgo.VoiceConnection
	HelloWorld        = "helloworld"
	Channels          = "channels"
	ChannelVoiceJoin  = "vcjoin"
	ChannelVoiceLeave = "vcleave"
)

func main() {
	loadEnv()

	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))

	if err != nil {
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)

	err = discord.Open()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")

	// channel作成
	stopBot := make(chan os.Signal, 1)

	// OSからのシグナルをキャッチする？
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// 処理を堰き止める
	<-stopBot

	fmt.Println("Interrupted.")

	err = discord.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Closed.")

}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err)
	}
}

func commandIs(command string, ses *discordgo.Session, mc *discordgo.MessageCreate) bool {
	return strings.HasPrefix(mc.Content, fmt.Sprintf("%s %s", fmt.Sprintf("<@%s>", os.Getenv("CLIENT_ID")), command))
}

// メッセージ受信ハンドラ
func onMessageCreate(ses *discordgo.Session, mc *discordgo.MessageCreate) {
	// if err != nil {
	// 	log.Println("Error getting channel: ", err)
	// 	return
	// }

	fmt.Printf("%20s %20s %20s > %s\n", mc.ChannelID, time.Now().Format(time.Stamp), mc.Author.Username, mc.Content)

	switch {
	case commandIs(HelloWorld, ses, mc):
		sendHelloWorld(ses, mc.ChannelID)
	case commandIs(Channels, ses, mc):
		sendChannels(ses, mc)
	}
}

// メッセージ送信関数
func sendMessage(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

// コマンド: helloworld
func sendHelloWorld(ses *discordgo.Session, channelId string) {
	sendMessage(ses, channelId, "Hello World!")
}

// コマンド: channels
func sendChannels(ses *discordgo.Session, mc *discordgo.MessageCreate) {
	st, err := ses.GuildChannels(mc.GuildID)
	if err != nil {
		fmt.Println("channels command error")
		fmt.Println(err)
	}

	var lines []string
	for _, v := range st {
		line := fmt.Sprintf("Name: %s(%s) - ID: %s", v.Name, v.Type, v.ID)
		lines = append(lines, line)
	}
	joinedLines := strings.Join(lines, "\n")
	fmt.Println(joinedLines)

	sendMessage(ses, mc.ChannelID, joinedLines)
}

// Cloud Text-to-Speech API呼び出し
// func fetchTextToSpeech(text string) {
// 	ctx := context.Background()

// 	client, err := texttospeech.NewClient(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	req := texttospeechpb.SynthesizeSpeechRequest{
// 		Input: &texttospeechpb.SynthesisInput{
// 			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
// 		},
// 		Voice: &texttospeechpb.VoiceSelectionParams{
// 			LanguageCode: "ja-JP",
// 			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
// 		},
// 		AudioConfig: &texttospeechpb.AudioConfig{
// 			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
// 		},
// 	}

// 	resp, err := client.SynthesizeSpeech(ctx, &req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	filename := "output.mp3"
// 	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Audio content written to file: %v\n", filename)
// }

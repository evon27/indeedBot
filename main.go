package main

import (
	"discord/lib"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	Prefix string = "/"
	Token string = "Nzc5MzA4NzY0NDYxMTM3OTgw.X7ep2A.SQO-MQLe1NnJ4w5OG9yvPMhoeww"
)

func main() {
	discord, err := discordgo.New("Bot " + Token); checkErr(err)

	discord.AddHandler(messageCreate)

	err = discord.Open(); checkErr(err)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func messageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID { return }
	if strings.HasPrefix(msg.Content, Prefix) {
		fmt.Println(msg.Content)
	}

	args := strings.Split(msg.Content, " ")
	command := strings.Trim(args[0], Prefix)

	if command == "indeed" {
		lib.Scrape(args[1])
		file, err := os.Open("jobs.csv"); checkErr(err)
		reply(session, msg, "Success!")
		session.ChannelFileSend(msg.ChannelID, args[1]+".csv", file)
		os.Remove("jobs.csv")
	}
}

func reply(session *discordgo.Session, msg *discordgo.MessageCreate, content string) {
	mention := "<@"+msg.Author.ID+">"
	session.ChannelMessageSend(msg.ChannelID, mention + " " + content)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
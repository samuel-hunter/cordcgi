/*
 * DiscordCGI: run commands from some cgi-bin/ directory as commands
 */
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const SettingsFile = "settings.json"

var logger = log.New(os.Stdout, "Discordcgi: ",
	log.Ldate|log.Ltime|log.Lshortfile)

var Settings struct {
	ClientID string
	Token    string
	Prefix   string
}

func loadSettingsOrPanic() {
	f, err := os.Open(SettingsFile)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&Settings); err != nil {
		panic(err)
	}
}

func botUrl() string {
	return fmt.Sprintf("https://discordapp.com/oauth2/authorize"+
		"?client_id=%s&scope=bot&permissions=2048",
		Settings.ClientID)
}

func main() {
	loadSettingsOrPanic()
	fmt.Println("Invite this bot at", botUrl())

	discord, err := discordgo.New("Bot " + Settings.Token)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(messageCreate)

	// Open websocket connection and begin listening
	if err = discord.Open(); err != nil {
		panic(err)
	}

	// Wait here or until ^C or other term signal is received.
	logger.Println("Bot is now running. Press ^C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the session gracefully.
	logger.Println("Closing gracefully...")
	discord.Close()
	logger.Println("Bye!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger.Println("New message: ", m.Content)
}

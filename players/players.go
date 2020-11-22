package players

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// Player struct
type Player struct {
	Connection       net.Conn                // TCP
	GameData         *map[string]interface{} // Vars of game
	NeedSendGameData bool                    // Flag to send game data to the player
	Online           bool                    // Is online
}

// HandleConnection is a main func for tcp connection
func (p *Player) HandleConnection() {
	defer p.Close()
	for p.Online {
		p.CheckOnline()
		if p.NeedSendGameData {
			p.SendGameData()
			fmt.Println(p.Read(2048))
			p.NeedSendGameData = false
		}
	}
}

func (p *Player) CheckOnline() {
	p.Connection.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	if _, err := p.Connection.Read(make([]byte, 1)); err == io.EOF {
		p.Online = false
	}
	p.Connection.SetReadDeadline(time.Time{})
}

// Write to socket
func (p *Player) Write(s string) {
	p.Connection.Write([]byte(s))
}

// SendGameData sends gamedata in json
func (p *Player) SendGameData() {
	bytes, _ := json.Marshal(*p.GameData)
	p.Connection.Write(bytes)
}

// Read from socket. bufLen = buffer length
func (p *Player) Read(bufLen int) string {
	buf := make([]byte, bufLen)
	_, err := p.Connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		p.Online = false
	}
	return string(buf)
}

// Close socket safely
func (p *Player) Close() {
	fmt.Println("Closed", p.Connection.RemoteAddr().String())
	p.Connection.Close()
}

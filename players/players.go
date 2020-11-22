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
	GameMap          *[][]int                // Hex grid
	NeedSendGameData bool                    // Flag to send game data to the player
	NeedSendGameMap  bool                    // Flag to send game map to the player
	Online           bool                    // Is online
	Trusted          bool                    // Protection against connections from outside the client application
}

// HandleConnection is a main func for tcp connection
func (p *Player) HandleConnection() {
	defer p.Close()
	for p.Online {
		p.CheckOnline()

		if p.NeedSendGameData {
			p.SendGameData()
			s, err := p.ReadBytes(2048)
			if !err {
				if !p.CheckAnswer(&s) {
					p.Online = false
				}
			}
			p.NeedSendGameData = false
		}

		if p.NeedSendGameMap {
			p.SendGameMap()
			s, err := p.ReadBytes(2048)
			if !err {
				if !p.CheckAnswer(&s) {
					p.Online = false
				}
			}
			p.NeedSendGameMap = false
		}

	}
}

// CheckAnswer checks if player answer is correct
func (p *Player) CheckAnswer(ans *[]byte) bool {
	var dat map[string]bool
	if err := json.Unmarshal(*ans, &dat); err != nil {
		return false
	}
	p.Trusted = true
	return dat["success"]
}

// CheckOnline is sending 1 byte to client for checking his online
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

// SendGameMap sends hex grid in json
func (p *Player) SendGameMap() {
	bytes, _ := json.Marshal(map[string]interface{}{
		"map": *p.GameMap,
	})
	p.Connection.Write(bytes)
}

// Read string from socket. bufLen = buffer length
func (p *Player) Read(bufLen int) (string, bool) {
	bytesAnsw, err := p.ReadBytes(bufLen)
	return string(bytesAnsw), err
}

// ReadBytes from socket. bufLen = buffer length
func (p *Player) ReadBytes(bufLen int) ([]byte, bool) {
	buf := make([]byte, bufLen)
	p.Connection.SetReadDeadline(time.Now().Add(2 * time.Second))
	ansLen, err := p.Connection.Read(buf)
	p.Connection.SetReadDeadline(time.Time{})
	if err != nil {
		p.Online = false
		return []byte("Error reading: " + string(err.Error())), true
	}
	return buf[:ansLen], false
}

// Close socket safely
func (p *Player) Close() {
	p.Trusted = false
	fmt.Println("Closed", p.Connection.RemoteAddr().String())
	p.Connection.Close()
}

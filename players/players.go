package players

import (
	"fmt"
	"net"
)

// Player struct
type Player struct {
	Connection net.Conn // TCP
}

// HandleConnection is a main func for tcp connection
func (p *Player) HandleConnection() {
	p.Write("Write message: ")
	fmt.Println(p.Read(2048))
	p.Write("Message received.")
	p.Connection.Close()
	fmt.Println("Closed", p.Connection.RemoteAddr().String())
}

func (p *Player) Write(s string) {
	p.Connection.Write([]byte(s))
}

func (p *Player) Read(bufLen int) string {
	buf := make([]byte, bufLen)
	_, err := p.Connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	return string(buf)
}

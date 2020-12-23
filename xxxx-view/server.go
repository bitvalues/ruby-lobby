package view

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/bitvalues/ruby-lobby/src/config"
	"github.com/bitvalues/ruby-lobby/src/data"
	"github.com/davecgh/go-spew/spew"
)

type Server struct {
	cfg    config.Config
	socket net.Listener
	data   *data.Server
}

func CreateServer(cfg config.Config, data *data.Server) *Server {
	return &Server{
		cfg:  cfg,
		data: data,
	}
}

func (s *Server) StartListening() {
	socket, err := net.Listen("tcp4", fmt.Sprintf(":%d", s.cfg.ViewPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("DEBUG: Socket is listening for connections...")
	s.socket = socket
	rand.Seed(time.Now().Unix())

	for {
		connection, err := s.socket.Accept()
		if err != nil {
			fmt.Println("DEBUG: Cannot accept connection...")
			log.Fatal(err)
		}

		// for every accepted connection, open a new go routine
		go s.HandleConnection(connection)
	}
}

func (s *Server) StopListening() {
	if s.socket != nil {
		s.socket.Close()
		s.socket = nil
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	fmt.Printf("View connection from: %s\n", conn.RemoteAddr().String())

	for {
		// try reading a chunk of data
		len, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("[view] eof")
				break
			} else if strings.Contains(err.Error(), "connection was forcibly closed") {
				fmt.Println("[view] closed forcibly")
				break
			}

			fmt.Println("DEBUG: Could not read buffer from connection")
			fmt.Println(err)
			break
		}

		fmt.Printf("[view] Read %d byte(s) of data\n", len)
		fmt.Printf("[view] received code: %d\n", buffer[8])

		switch buffer[8] {
		case 0x26:
			// TODO: Check version data
			// https://github.com/project-topaz/topaz/blob/f5fb1cddff2bd6f8c0dd1c49a9c4d7a180fb6b20/src/login/lobby.cpp#L527
			versionData := buffer[0x74:(0x74 + 10)]
			fmt.Printf("[view] Client is on version: %s\n", versionData)

			// sessionHash := buffer[12:28]

			packet := []byte{
				0x28, 0x00, 0x00, 0x00, 0x49, 0x58, 0x46, 0x46, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4f, 0xe0, 0x5d, 0xad,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}

			// expansion bitmask
			// 00000000 Bit0 - Not Used - Original FFXI bit
			// 00000010 Bit1 - Enables Rise of Zilart Icon
			// 00000100 Bit2 - Enables Chains of Promathia Icon
			// 00001000 Bit3 - Enables Treasures of Aht Urhgan Icon
			// 00010000 Bit4 - Enables Wings of The Goddess
			// 00100000 Bit5 - Enables A Crystalline Prophecy Icon
			// 01000000 Bit6 - Enables A Moogle Kupod'Etat Icon
			// 10000000 Bit7 - Enables A Shantotto Ascension Icon
			expansion := make([]byte, 2)
			binary.LittleEndian.PutUint16(expansion, 14)
			copy(packet[32:], expansion[:])

			// feature bitmask
			// 00000001 Bit0 - Enables Vision of Abyssea
			// 00000010 Bit1 - Enables Scars of Abyssea
			// 00000100 Bit2 - Enables Heroes of Abyssea
			// 00001000 Bit3 - Enables Seekers of Adoulin
			// 00010000 Bit4 - Not Used - Future expansion
			// 00100000 Bit5 - Not Used - Future expansion
			// 01000000 Bit6 - Not Used - Future expansion
			// 10000000 Bit7 - Not Used - Future expansion
			feature := make([]byte, 2)
			binary.LittleEndian.PutUint16(feature, 13)
			copy(packet[36:], feature[:])

			// hash the byte array
			hash := md5.Sum(packet)
			copy(packet[12:], hash[:])
			spew.Dump(packet)
			conn.Write(packet)
		case 0x1F:
			spew.Dump(buffer)
			packet := make([]byte, 5)
			packet[0] = 0x01

			s.data.GetConnection().Write(packet)
		}
	}
}

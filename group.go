package main

import (
	"fmt"
	"github.com/jes/go-ricochet/utils"
	"github.com/jes/ricochetbot"
	"log"
	"strings"
)

// avoidPeer can be nil to send a message to everyone
func SendToAll(bot *ricochetbot.RicochetBot, avoidPeer *ricochetbot.Peer, message string) {
	for _, p := range bot.Peers {
		if avoidPeer == nil || p.Onion != avoidPeer.Onion {
			p.SendMessage(message)
		}
	}
}

var ricochet2Nick map[string]string
var nick2Ricochet map[string]string

func main() {
	pk, err := utils.LoadPrivateKeyFromFile("./private_key")
	if err != nil {
		log.Fatalf("error reading private key file: %v", err)
	}

	ricochet2Nick = make(map[string]string)
	nick2Ricochet = make(map[string]string)

	commands := InitCommands()

	bot := new(ricochetbot.RicochetBot)
	bot.PrivateKey = pk

	bot.OnConnect = func(peer *ricochetbot.Peer) {
		fmt.Println("We connected to ", peer.Onion)
	}
	bot.OnNewPeer = func(peer *ricochetbot.Peer) bool {
		fmt.Println(peer.Onion, "connected to us")
		SendToAll(bot, peer, "*** "+peer.Onion+" has connected.")
		return true // true == already-known contact
	}
	bot.OnReadyToChat = func(peer *ricochetbot.Peer) {
		fmt.Println(peer.Onion, "ready to chat")
		peer.SendMessage("*** welcome to ricochet group chat.")
	}
	bot.OnMessage = func(peer *ricochetbot.Peer, message string) {
		if message[0] == '/' {
			words := strings.Fields(message)
			cmd, exists := commands[words[0]]
			if exists {
				cmd(peer, message, words)
			} else {
				peer.SendMessage("*** unrecognised command: " + words[0])
			}
		} else {
			name := peer.Onion
			nick, exists := ricochet2Nick[peer.Onion]
			if exists {
				name = nick
			}
			message = "<" + name + "> " + message
			SendToAll(bot, peer, message)
		}
	}
	bot.OnContactRequest = func(peer *ricochetbot.Peer, name string, desc string) bool {
		fmt.Println(peer.Onion, "wants to be our friend")
		return true // true == accept
	}
	bot.OnDisconnect = func(peer *ricochetbot.Peer) {
		fmt.Println(peer.Onion, "disconnected")
		SendToAll(bot, peer, "*** "+peer.Onion+" has disconnected.")

		nick, exists := ricochet2Nick[peer.Onion]
		if exists {
			delete(ricochet2Nick, peer.Onion)
			delete(nick2Ricochet, nick)
		}
	}

	bot.Run()
}

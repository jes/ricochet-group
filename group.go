package main

import (
	"fmt"
	"github.com/jes/go-ricochet/utils"
	"github.com/jes/ricochetbot"
	"log"
)

func SendToAll(bot *ricochetbot.RicochetBot, avoidPeer *ricochetbot.Peer, message string) {
	for _, p := range bot.Peers {
		if p.Onion != avoidPeer.Onion {
			p.SendMessage(message)
		}
	}
}

func main() {
	pk, err := utils.LoadPrivateKeyFromFile("./private_key")
	if err != nil {
		log.Fatalf("error reading private key file: %v", err)
	}

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
	bot.OnMessage = func(peer *ricochetbot.Peer, message string) {
		message = "<" + peer.Onion + "> " + message
		SendToAll(bot, peer, message)
	}
	bot.OnContactRequest = func(peer *ricochetbot.Peer, name string, desc string) bool {
		fmt.Println(peer.Onion, "wants to be our friend")
		return true // true == accept
	}
	bot.OnDisconnect = func(peer *ricochetbot.Peer) {
		fmt.Println(peer.Onion, "disconnected")
		SendToAll(bot, peer, "*** "+peer.Onion+" has disconnected.")
	}

	bot.Run()
}

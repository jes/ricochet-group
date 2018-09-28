package main

import (
	"github.com/jes/ricochetbot"
	"sort"
	"strings"
)

func InitCommands() map[string]func(*ricochetbot.Peer, string, []string) {
	m := make(map[string]func(*ricochetbot.Peer, string, []string))

	m["/help"] = func(peer *ricochetbot.Peer, message string, words []string) {
		peer.SendMessage(
			`Commands available:

  /help - Show this text
  /nick foo - Set your nickname
  /who - List connected peers
  /whois foo - Show the ricochet id for the given nickname
`)
	}

	m["/nick"] = func(peer *ricochetbot.Peer, message string, words []string) {
		if len(words) != 2 {
			peer.SendMessage("usage: /nick foo")
			return
		}

		curRicochet, exists := nick2Ricochet[words[1]]
		if exists {
			peer.SendMessage("The nick '" + words[1] + "' is already taken by " + curRicochet)
			return
		}

		oldnick, exists := ricochet2Nick[peer.Onion]
		if exists {
			delete(nick2Ricochet, oldnick)
		}

		// TODO: make sure new nick only has allowable characters and is short enough
		ricochet2Nick[peer.Onion] = words[1]
		nick2Ricochet[words[1]] = peer.Onion
		SendToAll(peer.Bot, nil, "*** "+peer.Onion+" is now known as "+words[1])
	}

	m["/who"] = func(peer *ricochetbot.Peer, message string, words []string) {
		peers := make([]string, 0)
		for _, p := range peer.Bot.Peers {
			text := p.Onion
			nick, exists := ricochet2Nick[p.Onion]
			if exists {
				text += " (" + nick + ")"
			}
			peers = append(peers, text)
		}
		sort.Slice(peers, func(a int, b int) bool {
			return strings.Compare(peers[a], peers[b]) < 0
		})
		peer.SendMessage("Connected peers:\n" + strings.Join(peers, "\n"))
	}

	m["/whois"] = func(peer *ricochetbot.Peer, message string, words []string) {
		if len(words) != 2 {
			peer.SendMessage("usage: /whois foo")
			return
		}

		onion, exists := nick2Ricochet[words[1]]
		if exists {
			peer.SendMessage(onion + " (" + words[1] + ")")
		} else {
			peer.SendMessage("no such nick: " + words[1])
		}
	}

	return m
}

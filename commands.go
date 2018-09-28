package main

import (
	"github.com/jes/ricochetbot"
	"regexp"
	"sort"
	"strings"
)

var IsAllowableNick = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).MatchString

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
			if curRicochet == peer.Onion {
				peer.SendMessage("But you're already called " + words[1] + "!")
			} else {
				peer.SendMessage("The nick '" + words[1] + "' is already taken by " + curRicochet)
			}
			return
		}

		oldnick, exists := ricochet2Nick[peer.Onion]
		if exists {
			delete(nick2Ricochet, oldnick)
		}

		if len(words[1]) > 16 {
			peer.SendMessage("Maximum of 16 characters for a nick")
			return
		}
		if !IsAllowableNick(words[1]) {
			peer.SendMessage("Nick can only contain letters, numbers, hyphen and underscore")
			return
		}
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

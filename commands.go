package main

import (
	"github.com/jes/ricochetbot"
	"github.com/spf13/viper"
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
  /welcome - Show the welcome message
  /who - List connected peers
  /whois foo - Show the ricochet id for the given nickname
`)

		if IsAdmin(peer.Onion) {
			peer.SendMessage(
				`Admin commands:

  /ban foo - Ban the given ricochet id or nickname
  /invite id - Invite the given ricochet id to the group
  /kick foo - Kick the given ricochet id or nickname
  /welcome [new message] - Set the welcome message
`)
		}
	}

	m["/kick"] = func(peer *ricochetbot.Peer, message string, words []string) {
		if !IsAdmin(peer.Onion) {
			peer.SendMessage("Sorry, only admins can kick.")
			return
		}

		if len(words) != 2 {
			peer.SendMessage("usage: /kick foo")
			return
		}

		kick := peer.Bot.LookupPeerByHostname(words[1])
		if kick == nil {
			peer.SendMessage("No such peer: " + words[1])
			return
		}

		kick.Disconnect()
		SendToAll(peer.Bot, nil, "*** "+words[1]+" was kicked by "+peer.Onion)
	}

	m["/nick"] = func(peer *ricochetbot.Peer, message string, words []string) {
		if len(words) != 2 {
			peer.SendMessage("usage: /nick foo")
			return
		}

		nickLock.Lock()
		defer nickLock.Unlock()

		curRicochet, exists := nick2Onion[words[1]]
		if exists {
			if curRicochet == peer.Onion {
				peer.SendMessage("But you're already called " + words[1] + "!")
			} else {
				peer.SendMessage("The nick '" + words[1] + "' is already taken by " + curRicochet)
			}
			return
		}

		oldnick, exists := onion2Nick[peer.Onion]
		if exists {
			delete(nick2Onion, oldnick)
		}

		if len(words[1]) > 15 {
			peer.SendMessage("Maximum of 15 characters for a nick")
			return
		}
		if !IsAllowableNick(words[1]) {
			peer.SendMessage("Nick can only contain letters, numbers, hyphen and underscore")
			return
		}
		onion2Nick[peer.Onion] = words[1]
		nick2Onion[words[1]] = peer.Onion
		SendToAll(peer.Bot, nil, "*** "+peer.Onion+" is now known as "+words[1])
	}

	m["/welcome"] = func(peer *ricochetbot.Peer, message string, words []string) {
		if IsAdmin(peer.Onion) && len(words) > 1 {
			peer.SendMessage("Sorry, changing the welcome message is not implemented yet")
		} else if len(words) == 1 {
			peer.SendMessage(viper.GetString("welcomemsg"))
		} else {
			peer.SendMessage("usage: /welcome")
		}
	}

	m["/who"] = func(peer *ricochetbot.Peer, message string, words []string) {
		peers := make([]string, 0)
		for _, p := range peer.Bot.Peers {
			text := p.Onion
			nick, exists := onion2Nick[p.Onion]
			if IsAdmin(p.Onion) {
				text = "*" + text
			}
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

		onion, exists := nick2Onion[words[1]]
		if exists {
			if IsAdmin(onion) {
				onion = "*" + onion
			}
			peer.SendMessage(onion + " (" + words[1] + ")")
		} else {
			peer.SendMessage("no such nick: " + words[1])
		}
	}

	return m
}

package main

import (
	"fmt"
	"github.com/jes/go-ricochet/utils"
	"github.com/jes/ricochetbot"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

var nickLock sync.Mutex
var onion2Nick map[string]string
var nick2Onion map[string]string

// avoidPeer can be nil to send a message to everyone
func SendToAll(bot *ricochetbot.RicochetBot, avoidPeer *ricochetbot.Peer, message string) {
	for _, p := range bot.Peers {
		if avoidPeer == nil || p.Onion != avoidPeer.Onion {
			p.SendMessage(message)
		}
	}
}

func IsAdmin(onion string) bool {
	return IsInList(onion, viper.GetStringSlice("admins")) || IsInList(onion, GetList("admins"))
}

func IsBanned(onion string) bool {
	return IsInList(onion, viper.GetStringSlice("bans")) || IsInList(onion, GetList("bans"))
}

func IsAllowedUser(onion string) bool {
	return IsInList(onion, viper.GetStringSlice("allowedusers")) || IsInList(onion, GetList("allowedusers"))
}

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	// paths that come first take priority
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/ricochet-group/")

	viper.SetDefault("welcomemsg", "*** welcome to ricochet group chat.")
	viper.SetDefault("allowedusers", []string{})
	viper.SetDefault("admins", []string{})
	viper.SetDefault("bans", []string{})
	viper.SetDefault("publicgroup", false)
	viper.SetDefault("datadir", ".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// basic sanity check...
	if viper.GetBool("publicgroup") == false && len(viper.GetStringSlice("allowedusers")) == 0 {
		log.Fatalf("Error: ricochet-group is configured to run a private group chat which no users are allowed to connect to")
	}

	// load a private key
	pkFilename := viper.GetString("datadir") + "/private_key"
	pk, err := utils.LoadPrivateKeyFromFile(pkFilename)
	if err != nil {
		// generate a new key if we can't load one
		pkNew, pkErr := utils.GeneratePrivateKey()
		if pkErr != nil {
			log.Fatalf("error reading private key file: %v, and error generating private key: %v", err, pkErr)
		}
		pk = pkNew

		err2 := ioutil.WriteFile(pkFilename, []byte(utils.PrivateKeyToString(pk)), 0600)
		if err2 != nil {
			log.Fatalf("error reading private key file: %v, and error writing private key file: %v", err, err2)
		}
	}

	onion, err := utils.GetOnionAddress(pk)
	if err != nil {
		log.Fatalf("can't get our onion address from our private key ??? %v", err)
	}
	fmt.Println("ricochet-group coming up at ricochet:" + onion + " ...")

	onion2Nick = make(map[string]string)
	nick2Onion = make(map[string]string)

	commands := InitCommands()

	bot := new(ricochetbot.RicochetBot)
	bot.PrivateKey = pk

	bot.OnConnect = func(peer *ricochetbot.Peer) {
		fmt.Println("We connected to ", peer.Onion)
	}
	bot.OnNewPeer = func(peer *ricochetbot.Peer) bool {
		fmt.Println(peer.Onion, "connected to us")
		if !IsInList(peer.Onion, GetList("peers")) {
			SendToAll(bot, peer, "*** "+peer.Onion+" has connected.")
		}
		if !IsBanned(peer.Onion) && (viper.GetBool("publicgroup") == true || IsAllowedUser(peer.Onion)) {
			return true
		} else {
			fmt.Println(peer.Onion + " not allowed! Refusing connection")
			return false
		}
	}
	bot.OnReadyToChat = func(peer *ricochetbot.Peer) {
		fmt.Println(peer.Onion, "ready to chat")
		if !IsInList(peer.Onion, GetList("peers")) {
			peer.SendMessage(viper.GetString("welcomemsg"))
		}
		AddToList("peers", peer.Onion)
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
			nick, exists := onion2Nick[peer.Onion]
			if exists {
				name = nick
			}
			message = "<" + name + "> " + message
			SendToAll(bot, peer, message)
		}
	}
	bot.OnContactRequest = func(peer *ricochetbot.Peer, name string, desc string) bool {
		fmt.Println(peer.Onion, "wants to be our friend")
		if viper.GetBool("publicgroup") == true || IsAllowedUser(peer.Onion) {
			return true
		} else {
			fmt.Println(peer.Onion + " not in allowed users list! Refusing contact request")
			return false
		}
	}
	bot.OnDisconnect = func(peer *ricochetbot.Peer) {
		fmt.Println(peer.Onion, "disconnected")

		nickLock.Lock()
		defer nickLock.Unlock()
		nick, exists := onion2Nick[peer.Onion]
		if exists {
			delete(onion2Nick, peer.Onion)
			delete(nick2Onion, nick)
		}
	}

	err = bot.ManageTor(viper.GetString("datadir"))
	if err != nil {
		log.Fatalf("can't start tor: %v", err)
	}
	fmt.Println("Started tor, we're controlling it at " + bot.TorControlAddress)

	go func() {
		// TODO: instead of sleeping 20 seconds, we should have a callback when tor is ready
		time.Sleep(20 * time.Second)
		for _, onion := range GetList("peers") {
			// don't connect out to banned users
			if IsBanned(onion) {
				RemoveFromList("peers", onion)
				continue
			}

			fmt.Println("Trying to connect out to", onion)
			go bot.Connect(onion)
		}
	}()

	bot.Run()
}

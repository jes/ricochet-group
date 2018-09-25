package main

import (
  "github.com/s-rah/go-ricochet/application"
  "github.com/s-rah/go-ricochet/channels"
  "github.com/s-rah/go-ricochet/utils"
  "log"
  "time"
  "fmt"
)

type BigBoy struct {
    bots []*ChatEchoBot
}

type ChatEchoBot struct {
	rai      *application.ApplicationInstance
    bigboy   *BigBoy
}

// We always want bidirectional chat channels
func (bot *ChatEchoBot) OpenInbound() {
	log.Println("OpenInbound() ChatChannel handler called...")
	outboutChatChannel := bot.rai.Connection.Channel("im.ricochet.chat", channels.Outbound)
	if outboutChatChannel == nil {
		bot.rai.Connection.Do(func() error {
			bot.rai.Connection.RequestOpenChannel("im.ricochet.chat",
				&channels.ChatChannel{
					Handler: bot,
				})
			return nil
		})
	}
}

func (bot *ChatEchoBot) ChatMessage(messageID uint32, when time.Time, message string) bool {
	log.Printf("ChatMessage(from: %v, %v", bot.rai.RemoteHostname, message)
    message = bot.rai.RemoteHostname + ": " + message
    MessageEveryone(bot.bigboy, bot, message)
	return true
}

// set "bot" to nil if it should go to everyone
func MessageEveryone(bigboy *BigBoy, bot *ChatEchoBot, message string) {
    for _, otherbot := range bigboy.bots {
        if otherbot != bot {
            SendMessage(otherbot.rai, message)
        }
    }
}

func SendMessage(rai *application.ApplicationInstance, message string) {
	log.Printf("SendMessage(to: %v, %v)\n", rai.RemoteHostname, message)
	rai.Connection.Do(func() error {

		log.Printf("Finding Chat Channel")
		channel := rai.Connection.Channel("im.ricochet.chat", channels.Outbound)
		if channel != nil {
			log.Printf("Found Chat Channel")
			chatchannel, ok := channel.Handler.(*channels.ChatChannel)
			if ok {
				chatchannel.SendMessage(message)
			}
		} else {
			log.Printf("Could not find chat channel")
		}
		return nil
	})
}

func (bot *ChatEchoBot) ChatMessageAck(messageID uint32, accepted bool) {

}



func main() {
  pk, err := utils.LoadPrivateKeyFromFile("./private_key")
  if err != nil {
    log.Fatalf("error reading private key file: %v", err)
  }

  bigboy := &BigBoy{bots: make([]*ChatEchoBot, 0)}


	fmt.Println("Initializing application factory...")
	af := application.ApplicationInstanceFactory{}
	af.Init()

	af.AddHandler("im.ricochet.contact.request", func(rai *application.ApplicationInstance) func() channels.Handler {
		return func() channels.Handler {
			contact := new(channels.ContactRequestChannel)
			contact.Handler = new(application.AcceptAllContactHandler)
			return contact
		}
	})

	fmt.Println("Starting alice...")
	alice := new(application.RicochetApplication)
	fmt.Println("Generating alice's pk...")
	aliceAddr, _ := utils.GetOnionAddress(pk)
	fmt.Println("Seting up alice's onion " + aliceAddr + "...")
	al, err := application.SetupOnion("127.0.0.1:9051", "tcp4", "", pk, 9878)
	if err != nil {
		fmt.Printf("Could not setup Onion for Alice: %v", err)
        return
	}

	fmt.Println("Initializing alice...")
	af.AddHandler("im.ricochet.chat", func(rai *application.ApplicationInstance) func() channels.Handler {
		return func() channels.Handler {
			chat := new(channels.ChatChannel)
			thisbot := &ChatEchoBot{rai: rai, bigboy: bigboy}
			chat.Handler = thisbot
            bigboy.bots = append(bigboy.bots, thisbot)
			return chat
		}
	})
	alice.Init("Alice", pk, af, new(application.AcceptAllContactManager))
	fmt.Println("Running alice...")
	alice.Run(al)
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var legals []string
var not_legals []string

var cname string
var crelease string
var cset string
var cmana string
var cstats string
var crarity string
var ctype string
var ctext string

var ctcg string
var cedh string
var cimg string

func main() {
	sess, err := discordgo.New("") // insert bot token in quotes
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			println("wef")
			return
		}
		if m.ChannelID == "" || m.ChannelID == "" { // commands only work in given channels
			latest := strings.Split(m.Content, " ")
			var toSearch string = ""
			if latest[0] == ".card" {
				for i := 1; i < len(latest); i++ {
					if i == 1 {
						toSearch = toSearch + latest[i]
					}
					if i > 1 {
						toSearch = toSearch + "_" + latest[i]
					}
				}
				var url string = "https://api.scryfall.com/cards/search?order=name&q=" + toSearch

				response, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()

				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				getInfo(string(body))

				var legStr string = "**Legal in:** "
				var nlegStr string = "**Not legal in:** "
				if len(legals) == 0 {
					legStr = legStr + "no formats"
				}
				for i := 0; i < len(legals); i++ {
					if i < len(legals)-1 {
						legStr = legStr + legals[i] + ", "
					}
					if i == len(legals)-1 {
						legStr = legStr + legals[i]
					}
				}

				if len(not_legals) == 0 {
					nlegStr = nlegStr + "no formats"
				}
				for i := 0; i < len(not_legals); i++ {
					if i < len(legals)-1 {
						nlegStr = nlegStr + not_legals[i] + ", "
					}
					if i == len(not_legals)-1 {
						nlegStr = nlegStr + not_legals[i]
					}
				}
				var msg string = cname + "\n" + crelease + "\n" + cset + "\n" + cmana + "\n" + cstats + "\n" +
					crarity + "\n" + ctype + "\n" + ctext + "\n" + legStr + "\n" + nlegStr + "\n" + cimg
				var msg2 = ctcg + "\n" + cedh
				s.ChannelMessageSend(m.ChannelID, msg)
				s.ChannelMessageSend(m.ChannelID, msg2)
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	println("The bot is online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func getInfo(input string) {
	parts := strings.Split(input, "\"")
	var temp string
	var info []string
	var legal []string
	var not_legal []string
	for i := 5; i < len(parts); i++ {
		switch parts[i] {
		case "name":
			temp = ("**Name**  :  " + parts[i+2])
			info = append(info, temp)
			cname = temp
		case "released_at":
			temp = ("**Released**  :  " + parts[i+2])
			info = append(info, temp)
			crelease = temp
		case "large":
			temp = ("[Image](" + parts[i+2] + ")")
			info = append(info, temp)
			cimg = temp
		case "mana_cost":
			rep1 := strings.Replace(parts[i+2], "B", "Bk", 999)
			rep2 := strings.Replace(rep1, "U", "Bu", 999)
			rep3 := strings.Replace(rep2, "W", "Wh", 999)
			rep4 := strings.Replace(rep3, "R", "Rd", 999)
			newStr := strings.Replace(rep4, "G", "Gr", 999)
			temp = ("**Mana cost**  :  " + newStr)
			info = append(info, temp)
			cmana = temp
		case "type_line":
			temp = ("**Card type**  :  " + parts[i+2])
			info = append(info, temp)
			ctype = temp
		case "oracle_text":
			rulesStr := strings.Replace(parts[i+2], "\\n", "    ", 999)
			temp = ("**Card text**  :\n" + rulesStr)
			info = append(info, temp)
			ctext = temp
		case "power":
			temp = ("**Stats**  :  " + parts[i+2] + " / " + parts[i+6])
			info = append(info, temp)
			cstats = temp
		case "legal":
			temp = (parts[i-2])
			legal = append(legal, temp)
		case "not_legal":
			temp = (parts[i-2])
			not_legal = append(not_legal, temp)
		case "set_name":
			temp = ("**Set**  :  " + parts[i+2])
			info = append(info, temp)
			cset = temp
		case "rarity":
			temp = ("**Rarity**  :  " + parts[i+2])
			info = append(info, temp)
			crarity = temp
		case "tcgplayer":
			temp = ("[TCGplayer](" + parts[i+2] + ")")
			info = append(info, temp)
			ctcg = temp
		case "edhrec":
			temp = ("[EDHREC](" + parts[i+2] + ")")
			info = append(info, temp)
			cedh = temp
		}
	}
	legals = legal
	not_legals = not_legal
}

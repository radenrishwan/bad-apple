package badapple

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var TOKEN = flag.String("token", os.Getenv("token"), "Bot token")

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "bad-apple",
		Description: "Running bad apple ",
	},
}

func RunDiscord(fps int, framePath string) {
	flag.Parse()

	// preparing bad apple frame
	fc, err := os.ReadDir(framePath)
	if err != nil {
		log.Fatalln("error while counting frame: ", err)
	}

	sort.Slice(fc, func(i, j int) bool {
		a := strings.Split(fc[i].Name(), "_")
		b := strings.Split(fc[j].Name(), "_")

		a = strings.Split(a[1], ".")
		b = strings.Split(b[1], ".")

		numA, _ := strconv.Atoi(a[0])
		numB, _ := strconv.Atoi(b[0])

		return numA < numB
	})

	frameCount := len(fc)

	discord, err := discordgo.New("Bot " + *TOKEN)
	if err != nil {
		log.Fatalln("Error while start bot: ", err)
	}

	defer discord.Close()

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "bad-apple" {
			// initial message
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Starting Bad Apple...",
				},
			})
			if err != nil {
				log.Printf("Error responding to interaction: %v", err)
				return
			}

			log.Println("Start run bad apple frame...")

			_, err = s.InteractionResponse(i.Interaction)
			if err != nil {
				log.Printf("Error getting interaction response: %v", err)
				return
			}

			for l := 0; l < frameCount; l++ {
				// get the frame path
				framePath := fmt.Sprintf("%s%s", framePath, fc[l].Name())
				log.Printf("Sending frame %v", framePath)

				frame, err := LoadFrame(framePath, 25, 72)
				if err != nil {
					log.Printf("Cannot load frame: %v", err)
					continue
				}

				// generate content
				content := "```\n"
				for _, row := range frame {
					for _, pixel := range row {
						if pixel == 0 {
							content += " "
						} else {
							content += "█"
						}
					}
					content += "\n"
				}
				content += "```" // End code block

				// edit the message with new frame
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &content,
				})

				if err != nil {
					log.Printf("Error editing message: %v", err)
					continue
				}
			}
		}
	})

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", discord.State.User.Username, discord.State.User.Discriminator)
	})

	err = discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for _, v := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Fatalf("Cannot create command: %v", err)
		}

		registeredCommands = append(registeredCommands, v)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")

	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}

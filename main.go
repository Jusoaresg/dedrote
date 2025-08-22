package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	twitch "github.com/gempir/go-twitch-irc/v4"
)

type Data struct {
	TwitchToken string `json:"twitch_token"`
	ChannelName string `json:"channel_name"`
	ObsHost     string `json:"obs_host"`
	ObsPort     int    `json:"obs_port"`
	ObsPassword string `json:"obs_password"`
	TextSource  string `json:"text_source"`
	Mortes      int    `json:mortes"`
}

func saveFile(data Data) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Erro ao gerar JSON: %v", err)
		return
	}

	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		log.Printf("Erro ao salvar JSON: %v", err)
	}
}

func main() {

	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Erro ao ler o arquivo de configuração:", err)
	}

	var data Data
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal("Erro ao transformar os dados do JSON:", err)
	}

	obs, err := goobs.New(fmt.Sprintf("%s:%d", data.ObsHost, data.ObsPort), goobs.WithPassword(data.ObsPassword))
	if err != nil {
		log.Fatalf("Erro ao conectar com o OBS: %v", err)
	}
	fmt.Println("✅ Conectado ao OBS")

	client := twitch.NewClient(data.ChannelName, data.TwitchToken)

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		if strings.HasPrefix(msg.Message, "!mortes") || strings.HasPrefix(msg.Message, "!morte") {
			if msg.User.Badges["moderator"] == 0 && msg.User.Badges["broadcaster"] == 0 {
				return
			}

			args := strings.Split(msg.Message, " ")
			if len(args) < 2 {
				client.Say(data.ChannelName, "Uso: !mortes <número>")
				return
			}

			deathCount := args[1]
			newText := fmt.Sprintf("Mortes: %s", deathCount)

			mortesInt, err := strconv.Atoi(deathCount)
			if err != nil {
				client.Say(data.ChannelName, "Número inválido")
				return
			}

			inputName := data.TextSource
			res, err := obs.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
				InputName: &inputName,
				InputSettings: map[string]interface{}{
					"text": newText,
				},
			})
			if err != nil {
				log.Printf("Erro ao alterar texto: %v", err)
				client.Say(data.ChannelName, "Erro ao alterar o texto no OBS")
				return
			}
			_ = res

			data.Mortes = mortesInt
			saveFile(data)
			client.Say(data.ChannelName, "✅ Contador mudado!")
		}
	})

	client.Join(data.ChannelName)

	err = client.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar na Twitch: %v", err)
	}
}

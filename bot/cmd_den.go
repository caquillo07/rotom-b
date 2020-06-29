package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleDenCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	fmt.Printf("%s", env.args)

	denData, err := b.pokemonRepo.den("1")
	if err != nil {
		return err
	}

	swordField := &discordgo.MessageEmbedField{}
	swordField.Inline = true
	swordField.Name += "HA in Sword"
	for i := 0; i < len(denData.Sword); i++ {
		if denData.Sword[i].Ability != "Standard" {
			swordField.Value += denData.Sword[i].Name + "\n"
		}
	}

	shieldField := &discordgo.MessageEmbedField{}
	shieldField.Inline = true
	shieldField.Name += "HA in Shield"
	for i := 0; i < len(denData.Shield); i++ {
		if denData.Shield[i].Ability != "Standard" {
			shieldField.Value += denData.Shield[i].Name + "\n"
		}
	}

	msgEmbedFields := []*discordgo.MessageEmbedField{swordField, shieldField}

	msgImage := discordgo.MessageEmbedImage{}
	msgImage.URL = "https://caquillo07.github.io/data/dens/den_1.png"

	msgEmbed := b.newEmbed()
	msgEmbed.URL = "https://www.serebii.net/swordshield/maxraidbattles/den1.shtml"
	msgEmbed.Title = "Pokémon found in Den " + "1" + ":"
	msgEmbed.Image = &msgImage
	// msgEmbed.Description = message
	msgEmbed.Fields = msgEmbedFields

	fmt.Printf("%+v\n\n", msgEmbed)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msgEmbed)
	return err
}

package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

// GuildSettings contains all the personal settings for a given guild
type GuildSettings struct {

	// ID internal unique ID
	ID int

	// Name is the name of the server in discord
	Name string `gorm:"column:guild_name"`

	// DiscordID is the server's unique ID given by discord
	DiscordID string `gorm:"column:guild_id"`

	// LastUpdatedBy is the username of the person who updated this config last
	LastUpdatedBy string

	// BotPrefix is the prefix used to check bot commands for this guild
	BotPrefix string

	// BotAdminRoles an array of role IDs that can admin the bot without the
	// guild administrator permission
	BotAdminRoles pq.StringArray

	// BotAdminUsers an array of user IDs that can admin the bot without a guild
	// administrator role
	BotAdminUsers pq.StringArray

	// ListeningChannels is a list of channels that the bot will listen and
	// respond to.
	ListeningChannels GuildSettingChannels

	// CreatedAt the date the user identity was created
	CreatedAt time.Time

	// UpdatedAt the date the user identity was last updated
	UpdatedAt time.Time
}

type GuildSettingChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GuildSettingChannels []*GuildSettingChannel

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *GuildSettingChannels) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := make([]*GuildSettingChannel, 0)
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Value return json value, implement driver.Valuer interface
func (j GuildSettingChannels) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	rawJSON, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return rawJSON, nil
}

// CreateGuildSettings creates a new config for a guild, and save a copy in the
// cache that does not expire.
func (r *Repository) CreateGuildSettings(config *GuildSettings) error {
	if err := r.db.Create(&config).Error; err != nil {
		return err
	}
	r.cache.Set(config.DiscordID, config, cache.NoExpiration)
	return nil
}

// GetGuildSettings returns the guild configuration for the given discord ID.
// This method will first check cache to prevent a DB call, and will hydrate
// cache on miss.
func (r *Repository) GetGuildSettings(guildID string) (*GuildSettings, error) {
	if cached, found := r.cache.Get(guildID); found {
		return cached.(*GuildSettings), nil
	}
	var c GuildSettings
	if err := r.db.Where("guild_id = ?", guildID).Take(&c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	// save in cache
	r.cache.Set(c.DiscordID, &c, cache.NoExpiration)
	return &c, nil
}

// UpdateGuildSettings takes the passed settings and updates the database.
// This method will also update the cache with the new settings.
func (r *Repository) UpdateGuildSettings(settings *GuildSettings) error {
	if err := r.db.Save(settings).Error; err != nil {
		return err
	}

	// update the cache to make sure we dont get off sync
	r.cache.Set(settings.DiscordID, settings, cache.NoExpiration)
	return nil
}

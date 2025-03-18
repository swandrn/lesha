package entity

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Channels    []Channel
	Name        string
	Description string
	Image       string
	UserID      uint
	User        User
}

type Channel struct {
	gorm.Model
	ServerID uint
	Server   Server
	Messages []Message
	Name     string
}

type User struct {
	gorm.Model
	Servers     []Server  `gorm:"many2many:user_servers;"`
	Channels    []Channel `gorm:"many2many:user_channels;"`
	Reactions   []Reaction
	Messages    []Message
	Name        string
	DisplayName string
	Email       string `gorm:"unique"`
	Password    string
	Status      bool
	Friends     []Friendship `gorm:"foreignKey:UserID"`
}

type Friendship struct {
	gorm.Model
	UserID   uint
	FriendID uint
	Status   string `gorm:"default:'pending'"` // Can be "accepted", "pending", "blocked"
	User     User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Friend   User   `gorm:"foreignKey:FriendID;constraint:OnDelete:CASCADE;"`
}

type Message struct {
	gorm.Model
	UserID    uint
	User      User
	Reactions []Reaction `gorm:"constraint:OnDelete:CASCADE;"`
	Medias    []Media    `gorm:"constraint:OnDelete:CASCADE;"`
	ChannelID uint
	Channel   Channel
	Content   string
	Pinned    bool
}

type Reaction struct {
	gorm.Model
	UserID    uint
	User      User
	MessageID uint
	Message   Message
	Emoji     string
}

type Media struct {
	gorm.Model
	MessageID uint
	Message   Message
	Type      string
	Extension string
	Url       string
}

type BlacklistedToken struct {
	gorm.Model
	Token string `gorm:"unique"`
}

package dto

type Config struct {
	Port          string        `mapstructure:"PORT"`
	JWTKey        string        `mapstructure:"JWTKEY"`
	Database      Database      `mapstructure:",squash"`
	MicrosoftAuth MicrosoftAuth `mapstructure:",squash"`
}
type Database struct {
	DBName string `mapstructure:"DATABASE_NAME"`
	DBUser string `mapstructure:"DATABASE_USER"`
}

type MicrosoftAuth struct {
	ClientSecret string `mapstructure:"MICROSOFT_AUTH_CLIENT_SECRET"`
	ClientID     string `mapstructure:"MICROSOFT_AUTH_CLIENT_ID"`
	RedirectURI  string `mapstructure:"MICROSOFT_AUTH_REDIRECT_URI"`
	TeamID       string `mapstructure:"MICROSOFT_AUTH_TEAM_ID"`
	ChannelID    string `mapstructure:"MICROSOFT_AUTH_CHANNEL_ID"`
}

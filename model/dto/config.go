package dto

type Config struct {
	Port     string   `mapstructure:"PORT"`
	Database Database `mapstructure:",squash"`
}
type Database struct {
	DBName string `mapstructure:"DATABASE_NAME"`
	DBUser string `mapstructure:"DATABASE_USER"`
}

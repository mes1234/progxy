package dto

type Config struct {
	Destintations map[string]Destination `mapstructure:"destinations"`
}

type Destination struct {
	Proxied Proxied `mapstructure:"proxied"`
	Port    int     `mapstructure:"port"`
	Name    string  `mapstructure:"name"`
}

type Proxied struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

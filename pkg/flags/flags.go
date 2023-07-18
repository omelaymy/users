package flags

import "flag"

type Flags struct {
	ConfigFile *string `json:"config"`
}

func New() (*Flags, error) {
	flags := &Flags{
		ConfigFile: flag.String("config", "config/config.yml", "config file path"),
	}
	flag.Parse()

	return flags, nil
}

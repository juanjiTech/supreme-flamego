package conf

type GlobalConfig struct {
	MODE        string       `yaml:"Mode"`
	ProgramName string       `yaml:"ProgramName"`
	BaseURL     string       `yaml:"BaseURL"`
	AUTHOR      string       `yaml:"Author"`
	Listen      string       `yaml:"Listen"`
	Port        string       `yaml:"Port"`
	LogPath     string       `yaml:"LogPath"`
	Databases   []Datasource `yaml:"Databases"`
	Caches      []Cache      `yaml:"Caches"`
	SentryDsn   string       `yaml:"SentryDsn"`
	Auth        struct {
		Secret string `yaml:"Secret"`
		Issuer string `yaml:"Issuer"`
	} `yaml:"Auth"`
}

type Datasource struct {
	Key      string `yaml:"Key"`
	Type     string `yaml:"Type"`
	IP       string `yaml:"Ip"`
	PORT     string `yaml:"Port"`
	USER     string `yaml:"User"`
	PASSWORD string `yaml:"Password"`
	DATABASE string `yaml:"Database"`
	Debug    bool   `yaml:"Debug"`
}

type Cache struct {
	Key      string `yaml:"Key"`
	Type     string `yaml:"Type"`
	IP       string `yaml:"Ip"`
	PORT     string `yaml:"Port"`
	PASSWORD string `yaml:"Password"`
	DB       int    `yaml:"Db"`
}

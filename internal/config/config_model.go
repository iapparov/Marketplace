package config


type Ad struct {
	MinLengthTitle       int    `yaml:"min_length_title" env-default:"3"`
	MaxLengthTitle       int    `yaml:"max_length_title" env-default:"100"`
	MinLengthDescription int    `yaml:"min_length_description" env-default:"10"`
	MaxLengthDescription int    `yaml:"max_length_description" env-default:"1000"`
	ImgType              []string `yaml:"img_type" env-default:"jpg,jpeg,png"`
	PriceMin             float64 `yaml:"price_min" env-default:"0.01"`

	AllowedImgTypesMap map[string]bool `yaml:"-"`
}

type Username struct{
	MinLength int `yaml:"min_length" env-default:"3"`
	MaxLength int `yaml:"max_length" env-default:"20"`
	AllowedCharacters string `yaml:"allowed_characters" env-default:"A-Za-z0-9_-"`
	CaseSensitive bool `yaml:"case_sensitive" env-default:"true"`
}

type Password struct {
	MinLength        int    `yaml:"min_length" env-default:"8"`
	MaxLength        int    `yaml:"max_length" env-default:"64"`
	RequireUpper     bool   `yaml:"require_upper" env-default:"true"`
	RequireLower     bool   `yaml:"require_lower" env-default:"true"`
	RequireDigit     bool   `yaml:"require_digit" env-default:"true"`
}

type Config struct {
    Env	string	`yaml:"env" env-default:"local"`
    Http_port	int	`yaml:"http_port" env-default:"8080"`
	Db	string	`yaml:"db" env-default:"./storage/marketplace.db"`
	JWT_ACCESS_SECRET	string	`yaml:"JWT_ACCESS_SECRET" env-default:"YOUR_JWT_SECRET"`
	JWT_REFRESH_SECRET	string	`yaml:"JWT_REFRESH_SECRET" env-default:"YOUR_JWT_SECRET"`
	JWT_EXP_ACCESS_TOKEN	int	`yaml:"JWT_EXP_ACCESS_TOKEN" env-default:"15"`
	JWT_EXP_REFRESH_TOKEN	int	`yaml:"JWT_EXP_REFRESH_TOKEN" env-default:"24"`
	Username  Username `yaml:"username"`
	Password  Password `yaml:"password"`
	Ad        Ad `yaml:"ad"`
}
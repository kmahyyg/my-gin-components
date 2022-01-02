package common_conf

type S3Config struct {
	Addr string `json:"addr"`
	Region string `json:"region,omitempty"`
	AKID string `json:"accesskey"`
	AKSK string `json:"secretkey"`
	Bucket string `json:"bucket"`
	FileEncryptPassword string `json:"encryption,omitempty"`
}

type WebConfig struct {
	ListenOn string `json:"listen"`
	Environment string `json:"environ"`   // prod OR dev
	DomainName string `json:"domain"`
	BizName string `json:"bizname"`
}

type ReCaptchaConfig struct {
	Secret string `json:"sitesecret"`
	SiteKey string `json:"sitekey"`
	Threshold float32 `json:"threshold"`
	InChina bool `json:"china_user"`
	DomainName string `json:"domain"`
}

type RedisConfig struct {
	DBNum int `json:"db_num,omitempty"`
	Addr string `json:"addr"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type TLSConfig struct {
	CertPath string `json:"cert_path,omitempty"`
	KeyPath string `json:"key_path,omitempty"`
}
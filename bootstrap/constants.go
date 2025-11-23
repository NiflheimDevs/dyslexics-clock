package bootstrap

import "time"

type Constants struct {
	Database    DBConst
	JWTKeysPath string
	SSLKeysPath string
	Port        string
	DevelopMode bool
	Kafka       Kafka
	Context     Context
}

type Context struct {
	UserID string
}

type Kafka struct {
	Cdctopics         []string
	GroupidForElastic string
}

type DBConst struct {
	MaxOpenDbConn int32
	MaxIdleDbConn time.Duration
	MaxDbLifeTime time.Duration
}

func NewConstant() *Constants {
	return &Constants{
		Database: DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
		JWTKeysPath: "./internal/jwt",
		SSLKeysPath: "./SSL",
		Port:        ":8080",
		DevelopMode: true,
		Kafka: Kafka{
			GroupidForElastic: "elastic-readers",
			Cdctopics:         []string{"postgres.public.users", "postgres.public.users_career_tag", "postgres.public.project", "postgres.public.project_tag", "postgres.public.team"},
		},
		Context: Context{
			UserID: "userID",
		},
	}
}

package config

import(
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/stackus/dotenv"

	"github.com/irononet/mallbots/internal/rpc"
	"github.com/irononet/mallbots/internal/web"
)

type (
	PGConfig struct{
		Conn string `required:"true"`
	}

	NatsConfig struct{
		URL string `required:"true"`
		Stream string `default:"mallbots"`
	}

	AppConfig struct{
		Environment string
		LogLevel string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG		PGConfig
		Nats	NatsConfig
		Rpc		rpc.RpcConfig
		Web		web.WebConfig
		ShutdownTimeout		time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg AppConfig, err error){
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil{
		return
	}

	err = envconfig.Process("", &cfg)

	return
}
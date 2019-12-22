package server

import (
	"net/http"

	ilogger "github.com/meateam/elasticsearch-logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	envPrefix               = "GWTMP"
	configPort              = "port"
	configAllowOrigins      = "allow_origins"
	configDownloadChromeURL = "chrome_download_url"
	configExternalApmURL    = "external_apm_url"
	configAuthURL           = "auth_url"
	configSupportLink       = "support_link"
	configService1          = "service1"
)

var (
	logger = ilogger.NewLogger()
)

func init() {
	viper.SetDefault(configPort, 8080)
	viper.SetDefault(configService1, "service1:8080")
	viper.SetDefault(configAllowOrigins, "http://localhost*")
	viper.SetDefault(configExternalApmURL, "http://localhost:8200")
	viper.SetDefault(configSupportLink, "https://open.rocket.chat")
	viper.SetDefault(configAuthURL, "http://localhost/auth/login")
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
}

// Server is a structure that holds the http server of the gateway-template.
type Server struct {
	server *http.Server
	conns  []*grpc.ClientConn
}

// NewServer creates a Server of the gateway-template.
func NewServer() *Server {
	router, conns := NewRouter(logger)

	s := &http.Server{
		Addr:           ":" + viper.GetString(configPort),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
	return &Server{server: s, conns: conns}
}

// Listen listens on configPort. Listen returns when listener is closed.
// Listener will be closed when this method returns, if listener is closed with non-nil
// error then it will be logged as fatal.
func (s *Server) Listen() {
	defer func() {
		for _, v := range s.conns {
			v.Close()
		}
	}()

	logger.Infof("server listening on port: %s", viper.GetString(configPort))
	if err := s.server.ListenAndServe(); err != nil {
		logger.Fatalf("%v", err)
	}
}

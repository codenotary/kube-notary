package config

import (
	"fmt"

	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Configuration variables names
const (
	LogLevel       = "log.level"
	WatchNamespace = "watch.namespace"
	WatchInterval  = "watch.interval"
	TrustKeys      = "trust.keys"
	TrustOrg       = "trust.org"
)

const (
	defaultConfigPath = "/etc/kube-notary/config.yaml"
)

// Interface is the kube-notary configuration
type Interface interface {
	LogLevel() log.Level
	Namespace() string
	Interval() time.Duration
	TrustedKeys() []string
	TrustedOrg() string
}

type cfg struct {
	v *viper.Viper
}

// New returns a kube-notary configuration instance
func New() (Interface, error) {
	v := viper.New()
	c := &cfg{
		v: v,
	}

	// Set defaults
	v.SetDefault(LogLevel, "info")
	v.SetDefault(WatchNamespace, "")
	v.SetDefault(WatchInterval, time.Second*60)
	v.SetDefault(TrustKeys, nil)
	v.SetDefault(TrustOrg, "")

	// Setup
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetTypeByDefaultValue(true)
	v.SetConfigFile(defaultConfigPath)

	// Find and read the config file
	err := v.ReadInConfig()
	// just use the default value(s) if the config file was not found
	if _, ok := err.(*os.PathError); ok {
		// logrus.Warnf("no config file '%s' not found. Using default values", defaultConfigPath)
	} else if err != nil { // Handle other errors that occurred while reading the config file
		return nil, fmt.Errorf("fatal error while reading the config file: %s", err)
	}

	// Monitor the changes in the config file
	v.WatchConfig()

	return c, nil

}

// LogLevel returns the log level
func (c cfg) LogLevel() log.Level {
	logLevel := c.v.GetString(LogLevel)
	l, err := log.ParseLevel(logLevel)
	if err != nil {
		l = log.InfoLevel
	}
	return l
}

// Namespace returns the namespace selector string
func (c cfg) Namespace() string {
	return c.v.GetString(WatchNamespace)
}

// Interval returns the watching cycle interval as time.Duration
func (c cfg) Interval() time.Duration {
	return c.v.GetDuration(WatchInterval)
}

// TrustedKeys returns the trusted keys list as a slice of strings
func (c cfg) TrustedKeys() []string {
	return c.v.GetStringSlice(TrustKeys)
}

// TrustedOrg returns the trusted organization ID as string
func (c cfg) TrustedOrg() string {
	return c.v.GetString(TrustOrg)
}

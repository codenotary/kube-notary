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
	LogLevel                   = "log.level"
	WatchNamespace             = "watch.namespace"
	WatchInterval              = "watch.interval"
	TrustKeys                  = "trust.keys"
	TrustOrg                   = "trust.org"
	LcHost                     = "cnc.host"
	LcPort                     = "cnc.port"
	LcCert                     = "cnc.cert"
	LcNoTls                    = "cnc.noTls"
	LcSkipTlsVerify            = "cnc.skipTlsVerify"
	LcCrossLedgerKeyLedgerName = "cnc.ledgerName"
	LcSignerID                 = "cnc.signerID"
	InternalMode               = "internal"
)

const (
	DefaultConfigPath = "/etc/kube-notary/config.yaml" // defined with mountPath in deployment.yaml by kubernetes
)

// Config populates required config values
type Config struct {
	v *viper.Viper
}

// New returns a kube-notary configuration instance
func New(filePath string) (*Config, error) {
	v := viper.New()
	c := &Config{
		v: v,
	}

	// Set defaults
	v.SetDefault(LogLevel, "info")
	v.SetDefault(WatchNamespace, "")
	v.SetDefault(WatchInterval, time.Second*60)
	v.SetDefault(TrustKeys, nil)
	v.SetDefault(TrustOrg, "")
	v.SetDefault(LcHost, "")
	v.SetDefault(LcPort, "3324")
	v.SetDefault(LcNoTls, true)
	v.SetDefault(LcSkipTlsVerify, true)
	v.SetDefault(LcCrossLedgerKeyLedgerName, "")
	v.SetDefault(LcSignerID, "")

	// Setup
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetTypeByDefaultValue(true)
	v.SetConfigFile(filePath)

	// Find and read the config file
	err := v.ReadInConfig()
	// just use the default value(s) if the config file was not found
	if _, ok := err.(*os.PathError); ok {
		log.Warnf("no config file '%s' not found. Using default values", DefaultConfigPath)
	} else if err != nil { // Handle other errors that occurred while reading the config file
		return nil, fmt.Errorf("fatal error while reading the config file: %s", err)
	}

	// Monitor the changes in the config file
	v.WatchConfig()

	return c, nil

}

// LogLevel returns the log level
func (c *Config) LogLevel() log.Level {
	logLevel := c.v.GetString(LogLevel)
	l, err := log.ParseLevel(logLevel)
	if err != nil {
		l = log.InfoLevel
	}
	return l
}

// Namespace returns the namespace selector string
func (c *Config) Namespace() string {
	return c.v.GetString(WatchNamespace)
}

// Interval returns the watching cycle interval as time.Duration
func (c *Config) Interval() time.Duration {
	return c.v.GetDuration(WatchInterval)
}

// TrustedKeys returns the trusted keys list as a slice of strings
func (c *Config) TrustedKeys() []string {
	return c.v.GetStringSlice(TrustKeys)
}

// TrustedOrg returns the trusted organization ID as string
func (c *Config) TrustedOrg() string {
	return c.v.GetString(TrustOrg)
}

// LcHost returns CNC connection host as a string
func (c *Config) LcHost() string {
	return c.v.GetString(LcHost)
}

// LcPort returns CNC connection port as a string
func (c *Config) LcPort() string {
	return c.v.GetString(LcPort)
}

// LcCert returns CNC connection certificate as a string
func (c *Config) LcCert() string {
	return c.v.GetString(LcCert)
}

// LcSkipTlsVerify returns the CNC LcSkipTlsVerify connection property as a bool
func (c *Config) LcSkipTlsVerify() bool {
	return c.v.GetBool(LcSkipTlsVerify)
}

// LcNoTls returns the CNC no tls connection property as a bool
func (c *Config) LcNoTls() bool {
	return c.v.GetBool(LcNoTls)
}

// LcCrossLedgerKeyLedgerName parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is possible
func (c *Config) LcCrossLedgerKeyLedgerName() string {
	return c.v.GetString(LcCrossLedgerKeyLedgerName)
}

// LcSignerID parameter is used to filter result on a specific signer ID.
func (c *Config) LcSignerID() string {
	return c.v.GetString(LcSignerID)
}

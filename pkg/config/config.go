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
	LcHost                     = "cnlc.host"
	LcPort                     = "cnlc.port"
	LcCert                     = "cnlc.cert"
	LcNoTls                    = "cnlc.noTls"
	LcSkipTlsVerify            = "cnlc.skipTlsVerify"
	LcCrossLedgerKeyLedgerName = "cnlc.crossLedgerKeyLedgerName"
	LcSignerID                 = "cnlc.signerID"
)

const (
	defaultConfigPath = "/etc/kube-notary/config.yaml" // defined with mountPath in deployment.yaml by kubernetes
)

// Interface is the kube-notary configuration
type Interface interface {
	LogLevel() log.Level
	Namespace() string
	Interval() time.Duration
	TrustedKeys() []string
	TrustedOrg() string
	LcHost() string
	LcPort() string
	LcCert() string
	LcSkipTlsVerify() bool
	LcNoTls() bool
	LcCrossLedgerKeyLedgerName() string
	LcSignerID() string
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

// LcHost returns CNLC connection host as a string
func (c cfg) LcHost() string {
	return c.v.GetString(LcHost)
}

// LcCert returns CNLC connection port as a string
func (c cfg) LcPort() string {
	return c.v.GetString(LcPort)
}

// LcCert returns CNLC connection certificate as a string
func (c cfg) LcCert() string {
	return c.v.GetString(LcCert)
}

// LcSkipTlsVerify returns the CNLC LcSkipTlsVerify connection property as a bool
func (c cfg) LcSkipTlsVerify() bool {
	return c.v.GetBool(LcSkipTlsVerify)
}

// LcNoTls returns the CNLC no tls connection property as a bool
func (c cfg) LcNoTls() bool {
	return c.v.GetBool(LcNoTls)
}

// LcCrossLedgerKeyLedgerName parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is possible
func (c cfg) LcCrossLedgerKeyLedgerName() string {
	return c.v.GetString(LcCrossLedgerKeyLedgerName)
}

// LcSignerID parameter is used to filter result on a specific signer ID.
func (c cfg) LcSignerID() string {
	return c.v.GetString(LcSignerID)
}

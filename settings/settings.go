package settings

import (
	"crypto/rand"
	"strings"

	"github.com/filebrowser/filebrowser/v2/rules"
)

// AuthMethod describes an authentication method.
type AuthMethod string

// Settings contain the main settings of the application.
type Settings struct {
	Key           []byte              `json:"key"`
	Signup        bool                `json:"signup"`
	CreateUserDir bool                `json:"createUserDir"`
	Defaults      UserDefaults        `json:"defaults"`
	AuthMethod    AuthMethod          `json:"authMethod"`
	Branding      Branding            `json:"branding"`
	Commands      map[string][]string `json:"commands"`
	Shell         []string            `json:"shell"`
	Rules         []rules.Rule        `json:"rules"`
	OnlyOffice    string              `json:"onlyoffice"`
}

// GetRules implements rules.Provider.
func (s *Settings) GetRules() []rules.Rule {
	return s.Rules
}

// Server specific settings.
type Server struct {
	Root                  string `json:"root"`
	BaseURL               string `json:"baseURL"`
	Socket                string `json:"socket"`
	TLSKey                string `json:"tlsKey"`
	TLSCert               string `json:"tlsCert"`
	Port                  string `json:"port"`
	Address               string `json:"address"`
	Log                   string `json:"log"`
	EnableThumbnails      bool   `json:"enableThumbnails"`
	ResizePreview         bool   `json:"resizePreview"`
	EnableExec            bool   `json:"enableExec"`
	TypeDetectionByHeader bool   `json:"typeDetectionByHeader"`
}

// Clean cleans any variables that might need cleaning.
func (s *Server) Clean() {
	s.BaseURL = strings.TrimSuffix(s.BaseURL, "/")
}

// GenerateKey generates a key of 256 bits.
func GenerateKey() ([]byte, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

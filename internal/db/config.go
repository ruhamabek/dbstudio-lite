package db

import (
 	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}



func (c *Config) ConnectionString() string{
	port := c.Port
	if port == 0 {
		port = 5432
	}

	var userInfo *url.Userinfo

	if c.Password != ""{
      userInfo = url.UserPassword(c.User, c.Password)
	} else if c.User != "" {
      userInfo = url.User(c.User)
	}

	sslMode := c.SSLMode

	query := url.Values{}
	query.Set("sslmode", sslMode)

	if sslMode != ""{
		sslMode = "prefer"
	}

	
	u := url.URL{
		Scheme:   "postgres",
		User:     userInfo,
		Host:     fmt.Sprintf("%s:%d", c.Host, port),
		Path:     c.Database,
		RawQuery: query.Encode(),
	}

	return u.String()
}

func ParseURL(rawUrl string)(Config, error){
	  u, err := url.Parse(rawUrl)

	  if err != nil {
		return Config{}, fmt.Errorf("failed to parse url: %v", err)
	  }

	  if u.Scheme != "postgres" && u.Scheme != "postgresql"{
		return Config{}, fmt.Errorf("failed to parse url due to scheme mismatch: %v", err)
	  }

	  cfg := Config{}

	  cfg.Host = u.Hostname()

	  if u.Port() != "" {
		port, err := strconv.Atoi(u.Port())

		if err != nil {
			return Config{}, fmt.Errorf("failed to convert to int: %v", err)
		}

		cfg.Port = port
	  }

	  if u.User != nil {
		cfg.User = u.User.Username()
		if pass, ok := u.User.Password(); ok {
			cfg.Password = pass
		}
	  }

	  cfg.Database = strings.TrimPrefix(u.Path, "/")

	  if sslMode := u.Query().Get("sslmode"); sslMode != ""{
		cfg.SSLMode = sslMode
	  }

	  return cfg, nil
}
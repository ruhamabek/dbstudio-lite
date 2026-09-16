package db

import (
 	"fmt"
	"net/url"
	"strconv"
	"strings"
)


type DriverType string

const (
	DriverPostgres DriverType = "postgres"
	DriverMySQL    DriverType = "mysql"
	DriverSQLite   DriverType = "sqlite"
)

type Config struct {
	Driver   DriverType `json:"driver"`
	Host     string     `json:"host"`
	Port     int        `json:"port"`
	User     string     `json:"user"`
	Password string     `json:"password"`
	Database string     `json:"database"`
	SSLMode  string     `json:"sslMode"`
	FilePath string     `json:"filePath"`
}



func (c *Config) ConnectionString() string{
	if c.Driver == DriverSQLite || c.FilePath != "" {
		return c.FilePath
	}

	if c.Driver == DriverMySQL {
		port := c.Port
		if port == 0 {
			port = 3306
		}

		userAuth := c.User
		if c.Password != "" {
			userAuth = fmt.Sprintf("%s:%s", c.User, c.Password)
		}

		if userAuth != "" {
			return fmt.Sprintf("%s@tcp(%s:%d)/%s", userAuth, c.Host, port, c.Database)
		}
		return fmt.Sprintf("tcp(%s:%d)/%s", c.Host, port, c.Database)
	}

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

func ParseURL(rawUrl string) (Config, error) {
 	if !strings.Contains(rawUrl, "://") && (strings.HasPrefix(rawUrl, "/") || strings.HasSuffix(rawUrl, ".db") || strings.HasSuffix(rawUrl, ".sqlite") || strings.HasSuffix(rawUrl, ".sqlite3")) {
		return Config{
			Driver:   DriverSQLite,
			FilePath: rawUrl,
		}, nil
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse url: %v", err)
	}

 	if u.Scheme == "sqlite" || u.Scheme == "file" {
		filePath := strings.TrimPrefix(rawUrl, u.Scheme+"://")
		return Config{
			Driver:   DriverSQLite,
			FilePath: filePath,
		}, nil
	}

 	if u.Scheme == "mysql" {
		cfg := Config{
			Driver: DriverMySQL,
			Host:   u.Hostname(),
		}

		if u.Port() != "" {
			port, err := strconv.Atoi(u.Port())
			if err != nil {
				return Config{}, fmt.Errorf("invalid port: %v", err)
			}
			cfg.Port = port
		} else {
			cfg.Port = 3306
		}

		if u.User != nil {
			cfg.User = u.User.Username()
			if pass, ok := u.User.Password(); ok {
				cfg.Password = pass
			}
		}

		cfg.Database = strings.TrimPrefix(u.Path, "/")
		return cfg, nil
	}

 	if u.Scheme == "postgres" || u.Scheme == "postgresql" {
		cfg := Config{
			Driver: DriverPostgres,
			Host:   u.Hostname(),
		}

		if u.Port() != "" {
			port, err := strconv.Atoi(u.Port())
			if err != nil {
				return Config{}, fmt.Errorf("invalid port: %v", err)
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

		if sslMode := u.Query().Get("sslmode"); sslMode != "" {
			cfg.SSLMode = sslMode
		}

		return cfg, nil
	}

	return Config{}, fmt.Errorf("unsupported database scheme: %q", u.Scheme)
}
package configs

import (
	"errors"
	"github.com/peterouob/golang_template/tools"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var Config *viper.Viper

func InitViper() {
	config := viper.New()
	wd, err := os.Getwd()
	tools.HandelError("error in os.Getwd()", err)

	rootDir := findRoot(wd)
	config.AddConfigPath(rootDir)
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = config.ReadInConfig()
	tools.HandelError("error in config.()", err)
	Config = config
}

func findRoot(path string) string {
	dir := path
	for {
		if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	tools.HandelError("error in findRoot", errors.New("config.yaml not found in any parent directories"))
	return ""
}

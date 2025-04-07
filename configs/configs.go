package configs

import (
	"errors"
	"github.com/peterouob/golang_template/utils"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var Config *viper.Viper

func InitViper() {
	config := viper.New()
	wd, err := os.Getwd()
	utils.HandelError("error in os.Getwd()", err)

	rootDir := findRoot(wd)
	config.AddConfigPath(rootDir)
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = config.ReadInConfig()
	utils.HandelError("error in config.()", err)
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
	utils.HandelError("error in findRoot", errors.New("config.yaml not found in any parent directories"))
	return ""
}

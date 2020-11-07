package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func loadConfig() *viper.Viper {
	v := viper.New()

	// Default values
	v.SetDefault("server.host", "example.com")
	v.SetDefault("server.port", 9000)

	// Env config
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Flags config
	pflag.Int("server.port", 1234, "Port name")
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// Config files
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	return v
}

func main() {
	v := loadConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	http.ListenAndServe(":"+v.GetString("server.port"), nil)

}

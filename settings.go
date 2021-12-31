package leikari

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var settingsInit bool

func initSettings() {
	if !settingsInit {
		viper.AddConfigPath(".")
		viper.AddConfigPath("conf")
		viper.AddConfigPath("config")
		viper.AddConfigPath("configs")

		viper.SetConfigName("config")

		viper.ReadInConfig()

		venv := viper.New()
		for _, env := range os.Environ() {
			if strings.HasPrefix(strings.ToLower(env), "leikari_") {
				pair := strings.SplitN(env, "=", 2)
				venv.Set(strings.ReplaceAll(strings.ToLower(pair[0]), "_", "."), pair[1])
			}
		}

		viper.MergeConfigMap(venv.AllSettings())
		
		settingsInit = true
	}
}

type Settings interface {
	Get(string) interface{}
	Set(string, interface{})
	GetBool(string) bool
	GetFloat64(string) float64
	GetInt(string) int
	GetIntSlice(string) []int
	GetString(string) string
	GetStringSlice(string) []string
	GetTime(string) time.Time
	GetDuration(string) time.Duration
	IsSet(string) bool

	AllSettings() map[string]interface{}

	GetDefault(string, interface{}) interface{}
	GetDefaultBool(string, bool) bool
	GetDefaultFloat64(string, float64) float64
	GetDefaultInt(string, int) int
	GetDefaultIntSlice(string, ...int) []int
	GetDefaultString(string, string) string
	GetDefaultStringSlice(string, ...string) []string
	GetDefaultTime(string, time.Time) time.Time
	GetDefaultDuration(string, time.Duration) time.Duration

	GetSub(string, ...Option) Settings
}

type SystemSettings interface {
	Settings

	NoSignature() bool
	GetActorSettings(string, ...Option) ActorSettings
}

type ActorSettings interface {
	Settings

	WorkerPoolSize() int
	MessageQueueSize() int
	Async() bool
}

type defaultWrapper struct {
	*viper.Viper
}

func (e *defaultWrapper) GetSub(key string, opts ...Option) Settings {
	if !e.IsSet(key) {
		e.Set(key, make(map[string]interface{}))
	}
	sub := e.Sub(key)
	for _, opt := range opts {
		sub.Set(opt.Name, opt.Value)
	}
	return &defaultWrapper{sub}
}

func (w *defaultWrapper) GetDefault(key string, v interface{}) interface{} {
	if w.IsSet(key) {
		return w.Get(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultBool(key string, v bool) bool {
	if w.IsSet(key) {
		return w.GetBool(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultFloat64(key string, v float64) float64 {
	if w.IsSet(key) {
		return w.GetFloat64(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultInt(key string, v int) int {
	if w.IsSet(key) {
		return w.GetInt(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultIntSlice(key string, v ...int) []int {
	if w.IsSet(key) {
		return w.GetIntSlice(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultString(key string, v string) string {
	if w.IsSet(key) {
		return w.GetString(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultStringSlice(key string, v ...string) []string {
	if w.IsSet(key) {
		return w.GetStringSlice(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultTime(key string, v time.Time) time.Time {
	if w.IsSet(key) {
		return w.GetTime(key)
	}
	return v
}

func (w *defaultWrapper) GetDefaultDuration(key string, v time.Duration) time.Duration {
	if w.IsSet(key) {
		return w.GetDuration(key)
	}
	return v
}

type systemSettings struct {
	*defaultWrapper
}

func newSystemSettings(opts ...Option) SystemSettings {
	initSettings()

	cfg := viper.GetViper().Sub("leikari")
	for _, opt := range opts {
		cfg.Set(opt.Name, opt.Value)
	}
	return &systemSettings{&defaultWrapper{cfg}}
}

func (s *systemSettings) GetActorSettings(name string, opts ...Option) ActorSettings {
	path := fmt.Sprintf("actor.%s", name)
	var cfg *viper.Viper
	if !s.IsSet(path) {
		s.Set(path, map[string]interface{}{
			"workerPool": 1,
			"messageQueue": 1000,
			"async": false,
		})
	}
	cfg = s.Sub(path)
	return newActorSettings(cfg, opts...)
}

func (s *systemSettings) NoSignature() bool {
	return s.GetBool("noSignature")
}

type actorSettings struct {
	*defaultWrapper
}

func newActorSettings(sub *viper.Viper, opts ...Option) ActorSettings {
	for _, opt := range opts {
		sub.Set(opt.Name, opt.Value)
	}
	return &actorSettings{&defaultWrapper{sub}}
}

func (as *actorSettings) WorkerPoolSize() int {
	if as.IsSet("workerPool") {
		wp := as.GetInt("workerPool")
		if wp > 0 {
			return wp
		}
	}
	return 1
}

func (as *actorSettings) MessageQueueSize() int {
	if as.IsSet("messageQueue") {
		mq := as.GetInt("messageQueue")
		if mq > 0 {
			return mq
		}
	}
	return 1000
}

func (as *actorSettings) Async() bool {
	return as.GetBool("async")
}

func init() {
	viper.SetDefault("leikari.loglevel", "INFO")
}
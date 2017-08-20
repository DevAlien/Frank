package managers

import (
	"frank/src/go/models"
	"frank/src/go/plugins"
)

type PluginsManager struct {
	Plugins map[string]Plugin
}

type Plugin struct {
	Name   string
	Type   string
	Plugin interface{}
}

type pluginI interface {
	ExecAction(models.CommandAction, map[string]string)
}

var ActivePlugins PluginsManager

func NewPlugins() {
	ActivePlugins = PluginsManager{
		Plugins: map[string]Plugin{},
	}
	p := plugins.NewPluginMusicStream()
	ActivePlugins.AddPlugin(Plugin{
		Name:   "music-stream",
		Type:   "stream",
		Plugin: &p,
	})

	p2 := plugins.NewPluginFirmata()
	ActivePlugins.AddPlugin(Plugin{
		Name:   "firmata",
		Type:   "device",
		Plugin: &p2,
	})

}

func (ctx *PluginsManager) AddPlugin(plugin Plugin) {
	ctx.Plugins[plugin.Name] = plugin
}

func (ctx *PluginsManager) ExecAction(action models.CommandAction, extraText map[string]string) {
	plugin := ctx.Plugins[action.Plugin].Plugin.(pluginI)
	plugin.ExecAction(action, extraText)
}

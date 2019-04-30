package managers

import (
	"frank/src/go/models"
	"frank/src/go/plugins"
)

//PluginsManager is the struct that defines the manager for the plugins.
type PluginsManager struct {
	Plugins map[string]Plugin
}

//Plugin is the struct that defines at Plugin
type Plugin struct {
	Name   string
	Type   string
	Plugin interface{}
}

//pluginI is an interface for all the plugins, they have to implement ExecAction and ExecReading
type pluginI interface {
	ExecAction(models.Action, map[string]string)
	ExecReading(models.Reading) (models.ReadingResponse, error)
}

//ActivePlugins stores the single instance for the Plugin Manager with all the active Plugins. This is initiated by the NewBlugin function.
var ActivePlugins PluginsManager

//NewPlugins is the part who loads the plugins and configures them to be used.
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

	p3 := plugins.NewPluginSonoff()
	ActivePlugins.AddPlugin(Plugin{
		Name:   "Sonoff",
		Type:   "device",
		Plugin: &p3,
	})

	p4 := plugins.NewPluginHTTP()
	ActivePlugins.AddPlugin(Plugin{
		Name:   "http",
		Type:   "device",
		Plugin: &p4,
	})

}

//AddPlugin is a simple way of adding a plugin into the ActivePlugins list.
func (ctx *PluginsManager) AddPlugin(plugin Plugin) {
	ctx.Plugins[plugin.Name] = plugin
}

//ExecAction is in charge of selecting the right plugin to Execute an Action.
func (ctx *PluginsManager) ExecAction(action models.Action, extraText map[string]string) {
	plugin := ctx.Plugins[action.Plugin].Plugin.(pluginI)
	plugin.ExecAction(action, extraText)
}

//ExecReading is in charge of the selecting the right plugin to Make  a read.
func (ctx *PluginsManager) ExecReading(action models.Reading) (models.ReadingResponse, error) {
	plugin := ctx.Plugins[action.Plugin].Plugin.(pluginI)
	return plugin.ExecReading(action)
}

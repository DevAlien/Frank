# Frank Assistant

Frank Assistant is a Voice Controlled assistant and not just that.
FA can run even on a Raspberry PI and can even control devices that are in your house, even outside.

It is build using GO because we wanted something that would just run the binary file and you are good to go.

## Why we built it?
We live in an italian speaking country and bots of this kind are all made for the english speakers.
Another reason is becacuse we did not find anything as configurable and as low level.

## How is it built?
- sox (http://sox.sourceforge.net/) is used to record the voice with the effect silence to cut the pieces
- snowboy - To recognize the keyword to make it listen to you.
- Google cloud voice recognition API to transform the voice to text.
- Gobot (https://gobot.io/) to interact with the devices

## How does it work?
Basically the assistant waits to hear the keyword (which is a snowboy file) and once it matches it will activate the voice recognition service.
The voice recognition service listens for your voice, transforms it into text and then it is passed to the commands matcher and then, if necessary will pass the matched commands to the plugins.
You can deactivate frank (so goes back at listening the keyword) using defined deactivation keywords or after X (can be defined) seconds that does not recognize any test it deactivates itself.

## Features
FA has a "plugin" system, basically you can attach commands to plugins and the plguins to the actions and they can connect to the defined actions. 

- [x] Music Stream [music-stream] - This plugin exposes 2 commands to play music from a stream and to stop it
- [x] Firmata [firmata] - Firmata plugin uses gobot to interact with the devices using the firmata protocol.
- [ ] Wather (weather) - You can ask for the weather around the world

## How do I build a plugin?
A plugin is a struct that exposes an `ExecAction(models.CommandAction, map[string]string)` and from that it will know what to do.
Example
```golang
package plugins

import (
	"frank/src/go/models"
)

var stations = map[string]string{
	"club": "URL,
}

type PluginMusicStreamer struct {
}

func NewPluginMusicStream() PluginMusicStreamer {
	pms := PluginMusicStreamer{}

	return pms
}

func (ctx *PluginMusicStreamer) ExecAction(action models.CommandAction, extraText map[string]string) {
	switch action.Action["action"].(string) {
	case "play":
		go ctx.startStream(stations[extraText["type"]], ctx.killCh)
	case "stop":
		go ctx.StopStream()
	}
}

func (ctx *PluginMusicStreamer) stopStream() {
    //STOP STREAM
}

func (ctx *PluginMusicStreamer) startStream(stream string, killChannel chan bool) {
	//START STREAM
}

```

## How is it configured?
FA has a config.json file which will be located at `~/.frank/config.json`. This file contains all the information and configuration for Frank to run.
```json
{
  "name": "Frank The Bot",
  "configs": {
    "google_api_key": ""
  },
  "deactivation": ["shut up"],
  "devices": [],
  "commands": []
}
```
##### Name
Name of the Assistant
##### Configs
Here we can have different configs, even for plugins. `google_api_key` is required at the moment.
##### Deactivation
You can deactivate the Assistant with some special keywords, you can define them here.
##### Devices
Here you can define devices. **(this config structure can change)** An example below of an `arduino` using `firmata` with 2 `interfaces` using `led drivers`.
So basically firmata can connect to it, using the defined PIN can change the value and, in this case, the led will switch on or off.
```json
{
      "name":"livingroom-light", // unique name of the device
      "type":"firmata", // type of the device
      "connection":{ //connection, required by the firmata plugin
        "type":"tcp",   //this means that the connection will be remote and not using a serial
        "address":"192.168.1.7:3030" //address of the arduino using firmata
      },
      "interfaces":[ //here you can define the interfaces, basically the GPIO you will be using
        {
          "name":"blueled", //unique name for the interface
          "driver":"led", // driver to use, 
          "pin":"12" // GPIO which is connected to
        }, {
          "name":"greenled",
          "driver":"led",
          "pin":"14"
        }
      ]
    }
```
##### Commands
The commands are basically the keywords you want to interact with and use with a plugin.
Example 1, using the `music-stream`. Basically `music-stream` expects to know the type of music you want to listen. So we can capture that in the command. `music {type}` will match `music club`, `music rock` and so on. And the `{type}` Will be passed to the plugin so it knows what to do.
```json
{
  "name":"Music On", //name of the command
  "commands":[ //list of commands
    "music {type}"
  ],
  "actions":[ // list of actions that have to be run when the command is matched
    {
      "plugin":"music-stream", //this action uses the music-stream plugin
      "action": { //we can define a json with a lot of data, if the plugin requires it
        "action": "play" //action is the only data required by the music-stream
      }
    }
  ]
}
```

The second command is the one using a device (the one we defined above).
This command basically has a matching interface, which is special, with the captured keword we can use it in the matchinInterface to map the word to an interface. This is because if we say `switch on light green` we want to switch on just the light on the green interface. using the matched `{color}` we can select the interface.
```json
{
      "name":"light on Special",
      "commands":[
        "switch on light {color}"
      ],
      "actions":[
        {
          "plugin": "firmata",
          "device":"livingroom-light",
          "matchingInterface":{
            "color":{
              "blue":"blueled",
              "green":"greenled",
            }
          },
          "action": {
            "action": "on"
          }
        }
      ]
    }
```

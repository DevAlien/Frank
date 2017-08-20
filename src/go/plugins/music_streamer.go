package plugins

import (
	"bytes"
	"os/exec"
	"strings"

	"frank/src/go/helpers/log"
	"frank/src/go/models"
)

var stations = map[string]string{
	"club": "http://club-high.rautemusik.fm",
	"top":  "http://ChartHits-High.RauteMusik.FM",
	"pop":  "http://ChartHits-High.RauteMusik.FM",
}

type PluginMusicStreamer struct {
	killCh    chan bool
	isPlaying bool
}

func NewPluginMusicStream() PluginMusicStreamer {
	pms := PluginMusicStreamer{}
	pms.killCh = make(chan bool, 1)
	pms.isPlaying = false

	return pms
}

func (ctx *PluginMusicStreamer) ExecAction(action models.CommandAction, extraText map[string]string) {
	switch action.Action["action"].(string) {
	case "play":
		if ctx.isPlaying == true {
			ctx.killCh <- true
		}
		go ctx.startStream(stations[extraText["type"]], ctx.killCh)
	case "stop":
		ctx.killCh <- true
	}
}

func (ctx *PluginMusicStreamer) stopStream() {
	ctx.isPlaying = false
}

func (ctx *PluginMusicStreamer) startStream(stream string, killChannel chan bool) {
	ctx.isPlaying = true
	defer ctx.stopStream()

	parts := strings.Fields("-t mp3 " + stream)
	cmd := exec.Command("play", parts...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-killChannel:
		if err := cmd.Process.Kill(); err != nil {
			log.Log.Error("Plugin[music-stream]", err)
			return
		}

		log.Log.Debug(out.String())
		return
	case err := <-done:
		log.Log.Error("Plugin[music-stream]", err)
		return
	}
	return
}

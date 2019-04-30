package services

import (
	"fmt"
	"frank/src/go/helpers/log"
	"frank/src/go/models"
	"github.com/imroc/req"
)

type Ddns struct {
	CurrentIp string
	Config    models.Ddns
}

var DdnsManager Ddns

func LoadDdns(config models.Ddns) Ddns {
	DdnsManager = Ddns{
		Config: config,
	}
	go DdnsManager.SetIp()
	return DdnsManager
}

func (d *Ddns) SetIp() {
	ip, err := GetPublicIp()
	if err != nil {
		log.Log.Error("Could not retrieve the Public Ip")
		log.Log.Error(err.Error())
		return
	}

	log.Log.Debugf("Public Ip Found: %s, Previous Ip: %s ", ip, d.CurrentIp)

	if d.CurrentIp == ip {
		log.Log.Debug("Not Updating")
		return
	}

	d.CurrentIp = ip
	switch d.Config.Type {
	case "noip":
		err = SetNoIp(d.Config, ip)
		if err != nil {
			d.CurrentIp = ""
		}
	}
}

func SetNoIp(cfg models.Ddns, ip string) error {
	url := fmt.Sprintf("http://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s", cfg.Username, cfg.Password, cfg.Hostname, ip)
	_, err := req.Get(url, req.Header{"User-Agent": "Frank/V0.1 g@margalho.info"})
	if err != nil {
		log.Log.Error("error updating NO-IP")
		log.Log.Error(err.Error())
		return err
	}

	log.Log.Debug("Updated Ip on NO-IP")

	return nil

}

func GetPublicIp() (string, error) {
	r, err := req.Get("https://api.ipify.org/")
	if err != nil {
		return "", err
	}
	return r.ToString()
}

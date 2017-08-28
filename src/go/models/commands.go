package models

type CommandAction struct {
	DeviceName        string                 `json:"device,omitempty"`
	MatchingInterface map[string]interface{} `json:"matchingInterface,omitempty"`
	InterfaceName     string                 `json:"interface,omitempty"`
	Plugin            string                 `json:"plugin,omitempty"`
	Action            map[string]interface{} `json:"action"`
}

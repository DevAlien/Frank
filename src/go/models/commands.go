package models

type CommandAction struct {
	DeviceName        string                 `json:"device"`
	MatchingInterface map[string]interface{} `json:"matchingInterface"`
	InterfaceName     string                 `json:"interface"`
	Plugin            string                 `json:"plugin"`
	Action            map[string]interface{} `json:"action"`
}

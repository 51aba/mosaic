package main

import (
  "fmt"
  "log"

  . "leadz/utils"
)

type leadplugin string

const codename string = "RAZDVAPUJCKA"

var configs_map = map[string]string {
  "api_url": "",
  "redirect_url": "https://ads.proficredit.cz/r/wVEYB/vKWRM?cid=",
  "username": "",
  "password": "",
}
var plugin_vars = []string{"link", "",}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "link": []func(map[string]any, map[string]string) (bool){
      redirect_action,},
    "": []func(map[string]any, map[string]string) (bool){
      redirect_action,},
  },
}

// ################################################################################################################################################################
func (p leadplugin) TestData(pPluginData map[string]any, is_paused bool) (ret bool) {
  return p.SendData(pPluginData, is_paused)
}

// ################################################################################################################################################################
func (p leadplugin) Validate(pPluginData map[string]any) ([]map[string]any) {
  return P_validate(codename, pPluginData, plugin_vars, validators_map)
}

// ################################################################################################################################################################
func (p leadplugin) SendData(pPluginData map[string]any, is_paused bool) (result bool) {
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func redirect_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = fmt.Sprintf("%v%v", config["redirect_url" + plugin_postfix], pPluginData["uid"])
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])

  return true
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

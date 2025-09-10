package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "SROVNEJNABIDKYLINK"

var configs_map = map[string]string {
  "redirect_url": "",
}
var plugin_vars = []string{"pd", "cl", "link"}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "link": []func(map[string]any, map[string]string) (bool){
      link_action,
    },
    "pd": []func(map[string]any, map[string]string) (bool){
      pd_action,
    },
    "cl": []func(map[string]any, map[string]string) (bool){
      cl_action,
    },
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
func link_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var variant int

  if "" == config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = fmt.Sprintf("https://track.sjednatpujcku.cz/81484708-796f-445d-9696-ebbbcad4c6ab?lead_id=%v&aff_id=%v&cascade=cl", pPluginData["uid"], pPluginData["affiliate_name"])
  } else {
    if strings.Contains(config["redirect_url" + plugin_postfix], "%v") {
      pPluginData["redirect_url"] = fmt.Sprintf(config["redirect_url" + plugin_postfix], pPluginData["uid"], pPluginData["affiliate_name"])
    } else {
      pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
    }
  }
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL [ %v ]: %v SALE_STATUS: %v", pPluginData["plugin_log"], variant, pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Printf("%v%v%v", YELLOW, pPluginData["description"], NC)

  return true
}

// ################################################################################################################################################################
func pd_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var variant int

  if "" == config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = fmt.Sprintf("https://track.sjednatpujcku.cz/81484708-796f-445d-9696-ebbbcad4c6ab?lead_id=%v&aff_id=%v&cascade=pd", pPluginData["uid"], pPluginData["affiliate_name"])
  } else {
    if strings.Contains(config["redirect_url" + plugin_postfix], "%v") {
      variant = 1
      pPluginData["redirect_url"] = fmt.Sprintf(config["redirect_url" + plugin_postfix], pPluginData["uid"], pPluginData["affiliate_name"])
    } else {
      variant = 2
      pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
    }
  }
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL [ %v ]: %v SALE_STATUS: %v", pPluginData["plugin_log"], variant, pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Printf("%v%v%v", YELLOW, pPluginData["description"], NC)

  return true
}

// ################################################################################################################################################################
func cl_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var variant int

  if "" == config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = fmt.Sprintf("https://track.sjednatpujcku.cz/81484708-796f-445d-9696-ebbbcad4c6ab?lead_id=%v&aff_id=%v&cascade=cl", pPluginData["uid"], pPluginData["affiliate_name"])
  } else {
    if strings.Contains(config["redirect_url" + plugin_postfix], "%v") {
      pPluginData["redirect_url"] = fmt.Sprintf(config["redirect_url" + plugin_postfix], pPluginData["uid"], pPluginData["affiliate_name"])
    } else {
      pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
    }
  }
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL [ %v ]: %v SALE_STATUS: %v", pPluginData["plugin_log"], variant, pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Printf("%v%v%v", YELLOW, pPluginData["description"], NC)

  return true
}

var LeadPlugin leadplugin
// ################################################################################################################################################################

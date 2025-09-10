package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "MMGLEAD"

var configs_map = map[string]string {
  "api_url": "https://api.mmgfg.cz/lead",
  "token": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
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
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "check"}

  pPluginData["map_data"] = map[string]any{
                       "birthnumber": pPluginData["birth_number"],
                       "name": pPluginData["first_name"],
                       "surname": pPluginData["last_name"],
                       "street": pPluginData["street"],
                       "postalCode": pPluginData["zip"],
                       "city": pPluginData["city"],
                     }
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Accept": "application/json", "Authorization": fmt.Sprintf("Bearer %v", config["token" + plugin_postfix])}

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "create"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_cell_phone = "+420" + GetString(pPluginData["cell_phone"])

  var fields = []Pair{{"birthnumber", "birth_number"},
                      {"name", "first_name"},
                      {"surname", "last_name"},
                      {"street", "street"},
                      {"postalCode", "zip"},
                      {"city", "city"},
                      {"email", "email"},
                      {"income", "monthly_income"},
                      {"requestedAmount", "requested_amount"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["token"] = strings.Replace(GetString(pPluginData["uid"]), "-", "", -1)
  data_map["phone"] = translated_cell_phone

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil != ret["errors"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "log"

  . "leadz/utils"
)

type leadplugin string

const codename string = "KVIKU"

var configs_map = map[string]string {
  "api_url": "https://kviku.es/api/cash/",
  "password": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      register_lead,
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
  pPluginData["headers"] = map[string]string{"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}

  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "create_Form"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["password"] = config["password" + plugin_postfix]

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_cell_phone string = "+34" + GetString(pPluginData["cell_phone"])
  var fields = []Pair{{"uid", "uid"},
                      {"amount", "requested_amount"},
                      {"user_ip", "ip_address"},
                      {"user_agent", "user_agent"},
                      {"user[email]", "email"},
                      {"user[name]", "first_name"},
                      {"user[surname]", "last_name"},
                      {"user[middle_name]", "last_name"},
                      {"user[passport]", "identity_card_number"},
                      {"user[birthday]", "birth_date"},
                      {"user[address]", "address"},
                      {"user[city]", "city"},
                      {"user[zipcode]", "zip"},
                    }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["market_id"] = 57
  data_map["user[phone]"] = translated_cell_phone

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  var data_map = GetMap(ret["data"])
  var redirect_url = GetString(data_map["deeplink"])

  if "" == redirect_url {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

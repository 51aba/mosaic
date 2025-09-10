package main

import (
  "fmt"
  "log"

  . "leadz/utils"
)

type leadplugin string

const codename string = "KOLIO"

var configs_map = map[string]string {
  "api_url": "http://system.koliocorporation.cz/www/api/",
  "token": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "birth_number", "func": "InsolvencyValidator"},},
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
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var config map[string]string = P_init_named(codename, pPluginData, configs_map, plugin_postfix)

  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": config["token" + plugin_postfix]}

  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "ping-lead/"}

  pPluginData["map_data"] = map[string]any {
                        "firstName": pPluginData["first_name"],
                        "lastName": pPluginData["last_name"],
                        "subjectCode": pPluginData["birth_number"],
                      }
  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "import-lead/"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

                      /*
                      // {"employerName", "employer"},
                      // {"employerAddress", "employer_address"},
                      // {"employerPhone", "work_phone"},
                      */
  var fields = []Pair{
                      {"firstName", "first_name"},
                      {"lastName", "last_name"},
                      {"email", "email"},
                      {"phone", "cell_phone"},
                      {"subjectCode", "birth_number"},
                      {"zipCode", "zip"},
                      {"amount", "requested_amount"},
                      {"street", "street"},
                      {"housenumber", "house_number"},
                      {"idcard", "identity_card_number"},
                      {"monthlyIncome", "monthly_income"},
                      {"expenses", "monthly_expenses"},
                      {"bankCode", "bank_code"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["address"] = fmt.Sprintf("%v %v, %v", pPluginData["street"], pPluginData["house_number"], pPluginData["city"])
  data_map["account_number"] = pPluginData["iban"]

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if 201 != GetInt(ret["code"]) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

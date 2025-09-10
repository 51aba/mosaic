package main

import (
  "log"
  "fmt"
  "encoding/base64"

  . "leadz/utils"
)

type leadplugin string

const codename string = "SVYCARSKAPUJCKA"

var configs_map = map[string]string {
  "api_url": "https://portal.svycarskapujcka.cz",
  "username": "",
  "password": "",
}
var plugin_vars = []string{"cpl", "",}
var validators_map = map[string][]map[string]any {
  "cpl": {
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 8000},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "identity_card_number", "func": "IdentityCardValidator",},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED", "EMPLOYED", "PENSION", "PART_TIME_EMPLOYMENT",}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 74},},
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpl": []func(map[string]any, map[string]string) (bool){
      register_lead,
    },
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
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "api/lead"}
  var basic = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", config["username" + plugin_postfix], config["password" + plugin_postfix]))))
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": basic}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_income_type string
  var translated_requested_amount (int) = GetInt(pPluginData["requested_amount"])

  switch GetString(pPluginData["income_type"]) {
    case "EMPLOYED":
      translated_income_type = "employed"
    case "UNEMPLOYED":
      translated_income_type = "unemployed"
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = "part-time-employed"
    case "MATERNITY_LEAVE":
      translated_income_type = "maternity"
    case "PENSION":
      translated_income_type = "pensioner"
    case "OTHER":
      translated_income_type = "other"
    default:
      translated_income_type = "other"
  }

  if translated_requested_amount > 5000 {
    translated_requested_amount = 5000
  }
  var client_map = map[string]any {"firstName": pPluginData["first_name"],
                                   "lastName": pPluginData["last_name"],
                                   "birthNumber": pPluginData["birth_number"],
                                   "identityNumber": pPluginData["identity_card_number"],
                                   "phone": pPluginData["cell_phone"],
                                   "mail": pPluginData["email"],
                                   "address": map[string]any{"street": pPluginData["street"], "number": pPluginData["house_number"],
                                                             "city": pPluginData["city"], "zip": pPluginData["zip"],},
                                   "work": map[string]any{"type": translated_income_type, "additionInfo": GetString(pPluginData["monthly_income"])},
                                  }
  if "employed" == translated_income_type {
    GetMap(client_map["work"])["employer"] = GetString(pPluginData["income_type"])
  }
  data_map["sourceId"] = GetString(pPluginData["uid"])
  data_map["client"] = client_map
  data_map["loanRequest"] = map[string]any{"amount": translated_requested_amount, "currency": "CZK", "days": pPluginData["period"]}

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  var redirect_url = GetString(ret["loanCreateUrl"])

  if nil != ret["errors"] || "" == redirect_url {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = ret["leadId"]
  pPluginData["redirect_url"] = redirect_url

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

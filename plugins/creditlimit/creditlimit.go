package main

import (
  "fmt"
  "log"

  . "leadz/utils"
)

type leadplugin string

const codename string = "CREDITLIMIT"

var configs_map = map[string]string{
  "api_url": "https://proxy.ext.ferratum.com/czech/cls/",
}
var plugin_vars = []string{"cpl"}
var validators_map = map[string][]map[string]any {
  "cpl": {
    {},},
  "": {
    {"field": "birth_number", "func": "InsolvencyValidator",},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpl": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "": []func(map[string]any, map[string]string) (bool){
      register_lead,},
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
  var call_config = map[string]any{"command": "Affiliate"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var translated_requested_amount = GetInt(pPluginData["requested_amount"])

  if translated_requested_amount > 8000 {
    translated_requested_amount = 8000
  }
  var params = map[string]any {
    "name_first": GetString(pPluginData["first_name"]),
    "name_last": GetString(pPluginData["last_name"]),
    "gsm": GetString(pPluginData["cell_phone"]),
    "ssn": GetString(pPluginData["birth_number"]),
    "email": GetString(pPluginData["email"]),
    "loan_amount": fmt.Sprintf("%d", translated_requested_amount),
    "loan_term": "30",
    "id_number": GetString(pPluginData["identity_card_number"]),
    "bank_account": fmt.Sprintf("%v/%v", GetString(pPluginData["bank_account_number"]), GetString(pPluginData["bank_code"])),
    "address": GetString(pPluginData["street"]),
    "city": GetString(pPluginData["city"]),
    "zip": GetString(pPluginData["zip"]),
    "referrer": "PARTNER_REF",
    "monthly_income": fmt.Sprintf("%v.00", GetInt(pPluginData["monthly_income"])),
    "income_source": "salary",
    "user_agent": USERAGENT,
    "google_adwords" : false,
    "country_code": "cz",
    "partner_values": map[string]any {
      "affilbox_code": "eh83gbp6",
      "some_partner_value": "partner",
      "other_partner_value": "partner",},}
  data_map["params"] = params
  data_map["id"] = "123"
  data_map["method"] = "submitApplication"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil == ret["result"] {
    return
  }
  var result_map = GetMap(ret["result"])
  log.Printf("%s RESPONSE_DATA: %v [%v]", pPluginData["plugin_log"], ret["result"], result_map)

  if nil == result_map {
    return
  }
  var lead_type int = GetInt(result_map["new"])

  if 1 == lead_type {
    result = true
    pPluginData["sale_status"] = "UNCONFIRMED"
    pPluginData["external_id"] = result_map["order_id"]

    return
  }

  if 0 == lead_type || 1 == GetInt(result_map["duplicate"]) {
    pPluginData["sale_status"] = "DUPLICATE"
  }
  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

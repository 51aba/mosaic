package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "CFG"

var configs_map = map[string]string{
  "api_url": "https://njf51c4adm.affil.comfortfinancegroup.com/v1/PARTNER_NAME/",
  "token": "",
  "channel_id": "",
}
var plugin_vars = []string{"unique", "duplicate", "uniquefl2", "duplicatefl2"}
var validators_map = map[string][]map[string]any{
  "unique": {
    {"field": "home_status", "func": "AllowedValuesValidator", "param1": []string{"HOME_OWNER",}},},
  "duplicate": {
    {"field": "home_status", "func": "AllowedValuesValidator", "param1": []string{"HOME_OWNER",}},},
  "uniquefl2": {
    {"field": "home_status", "func": "AllowedValuesValidator", "param1": []string{"HOME_OWNER",}},},
  "duplicatefl2": {
    {"field": "home_status", "func": "AllowedValuesValidator", "param1": []string{"HOME_OWNER",}},},
  "": {
    {"field": "home_status", "func": "AllowedValuesValidator", "param1": []string{"HOME_OWNER",}},
    {"field": "birth_number", "func": "InsolvencyValidator",},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "unique": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "duplicate": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "uniquefl2": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "duplicatefl2": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
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
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var config map[string]string = P_init_named(codename, pPluginData, configs_map, plugin_postfix)

  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Token %v", config["token" + plugin_postfix])}

  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": ""}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["channel_id"] = config["channel_id" + plugin_postfix]

  pPluginData["map_data"] = data_map

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "contact"}

  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    data_map["channel_id"] = config["channel_id" + plugin_postfix]
    pPluginData["map_data"] = data_map
  }
  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

                      //{"marital_status", "marital_status"},
                      // {"dependent_children", "dependent_children"},
                      /*
                      // {"employed_time", "employed_time"},
                      // {"employer", "employer"},
                      // {"employer_address", "employer_address"},
                      // {"employer_phone", "employer_phone"},
                      // {"job_title", "job_title"},
                      */
  var fields = []Pair{
                      {"first_name", "first_name"},
                      {"last_name", "last_name"},
                      {"birth_number", "birth_number"},
                      {"cell_phone", "cell_phone"},
                      {"email", "email"},
                      {"bank_code", "bank_code"},
                      {"city", "city"},
                      {"contact_city", "city"},
                      {"contact_home_status", "home_status"},
                      {"contact_house_number", "house_number"},
                      {"contact_street", "street"},
                      {"contact_zip", "zip"},
                      {"distraint", "distraint"},
                      {"education", "education"},
                      {"expenses", "monthly_expenses"},
                      {"gender", "gender"},
                      {"home_status", "home_status"},
                      {"house_number", "house_number"},
                      {"identity_card_number", "identity_card_number"},
                      {"income_type", "income_type"},
                      {"insolvency", "insolvency"},
                      {"ip_address", "ip_address"},
                      {"monthly_income", "monthly_income"},
                      {"period", "period"},
                      {"requested_amount", "requested_amount"},
                      {"street", "street"},
                      {"user_agent", "user_agent"},
                      {"zip", "zip"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["bank_account"] = GetString(pPluginData["bank_account_number"])
  data_map["unique_key"] = fmt.Sprintf("%v", strings.Replace(GetString(pPluginData["uid"]), "-", "", -1))
  data_map["address_time"] = "1"
  data_map["company_number"] = nil
  data_map["contract_duration"] = "0"
  data_map["has_guarantee"] = "0"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil == ret["ACCEPTED"] && nil == ret["id"] {
    return
  }
  var accepted  = GetString(ret["ACCEPTED"])
  var external_id = GetString(ret["id"])

  if "" != accepted {
    result = true
    pPluginData["sale_status"] = "UNCONFIRMED"
  } else if "" != external_id {
    result = true
    pPluginData["sale_status"] = "UNCONFIRMED"

    pPluginData["external_id"] = external_id
    pPluginData["redirect_url"] = ret["redirect_url"]
  }
  log.Printf("%s SET_RESPONSE_DATA: %v %v", pPluginData["plugin_log"], accepted, external_id)

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "KIMBI"
var loan_period_variants = []int {6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36}

var configs_map = map[string]string {
  "api_url": "https://www.zaplonasplatky.cz/api/affiliate-api",
  "username": "",
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
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": ""}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_birth_number = GetString(pPluginData["birth_number"])
  var translated_requested_amount (int) = GetInt(pPluginData["requested_amount"])
  var translated_home_status (string) = GetString(pPluginData["home_status"])
  var translated_income_type (string) = "EMPLOYEE"
  var fields = []Pair{{"firstName", "first_name"},
                      {"lastName","last_name"},
                      {"cardId","identity_card_number"},
                      {"mobilePhone","cell_phone"},
                      {"email", "email"},
                      {"street", "street"},
                      {"houseFlatNumber", "house_number"},
                      {"city", "city"},
                      {"postalCode", "zip"},
                      {"password", "user_password"},
                      {"passwordRepeat", "user_password"},
                      {"clientIncome", "monthly_income"},
                      {"clientExpenses", "monthly_expenses"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "OWNED"
    case "TENANT":
      translated_home_status = "RENT"
    case "OTHER":
      translated_home_status = "EMPLOYER"
    default:
      translated_home_status = "OWNED"
  }

  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%s/%s", translated_birth_number[:6], translated_birth_number[6:])
  }
  data_map["amount"] = translated_requested_amount
  data_map["term"] = FindClosest(GetInt(pPluginData["period"]), loan_period_variants)
  data_map["termType"] = "MONTHS"
  data_map["personalId"] = translated_birth_number
  data_map["accommodationType"] = translated_home_status
  data_map["numberOfEmployedHouseholdMembers"] = 2
  data_map["employmentType"] = translated_income_type
  data_map["acceptAgreement"] = true
  data_map["affiliateProvider"] = "54_API"
  data_map["affiliateToken"] = ""
  data_map["newsAccepted"] = true
  data_map["applicationType"] = "WEB"
  data_map["bankAccountNumber"] = fmt.Sprintf("%v/%v", pPluginData["bank_account_number"], pPluginData["bank_code"])

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil == ret["resolution"] || "success" != strings.ToLower(GetString(ret["resolution"])) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = fmt.Sprintf("https://exc4finance.com/?a=&c=&s1=API&s2=%v&ckmrdr=%v", pPluginData["uid"], ret["url"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

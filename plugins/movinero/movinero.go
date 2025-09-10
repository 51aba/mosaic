package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "MOVINERO"

var configs_map = map[string]string {
  "api_url": "https://movinero.cz/webapi/v1/loan/",
  "token": "",
}
var plugin_vars = []string{"cpl", ""}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT", "SELF_EMPLOYED", "PENSION", "MATERNITY_LEAVE"}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 20, "param2": 100},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 7300},},
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
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": fmt.Sprintf("registrations/check/%v", pPluginData["birth_number"])}

  pPluginData["map_data"] = map[string]any{}

  return P_check_unique_get(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "register"}
  var data_map = map[string]any{}
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Accept": "application/json", "Authorization": config["token" + plugin_postfix]}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var fields = []Pair{{"first_name", "first_name"},
                      {"last_name", "last_name"},
                      {"email", "email"},
                      {"phone", "cell_phone"},
                      {"personal_id", "birth_number"},
                      {"income_type", "income_type"},
                      {"neto_income", "monthly_income"},
                      {"next_payday", "next_pay_date"},
                      {"postal_index", "zip"},
                      {"city", "city"},
                      {"address", "address"},
                      {"house_number", "house_number"},
                      {"mailing_postal_index", "zip"},
                      {"mailing_city", "city"},
                      {"mailing_address", "address"},
                      {"mailing_house_number", "house_number"},
                      {"cz_bank_code", "bank_code"},
                      {"password", "111"},
                      {"password_confirmation", "111"},
                      {"origin_ip", "ip_address"},
                     }
  var inner_data_map = map[string]any{}

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      inner_data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      inner_data_map[GetString(f.A)] = f.B
    }
  }
  var translated_bank_account_number string = GetString(pPluginData["bank_account_number"])
  var bank_account_prefix string = "000000"
  var parts = strings.Split(translated_bank_account_number, "-")

  if nil != parts && len(parts) > 1 {
    bank_account_prefix = parts[0]
    translated_bank_account_number = parts[1]
  }
  var idc = GetString(pPluginData["identity_card_number"])
  var ret int

  if len(idc) > 8 {
    ret = ((GetInt(idc[:1]) + GetInt(idc[3:4]) + GetInt(idc[6:7])) * 7 + (GetInt(idc[1:2]) + GetInt(idc[4:5]) + GetInt(idc[7:8])) * 3 +
            GetInt(idc[2:3]) + GetInt(idc[5:6]) + GetInt(idc[8:9])) % 10
  }
  inner_data_map["id_card_number"] = fmt.Sprintf("%v%v", idc, ret)
  inner_data_map["cz_bank_account_prefix"] = bank_account_prefix
  inner_data_map["bank_account"] = translated_bank_account_number

  var sum = GetInt(pPluginData["requested_amount"])
  var period = GetInt(pPluginData["period"])

  inner_data_map["employer"] = "Zaměstnání1" // LEAD-877

  if sum > 20000{
    sum = 20000
  }

  if period > 30{
    period = 30
  }
  var loan_map = map[string]any{}

  loan_map["loan_sum"]= sum
  loan_map["loan_period"] = period

  data_map["client"] = inner_data_map
  data_map["loan"] = loan_map

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if ! GetBool(ret["accepted"]) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = ret["redirect_uri"]
  pPluginData["external_id"] = ret["customer_number"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

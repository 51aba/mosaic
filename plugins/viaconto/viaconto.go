package main

import (
  "fmt"
  "log"
  "encoding/base64"

  . "leadz/utils"
)

type leadplugin string

const codename string = "VIACONTO"

var requested_sum_variants = []int{1000, 1500, 2000, 2500, 3000, 3500, 4000, 4500, 5000, 5500, 6000, 6500, 7000, 7500, 8000, 8500, 9000, 9500,
                                   10000, 10500, 11000, 11500, 12000, 12500, 13000, 13500, 14000, 14500, 15000, 15500, 16000}
var configs_map = map[string]string {
  "api_url": "https://www.viasms.cz/lead/api/v2/",
  "login": "",
  "password": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 21, "param2": 70},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 10000},
  },
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      register_lead,},
  },
  true: map[string][]func(map[string]any, map[string]string) (bool){
    "": []func(map[string]any, map[string]string) (bool){
      check_sms,},
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
  var call_config = map[string]any{"command": "application/create"}
  var basic = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(config["login" + plugin_postfix] + ":" + config["password" + plugin_postfix])))
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": basic}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func check_sms(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "checkPin"}
  var sale_data_map = GetMap(pPluginData["sale_data"])

  if nil == sale_data_map {
    log.Printf("%s %v CHECK_SMS: LEAD_MAP_ISNULL", pPluginData["plugin_log"], codename)

    return
  }
  pPluginData["map_data"] = map[string]any {"customer-id": pPluginData["external_id"], "sms-pin": sale_data_map["sms_code"]}

  return P_check_sms(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  var translated_cell_phone string = "+420" + GetString(pPluginData["cell_phone"])
  var translated_monthly_income int = GetInt(pPluginData["monthly_income"])
  var translated_period int = 30
  var translated_income_type int

  switch GetString(pPluginData["income_type"]) {
    case "EMPLOYED":
      translated_income_type = 6
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = 2
    case "SELF_EMPLOYED":
      translated_income_type = 7
    case "STUDENT":
      translated_income_type = 5
    default:
      translated_income_type = 6
  }
  var fields = []Pair{{"Name", "first_name"},
                      {"Surname", "last_name"},
                      {"PCode", "birth_number"},
                      {"IdCard", "identity_card_number"},
                      {"Email", "email"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["IncomeAmount"] = translated_monthly_income
  data_map["BasicNeeds"] = translated_monthly_income / 4
  data_map["Phone"] = translated_cell_phone
  data_map["Amount"] = FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants)
  data_map["Days"] = translated_period
  data_map["IncomeType"] = translated_income_type
  data_map["PersonsWithIncome"] = 2
  data_map["OtherCreditsTotal"] = 0
  data_map["Dependants"] = 0
  /*
  "AccommodationExpenses": "2500",
  /**/

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var app_map = GetMap(ret["Application"])
  var msg = GetString(app_map["Message"])

  if "Returning customer, accepted" == msg {
    pPluginData["sale_status"] = "DUPLICATE"

    return
  }

  if nil == app_map || nil == app_map["Id"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = app_map["Id"]

  log.Printf("%s RESPONSE_DATA: %v\n", pPluginData["plugin_log"], app_map)

  if nil != app_map["Redirect"] {
    pPluginData["redirect_url"] = app_map["Redirect"]
  }

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

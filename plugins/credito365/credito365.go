package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "CREDITO365"

var configs_map = map[string]string{
  "api_url":  "https://api.credito-365.mx/api/partner-leads",
  "partner_id": "",
  "api_key": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any{}
var send_map = SEND_MAP{
  false: VAR_FUNCS_MAP{
    "": []func(map[string]any, map[string]string) bool{
      register_lead,
    },
  },
}

// ################################################################################################################################################################
func (p leadplugin) TestData(pPluginData map[string]any, is_paused bool) (ret bool) {
  return p.SendData(pPluginData, is_paused)
}

// ################################################################################################################################################################
func (p leadplugin) Validate(pPluginData map[string]any) []map[string]any {
  return P_validate(codename, pPluginData, plugin_vars, validators_map)
}

// ################################################################################################################################################################
func (p leadplugin) SendData(pPluginData map[string]any, is_paused bool) (result bool) {
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{
    "Content-Type": "application/json",
    "X-Parther-id": config["partner_id" + plugin_postfix],
    "X-Partner-Token": config["api_key" + plugin_postfix],
  }

  return P_register_lead(pPluginData, nil, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_income_type = GetString(pPluginData[INCOME_TYPE])
  var translated_marital_status = GetString(pPluginData[MARITAL_STATUS])
  var translated_home_status = GetString(pPluginData[HOME_STATUS])

  switch translated_marital_status {
    case "MARRIED":
      translated_marital_status = "married"
    case "DIVORCED":
      translated_marital_status = "divorced"
    case "SINGLE":
      translated_marital_status = "single"
    case "PARTNERSHIP":
      translated_marital_status = "in_relationship"
    case "OTHER":
      translated_marital_status = "other"
    default:
      translated_marital_status = "married"
  }

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "own"
    case "TENANT":
      translated_home_status = "rent"
    case "OTHER":
      translated_home_status = "family_house"
    default:
      translated_home_status = "own"
  }

  switch translated_income_type {
    case "EMPLOYED":
      translated_income_type = "employed"
    case "PENSION":
      translated_income_type = "pensioner"
    case "SELF_EMPLOYED":
      translated_income_type = "self__employed"
    case "STUDENT":
      translated_income_type = "student"
    case "UNEMPLOYED":
      translated_income_type = "temporarily_unemployed"
    default:
      translated_income_type = "employed"
  }

  var fields = []Pair{
    {"personal_id", IDENTITY_CARD_NUMBER},
    {"requested_amount", REQUESTED_AMOUNT},
    {"requested_tenor", PERIOD},
    {"first_name", FIRST_NAME},
    {"first_last_name", LAST_NAME},
    {"second_last_name", LAST_NAME_2},
    {"gender", GENDER},
    {"birthday", BIRTH_DATE},
    {"phone_number", CELL_PHONE},
    {"email", EMAIL},
    {"residence_postal_code", ZIP},
    {"residence_city", CITY},
    {"residence_region", STATE},
    {"residence_street", STREET},
    {"residence_internal_house", HOUSE_NUMBER},
    {"residence_external_house", HOUSE_NUMBER},
    {"residence_municipality", DISTRICT},
    {"monthly_income", MONTHLY_INCOME},
  }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      var value any = pPluginData[GetString(f.B)]

      if nil == value {
        data_map[GetString(f.A)] = ""
      } else {
        data_map[GetString(f.A)] = value
      }
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["residence_country"] = "MX"
  data_map["nationality"] = "MX"
  data_map["residence_type"] = translated_home_status
  data_map["employment_type"] = translated_income_type
  data_map["marital_status"] = translated_marital_status
  // data_map["rfc"] = ""
  // data_map["subid1"] = "subid1"
  // data_map["subid2"] = "subid2"
  // data_map["education"] = "advanced_level"
  // data_map["residence_colony"] = ""
  // data_map["residence_duration"] = ""

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var status = strings.ToLower(GetString(ret["status"]))

  if "accept" != status {
    log.Printf("%v%v SET_RESPONSE_DATA: STATUS_ERROR: %v%v", RED, pPluginData["plugin_log"], ret, NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  fmt.Printf("%s RESPONSE_DATA: %v [%v]\n", pPluginData["plugin_log"], ret["offer"], status)

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

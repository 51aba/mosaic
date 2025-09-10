package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "SILVERSIDE"

var configs_map = map[string]string {
  "api_url": "https://njf51c4adm.affil.comfortfinancegroup.com/v1/PARTNER_NAME/",
  "api_token": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 21, "param2": 65},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED",}},
    {"field": "birth_number", "func": "InsolvencyValidator",},},
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

  if "" == config["api_token" + plugin_postfix] {
    pPluginData["description"] = fmt.Sprintf("%v REQUEST: TOKEN_ISNULL_ERROR SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["sale_status"])
    pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
    log.Printf("%v%v%v", RED, pPluginData["description"], NC)

    return
  }
  pPluginData["map_data"] = map[string]any{"unique_key": fmt.Sprintf("%v", strings.Replace(GetString(pPluginData["uid"]), "-", "", -1)),
                                     "birth_number": pPluginData["birth_number"],
                                     "email": pPluginData["email"],}
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": config["api_token" + plugin_postfix]}

  return P_check_unique(pPluginData, call_config, set_response_data)
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

                      // {"dependent_children", "dependent_children"},
  var fields = []Pair{
                      {"birth_number", "birth_number"},
                      {"email", "email"},
                      {"cell_phone", "cell_phone"},
                      {"bank_code", "bank_code"},
                      {"city", "city"},
                      {"contact_city", "city"},
                      {"contact_house_number", "house_number"},
                      {"contact_street", "street"},
                      {"contact_zip", "zip"},
                      {"expenses", "monthly_expenses"},
                      {"first_name", "first_name"},
                      {"last_name", "last_name"},
                      {"home_status", "home_status"},
                      {"house_number", "house_number"},
                      {"ip_address", "ip_address"},
                      {"identity_card_number", "identity_card_number"},
                      {"monthly_income", "monthly_income"},
                      {"period", "period"},
                      {"requested_amount", "requested_amount"},
                      {"street", "street"},
                      {"zip", "zip"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  if "M" == GetString(pPluginData["gender"]) {
    data_map["gender"] = "male"
  } else {
    data_map["gender"] = "female"
  }
  data_map["unique_key"] = fmt.Sprintf("%v", strings.Replace(GetString(pPluginData["uid"]), "-", "", -1))
  data_map["income_type"] = strings.ToLower(GetString(pPluginData["income_type"]))
  // data_map["marital_status"] = strings.ToLower(GetString(pPluginData["marital_status"]))

  var translated_bank_account_number string = strings.Replace(strings.Replace(GetString(pPluginData["bank_account_number"]), "000000-", "", -1), "-", "", -1)
  var bank_code string = GetString(pPluginData["bank_code"])
  var bank_account_prefix string = "000000"
  data_map["bank_account"] = fmt.Sprintf("%v-%v/%v", bank_account_prefix, translated_bank_account_number, bank_code)

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if "accepted" != strings.ToLower(GetString(ret["status"])) {
    log.Printf("%v%v SET_RESPONSE_DATA: ERROR: %v%v", RED, pPluginData["plugin_log"], ret["status"], NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

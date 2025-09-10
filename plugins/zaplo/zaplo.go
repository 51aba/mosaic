package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "ZAPLO"

var configs_map = map[string]string {
  "api_url": "https://api.zaplo.cz/affiliate-api",
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
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %v", config["auth_token"])}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  var translated_bn string
  var translated_bank_acc = fmt.Sprintf("%v/%v", pPluginData["bank_account_number"], pPluginData["bank_code"])
  var bn = GetString(pPluginData["birth_number"])

                      // {"numberOfHouseholdMembers", "dependent_children"},
  var fields = []Pair{{"amount", "requested_amount"},
                      {"term", "period"},
                      {"firstName", "first_name"},
                      {"lastName", "last_name"},
                      {"cardId", "identity_card_number"},
                      {"mobilePhone", "cell_phone"},
                      {"email", "email"},
                      {"city", "city"},
                      {"street", "street"},
                      {"houseFlatNumber", "house_number"},
                      {"postalCode", "zip"},
                      {"acceptAgreement", true},
                      {"affiliateProvider", "54_API"},
                      {"affiliateToken", "96080052"},
                      {"newsAccepted", true},
                      {"clientIncome", "monthly_income"},
                      {"clientExpenses", "monthly_expenses"},
                      {"applicationType", "WEB"},
                      {"employmentType", "income_type"},
                      {"numberOfEmployedHouseholdMembers", 2},
                      {"password", "111"},
                      {"passwordRepeat", "111"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }

  if len(bn) > 6 {
    translated_bn = fmt.Sprintf("%v/%v", bn[:6], bn[6:])
  }
  data_map["personalId"] = translated_bn
  data_map["bankAccountNumber"] = translated_bank_acc

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if "success" != strings.ToLower(GetString(ret["resolution"])) {
    return
  }
  var redirect_url string

  if nil != ret["url"] {
    redirect_url = GetString(ret["url"])
  }
  redirect_url = fmt.Sprintf("https://exc4finance.com/?a=&c=&s1=API&s2=%v&ckmrdr=%v", pPluginData["uid"], redirect_url)
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = redirect_url

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

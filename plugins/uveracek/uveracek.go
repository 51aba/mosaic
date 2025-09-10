package main

import (
  "fmt"
  "log"
  "time"

  . "leadz/utils"
)

type leadplugin string

const codename string = "UVERACEK"

var configs_map = map[string]string {
  "api_url": "https://api.jncg.eu/api/v1/",
  "api_key": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "birth_number", "func": "InsolvencyValidator",},},
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
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": fmt.Sprintf("applicant-uveracek/?apiKey=%v", config["api_key" + plugin_postfix])}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["method"] = "PUT"

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var fields = []Pair{
                      {"idCardNumber", "identity_card_number"},
                      {"firstName", "first_name"},
                      {"lastName", "last_name"},
                      {"birthNumber", "birth_number"},
                      {"email", "email"},
                      {"phone", "cell_phone"},
                      {"loanAmount", "requested_amount"},
                      {"loanPeriod", "period"},
                      {"ipAddress", "ip_address"},
                      {"monthlyIncome", "monthly_income"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  var translated_monthly_expenses = GetInt(pPluginData["monthly_expenses"])
  var translated_income_type = GetString(pPluginData["income_type"])

  switch GetString(pPluginData["income_type"]) {
    case "EMPLOYED":
      translated_income_type = "EMPLOYEE"
    case "PENSION":
      translated_income_type = "RETIREE"
    case "SELF_EMPLOYED":
      translated_income_type = "ENTREPRENEUR"
    default:
      translated_income_type = "OTHER"
  }

  if 0 == translated_monthly_expenses {
    translated_monthly_expenses = GetInt(data_map["monthlyIncome"]) / 2
  }

  /*
  var translated_marital_status = GetString(pPluginData["marital_status"])
  if "" == translated_marital_status {
    translated_marital_status = "SINGLE"
  }
  */
  var translated_marital_status string = "SINGLE"
  data_map["maritalStatus"] = translated_marital_status
  data_map["monthlyPayment"] = translated_monthly_expenses
  data_map["incomeType"] = translated_income_type
  data_map["childrenCount"] = 1
  data_map["dateOfAgreement"] = time.Now().Format("2006-01-02 15:04:05")
  data_map["propertyOwner"] = true
  data_map["permanentAddress"] = map[string]any{
                                 "street": pPluginData["street"],
                                 "zip": pPluginData["zip"],
                                 "city": pPluginData["city"],}
  data_map["contactAddress"] = map[string]any{
                                 "street": pPluginData["street"],
                                 "zip": pPluginData["zip"],
                                 "city": pPluginData["city"],}

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if GetInt(ret["code"]) > 299 {
    return
  }
  result = true
  pPluginData["redirect_url"] = ret["redirectUrl"]
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)
type leadplugin string

const codename string = "FLEXIFIN"

var configs_map = map[string]string{
  "api_url": "https://lead.flexifin.cz/api/v1/leads/",
  "auth": "",
  "code": "",
  "partner": "",
  "redirect_url": "https://flexifin.cz/?utm_source=PARTNER_NAME&utm_parametr1=hodnota1&utm_parametr=hodnota2#PARTNER_NAME/",
}
var plugin_vars = []string{"cpl2", "cpl", "link", ""}
var validators_map = map[string][]map[string]any{
  "cpl2": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
  },
  "cpl": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
  },
  "link": {
    {},},
  "": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
  },
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpl2": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "cpl": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "link": []func(map[string]any, map[string]string) (bool){
      link_action,
    },
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
func link_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = fmt.Sprintf("%v%v", config["redirect_url" + plugin_postfix], pPluginData["uid"])

  pPluginData["description"] = fmt.Sprintf("%v REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])
  result = true

  return
}

/*
// ################################################################################################################################################################
func check_unique_cpl(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["partner"] = config["partner" + plugin_postfix]
  data_map["token"] = config["auth" + plugin_postfix]
  data_map["code"] = config["code" + plugin_postfix]

  pPluginData["map_data"] = data_map

  result = true

  return
}
*/

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "check"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["token"] = config["auth" + plugin_postfix]
  data_map["code"] = config["code" + plugin_postfix]

  if "" != GetString(config["partner" + plugin_postfix]) {
    data_map["partner"] = config["partner" + plugin_postfix]
  }
  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json",}

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "signup"}

  if _, dok := pPluginData["map_data"]; ! dok {
    var plugin_postfix = GetString(pPluginData["plugin_postfix"])
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    data_map["token"] = config["auth" + plugin_postfix]
    data_map["code"] = config["code" + plugin_postfix]

    if "" != GetString(config["partner" + plugin_postfix]) {
      data_map["partner"] = config["partner" + plugin_postfix]
    }
    pPluginData["map_data"] = data_map
  }
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json",}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var translated_home_status string = GetString(pPluginData["home_status"])

  if "" == translated_home_status {
    translated_home_status = "TENANT"
  }
  data_map["ipAddress"] = GetString(pPluginData["ip_address"])
  data_map["requestedAmount"] = GetInt(pPluginData["requested_amount"])
  data_map["period"] = GetInt(pPluginData["period"])
  data_map["sameAddress"] = true

  data_map["personalInfo"] = map[string]any{
    "firstName": GetString(pPluginData["first_name"]),
    "lastName": GetString(pPluginData["last_name"]),
    "birthNumber": GetString(pPluginData["birth_number"]),
    "identityCardNumber": GetString(pPluginData["identity_card_number"]),
    /*
    "dependentChildren": GetString(pPluginData["dependent_children"]),
    "maritalStatus": GetString(pPluginData["marital_status"]),
    */
    "email": GetString(pPluginData["email"]),
    "cellPhone": "+420" + GetString(pPluginData["cell_phone"]),
    "incomeType": GetString(pPluginData["income_type"]),}

  data_map["residenceAddress"] = map[string]any{
    "zip": GetString(pPluginData["zip"]),
    "city": GetString(pPluginData["city"]),
    "street": GetString(pPluginData["street"]),
    "houseNumber": GetString(pPluginData["house_number"]),
    "homeStatus": translated_home_status,}

  /*
  data_map["sourceDetail"] = ""
  data_map["product"] = ""
  data_map["contactAddress"] = map[string]any{
    "zip": GetString(pPluginData["zip"]),
    "city": GetString(pPluginData["city"]),
    "street": GetString(pPluginData["street"]),
    "houseNumber": GetString(pPluginData["house_number"]),
    "homeStatus": GetString(pPluginData["home_status"])}},}
  */

  /*
  // if "EMPLOYED" == pPluginData["income_type"] {
  //   data_map["employerData"] = map[string]any{
  //     "employer": GetString(pPluginData["employer"]),
  //     "employerAddress" : map[string]any{
  //       "city": GetString(pPluginData["employer_city"]),
  //       "street": GetString(pPluginData["employer_address"]),
  //     },
  //     "employerPhone": "+420" + GetString(pPluginData["employer_phone"]),
  //     "jobTitle": GetString(pPluginData["job_title"]),
  //     "employedTime": 0,
  //     "employerDescription": "",
  //     "companyNumber": "company_number",}
  // } else {
  // }
  */
  var translated_bank_account_number string = GetString(pPluginData["bank_account_number"])
  var bank_account_prefix string = "000000"
  var parts = strings.Split(translated_bank_account_number, "-")

  if nil != parts && len(parts) > 1 {
    bank_account_prefix = parts[0]
    translated_bank_account_number = parts[1]
  }

  data_map["bankData"] = map[string]any{
    "bankAccountPrefix": bank_account_prefix,
    "bank_account_number": translated_bank_account_number,
    "bank_code": GetString(pPluginData["job_title"]),
  }

  data_map["incomeData"] = map[string]any{
    "monthlyIncome": GetInt(pPluginData["monthly_income"]),
    "monthlyExpenses": GetInt(pPluginData["monthly_expenses"]),
  }

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var data_map = GetMap(ret["data"])

  if nil == data_map {
    log.Printf("%v SET_RESPONSE_DATA: NO_DATA_RET: %v", pPluginData["plugin_log"], ret)

    return
  }
  var response_status string = strings.ToLower(GetString(data_map["status"]))

  if "accepted" != response_status {
    return
  }
  log.Printf("%v SET_RESPONSE_DATA: STATUS: %v [ %v ]", pPluginData["plugin_log"], code, response_status)

  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = data_map["url"]
  pPluginData["external_id"] = data_map["refNo"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

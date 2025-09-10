package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "CREDITON"

var configs_map = map[string]string{
  "api_url": "https://proxy.ext.ferratum.com/czech/cls/",
  "login": "",
  "password": "",
}
var plugin_vars = []string{"mainhigh", "maintime", "new", "duplicatecps", "plus", "link"}
var validators_map = map[string][]map[string]any {
  "mainhigh": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},},
  "maintime": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},},
  "new": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},},
  "duplicatecps": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},},
  "plus": {
    {},},
  "link": {
    {},},
  "": {
    {},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "mainhigh": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "maintime": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "new": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "duplicatecps": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "plus": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "link": []func(map[string]any, map[string]string) (bool){
      link_action,},
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
func link_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  if "" != config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
  } else {
    pPluginData["redirect_url"] = fmt.Sprintf("https://www.crediton.cz/?utm_source=PARTNER_NAME&utm_medium=affiliate&utm_campaign=PARTNER_NAME")
  }
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])

  return true
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "checkDuplicate"}

  pPluginData["map_data"] = map[string]any{"login": config["login" + plugin_postfix], "password": config["password" + plugin_postfix],
                                     "data": map[string]any {"new-registration": true, "customer-personcode": pPluginData["birth_number"]},}
  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "registerCustomer"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["ref"] = config["login" + plugin_postfix]

  pPluginData["map_data"] = map[string]any{"login": config["login" + plugin_postfix], "password": config["password" + plugin_postfix], "data": data_map}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func check_sms(pPluginData map[string]any, config map[string]string) (result bool) {
  var sale_data_map = GetMap(pPluginData["sale_data"])

  if nil == sale_data_map {
    log.Printf("%s SEND_DATA: [%v] CHECK_SMS: SALE_STATUS: %v SALE_MAP_ISNULL_ERROR", pPluginData["plugin_log"], pPluginData["function_idx"], pPluginData["sale_status"])

    return
  }
  var call_config = map[string]any{"command": "checkPin"}

  pPluginData["map_data"] = map[string]any {"customer-id": pPluginData["external_id"], "sms-pin": sale_data_map["sms_code"]}

  return P_check_sms(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_gender, translated_income_type, translated_income_level, translated_working_time int
  var translated_requested_amount int = GetInt(pPluginData["requested_amount"])
  var monthly_income int = GetInt(pPluginData["monthly_income"])
  var translated_home_status string = GetString(pPluginData["home_status"])
  var translated_address_time string = ""
  var translated_bank_account string = fmt.Sprintf("%v/%v", pPluginData["bank_account_number"], pPluginData["bank_code"])
  var translated_street string = fmt.Sprintf("%v %v", pPluginData["street"], pPluginData["house_number"])
  var income_type string = GetString(pPluginData["income_type"])

  if "F" == GetString(pPluginData["gender"]) {
    translated_gender = 0
  } else {
    translated_gender = 1
  }

  if translated_requested_amount < 1000 { // map[err:map[loan-amount:[invalid-value error-loan-limit] loan-limit:[error-loan-limit] min-loan-amount:1000]]
    translated_requested_amount = 1000
  } else if translated_requested_amount > 40000 { // *** LEAD-890 ***
    translated_requested_amount = 40000
  }

  if monthly_income <= 5000 {
    translated_income_level = 0
  } else if 5000 < monthly_income && monthly_income <= 8000 {
    translated_income_level = 1
  } else if 8000 < monthly_income && monthly_income <= 10000 {
    translated_income_level = 2
  } else if 10000 < monthly_income && monthly_income <= 13000 {
    translated_income_level = 3
  } else if 13000 < monthly_income && monthly_income <= 16000 {
    translated_income_level = 4
  } else if 16000 < monthly_income && monthly_income <= 19000 {
    translated_income_level = 5
  } else if 19000 < monthly_income && monthly_income <= 22000 {
    translated_income_level = 6
  } else if 22000 < monthly_income && monthly_income <= 25000 {
    translated_income_level = 7
  } else if 25000 < monthly_income && monthly_income <= 30000 {
    translated_income_level = 8
  } else if 30000 < monthly_income && monthly_income <= 35000 {
    translated_income_level = 9
  } else if 35000 < monthly_income && monthly_income <= 40000 {
    translated_income_level = 10
  } else {
    translated_income_level = 11
  }

  switch income_type {
    case "EMPLOYED":
      translated_income_type = 1
      translated_working_time = 1
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = 1
      translated_working_time = 2
    case "SELF_EMPLOYED":
      translated_income_type = 2
      translated_working_time = 1
    case "MATERNITY_LEAVE":
      translated_income_type = 3
      translated_working_time = 5
    case "STUDENT":
      translated_income_type = 5
      translated_working_time = 5
    case "PENSION":
      translated_income_type = 4
      translated_working_time = 5
    case "SAVINGS":
      translated_income_type = 7
      translated_working_time = 5
    case "UNEMPLOYED":
      translated_income_type = 8
      translated_working_time = 5
    case "BENEFITS":
      translated_income_type = 8
      translated_working_time = 4
    case "OTHER":
      translated_income_type = 8
      translated_working_time = 5
    default:
      translated_income_type = 8
      translated_working_time = 5
  }

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "1"
    case "CO_OWNED":
      translated_home_status = "1"
    case "HOSTEL":
      translated_home_status = "4"
    default:
      translated_home_status = "3"
  }
  var fields = []Pair {
                      {"customer-firstname", "first_name"},
                      {"customer-lastname", "last_name"},
                      {"customer-personcode", "birth_number"},
                      {"customer-birthday", "birth_date"},
                      {"customer-email", "email"},
                      {"customer-password", "user_password"},
                      {"customer-phone", "cell_phone"},
                      {"customer-city", "city"},
                      {"customer-zip", "zip"},
                      {"document-number", "identity_card_number"},
                      {"total-expenses", "monthly_expenses"},
                      {"ipaddress", "ip_address"},
                      }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }

  if "SELF_EMPLOYED" == income_type {
    if _, ok := pPluginData["company_number"]; ! ok {
      data_map["job-ico"] = 24312134
    } else {
      data_map["job-ico"] = pPluginData["company_number"]
    }
  }
  var job_title string = GetString(pPluginData["job_title"])

  if "" == job_title {
    job_title = "XXX"
  }
  data_map["job-name"] = job_title
  data_map["customer-monthly-income"] = monthly_income
  data_map["use-products"] = 1
  data_map["new-registration"] = 1
  data_map["product"] = 3
  data_map["loan-term"] = 12
  data_map["language"] = "cz"
  data_map["bank-account"] = translated_bank_account
  data_map["customer-gender"] = translated_gender
  data_map["customer-street"] = translated_street
  data_map["income-level-selection"] = translated_income_level
  data_map["type-of-living"] = translated_home_status
  data_map["living-at-address"] = translated_address_time
  data_map["data-usage"] = "1"
  data_map["loan-amount"] = translated_requested_amount
  data_map["customer-employment-type"] = translated_income_type
  data_map["customer-working-time"] = translated_working_time
  data_map["customer-employment-length"] = "1"
  data_map["agb"] = 1
  data_map["newsletter"] = 1
  data_map["other"] = "PARTNER_NAME"
  data_map["telco-consent"] = "1"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil != ret["err"] {
    if strings.Contains(GetString(ret["err"]), "duplicate") {
      pPluginData["sale_status"] = "DUPLICATE"
    }
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  var ret_data_map = GetMap(ret["data"])

  log.Printf("%v SET_RESPONSE_DATA: DATA_MAP: %v [%v]", pPluginData["plugin_log"], ret["data"], ret_data_map)

  if nil != ret_data_map {
    pPluginData["redirect_url"] = ret_data_map["forward-url"]
    pPluginData["external_id"] = ret_data_map["customer-id"]
  }
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var config map[string]string = GetMapStrings(pPluginData["config"])

  if "" != GetString(config["form_context" + plugin_postfix]) {
    pPluginData["form_context"] = config["form_context" + plugin_postfix]
  } else {
    pPluginData["form_context"] = `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`
  }

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

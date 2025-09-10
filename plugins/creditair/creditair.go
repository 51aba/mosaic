package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)
type leadplugin string

const codename string = "CREDITAIR"

var configs_map = map[string]string{
  "api_url": "https://cz.platon.avafin.com/affiliate/",
  "redirect_url": "https://www.creditair.cz/?utm_source=pap&utm_medium=affiliate&utm_campaign=cps&utm_content=970x310&ref=57715a15a6abf&bid=e492f0b6",
  "form_context": `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`,
  "login": "",
  "password": "",
}
var plugin_vars = []string{"new", "duplicatecps", "link"}
var validators_map = map[string][]map[string]any{
  "new": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},
  },
  "duplicatecps": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},
  },
  "link": {
    {},
  },
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 79},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},
  },
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "new": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "duplicatecps": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "link": []func(map[string]any, map[string]string) (bool){
      link_action,},
  },
  true: map[string][]func(map[string]any, map[string]string) (bool){
    "new": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "duplicatecps": []func(map[string]any, map[string]string) (bool){
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
func link_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  pPluginData["sale_status"] = "UNCONFIRMED"

  if "" != config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
  } else {
    pPluginData["redirect_url"] = fmt.Sprintf("https://www.creditair.cz/?utm_source=pap&utm_medium=affiliate&utm_campaign=cps&utm_content=970x310&ref=57715a15a6abf&bid=e492f0b6")
  }
  pPluginData["description"] = fmt.Sprintf("%v REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])

  return true
}

/*
// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "checkDuplicate"}

  pPluginData["map_data"] = map[string]any{"login": config["login" + plugin_postfix], "password": config["password" + plugin_postfix],
                                     "data": map[string]any {"new-registration": true, "customer-personcode": pPluginData["birth_number"]},}
  return P_check_unique(pPluginData, call_config, set_response_data)
}
*/

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "registerCustomer"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
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
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "checkPin"}

  pPluginData["map_data"] = map[string]any{"login": config["login" + plugin_postfix], "password": config["password" + plugin_postfix],
                                     "data": map[string]any {"customer-id": pPluginData["external_id"], "sms-pin": sale_data_map["sms_code"]},}
  return P_check_sms(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var translated_gender, translated_income_type, translated_income_level, translated_working_time int
  var translated_requested_amount = GetInt(pPluginData["requested_amount"])
  var translated_home_status string = GetString(pPluginData["home_status"])
  var translated_address_time string
  var translated_bank_account string = fmt.Sprintf("%v/%v", pPluginData["bank_account_number"], pPluginData["bank_code"])
  var monthly_income = GetInt(pPluginData["monthly_income"])
  var income_type = GetString(pPluginData["income_type"])

  if "F" == GetString(pPluginData["gender"]) {
    translated_gender = 0
  } else {
    translated_gender = 1
  }

  if translated_requested_amount < 1000 { // *** map[err:map[0:service-error loan-amount:[invalid-value error-loan-limit] min-loan-amount:1000]] ***
    translated_requested_amount = 1000
  } else if translated_requested_amount > 20000 { // *** LEAD-890 ***
    translated_requested_amount = 20000
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
      translated_working_time = 1
    case "SELF_EMPLOYED":
      translated_working_time = 1
    case "PART_TIME_EMPLOYMENT":
      translated_working_time = 2
    case "OTHER":
      translated_working_time = 5
    case "BENEFITS":
      translated_working_time = 4
    default:
      translated_working_time = 5
  }

  switch income_type {
    case "EMPLOYED":
      translated_income_type = 1
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = 1
    case "SELF_EMPLOYED":
      translated_income_type = 2
    case "MATERNITY_LEAVE":
      translated_income_type = 3
    case "STUDENT":
      translated_income_type = 5
    case "PENSION":
      translated_income_type = 4
    case "SAVINGS":
      translated_income_type = 7
    case "UNEMPLOYED":
      translated_income_type = 8
    default:
      translated_income_type = 8
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

                      /*
                      // {"job-name", "employer"},
                      // {"job-phone", "employer_phone"},
                      // {"job-address", "employer_address"},
                      // {"job-position", "job_title"},
                      */
  var fields = []Pair{
                      {"customer-firstname", "first_name"},
                      {"customer-lastname", "last_name"},
                      {"customer-personcode", "birth_number"},
                      {"customer-birthday", "birth_date"},
                      {"customer-email", "email"},
                      {"customer-password", "user_password"},
                      {"customer-phone", "cell_phone"},
                      {"customer-city", "city"},
                      {"customer-street", "street"},
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
  /*
  } else if "EMPLOYED" == income_type {
    data_map["job-name"] = "Zaměstnání1"
  */
  }
  var job_title string = GetString(pPluginData["job_title"])

  if "" == job_title {
    job_title = "Zaměstnání1"
  }
  data_map["job-name"] = job_title
  data_map["customer-monthly-income"] = monthly_income
  data_map["income-level-selection"] = translated_income_level
  data_map["use-products"] = 1
  data_map["new-registration"] = 1
  data_map["product"] = 3
  data_map["loan-term"] = 12
  data_map["bank-account"] = translated_bank_account
  data_map["customer-gender"] = translated_gender
  data_map["type-of-living"] = translated_home_status
  data_map["living-at-address"] = translated_address_time
  data_map["data-usage"] = "1"
  data_map["language"] = "cz"
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
    log.Printf("%v%v SET_RESPONSE_DATA: ERROR: %v%v", RED, pPluginData["plugin_log"], ret, NC)

    if strings.Contains(GetString(ret["err"]), "duplicate") {
      pPluginData["sale_status"] = "DUPLICATE"
    }
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  var ret_data_map = GetMap(ret["data"])

  log.Printf("%s SET_RESPONSE_DATA: RET_DATA_TEST: %v [%v]", pPluginData["plugin_log"], ret["data"], ret_data_map)

  if nil != ret_data_map && "registerCustomer" == command {
    // pPluginData["sale_status"] = "PAUSED"
    var plugin_postfix = GetString(pPluginData["plugin_postfix"])
    var config map[string]string = GetMapStrings(pPluginData["config"])

    if "" != GetString(config["form_context" + plugin_postfix]) {
      pPluginData["form_context"] = config["form_context" + plugin_postfix]
    } else {
      pPluginData["form_context"] = `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`
    }
    pPluginData["redirect_url"] = ret_data_map["forward-url"]
    pPluginData["external_id"] = ret_data_map["customer-id"]
    log.Printf("%s SET_RESPONSE_DATA: %v", pPluginData["plugin_log"], ret_data_map)
  }

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

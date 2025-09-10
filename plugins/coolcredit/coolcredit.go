package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)
type leadplugin string

const codename string = "COOLCREDIT"

var configs_map = map[string]string{
  "api_url": "https://coolcredit.cz/api/v2/",
  "form_context": `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "validator": "number", "label": "SMS Code"}]}`,
  "token": "",
  "redirect_url": "https://coolcredit.cz/?utm_source=PARTNER_NAME&utm_medium=cpc&utm_campaign=klikacka&utm_content=cerven&aid=138",
}
var plugin_vars = []string{"newhq", "newlq", "newnormal", "cps", "big", "redirect", "sms", ""}
var validators_map = map[string][]map[string]any{
  "newhq": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 74},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 10000},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED", "EMPLOYED", "PENSION", "PART_TIME_EMPLOYMENT",}},},
  "newlq": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 74},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 8000},},
  "newnormal": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 18, "param2": 74},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 10000},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED", "EMPLOYED", "PENSION", "PART_TIME_EMPLOYMENT",}},},
  "cps": {
    {},},
  "big": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 20, "param2": 75},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 20000},
    {"field": "home_status", "func": "AllowedValuesValidator", "param1": []string{"TENANT", "HOME_OWNER",}},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED",}},
    {"field": "insolvency", "func": "AllowedValuesValidator", "param1": []string{"NO",}},},
  "redirect": {
    {},},
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cps": []func(map[string]any, map[string]string) (bool){
      cps_check_unique, cps_register_lead,},
    "newhq": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "newlq": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "newnormal": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "big": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "redirect": []func(map[string]any, map[string]string) (bool){
      redirect_action,},
    "sms": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "": []func(map[string]any, map[string]string) (bool){
      register_lead,},
  },
  true: map[string][]func(map[string]any, map[string]string) (bool){
    "sms": []func(map[string]any, map[string]string) (bool){
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
func redirect_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]

  pPluginData["description"] = fmt.Sprintf("%v REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])

  return true
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": fmt.Sprintf("registrations/check?token=%v", config["token" + plugin_postfix])}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  pPluginData["map_data"] = data_map

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func cps_check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": fmt.Sprintf("pre-registrations?token=%v", config["token" + plugin_postfix])}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  pPluginData["map_data"] = data_map

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": fmt.Sprintf("registrations?token=%v", config["token" + plugin_postfix])}

  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    pPluginData["map_data"] = data_map
  }
  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func cps_register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": fmt.Sprintf("pre-registrations?token=%v", config["token" + plugin_postfix]), "description": "CPS"}

  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    pPluginData["map_data"] = data_map
  }
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
  var call_config = map[string]any{"command": fmt.Sprintf("registrations/activation?token=%v", config["token" + plugin_postfix])}

  pPluginData["map_data"] = map[string]any {"pincode": sale_data_map["sms_code"]}

  return P_check_sms(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
                      /*
                      // {"employer", "employer"},
                      // {"employerPosition", "job_title"},
                      // {"employerPhoneNumber", "work_phone"},
                      */
  var fields = []Pair{{"nin", "birth_number"},
                      {"name", "first_name"},
                      {"surname", "last_name"},
                      {"idCardNumber", "identity_card_number"},
                      {"phoneNumber", "cell_phone"},
                      {"email", "email"},
                      {"street", "street"},
                      {"streetNumber", "house_number"},
                      {"city", "city"},
                      {"zipCode", "zip"},
                      {"netIncome", "monthly_income"},
                      {"accountNumber", "bank_account_number"},
                      {"costs", "monthly_expenses"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  var translated_requested_amount = GetInt(pPluginData["requested_amount"])

  if translated_requested_amount > 30000 {
    translated_requested_amount = 30000
  } else if translated_requested_amount > 10000 {
    translated_requested_amount = 10000
  } else if translated_requested_amount < 1000 {
    translated_requested_amount = 1000
  }
  data_map["amount"] = translated_requested_amount
  data_map["workingYears"] = "1to2years"
  data_map["maritalStatus"] = "free"
  data_map["housing"] = "owner"
  data_map["moneySource"] = 1
  data_map["bank"] = 1
  data_map["days"] = 30

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  if "_cps" == plugin_postfix {
    var data_map = GetMap(ret["data"])
    var id_old = ret["id"]
    var id =  data_map["id"]

    if nil == id_old && nil == id {
      log.Printf("%v%v SET_RESPONSE_DATA: CPS_ID_ISNULL_ERROR: [ %v : %v ] %v%v", RED, pPluginData["plugin_log"], id_old, id, ret, NC)

      return
    }
    pPluginData["external_id"] = ret["id"]
    result = true
    pPluginData["sale_status"] = "UNCONFIRMED"
  } else {
    if nil == ret["data"] {
      log.Printf("%v%v SET_RESPONSE_DATA: DATA_ISNULL_ERROR: [ %v ]%v", RED, pPluginData["plugin_log"], ret["data"], NC)

      return
    }
    var data_map = GetMap(ret["data"])
    log.Printf("%s SET_RESPONSE_DATA: %v [%v]", pPluginData["plugin_log"], ret["data"], data_map)

    if nil == data_map {
      log.Printf("%v%v SET_RESPONSE_DATA: DATA_MAP_ISNULL_ERROR: [ %v ]%v", RED, pPluginData["plugin_log"], ret["data"], NC)

      return
    }

    if "accepted" != strings.ToLower(GetString(data_map["result"])) { // *** NEW_VERSION v2 ***
      log.Printf("%v%v SET_RESPONSE_DATA: RESULT_ISNULL_WARNING: [ %v ]%v", YELLOW, pPluginData["plugin_log"], data_map["result"], NC)

      return
    }

    if nil != data_map["redirect_url"] {
      pPluginData["redirect_url"] = data_map["redirect_url"]
    } else if nil != data_map["redirectUrl"] {
      pPluginData["redirect_url"] = data_map["redirectUrl"]
    }
    result = true

    if strings.Contains(command, "activation") { // *** sms_check ***
      pPluginData["sale_status"] = "UNCONFIRMED"

      return
    }

    if StringInArray([]string{"_cps", "_newlq", "_newhq", "_newnormal", "_big"}, plugin_postfix) {
      pPluginData["sale_status"] = "UNCONFIRMED"
    } else { // *** _api ***
      pPluginData["sale_status"] = "PAUSED"
    }
    var config map[string]string = GetMapStrings(pPluginData["config"])

    if "" != GetString(config["form_context" + plugin_postfix]) {
      pPluginData["form_context"] = config["form_context" + plugin_postfix]
    } else {
      pPluginData["form_context"] = `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "validator": "number", "label": "SMS Code"}]}`
    }
  } // *** NON CPS ***

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

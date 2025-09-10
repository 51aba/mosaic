package main

import (
  "fmt"
  "log"
  "strings"
  "encoding/base64"

  . "leadz/utils"
)

type leadplugin string

const codename string = "MOBILPUJCKA"

var configs_map = map[string]string {
  "api_url": "https://sys.mobilpujcka.cz/api-loan/",
  "login": "",
  "password": "",
  "companyID": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
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
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "check-customer"}
  var basic = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(config["login" + plugin_postfix] + ":" + config["password" + plugin_postfix])))
  var data_map = map[string]any{}

  data_map["companyId"] = config["companyID" + plugin_postfix]

  translate(pPluginData, data_map)
  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": basic}

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "loan-request"}

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
  var call_config = map[string]any{"command": "checkPin"}

  return P_check_sms(pPluginData, call_config, set_response_data)

  return
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
                      /*
                      // {"employer", "employer"},
                      // {"employerPosition", "job_title"},
                      // {"position", "job_title"},
                      */
  var fields = []Pair{{"personalIdentificationNumber", "birth_number"},
                      {"identityCardNumber", "identity_card_number"},
                      {"firstname", "first_name"},
                      {"lastname", "last_name"},
                      {"email", "email"},
                      {"mobil", "cell_phone"},
                      {"loanAmount", 1000},
                      {"accountNumber", "bank_account_number"},
                      {"bankCode", 1},
                      {"city", "city"},
                      {"street", "street"},
                      {"descriptiveNumber", "10"},
                      {"orientationNumber", ""},
                      {"postalCode", "zip"},
                     }
  var inner_data_map = make(map[string]any)

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      inner_data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      inner_data_map[GetString(f.A)] = f.B
    }
  }
  data_map["data"] = inner_data_map

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if code > 299 {
    if strings.Contains(fmt.Sprintf("%v", ret), "duplicate") {
      pPluginData["sale_status"] = "DUPLICATE"
    }
    return
  }

  if nil == ret["status"] || "ok" != ret["status"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = ret["customerId"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

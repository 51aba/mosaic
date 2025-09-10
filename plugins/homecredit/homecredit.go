package main

import (
  "fmt"
  "log"

  . "leadz/utils"
)

type leadplugin string

// *** README FIXME server ip whitelist required *** OR {"reason": "Blocked by WAF", "response_type": "General", "supportID": 17280053378626552529}
const codename string = "HOMECREDIT"

var configs_map = map[string]string {
  "api_url": "https://api.homecredit.cz/",
  "auth_token": "",
  "username": "",
  "password": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      get_token, check_unique, register_lead,},
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
func get_token(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "authentication/v1/partner", "description": "GET_TOKEN"}

  pPluginData["map_data"] = map[string]any{"username": config["username" + plugin_postfix], "password": config["password" + plugin_postfix]}

  return P_check_unique(pPluginData, call_config, set_response_data_token)
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  if "" == config["auth_token"] {
    return
  }
  var call_config = map[string]any{"command": "financing/v2/scorings/precheck"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %v", config["auth_token"])}

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "financing/v2/applications"}

  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    pPluginData["map_data"] = data_map
  }

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  var translated_requested_amount = GetInt(pPluginData["requested_amount"])

  if translated_requested_amount < 10000 {
    translated_requested_amount = 1000000
  } else if translated_requested_amount > 250000 {
    translated_requested_amount = 25000000
  } else {
    translated_requested_amount *= 100
  }
  data_map["customer"] = map[string]any {
                                 "firstName": pPluginData["first_name"],
                                 "lastName": pPluginData["last_name"],
                                 "email": pPluginData["email"],
                                 "phone": "+420" + GetString(pPluginData["cell_phone"]),
                                 "identificationNumber": pPluginData["birth_number"],
                                 "addresses": []map[string]any{
                                     map[string]any{
                                     "streetAddress": pPluginData["street"],
                                     "streetNumber": pPluginData["house_number"],
                                     "zip": pPluginData["zip"],
                                     "city": pPluginData["city"],
                                     "flags": []string{ "PERMANENT"},},},}
  data_map["financedSubject"] = map[string]any{"type": "CASH",
                                 "cash": map[string]any{
                                         "amount": map[string]any{
                                                   "value": translated_requested_amount,
                                                   "currency": "CZK",},
                                         "financing": map[string]any{
                                                      "type": "LOAN_WALLET",
                                                      "loanWallet": map[string]any{
                                                                    "preferredInstallment": map[string]any{
                                                                                            "value": translated_requested_amount,
                                                                                            "currency": "CZK",},},},},}
  data_map["originationContext"] = map[string]any{
     "channel": map[string]any{
                "type": "ONLINE",
                "online": map[string]any{
                          "approvedRedirect": "https://www.creditsor.cz/done",
                          "rejectedRedirect": "https://www.creditsor.cz/",
                          "notificationEndpoint": fmt.Sprintf("https://www.PARTNER_NAME.cz/api/leads/tracking/hc/%v/", pPluginData["uid"]),},},
      "sourceId": "PARTNER_NAME",
      "extendedData": []map[string]any{map[string]any{
                           "key": "utm_source",
                           "value": "source"},
                           map[string]any{
                           "key": "utm_campaign",
                           "value": "PARTNER_NAME",},},}
  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data_token(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  log.Printf("%v SET_RESPONSE_DATA_TOKEN: %v", pPluginData["plugin_log"], ret["accessToken"])
  var config = GetMapStrings(pPluginData["config"])

  config["auth_token"] = GetString(ret["accessToken"])
  pPluginData["config"] = config

  if "" == config["auth_token"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], ret == nil, code > 299, NC)

    return
  }
  log.Printf("%v SET_RESPONSE_DATA: RET_RESULT: %v RET_STATUS: %v", pPluginData["plugin_log"], ret["result"], ret["status"])

  if nil == ret["result"] && nil == ret["status"] {
    return
  }

  if "NEW_CONTRACT" == GetString(ret["result"]) || "CREATED" == GetString(ret["status"]) {
  } else {
    return
  }
  var steps = GetArray(ret["workflowSteps"])
  log.Printf("%v SET_RESPONSE_DATA: STEPS: %v", pPluginData["plugin_log"], steps)

  if nil == steps || 0 == len(steps) {
  } else {
    var redirects_map = GetMap(GetMap(steps[0])["redirect"])

    log.Printf("%v SET_RESPONSE_DATA: REDIRECT: %v", pPluginData["plugin_log"], redirects_map)

    if nil != redirects_map {
      pPluginData["redirect_url"] = redirects_map["gatewayRedirectUrl"]
    }
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "fmt"
  "log"
  "encoding/base64"

  . "leadz/utils"
)
type leadplugin string

const codename string = "CREDITEA"

var configs_map = map[string]string{
  "api_url": "https://production-esb.ipfdigital.io/api/partners/v1/",
  "correlation": "",
  "partner_id": "",
  "password": "",
}
var plugin_vars = []string{"cpl"}
var validators_map = map[string][]map[string]any{
  "cpl": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpl": []func(map[string]any, map[string]string) (bool){
      register_lead,
    },
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
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "applications"}
  var basic = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", config["correlation" + plugin_postfix], config["password" + plugin_postfix])))

  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Basic %v", basic),
                                       "Correlation-ID": config["correlation" + plugin_postfix], "CountryCode": "CZ", "Partner-ID": config["partner_id" + plugin_postfix],}
  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    pPluginData["map_data"] = data_map
  }
  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var app_map = []map[string]any {{ "role": "MAIN",
                                    "firstName": GetString(pPluginData["first_name"]),
                                    "lastName": GetString(pPluginData["last_name"]),
                                    "ip": GetString(pPluginData["ip_address"]),
                                    "contact" : map[string]any{"email": GetString(pPluginData["email"]), "phone": GetString(pPluginData["cell_phone"])},
                                    "identity": map[string]any{
                                                "personalIdentificationNumber": GetString(pPluginData["birth_number"]),
                                                "documentNumber": GetString(pPluginData["identity_card_number"])},
                                    "addresses": []map[string]any{{
                                                 "purpose": []string{"RESIDENCE"},
                                                 "city": GetString(pPluginData["city"]),
                                                 "postalCode": GetString(pPluginData["zip"]),
                                                 "street": GetString(pPluginData["street"]),
                                                 "houseNumber": GetString(pPluginData["house_number"])}},
                                    "financials": map[string]any{
                                                  "incomes": []map[string]any{{
                                                             "type": "SALARY",
                                                             "netAmount": GetFloat(pPluginData["monthly_income"])}},
                                                             "consents": map[string]bool{"personalDataProcessing": true,
                                                             "debtRegistryCheck": true}},
                                 },}
   data_map["application"] = map[string]any{
                             "applicants": app_map,
                             "loan": map[string]any{
                                     "amount": GetFloat(pPluginData["requested_amount"])},
                             }
  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil == ret["result"] {
    return
  }
  var result_map = GetMap(ret["result"])

  fmt.Printf("%s RESPONSE_DATA: %v [%v]\n", pPluginData["plugin_log"], ret["result"], result_map)

  if nil == result_map || ! GetBool(result_map["success"]) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

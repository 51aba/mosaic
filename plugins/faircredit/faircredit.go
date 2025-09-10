package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "FAIRCREDIT"

var configs_map = map[string]string {
  "api_url": "https://fis.faircredit.cz/api/ContactCentrum/",
  "token": "",
  "redirect_url": "https://www.faircredit.cz/pujcka-v-hotovosti/?utm_source=PARTNER_NAME&utm_medium=cpc&utm_campaign=PARTNER_NAME_s",
}
var plugin_vars = []string{"link", ""}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "birth_date", "func": "MinAgeValidator", "param1": 18},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED", "STUDENT",}},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "link": []func(map[string]any, map[string]string) (bool){
      link_action,
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
  // pPluginData["redirect_url"] = fmt.Sprintf("%v%v", config["redirect_url" + plugin_postfix], pPluginData["uid"])
  pPluginData["redirect_url"] = fmt.Sprintf("%v", config["redirect_url" + plugin_postfix])

  pPluginData["description"] = fmt.Sprintf("%v REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])
  result = true

  return
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "AddContacts"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var fields = []Pair{{"ExternalIdentifier", "uid"},
                      {"FirstName", "first_name"},
                      {"LastName", "last_name"},
                      {"SocialSecurityNumber", "birth_number"},
                      {"IdCardNumber", "identity_card_number"},
                      {"Street", "street"},
                      {"HouseNumber", "house_number"},
                      {"City", "city"},
                      {"Zip", "zip"},
                      {"Email", "email"},
                     }

  for _, f := range fields {
    var key = GetString(f.A)

    if nil != pPluginData[GetString(f.B)] {
      data_map[key] = pPluginData[GetString(f.B)]
    } else {
      data_map[key] = f.B
    }
  }
  data_map["ContactType"] = 1
  data_map["StreetNumber"] = "0001"
  data_map["Phone"] = "+420" + GetString(pPluginData["cell_phone"])

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret {
    return
  }

  if code > 299 {
    if strings.Contains(GetString(ret), "Register") {
      pPluginData["sale_status"] = "DUPLICATE"
    }
    return
  }
  var result_array = GetArray(ret["AddContactResults"])
  var result_map = GetMap(result_array[0])

  log.Printf("%v RESPONSE_DATA_CHECK: %v [ %v ]", pPluginData["plugin_log"], result_map, result_map["ExternalIdentifier"])

  if nil == result_map || nil == result_map["ExternalIdentifier"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = GetString(result_map["ExternalIdentifier"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

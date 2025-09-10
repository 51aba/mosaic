package main

import (
  "log"
  "strings"

  . "leadz/utils"
)

// *** README FIXME server ip whitelist required *** RESPONSE_DATA: <html> <head> <title>Request Rejected/Přístup zamítnut</title>
type leadplugin string

const codename string = "KAMALI"

var configs_map = map[string]string {
  "api_url": "https://www.kamali.cz/apiaffil/v01",
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
  var call_config = map[string]any{"command": "Clients"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var translated_cell_phone = "+420" + GetString(pPluginData["cell_phone"])

  var fields = []Pair{{"subjectCode", "PARTNER_NAME"},
                      {"discriminator", "PARTNER_NAME"},
                      {"campaign", "PARTNER_NAME"},
                      {"name", "first_name"},
                      {"surname","last_name"},
                      {"email", "email"},
                      {"personalNumber", "birth_number"},
                      {"nationality", "1"},
                      {"birthplace", "city"},
                      {"addressPostalCode", "zip"},
                      {"addressCity", "city"},
                      {"addressDistrict", "district"},
                      {"addressStreet", "street"},
                      {"addressStreetNumber", "house_number"},
                      {"addressMailingPostalCode", "zip"},
                      {"addressMailingCity", "city"},
                      {"addressMailingDistrict", "district"},
                      {"addressMailingStreet", "street"},
                      {"addressMailingStreetNumber", "house_number"},
                      {"amtClientIncomes", "monthly_income"},
                      {"amtClientExpenses", "monthly_expenses"},
                      {"clientIncomesSourceCode", "1"},
                      {"bankAccountNumber", "bank_account_number"},
                      {"bankAccountBankCode", "bank_code"},
                      {"amtLoanAmount", "requested_amount"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      var value any = pPluginData[GetString(f.B)]

      if nil == value {
        data_map[GetString(f.A)] = ""
      } else {
        data_map[GetString(f.A)] = value
      }
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["phoneNumber"] = translated_cell_phone

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if _, eok := ret["errors"]; eok {
    return
  }
  var data string = GetString(ret["data"])

  if "" != data {
    log.Printf("%v%v SET_RESPONSE_DATA: DATA: %v%v", RED, pPluginData["plugin_log"], data, NC)

    if strings.Contains(data, "Request Rejected") {
      return
    }
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = "https://www.kamali.cz/"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

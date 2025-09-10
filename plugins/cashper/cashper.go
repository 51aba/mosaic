package main

import (
  "fmt"
  "log"

  . "leadz/utils"
)

type leadplugin string

const codename string = "CASHPER"

var requested_sum_variants = []int{50, 100, 150, 200, 300}
var period_variants = []int{15, 30}

var configs_map = map[string]string{
  "api_url": "https://partnerservice.testing.novumbankgroup.com/Leads/",
  "username": "",
  "password": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {},}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool) {
      register_lead,},
  },
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
  var call_config = map[string]any{"command": "NewLeadInstallment"}

  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    data_map["UserName"] = config["username" + plugin_postfix]
    data_map["Password"] = config["password" + plugin_postfix]

    translate(pPluginData, data_map)
    pPluginData["map_data"] = data_map
  }
  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_home_status int
  var translate_requested_amount = fmt.Sprintf("%v", FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants))
  /*
  var translated_marital_status int

  switch GetString(pPluginData["marital_status"]) {
    case "SINGLE":
      translated_marital_status = 0
    case "MARRIED":
      translated_marital_status = 1
    case "DIVORCED":
      translated_marital_status = 2
    case "WIDOWED":
      translated_marital_status = 3
    case "PARTNERSHIP":
      translated_marital_status = 4
    default:
      translated_marital_status = 0
  }
  */

  switch GetString(pPluginData["home_status"]) {
    case "HOME_OWNER":
      translated_home_status = 2
    case "TENANT":
      translated_home_status = 1
    case "OTHER":
      translated_home_status = 3
    default:
      translated_home_status = 1
   }
  var fields = []Pair{{"ReferenceId", "id"},
                      {"Duration", "period"},
                      {"IpAddress", "ip_address"},
                      {"FirstName", "first_name"},
                      {"SurName", "last_name"},
                      {"DateOfBirth", "birth_date"},
                      {"PhoneMobile", "cell_phone"},
                      {"Street", "street"},
                      {"HouseNumber", "house_number"},
                      {"PostalCode", "zip"},
                      {"City", "city"},
                      {"Province", "province_name"},
                      {"Email", "email"},
                      {"Gender", "gender"},
                      {"IncomeNetto", "monthly_income"},
                      {"IBAN", "iban_number"},
                      {"BankAccountNumber", "bank_account_number"},
                      {"CreatedOn", "created"},
                      {"DNI", "identity_card_number"},
                      {"SecondSurname", "last_name_2"},
                    }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["CountryCode"] = "ES"
  data_map["Version"] = "1"
  data_map["Amount"] = translate_requested_amount
  data_map["Interest"] = "10"
  data_map["IncomeSourceTypeId"] = 2005
  data_map["IncomeDay"] = 10
  // data_map["MaritalStatusId"] = translated_marital_status
  // data_map["AmountOfChildren"] = GetString(pPluginData["dependent_children"])
  data_map["HousingSituationId"] = translated_home_status
  data_map["EducationId"] = 9004
  data_map["LoanReasonTypeId"] = 10
  data_map["ProductTypeIds"] = "0"
  data_map["RoadTypeId"] = 5

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  var ret_status = GetInt(ret["Status"])

  if 3 != ret_status {
    return
  }
  // var referenceId = GetInt(ret["ReferenceId"])
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = ret["InternalId"]

  log.Printf("%s SET_RESPONSE_DATA: %v", pPluginData["plugin_log"], ret["InternalId"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

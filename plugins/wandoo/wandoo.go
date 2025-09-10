package main

import (
  "fmt"
  "log"
  "strings"
  "encoding/base64"

  . "leadz/utils"
)

type leadplugin string

const codename string = "WANDOO"

var configs_map = map[string]string {
  "api_url": "https://affiliates.wandoo.eu/partners-api/spain/wandoo/PARTNER_NAME",
  "username": "",
  "password": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "redirect": {
    {},},
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 23, "param2": 65},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 500},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"STATIONARY",}},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "redirect": []func(map[string]any, map[string]string) (bool){
      redirect_action,
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
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "lead"}
  var basic = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", config["username" + plugin_postfix], config["password" + plugin_postfix])))
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Basic %v", basic)}
  pPluginData["method"] = "PUT"
  // req, err = http.NewRequest("PUT", request_url, bytes.NewBuffer(data)) // FIXME TODO *** PUT method TO ADD ***

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func redirect_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  pPluginData["redirect_url"] = config["api_url" + plugin_postfix]
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])

  return true
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  var translated_income_type = GetString(pPluginData["income_type"])

  switch translated_income_type {
    case "FULL_TIME_EMPLOYMENT":
      translated_income_type = "EMPLOYED_WORKER"
    case "OFFICIAL":
      translated_income_type = "OFFICIAL"
    case "INDIVIDUAL_ENTREPRENEUR":
      translated_income_type = "SELF_EMPLOYED"
    case "PENSION":
      translated_income_type = "PENSIONER"
    case "STUDENT":
      translated_income_type = "STUDENT"
    case "UNEMPLOYED":
      translated_income_type = "UNEMPLOYED"
    default:
      translated_income_type = "EMPLOYED_WORKER"
  }
  var address_map = map[string]any{
                    "apartmentLetter": "",
                    "apartmentNumber": pPluginData["house_ext"],
                    "city": pPluginData["city"],
                    "postalCode": GetInt(pPluginData["zip"]),
                    "province": pPluginData["province"],
                    "street": pPluginData["street"],
                    "streetNumber": pPluginData["house_number"],
                    }
  var contact_map = map[string]any{
                    "email": pPluginData["email"],
                    "phone": pPluginData["cell_phone"],
                    }
  var person_map = map[string]any{
                    "birthDate": pPluginData["birth_date"],
                    "firstName": pPluginData["first_name"],
                    "lastName": pPluginData["last_name"],
                    "nationality": "ES",
                    "personalId": pPluginData["identity_card_number"],
                    "secondLastName": pPluginData["last_name_2"],
                    }

  var employ_map = map[string]any{
                    "employmentStatus": translated_income_type,
                    "monthlyIncome": pPluginData["monthly_income"],
                    }
  data_map["customer"] = map[string]any{
                         "address": address_map,
                         "contactInformation": contact_map,
                         "employmentInformation": employ_map,
                         "personalInformation": person_map,
                         }
  data_map["lead"] = map[string]any{
                    "ipAddress": pPluginData["ip_address"],
                    "loanAmount": pPluginData["requested_amount"],
                    "loanTermDays": pPluginData["period"],
                    "uuid": pPluginData["uid"],
                    }

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil == ret["status"] || "accepted" != strings.ToLower(GetString(ret["status"])) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = ret["landingPageUri"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "fmt"
  "log"
  "strings"
  "encoding/base64"

  . "leadz/utils"
)
type leadplugin string

const codename string = "AVINTO"

var configs_map = map[string]string{
  "api_url": "https://platon.creamfinance.com/affiliate/",
  "username": "",
  "password": "",
}
var plugin_vars = []string{"redirect"}// *** README : if spec var added here it must be added to "send_map" in order request to be executed
var validators_map = map[string][]map[string]any{
  "redirect": {
    {},},
    "": {
      {"field": "birth_date", "func": "AgeRangeValidator", "param1": 23, "param2": 65},
      {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 500},
      {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"STATIONARY",}, },
  },
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "redirect": []func(map[string]any, map[string]string) (bool){
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
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var config map[string]string = P_init_named(codename, pPluginData, configs_map, plugin_postfix)
  var base64value = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", config["username" + plugin_postfix], config["password" + plugin_postfix])))

  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Basic %v", base64value)}

  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "lead"}

  if _, dok := pPluginData["map_data"]; ! dok {
    var data_map = map[string]any{}

    translate(pPluginData, data_map)
    pPluginData["map_data"] = data_map
  }
  return P_register_lead(pPluginData, call_config, set_response_data)
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
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Println(pPluginData["description"])

  return true
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
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
                    "apartmentNumber": GetString(pPluginData["house_ext"]),
                    "city": GetString(pPluginData["city"]),
                    "postalCode": GetInt(GetString(pPluginData["zip"])),
                    "province": GetString(pPluginData["province"]),
                    "street": GetString(pPluginData["street"]),
                    "streetNumber": GetString(pPluginData["house_number"]),
                    }
  var contact_map = map[string]any{
                    "email": GetString(pPluginData["email"]),
                    "phone": GetString(pPluginData["cell_phone"]),
                    }
  var person_map = map[string]any{
                    "birthDate": GetString(pPluginData["birth_date"]),
                    "firstName": GetString(pPluginData["first_name"]),
                    "lastName": GetString(pPluginData["last_name"]),
                    "nationality": "ES",
                    "personalId": GetString(pPluginData["identity_card_number"]),
                    "secondLastName": GetString(pPluginData["last_name_2"]),
                    }

  var employ_map = map[string]any{
                    "employmentStatus": translated_income_type,
                    "monthlyIncome": GetFloat(pPluginData["monthly_income"]),
                    }
  data_map["customer"] = map[string]any{
                         "address": address_map,
                         "contactInformation": contact_map,
                         "employmentInformation": employ_map,
                         "personalInformation": person_map,
                         }
  data_map["lead"] = map[string]any{
                    "ipAddress": GetString(pPluginData["ip_address"]),
                    "loanAmount": GetFloat(pPluginData["requested_amount"]),
                    "loanTermDays": GetInt(pPluginData["period"]),
                    "uuid": GetString(pPluginData["uid"]),
                    }
  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}
/*
{
  "customer": {
    "address": {
      "apartmentLetter": "a",
      "apartmentNumber": 33,
      "city": "Madrid",
      "postalCode": 24989,
      "province": "Madrid",
      "street": "Ronda Raúl",
      "streetNumber": 66
    },
    "contactInformation": {
      "email": "somemail@somewhere.es",
      "phone": "+34929299489"
    },
    "employmentInformation": {
      "employmentStatus": "Allowed values: EMPLOYED_WORKER, OFFICIAL, PENSIONER, SELF_EMPLOYED, STUDENT, UNEMPLOYED",
      "monthlyIncome": 800
    },
    "personalInformation": {
      "birthDate": "1991-07-18",
      "firstName": "Alejandra",
      "lastName": "Solorio",
      "nationality": "ES",
      "personalId": "19882079M",
      "secondLastName": "Garcia"
    }
  },
  "lead": {
    "ipAddress": "string",
    "loanAmount": 200,
    "loanTermDays": 25,
    "uuid": "string"
  }
}
*/

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil == ret["status"] || "accepted" != strings.ToLower(GetString(ret["status"])) {
    log.Printf("%v%v SET_RESPONSE_DATA: REJECTED: %v%v", RED, pPluginData["plugin_log"], ret, NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = ret["landingPageUri"]

  log.Printf("%s SET_RESPONSE_DATA: %v", pPluginData["plugin_log"], pPluginData["redirect_url"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

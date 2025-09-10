package main

import (
  "fmt"
  "log"
  "strings"
  "encoding/base64"

  . "leadz/utils"
)

type leadplugin string

const codename string = "CREDITKASA"

var configs_map = map[string]string{
  "api_url": "https://nazivnost.cz/leadpartnerapi",
  "username": "",
  "password": "",
  "external_session_id": "",
}
var plugin_vars = []string{"osvc", "cps"}
var validators_map = map[string][]map[string]any {
  "osvc": {
    {},},
  "": {
    {},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "osvc": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "cps": []func(map[string]any, map[string]string) (bool){
      register_lead_cps,},
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
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
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var translated_birth_number = GetString(pPluginData["birth_number"])

  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%v/%v", translated_birth_number[:6], translated_birth_number[6:])
  }
  var call_config = map[string]any{"command": ""}
  var data_map = map[string]any{"method": "checkDuplicate", "customer-personcode": translated_birth_number}

  pPluginData["map_data"] = data_map

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": ""}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["method"] = "registerCustomer"

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead_cps(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "service/lead/tiny"}
  var data_map = map[string]any{}

  translate_cps(pPluginData, data_map)

  data_map["external_session_id"] = config["external_session_id" + plugin_postfix]
  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data_cps)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_requested_amount (int) = GetInt(pPluginData["requested_amount"])
  var translated_birth_number = GetString(pPluginData["birth_number"])

  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%v/%v", translated_birth_number[:6], translated_birth_number[6:])
  }
                      //{"marital-status", "marital_status"},
                      // {"dependents-count-status", "dependent_children"},
                      /*
                      // {"company-name", "employer"},
                      // {"employed-for", "employed_time"},
                      // {"employer-phone", "employer_phone"},
                      // {"profession", "job_title"},
                      */
  var fields = []Pair{{"loan-term", "period"},
                      {"firstname", "first_name"},
                      {"lastname", "last_name"},
                      {"passport", "identity_card_number"},
                      {"sex", "gender"},
                      {"email", "email"},
                      {"income-source", "income_type"},
                      {"industry", "specialization"},
                      {"company-ico", "company_number"},
                      {"monthly-income", "monthly_income"},
                      {"address-zip", "zip"},
                      {"address-city", "city"},
                      {"address-street", "street"},
                      {"address-number", "house_number"},
                      {"housing", "home_status"},
                      {"contact-address-zip", "contact_zip"},
                      {"contact-address-city", "contact_city"},
                      {"contact-address-street", "contact_street"},
                      {"contact-address-number", "contact_house_number"},
                      {"contact-address-housing", "contact_home_status"},
                      {"contact-address-live-for", "address_time"},
                      {"execution", "distraint"},
                      {"insolvency", "insolvency"},
                      {"education", "education"},
                      {"third-parties-debt", "debts_to_other"},
                      {"monthly-expenses", "monthly_expenses"},
                     }
  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%s/%s", translated_birth_number[:6], translated_birth_number[6:])
  }

  if translated_requested_amount > 50000 {
    translated_requested_amount = 50000
  } else if translated_requested_amount < 5000 {
    translated_requested_amount = 5000
  }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["loan-amount"] = translated_requested_amount
  data_map["pin"] = translated_birth_number
  data_map["phone"] = "+420" + GetString(pPluginData["cell_phone"])
  data_map["type-of-income"] = "salary"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func translate_cps(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE_CPS: STARTED", pPluginData["plugin_log"])

  var fields = []Pair{{"surname", "last_name"},
                      {"name", "first_name"},
                      {"phone", "cell_phone"},
                      {"birthdate", "birth_date"},
                     }
  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["loan-type"] = "2"
  log.Printf("%v TRANSLATE_CPS: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if "" != GetString(ret["err"]) || "1" != ret["result"] {
    log.Printf("%v%v SET_RESPONSE_DATA: REJECTED: [ %v ] %v%v", RED, pPluginData["plugin_log"], GetString(ret["err"]), ret, NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = GetString(ret["lead-id"])

  if nil != ret["return-url"] {
    pPluginData["redirect_url"] = ret["return-url"]
  }
  log.Printf("%s SET_RESPONSE_DATA: %v %v", pPluginData["plugin_log"], ret["lead-id"], ret["return-url"])

  return
}

// ################################################################################################################################################################
func set_response_data_cps(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if errval := GetString(ret["err"]); "" != errval {
    log.Printf("%v%v SET_RESPONSE_DATA: REJECTED: [ %v ] %v%v", RED, pPluginData["plugin_log"], errval, ret, NC)

    if "success" != strings.ToLower(errval) {
      return
    }
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  if nil != ret["lead_tiny_id"] {
    pPluginData["external_id"] = ret["lead_tiny_id"]
  }

  if nil != ret["return-url"] {
    pPluginData["redirect_url"] = ret["return-url"]
  }
  log.Printf("%s SET_RESPONSE_DATA_CPS: %v %v", pPluginData["plugin_log"], ret["lead_tiny_id"], ret["return-url"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

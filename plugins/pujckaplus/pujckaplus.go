package main

import (
  "fmt"
  "log"
  "strings"
  "net/url"
  "encode/json"

  . "leadz/utils"
)

type leadplugin string

const codename string = "PUJCKAPLUS"

var requested_sum_variants = []int{500, 1000, 1500, 2000, 2500, 3000, 3500, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 11000, 12000, 13000, 15000}
var loan_period_variants = []int{5, 10, 15, 20, 25, 30}

var configs_map = map[string]string {
  "api_url": "https://leaderfin-cz.creditonline.eu/",
  "form_context": `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`,
  "login": "",
  "password": "",
}
var plugin_vars = []string{"cpl2", "cpl", "", }
var validators_map = map[string][]map[string]any {
  "cpl2": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "SELF_EMPLOYED", "PART_TIME_EMPLOYMENT",}},},
  "cpl": {
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "SELF_EMPLOYED", "PART_TIME_EMPLOYMENT",}},},
  "": {
    {"field": "distraint", "func": "AllowedValuesValidator", "param1": []string{"NO",}},
    {"field": "birth_number", "func": "InsolvencyValidator",},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"UNEMPLOYED",}},
    {"field": "home_status", "func": "DisallowedValuesValidator", "param1": []string{"MINISTRY",}},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpl2": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "cpl": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
  },
  true: VAR_FUNCS_MAP {
    "cpl2": []func(map[string]any, map[string]string) (bool){
      check_sms,
    },
    "cpl": []func(map[string]any, map[string]string) (bool){
      check_sms,
    },
    "": []func(map[string]any, map[string]string) (bool){
      check_sms,
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
func prepare_data(plugin_postfix string, data_map map[string]any, config map[string]string) (urlencoded_data url.Values) {
  var data []byte
  var err error

  urlencoded_data = url.Values{}
  data, err = json.Marshal(data_map)

  if nil != err {
    return
  }
  urlencoded_data.Set("login", config["login" + plugin_postfix])
  urlencoded_data.Set("pass", config["password" + plugin_postfix])
  urlencoded_data.Set("data", string(data))

  return
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "?RPC=agents.checkPersonalData"}
  var data_map = map[string]any {
                 "person_code": pPluginData["birth_number"],
                 "email": pPluginData["email"],
                 "mob_phone": pPluginData["cell_phone"],
                 "surname": pPluginData["last_name"],
                 "realname": pPluginData["first_name"],
                 "id_number": pPluginData["identity_card_number"],
                 "address": pPluginData["street"],
                 "city": pPluginData["city"],
                 "house": pPluginData["house_number"],
                 "zipcode": pPluginData["zip"],
                 "account_number": fmt.Sprintf("%v/%v", pPluginData["bank_account_number"], pPluginData["bank_code"]),
                }
  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, data_map, config)

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var data_map = map[string]any{}
  var call_config = map[string]any{"command": "?RPC=agents.registerClient"}

  translate(pPluginData, data_map)

  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, data_map, config)

  return P_register_lead(pPluginData, call_config, set_response_data_register)
}

// ################################################################################################################################################################
func check_sms(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "?RPC=agents.checkPin"}
  var sale_data_map = GetMap(pPluginData["sale_data"])

  if nil == sale_data_map {
    log.Printf("%s %v CHECK_SMS: LEAD_MAP_ISNULL", pPluginData["plugin_log"], codename)

    return
  }
  var data_map = map[string]any {"customer-id": pPluginData["external_id"], "sms-pin": sale_data_map["sms_code"]}

  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, data_map, config)

  return P_check_sms(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_income_type string
  var fields = []Pair{
                      {"email", "email"},
                      {"mob_phone", "cell_phone"},
                      {"contract_c", "1"},
                      {"marketing", "1"},
                      {"realname", "first_name"},
                      {"surname", "last_name"},
                      {"person_code", "birth_number"},
                      {"address", "street"},
                      {"house", "house_number"},
                      {"zipcode", "zip"},
                      {"address_different", "1"},
                      {"income", "monthly_income"},
                      {"ico", ""},
                      {"housing_type", ""},
                      {"monthly_mortgage_repayments", ""},
                      {"monthly_other_loans_repayments", ""},
                      {"other_costs", ""},
                      {"id_number", "identity_card_number"},
                      {"ref", "PARTNER_NAME"},
                      {"ip", "ip_address"},
                      {"chk9_formular", "1"},
                      {"chk10_vop_smlouva", "1"},
                      {"chk11_politicky_ex_osoba", "1"},
                      {"chk12_osobni_udaje_registry", "1"},
                      {"chk_osobni_udaje", "1"},
                      {"chk2", "1"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["create_credit"] = fmt.Sprintf("%v-%v", FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants),
                                                   FindClosest(GetInt(pPluginData["period"]), loan_period_variants))

  switch GetString(pPluginData["income_type"]) {
    case "EMPLOYED":
      translated_income_type = "Zaměstnanec"
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = "Personální agentura"
    case "SELF_EMPLOYED":
      translated_income_type = "OSVČ"
      data_map["workplace"] = fmt.Sprintf("OSVČ, IČ: %v", pPluginData["company_number"])
      data_map["workplace_position"] = pPluginData["company_number"]
    case "PENSION":
      translated_income_type = "Důchodce"
      data_map["workplace"] = "Důchod"
      data_map["workplace_position"] = "Důchod"
    case "MATERNITY_LEAVE":
      translated_income_type = "Mateřská dovolena"
    case "STUDENT":
      translated_income_type = "Student"
    case "BENEFITS":
      translated_income_type = "Invalidní důchodce"
    case "UNEMPLOYED":
      translated_income_type = "Nezaměstnaný"
    default:
      translated_income_type = "Ostatní"
  }
  data_map["income_source"] = translated_income_type

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data_register(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if strings.Contains(fmt.Sprintf("%v", ret), "registered") {
    pPluginData["sale_status"] = "DUPLICATE"

    return
  }

  if nil != ret["err"] || 1 != GetInt(ret["ok"]) {
    log.Printf("%v SET_RESPONSE_DATA: ERROR: %v %v", pPluginData["plugin_log"], ret["err"], ret["ok"])

    return
  }
  result = true

  pPluginData["sale_status"] = "UNCONFIRMED"

  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var config map[string]string = GetMapStrings(pPluginData["config"])

  if "" != GetString(config["form_context" + plugin_postfix]) {
    pPluginData["form_context"] = config["form_context" + plugin_postfix]
  } else {
    pPluginData["form_context"] = `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`
  }
  pPluginData["redirect_url"] = ret["url"]

  var ret_data_map = GetMap(ret["data"])

  if nil != ret_data_map {
    log.Printf("%s DATA: %v", pPluginData["plugin_log"], ret_data_map)
    pPluginData["external_id"] = ret_data_map["customer_id"]
  }

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if strings.Contains(fmt.Sprintf("%v", ret), "registered") {
    pPluginData["sale_status"] = "DUPLICATE"

    return
  }

  if nil != ret["err"] || 1 != GetInt(ret["ok"]) {
    log.Printf("%v SET_RESPONSE_DATA: ERROR: %v %v", pPluginData["plugin_log"], ret["err"], ret["ok"])

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = ret["url"]

  var ret_data_map = GetMap(ret["data"])

  if nil != ret_data_map {
    log.Printf("%s DATA: %v", pPluginData["plugin_log"], ret_data_map)
    pPluginData["external_id"] = ret_data_map["customer_id"]
  }

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

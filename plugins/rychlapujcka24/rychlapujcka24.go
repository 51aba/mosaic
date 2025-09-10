package main

import (
  "fmt"
  "log"
  "net/url"
  "github.com/bytedance/sonic"

  . "leadz/utils"
)

type leadplugin string

// *** README FIXME server ip whitelist required *** OR {"err":"Authentification was not successful."}
const codename string = "RYCHLAPUJCKA24"

var requested_sum_variants = []int{3000, 4000, 5000, 6000, 7000}
var loan_period_variants = []int{10, 15, 20, 25, 30}

var configs_map = map[string]string {
  "api_url": "https://rychlapujcka24-cz.creditonline.eu/",
  "login": "",
  "password": "",
  "form_context": `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`,
}
var plugin_vars = []string{"fl2", "", }
var validators_map = map[string][]map[string]any {
  "fl2": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 22, "param2": 100},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED",}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 8000},},
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 22, "param2": 100},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED",}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 8000},
    {"field": "requested_amount", "func": "MaxMoneyValidator", "param1": 15000},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "fl2": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
  },
  true: VAR_FUNCS_MAP {
    "fl2": []func(map[string]any, map[string]string) (bool){
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
  data, err = sonic.Marshal(data_map)

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
    "mob_phone": pPluginData["cell_phone"],}

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

  return P_register_lead(pPluginData, call_config, set_response_data)
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

  var fields = []Pair{
                      {"person_code", "birth_number"},
                      {"id_number", "identity_card_number"},
                      {"email", "email"},
                      {"mob_phone", "cell_phone"},
                      {"city", "city"},
                      {"address", "street"},
                      {"zipcode", "zip"},
                      {"ip", "ip_address"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  switch GetString(pPluginData["income_type"]) {
    case "PENSION":
    default:
      data_map["income"] = pPluginData["monthly_income"]
  }
  data_map["accept_contract4"] = 1
  data_map["marketing"] = 1
  data_map["realname"] =  fmt.Sprintf("%v %v", pPluginData["first_name"], pPluginData["last_name"])
  data_map["initial_amount"] = FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants)
  data_map["initial_term"] = FindClosest(GetInt(pPluginData["period"]), loan_period_variants)
  data_map["client_type"] = "business"
  data_map["finished"] = "1"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil != ret["err"] || nil != ret["errors"] {
    return
  }
  result = true

  if "?RPC=agents.registerClient" == command { // *** REGISTER_LEAD ***
    // pPluginData["sale_status"] = "PAUSED"
    pPluginData["sale_status"] = "UNCONFIRMED"

    var plugin_postfix = GetString(pPluginData["plugin_postfix"])
    var config map[string]string = GetMapStrings(pPluginData["config"])

    if "" != GetString(config["form_context" + plugin_postfix]) {
      pPluginData["form_context"] = config["form_context" + plugin_postfix]
    } else {
      pPluginData["form_context"] = `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "label": "SMS Code"}]}`
    }

    if nil != ret["url"] {
      pPluginData["redirect_url"] = ret["url"]
    }
  } else {
    pPluginData["sale_status"] = "UNCONFIRMED"
  } // *** CHECK_UNIQUE ***

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "fmt"
  "log"
  "strings"
  "net/url"
  "github.com/bytedance/sonic"

  . "leadz/utils"
)

type leadplugin string

// *** README FIXME server ip whitelist required *** OR {"err":"Authentification was not successful."}
const codename string = "PUJCKA7"

var requested_sum_variants = []int{1000, 1500, 2000, 2500, 3000, 3500, 4000}
var loan_period_variants = []int{7, 14, 21, 28}

var configs_map = map[string]string {
  "api_url": "https://pujcka7-cz.creditonline.eu/",
  "login": "",
  "password": "",
}
var plugin_vars = []string{"apinosms", "apisms", ""}
var validators_map = map[string][]map[string]any {
  "apinosms": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 21, "param2": 65},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT", "PENSION", "SELF_EMPLOYED"}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 9500},
    {"field": "identity_card_number", "func": "IdentityCardValidator"},},
  "apisms": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 21, "param2": 65},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT", "PENSION", "SELF_EMPLOYED"}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 9500},
    {"field": "identity_card_number", "func": "IdentityCardValidator"},},
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 21, "param2": 65},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT", "PENSION", "SELF_EMPLOYED"}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 9500},
    {"field": "identity_card_number", "func": "IdentityCardValidator"},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "apinosms": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "apisms": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
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
  var data_map = map[string]any{}
  var call_config = map[string]any{"command": "?RPC=agents.checkPersonalData"}

  translate(pPluginData, data_map)

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
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var income_type = GetString(pPluginData["income_type"])
  var translated_income_type string
  /*
  // var translated_employed_time string
  // var employed_time = GetString(pPluginData["employed_time"])

  // switch employed_time {
  //   case "0":
  //     translated_employed_time= "0 měsíců"
  //   default:
  //     translated_employed_time= "více než 5 let"
  //   case "PART_TIME_EMPLOYMENT":
  // }
  // data_map["workplace_period"] = translated_employed_time

  //   {"workplace", "employer"},
  //   {"workplace_address", "employer_city"},
  //   {"workplace_position", "job_title"},
  */

  var fields = []Pair{
    {"person_code", "birth_number"},
    {"email", "email"},
    {"mob_phone", "cell_phone"},
    {"id_number", "identity_card_number"},
    {"house_number2", "house_number"},
    {"city2", "city"},
    {"zipcode2", "zip"},
    {"income3", "monthly_income"},
    {"monthly_payments", "monthly_expenses"},
    {"ip", "ip_address"},
   }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }

  switch income_type {
    case "EMPLOYED":
      translated_income_type = "Zaměstnaný"
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = "Personální agentura"
    case "MATERNITY_LEAVE":
      translated_income_type = "Mateřská dovolena"
    case "SELF_EMPLOYED":
      translated_income_type = "OSVČ"
      data_map["workplace_position"] = fmt.Sprintf("OSVČ, IČ: %v", pPluginData["company_number"])
    case "PENSION":
      translated_income_type = "Důchodce"
      data_map["workplace"] = "Důchod"
      data_map["workplace_position"] = "Důchod"
      data_map["workplace_address"] = "Důchod"
    case "STUDENT":
      translated_income_type = "Student"
    case "BENEFITS":
      translated_income_type = "Invalidní důchodce"
    case "UNEMPLOYED":
      translated_income_type = "Nezaměstnaný"
    case "OTHER":
      translated_income_type = "Ostatní"
    default:
      translated_income_type = "Ostatní"
  }
  data_map["realname"] = fmt.Sprintf("%v%v", pPluginData["first_name"], pPluginData["last_name"])
  data_map["account_number_owner"] = fmt.Sprintf("%v%v", pPluginData["first_name"], pPluginData["last_name"])
  data_map["chk2"] = "1"
  data_map["chk6"] = "1"
  data_map["chk8"] = "1"
  data_map["address2"] = fmt.Sprintf("%v, %v, %v", pPluginData["city"], pPluginData["street"], pPluginData["house_number"])
  data_map["work_activities"] = translated_income_type
  data_map["account_number"] = fmt.Sprintf("%v/%v", pPluginData["bank_account_number"], pPluginData["bank_code"])
  data_map["accept_contract"] = "1"
  data_map["accept_contract2"] = 1
  data_map["marketing"] = 1
  data_map["chk3"] = "2"
  data_map["ref"] = "PARTNER_NAME"
  data_map["create_credit"] = fmt.Sprintf("%v-%v", FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants),
                                                   FindClosest(GetInt(pPluginData["period"]), loan_period_variants))

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if strings.Contains(GetString(ret), "registered") {
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

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

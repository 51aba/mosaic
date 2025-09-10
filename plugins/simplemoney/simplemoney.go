package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "SIMPLEMONEY"

var requested_sum_variants = []int{1000, 2000, 3000, 4000, 5000, 6000}
var period_variants = []int{7, 10, 12, 14, 21, 25, 30}

var configs_map = map[string]string {
  "api_url": "http://api.simplepujcka.cz/",
  "api_id": "",
}
var plugin_vars = []string{"cpa2", "cpa", "cpl", "",}
var validators_map = map[string][]map[string]any {
  "cpa2": {
    {"field": "birth_number", "func": "InsolvencyValidator"},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 25, "param2": 60},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT", "SELF_EMPLOYED", "MATERNITY_LEAVE", "STUDENT", "OTHER"}},
    {"field": "distraint", "func": "AllowedValuesValidator", "param1" : []string {"NO"}},},
  "cpa": {
    {},},
  "cpl": {
    {},},
  "": {
    {"field": "birth_number", "func": "InsolvencyValidator"},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 25, "param2": 60},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT", "SELF_EMPLOYED", "MATERNITY_LEAVE", "STUDENT", "OTHER"}},
    {"field": "distraint", "func": "AllowedValuesValidator", "param1" : []string {"NO"}},
    {"field": "email", "func": "WeekdayOffValidator", "param1": 5, "param2": 7},
    {"field": "first_name", "func": "Timeslot8to18Validator", "param1": true, "param2": 900, "param3": 1700},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpa": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,
    },
    "cpl": []func(map[string]any, map[string]string) (bool){
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
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "validate/validate.php"}
  var data_map = map[string]any{}

  data_map["afil_alias"] = config["api_id" + plugin_postfix]

  translate(pPluginData, data_map)
  pPluginData["map_data"] = data_map

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "save/save.php"}

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

  var translated_income_type string
  var income_type = GetString(pPluginData["income_type"])

  switch income_type {
    case "EMPLOYED":
      translated_income_type = "zaměstnanec"
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = "zaměstnanec"
    case "MATERNITY_LEAVE":
      translated_income_type = "bez zaměstnání"
    case "SELF_EMPLOYED":
      translated_income_type = "OSVČ"
    case "PENSION":
      translated_income_type = "důchodce"
    case "STUDENT":
      translated_income_type = "student"
    case "OTHER":
      translated_income_type = "bez zaměstnání"
    default:
      translated_income_type = "bez zaměstnání"
  }
  var translate_requested_amount = fmt.Sprintf("%v", FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants))
  var translate_period = fmt.Sprintf("%v", FindClosest(GetInt(pPluginData["period"]), period_variants))
  var fields = []Pair{
                      {"jmeno", "first_name"},
                      {"prijmeni", "last_name"},
                      {"rodne_cislo", "birth_number"},
                      {"cislo_uctu", "bank_account_number"},
                      {"cislo_op", "identity_card_number"},
                      {"email", "email"},
                      {"mobil", "cell_phone"},
                      {"prijem", "monthly_income"},
                      {"ulice", "street"},
                      {"cislo_popisne", "house_number"},
                      {"mesto", "city"},
                      {"psc", "zip"},
                      {"trvaly_pobyt_ulice", "contact_street"},
                      {"trvaly_pobyt_cislo_popisne", "contact_house_number"},
                      {"trvaly_pobyt_mesto", "contact_city"},
                      {"trvaly_pobyt_psc", "contact_zip"},
                      {"blizka_osoba_jmeno", "first_name"},
                      {"blizka_osoba_prijmeni", "last_name"},
                      {"blizka_osoba_telefon", "cell_phone"},
                      {"ip_adresa", "ip_address"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  data_map["castka"] = translate_requested_amount
  data_map["doba"] = translate_period
  data_map["povolani"] = translated_income_type
  data_map["majitel_auta"] = "ne"
  data_map["blizka_osoba_vztah"] = "jiný rodinný příslušník"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  log.Printf("%v SET_RESPONSE_DATA: COMMAND: %v", pPluginData["plugin_log"], command)

  if "validate/validate.php" == command {
    if strings.Contains(strings.ToLower(fmt.Sprintf("%v", ret)), "ano") {
      result = true
      pPluginData["sale_status"] = "UNCONFIRMED"
    }
    return
  }

  if nil == ret["linkSouhlas"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = ret["linkSouhlas"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

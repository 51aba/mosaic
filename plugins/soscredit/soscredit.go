package main

import (
  "fmt"
  "log"
  "net/url"
  "strings"
  "github.com/bytedance/sonic"

  . "leadz/utils"
)

type leadplugin string

const codename string = "SOSCREDIT"

var requested_sum_variants = []int{1000, 1500, 2000, 2500, 3000, 3500, 4000, 4500, 5000, 6000, 7000}
var loan_period_variants = []int{7, 14, 21, 28}

var configs_map = map[string]string{
  "api_url":  "https://www.soscredit.cz/api/broker/action/",
  "referer":  "PARTNER_NAME.cz",
  "password": "",
  "unique_id": "",
}
var plugin_vars = []string{""}
var validators_map = map[string][]map[string]any{
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 21, "param2": 68},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "PART_TIME_EMPLOYMENT"}},
    {"field": "insolvency", "func": "DisallowedValuesValidator", "param1": []string{"YES"}},
    {"field": "distraint", "func": "AllowedValuesValidator", "param1": []string{"NO"}}},
}
var send_map = SEND_MAP{
  false: VAR_FUNCS_MAP{
    "": []func(map[string]any, map[string]string) bool{
      check_unique, register_lead,},
  },
}

// ################################################################################################################################################################
func (p leadplugin) TestData(pPluginData map[string]any, is_paused bool) (ret bool) {
  return p.SendData(pPluginData, is_paused)
}

// ################################################################################################################################################################
func (p leadplugin) Validate(pPluginData map[string]any) []map[string]any {
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
    log.Printf("%v PREPARE_DATA: MARSHAL_ERROR: %v", codename, err)

    return
  }
  urlencoded_data.Set("uniqueId", config["unique_id" + plugin_postfix])
  urlencoded_data.Set("password", config["password" + plugin_postfix])
  urlencoded_data.Set("data", string(data))

  return
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var call_config = map[string]any{"command": "checkUserDuplicate"}
  var user_map = map[string]any{
    "persId":      pPluginData[BIRTH_NUMBER],
    "persId2":     pPluginData[IDENTITY_CARD_NUMBER],
    "phone":       pPluginData[CELL_PHONE],
    "email":       pPluginData[EMAIL],
    "bankAccount": pPluginData[BANK_ACCOUNT_NUMBER],
  }
  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, user_map, config)

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var call_config = map[string]any{"command": "userRegister"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  data_map["referer"] = config["referer"]
  data_map["refererHost"] = config["referer"]
  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, data_map, config)

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_income_type = GetString(pPluginData[INCOME_TYPE])

  switch translated_income_type {
    case "EMPLOYED", "PART_TIME_EMPLOYMENT", "TEMPORARY_EMPLOYMENT", "STUDENT":
      translated_income_type = "ZAM"
    case "SELF_EMPLOYED":
      translated_income_type = "OSVČ"
    case "UNEMPLOYED":
      translated_income_type = "BEZ"
    case "PENSION", "RETIRED":
      translated_income_type = "DUCH"
    case "MATERNITY_LEAVE":
      translated_income_type = "MAT"
    default:
      translated_income_type = "ZAM"
  }
  var translate_period = fmt.Sprintf("%v", FindClosest(GetInt(pPluginData[PERIOD]), loan_period_variants))
  var translate_requested_amount = fmt.Sprintf("%v", FindClosest(GetInt(pPluginData[REQUESTED_AMOUNT]), requested_sum_variants) * 100)

  data_map["loanData"] = map[string]any{
    "amount": translate_requested_amount,
    "days":   translate_period,
  }

  data_map["userData"] = map[string]any{
    "name":              pPluginData[FIRST_NAME],
    "surname":           pPluginData[LAST_NAME],
    "persId":            pPluginData[BIRTH_NUMBER],
    "persId2":           pPluginData[IDENTITY_CARD_NUMBER],
    "phone":             "420" + GetString(pPluginData[CELL_PHONE]),
    "gender":            pPluginData[GENDER],
    "email":             pPluginData[EMAIL],
    // "workplace":         pPluginData["employer"],
    // "workplaceAddress":  pPluginData["employer_city"],
    // "job":               pPluginData["job_title"],
    // "employmentTime":    pPluginData[EMPLOYED_TIME],
    // "liabilities":       GetString(pPluginData["debts_to_other"]
    "occupation":        translated_income_type,
    "birthDate":         pPluginData[BIRTH_DATE],
    "actualAddrStreet":  pPluginData[STREET],
    "actualAddrCity":    pPluginData[CITY],
    "actualAddrPostal":  pPluginData[ZIP],
    "actualAddrNumber":  pPluginData[HOUSE_NUMBER],
    "legalAddrStreet":   pPluginData[STREET],
    "legalAddrCity":     pPluginData[CITY],
    "legalAddrPostal":   pPluginData[ZIP],
    "legalAddrNumber":   pPluginData[HOUSE_NUMBER],
    "salary":            pPluginData[MONTHLY_INCOME],
    "registrationIp":    pPluginData[IP_ADDRESS],
    "bankAccount":       fmt.Sprintf("%v/%v", pPluginData[BANK_ACCOUNT_NUMBER], pPluginData[BANK_CODE]),
    "bankAccountHolder": fmt.Sprintf("%v %v", pPluginData[FIRST_NAME], pPluginData[LAST_NAME]),
  }

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var status = strings.ToLower(GetString(ret["status"]))

  if "success" != status {
    log.Printf("%v%v SET_RESPONSE_DATA: STATUS_ERROR: %v%v", RED, pPluginData["plugin_log"], ret, NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

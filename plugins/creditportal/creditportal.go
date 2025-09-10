package main

import (
  "fmt"
  "log"
  "sort"
  "time"
  "strings"
  "net/url"
  "crypto/md5"
  "github.com/leekchan/timeutil"

  . "leadz/utils"
)
type leadplugin string

const codename string = "CREDITPORTAL"

var configs_map = map[string]string{
  "api_url": "https://www.creditportal.cz/cs/api",
  "partner": "",
  "partner_secret": "",
  "channelid": "",
  "form_context": `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "validator": "number", "label": "SMS Code"}]}`,
}
var plugin_vars = []string{"301", "302", "303", "304", "305", "3651", "3652", "3653", "short"}// *** README : if spec var added here it must be added to "send_map" in order request to be executed
var validators_map = map[string][]map[string]any{
  "short": {
    {},},
  "301": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED",}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 26, "param2": 100},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 24000},},
  "302": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "SELF_EMPLOYED",}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 23, "param2": 100},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 19000},},
  "303": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "SELF_EMPLOYED", "PENSION",}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 15000},},
  "304": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "SELF_EMPLOYED", "PENSION",}},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 8000},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 19, "param2": 65},
    {"field": "contactPersonType", "func": "AllowedValuesValidator", "param1": []string{"relative", "friend", "employee",}}, // FIXME TODO not applicable
    {"field": "bank_code", "func": "AllowedValuesValidator", "param1": []string{"0800", "0300", "0600", "0100", "5500"}}, // FIXME TODO not applicable
    {"field": "requested_amount", "func": "MinMoneyValidator", "param1": 500},
    {"field": "requested_amount", "func": "MaxMoneyValidator", "param1": 4900},},
  "305": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED", "SELF_EMPLOYED", "PENSION",}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 20, "param2": 100},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 15000},},
  "": {
    {"field": "birth_number", "func": "AgeRangeValidator", "param1": 21, "param2": 65, "param3": "test"},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "short": []func(map[string]any, map[string]string) (bool){
      register_lead_short,},
    "301": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "302": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "303": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "304": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "305": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
    "3651": []func(map[string]any, map[string]string) (bool){
      check_unique_365, register_lead_365,},
    "3652": []func(map[string]any, map[string]string) (bool){
      check_unique_365, register_lead_365,},
    "3653": []func(map[string]any, map[string]string) (bool){
      check_unique_365, register_lead_365,},
    "": []func(map[string]any, map[string]string) (bool){
      check_unique, register_lead,},
  },
  true: VAR_FUNCS_MAP {
    "301": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "302": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "303": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "304": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "305": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "3651": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "3652": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "3653": []func(map[string]any, map[string]string) (bool){
      check_sms,},
    "": []func(map[string]any, map[string]string) (bool){
      check_sms,},
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
func prepare_data(command string, plugin_postfix string, data_map map[string]any, config map[string]string) (urlencoded_data url.Values) {
  var hash_sum string
  var keys []string
  var locale = time.FixedZone("Europe/Warsaw", 6*60*60)
  var tp = time.Date(1970, 1, 1, 0, 0, 0, 0, locale)

  urlencoded_data = url.Values{}
  data_map["partner"] = config["partner" + plugin_postfix]
  data_map["version"] = "1.2"
  data_map["command"] = command
  data_map["timestamp"] = time.Now().Sub(tp).Seconds()

  keys = make([]string, 0, len(data_map))

  for k := range data_map {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  for _, key := range keys {
    var value = GetString(data_map[key])

    if "" == hash_sum {
      hash_sum = value
    } else {
      hash_sum = hash_sum + "|" + value
    }
    urlencoded_data.Set(key, value)
  }
  hash_sum = hash_sum + "|" + config["partner_secret" + plugin_postfix]
  urlencoded_data.Set("signature", fmt.Sprintf("%x", md5.Sum([]byte(hash_sum))))

  return
}

// ################################################################################################################################################################
func register_lead_short(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var call_config = map[string]any{"description": "SHORT"}
  var td = timeutil.Timedelta{Days: time.Duration(30)}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["borrowUntil"] = time.Now().Add(td.Duration()).Format("02.01.2006")

  pPluginData["urlencoded_data"] = prepare_data("submitShortLead", plugin_postfix, data_map, config)

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var data_map = map[string]any{"birthId": pPluginData["birth_number"]}
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)

  pPluginData["urlencoded_data"] = prepare_data("checkBirthId", plugin_postfix, data_map, config)

  return P_check_unique(pPluginData, nil, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var td = timeutil.Timedelta{Days: time.Duration(30)}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["date"] = time.Now().Add(td.Duration()).Format("02.01.2006")

  //=========================================================
  var call_config = map[string]any{"description": "DOC"}
  pPluginData["urlencoded_data"] = prepare_data("documentation", plugin_postfix, data_map, config)

  if result = P_register_lead(pPluginData, call_config, set_response_data); ! result {
  }
  //---------------------------------------------------------
  var command string = "submitLead"

  //=========================================================
  data_map["channelid"] = config["channelid" + plugin_postfix]
  data_map["channelId"] = config["channelid" + plugin_postfix]
  pPluginData["urlencoded_data"] = prepare_data(command, plugin_postfix, data_map, config)
  pPluginData["command"] = command

  return P_register_lead(pPluginData, nil, set_response_data)
  //---------------------------------------------------------
}

// ################################################################################################################################################################
func check_sms(pPluginData map[string]any, config map[string]string) (result bool) {
  var sale_data_map = GetMap(pPluginData["sale_data"])
  var form_path string = GetString(pPluginData["form_path"])

  if nil == sale_data_map || "" == form_path {
    log.Printf("%v SEND_DATA: [%v] CHECK_SMS: SALE_STATUS: %v SALE_MAP_OR_SMS_DATA_ISNULL_ERROR: %v, %v", pPluginData["plugin_log"], pPluginData["function_idx"], pPluginData["sale_status"], sale_data_map, form_path)

    return
  }
  //=========================================================
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var command string = "submitToken"
  var data_map = map[string]any{"birthId": pPluginData["birth_number"]}
  var token_code, token_id string

  var parts = strings.Split(form_path, ",")

  if len(parts) > 0 {
    token_code = parts[0]

    if len(parts) > 1 {
      token_id = parts[1]
    }
  }
  data_map["tokenValue"] = sale_data_map["sms_code"]
  data_map["tokenId"] = token_id
  data_map["tokenCode"] = token_code
  pPluginData["urlencoded_data"] = prepare_data(command, plugin_postfix, data_map, config)
  pPluginData["command"] = command

  return P_check_sms(pPluginData, nil, set_response_data)
  //---------------------------------------------------------
}

// ################################################################################################################################################################
func check_unique_365(pPluginData map[string]any, config map[string]string) (result bool) {
  var data_map = map[string]any{"birthId": pPluginData["birth_number"]}
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)

  pPluginData["urlencoded_data"] = prepare_data("checkInstallmentBirthId", plugin_postfix, data_map, config)

  return P_check_unique(pPluginData, nil, set_response_data)
}

// ################################################################################################################################################################
func register_lead_365(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var td = timeutil.Timedelta{Days: time.Duration(30)}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)
  data_map["date"] = time.Now().Add(td.Duration()).Format("02.01.2006")

  //=========================================================
  var call_config = map[string]any{"description": "INSTALLMENT_DOC"}
  pPluginData["urlencoded_data"] = prepare_data("installmentDocumentation", plugin_postfix, data_map, config)

  if result = P_register_lead(pPluginData, call_config, set_response_data); ! result {
  }
  //=========================================================
  var command string = "submitInstallmentLead"

  //=========================================================
  data_map["channelid"] = config["channelid" + plugin_postfix]
  data_map["channelId"] = config["channelid" + plugin_postfix]
  pPluginData["urlencoded_data"] = prepare_data(command, plugin_postfix, data_map, config)
  pPluginData["command"] = command

  return P_register_lead(pPluginData, nil, set_response_data)
  //=========================================================
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  if nil == data_map {
    data_map = map[string]any{}
  }
  var fields = []Pair{{"birthId", "birth_number"},
                      {"amount", "requested_amount"},
                      {"name", "first_name"},
                      {"surname", "last_name"},
                      {"citizenId", "identity_card_number"},
                      {"phone", "cell_phone"},
                      {"email", "email"},
                      {"town", "city"},
                      {"townContact", "city"},
                      {"contactPersonPhone", "cell_phone"},
                      {"incomeAmount", "monthly_income"},
                      {"bankId", "bank_code"},
                      {"bankAccount", "bank_account_number"},
                      {"contactPerson", "last_name"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      var value = pPluginData[GetString(f.B)]

      if nil != value {
        data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
      } else {
        data_map[GetString(f.A)] = ""
      }
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  var translated_income_type = GetString(pPluginData["income_type"])
  var translated_home_status = GetString(pPluginData["home_status"])

  switch translated_income_type {
    case "EMPLOYED":
      translated_income_type = "full-time"
    case "PART_TIME_EMPLOYMENT":
      translated_income_type = "part-time"
    case "SELF_EMPLOYED":
      translated_income_type = "self-employed"
    case "STUDENT":
      translated_income_type = "student"
    case "MATERNITY_LEAVE":
      translated_income_type = "maternity"
    default:
      translated_income_type = "self-employed"
  }

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "home owner"
    case "TENANT":
      translated_home_status = "tenant"
    case "DORMITORY":
      translated_home_status = "dormitory"
    default:
      translated_home_status = "tenant"
  }
  data_map["street"] = fmt.Sprintf("%v %v", pPluginData["street"], pPluginData["house_number"])
  data_map["streetContact"] = data_map["street"]
  data_map["incomeSource"] = translated_income_type
  data_map["product"] = "online"
  data_map["contactPersonType"] = "employee"
  data_map["universityName"] = "university"
  data_map["homeStatus"] = translated_home_status
  data_map["homeType"] = "flat"
  data_map["postalCode"] = pPluginData["zip"]
  data_map["postalCodeContact"] = pPluginData["zip"]

  var employer string = GetString(pPluginData["employer"])
  var employer_phone string = GetString(pPluginData["employer_phone"])

  if "" == employer {
     employer = "Hrbota"
  }

  if "" == employer_phone {
    employer_phone = "777777777"
  }
  data_map["companyName"] = employer
  data_map["companyPhone"] = employer_phone

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  // *** {"status":"ok","version":"1.2","data":{"accepted":true}}

  if nil == ret["data"] {
    return
  }
  var ret_data_map = GetMap(ret["data"])
  log.Printf("%v RESPONSE_DATA: %v [%v] COMMAND: %v", pPluginData["plugin_log"], ret["data"], ret_data_map, pPluginData["command"])

  if nil == ret_data_map || ! GetBool(ret_data_map["accepted"]) {
    return
  }
  result = true
  pPluginData["redirect_url"] = ret_data_map["redirectUri"]

  if strings.Contains(GetString(pPluginData["command"]), "submitLead") || strings.Contains(GetString(pPluginData["command"]), "submitInstallmentLead") {
    pPluginData["sale_status"] = "PAUSED"
    pPluginData["external_id"] = ret_data_map["loanNumber"]

    var form_path string = fmt.Sprintf("%v,%v", ret_data_map["tokenCode"], ret_data_map["tokenId"])

    if len(form_path) > 128 {
      form_path = form_path[:128]
    }
    pPluginData["form_path"] = form_path
    var plugin_postfix = GetString(pPluginData["plugin_postfix"])
    var config map[string]string = GetMapStrings(pPluginData["config"])

    if "" != GetString(config["form_context" + plugin_postfix]) {
      pPluginData["form_context"] = config["form_context" + plugin_postfix]
    } else {
      pPluginData["form_context"] = `{"items": [{"id": 0, "name": "sms_code", "input_type": "text", "validator": "number", "label": "SMS Code"}]}`
    }
  } else { // *** _365 _short ***
    pPluginData["sale_status"] = "UNCONFIRMED"
  }
  log.Printf("%s SET_RESPONSE_DATA: FINAL_SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["sale_status"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

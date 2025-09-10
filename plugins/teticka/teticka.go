package main

import (
  "fmt"
  "log"
  "time"
  "github.com/leekchan/timeutil"

  . "leadz/utils"
)

type leadplugin string

// *** README FIXME server ip whitelist required *** result:map[errorMessage:IP address 46.32.86.143 is not whitelisted! isOk:false token:<nil>] success:true targetUrl:<nil> unAuthorizedRequest:false]
const codename string = "TETICKA"

var requested_sum_variants = []int{1000, 1500, 2000, 2500, 3000, 3500, 4000, 4500, 5000, 6000, 7000, 8000, 9000, 10000}
var loan_period_variants = []int{7, 14, 21, 28}

var configs_map = map[string]string {
  "api_url": "https://orbis.faircredit.cz/api/",
  "apikey": "",
}
var plugin_vars = []string{"cpl", "new", "",}
var validators_map = map[string][]map[string]any {
  "cpl": {
    {},},
  "new": {
    {},},
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 19, "param2": 75},
    {"field": "birth_number", "func": "InsolvencyValidator"},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 14000},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"EMPLOYED",}},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "cpl": []func(map[string]any, map[string]string) (bool){
      get_token, check_unique, register_lead,},
    "new": []func(map[string]any, map[string]string) (bool){
      get_token, register_lead,},
    "": []func(map[string]any, map[string]string) (bool){
      get_token, register_lead,},
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
func get_token(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "leads/login-with-api-key", "description": "GET_TOKEN"}

  pPluginData["map_data"] = map[string]any{"apiKey": config["apikey" + plugin_postfix]}

  return P_check_unique(pPluginData, call_config, set_response_data_token)
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  if "" == config["auth_token"] {
    return
  }
  var call_config = map[string]any{"command": fmt.Sprintf("leads/pre-check/%v", pPluginData["birth_number"])}

  pPluginData["map_data"] = map[string]any{}
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %v", config["auth_token"])}

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "leads/import"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map
  pPluginData["headers"] = map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %v", config["auth_token"])}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_income_type string = GetString(pPluginData["income_type"])
  var translated_home_status string = GetString(pPluginData["home_status"])
  var translated_marital_status string = GetString(pPluginData["marital_status"])

  switch translated_income_type {
    case "EMPLOYED":
      translated_income_type = "FullEmployment"
    case "SELF_EMPLOYED":
      translated_income_type = "PartialEmployment"
    case "UNEMPLOYED":
      translated_income_type = "Unemployed"
    case "PENSION":
      translated_income_type = "Pension"
    case "MATERNITY":
      translated_income_type = "MaternityLeave"
    case "OTHER":
      translated_income_type = "OwnAccountWorker"
    default:
      translated_income_type = "FullEmployment"
  }

  /*
  // var translated_employed_time string = GetString(pPluginData["employed_time"])
  // switch translated_employed_time {
  //   case "PRIMARY":
  //     translated_employed_time = "Elementary"
  //   case "APPRENTICESHIP":
  //     translated_employed_time = "Apprenticeship"
  //   case "SECONDARY_PROFESSIONAL":
  //     translated_employed_time = "HighSchool"
  //   case "UNIVERSITY_BACHELOR":
  //     translated_employed_time = "HigherEducation"
  //   case "GRADUATE":
  //     translated_employed_time = "University"
  //   default:
  //     translated_employed_time = "Elementary"
  // }
  // inner_data_map["educationType"] = translated_employed_time
  // inner_data_map["employmentUntil"] = time.Now().Add(td.Duration()).Format("2006-01-02")
  */

  switch translated_marital_status {
    case "MARRIED":
      translated_marital_status = "Married"
    case "DIVORCED":
      translated_marital_status = "Divorced"
    case "SINGLE":
      translated_marital_status = "Single"
    case "PARTNERSHIP":
      translated_marital_status = "CommonHousehold"
    case "WIDOWED":
      translated_marital_status = "Widow"
    default:
      translated_marital_status = "Married"
  }

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "OwnHouse"
    case "TENANT":
      translated_home_status = "InRent"
    case "CO_OWNED":
      translated_home_status = "ParentsHouse"
    case "EMPLOYEE_BENEFIT":
      translated_home_status = "Hostel"
    case "OTHER":
      translated_home_status = "Hostel"
    default:
      translated_home_status = "OwnHouse"
  }
                      /*
                      // {"workplace_position", "job_title"},
                      // {"employerName", "employer"},
                      // {"workplace", "employer"},
                      // {"workplace_address", "employer_city"},
                      */
  var fields = []Pair{
                      {"externalId", "uid"},
                      {"permanentAddressCity", "city"},
                      {"permanentAddressStreet", "street"},
                      {"permanentAddressZipCode", "zip"},
                      {"identificationDocumentNumber", "identity_card_number"},
                      {"bankAccountNumber", "bank_account_number"},
                      {"bankAccountBankCode", "bank_code"},
                      {"userEmailAddress", "email"},
                      {"userFirstName", "first_name"},
                      {"userLastName", "last_name"},
                      {"userPhoneNumber", "cell_phone"},
                      {"personalIdentificationNumber", "birth_number"},
                      {"house_number2", "house_number"},
                      {"zipcode2", "zip"},
                      {"income3", "monthly_income"},
                      {"monthly_payments", "monthly_expenses"},
                      {"ip", "ip_address"},
                     }

  var inner_data_map = make(map[string]any)

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      inner_data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      inner_data_map[GetString(f.A)] = f.B
    }
  }
  var td = timeutil.Timedelta{Days: time.Duration(2 * 365)}

  inner_data_map["identificationDocumentExpirationDate"] = time.Now().Add(td.Duration()).Format("2006-01-02")
  td = timeutil.Timedelta{Days: time.Duration(10 * 365)}
  inner_data_map["permanentAddressRuianCode"] = ""
  inner_data_map["incomeSource"] = "IndefinitePeriod"
  inner_data_map["employmentType"] = translated_income_type
  inner_data_map["daysToDueDate"] = 365
  inner_data_map["housingType"] = translated_home_status
  inner_data_map["maritalStatus"] = translated_marital_status
  inner_data_map["incomeNet"] = 0
  inner_data_map["expensesOnLoans"] = 0
  inner_data_map["expensesOnLiving"] = 0
  inner_data_map["expenses"] = 0
  inner_data_map["dependentChildCount"] = 0
  inner_data_map["note"] = ""
  inner_data_map["loanAmount"] = FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants)
  data_map["leads"] = []map[string]any{inner_data_map}

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data_token(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  log.Printf("%v SET_RESPONSE_DATA_TOKEN: %v", pPluginData["plugin_log"], ret["accessToken"])
  var config = GetMapStrings(pPluginData["config"])

  config["auth_token"] = GetString(ret["accessToken"])
  pPluginData["config"] = config

  if "" == config["auth_token"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil == ret["result"] && nil == ret["status"] {
    return
  }
  log.Printf("%s RESPONSE_DATA: [ %v ] [ %v ]\n", pPluginData["plugin_log"], ret["result"], ret["status"])

  if "NEW_CONTRACT" != GetString(ret["result"]) && "CREATED" != GetString(ret["status"]) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

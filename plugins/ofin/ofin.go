package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "net/url"
  "github.com/bytedance/sonic"
  "github.com/leekchan/timeutil"

  . "leadz/utils"
)

type leadplugin string

const codename string = "OFIN"

const (
  min_amount int = 1000
  mid_amount int = 20000
  max_amount int = 50000
  period_max int = 30
)

var configs_map = map[string]string {
  "api_url": "https://ofin-cz.creditonline.eu/",
  "login": "",
  "password": "",
  "broker_id": "",
  "api_key": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "fl2": {
    {"field": "requested_amount", "func": "MaxMoneyValidator", "param1": max_amount},},
  "": {
    nil,},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
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
  pPluginData["headers"] = map[string]string{"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}

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
func loan_application(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "api/v1.0/creditLine/application", "description": "LOAN_APPLICATION"}
  var period int = GetInt(pPluginData["period"])

  if period > period_max {
  }
  period = period_max // *** FIXME ***

  var paydays_td = timeutil.Timedelta{Days: time.Duration(period)}
  var amount int = GetInt(pPluginData["requested_amount"])

  if amount > mid_amount {
    amount = mid_amount
  } else if amount < min_amount {
    amount = min_amount
  }
  var bank_number string = GetString(pPluginData["bank_account_number"])

  if "" == bank_number {
  } else {
    bank_number = fmt.Sprintf("%v/%v", bank_number, pPluginData["bank_code"])
  }
  pPluginData["map_data"] = map[string]any{
                        "customerId": pPluginData["external_id"],
                        "amount": amount,
                        "maxAmount": mid_amount,
                        "account_number": bank_number,
                        "paydays":  time.Now().Add(paydays_td.Duration()).Format("2006-01-02"),
                        "brokerId": config["broker_id" + plugin_postfix],
                        "_method": "post",
                        "apiLan": "cz",
                        "apiKey": config["api_key" + plugin_postfix],
                       }

  return P_register_lead_get(pPluginData, call_config, set_response_data_get)
}

// ################################################################################################################################################################
func check_unique(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var data_map = map[string]any{}
  var call_config = map[string]any{"command": "?RPC=agents.checkPersonalData"}

  translate(pPluginData, data_map, true)

  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, data_map, config)

  return P_check_unique(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "?RPC=agents.registerClient"}

  log.Printf("%v%v REGISTER_LEAD: OFIN_TYPE: %v%v", YELLOW, pPluginData["plugin_log"], pPluginData["ofin_type"], NC)

  //-------------------------------------------------------------------------------------------------
  if "klient" == pPluginData["ofin_type"] || "registrator" == pPluginData["ofin_type"] {
    pPluginData["config"] = nil
    plugin_postfix = "_" + GetString(pPluginData["ofin_type"])
    pPluginData["plugin_postfix"] = plugin_postfix
    pPluginData["plugin_log"] = fmt.Sprintf("%v%v %v", codename, plugin_postfix, pPluginData["uid"])

    config = P_init_named(codename, pPluginData, configs_map, "_" + GetString(pPluginData["ofin_type"]))
    pPluginData["config"] = config

    if result = check_unique(pPluginData, config); ! result {
      return
    }
    result = loan_application(pPluginData, config)

    return
  }
  //=================================================================================================
  var data_map = map[string]any{}

  translate(pPluginData, data_map, false)
  data_map["ref"] = config["login" + plugin_postfix]

  pPluginData["urlencoded_data"] = prepare_data(plugin_postfix, data_map, config)

  if result = P_register_lead(pPluginData, call_config, set_response_data); ! result {
    return
  }
  result = loan_application(pPluginData, config)

  return
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any, unique bool) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var fields = []Pair{
                      {"realname", "first_name"},
                      {"surname", "last_name"},
                      {"email", "email"},
                      {"mob_phone", "cell_phone"},
                      {"person_code", "birth_number"},
                      {"id_number", "identity_card_number"},
                      {"address", "street"},
                      {"house", "house_number"},
                      {"city", "city"},
                      {"zipcode", "zip"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  var bank_number string = GetString(pPluginData["bank_account_number"])

  if "" == bank_number {
    data_map["account_number"] = ""
  } else {
    data_map["account_number"] = fmt.Sprintf("%v/%v", bank_number, pPluginData["bank_code"])
  }

  if unique {
    return
  }
  var translated_home_status = GetString(pPluginData["home_status"])
  var translated_education = GetString(pPluginData["education"])
  var translated_income_type = GetString(pPluginData["income_type"])

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "Vlastní byt bez hypotéky"
    case "CO_OWNED":
      translated_home_status = "Sdílené bydlení"
    case "HOSTEL":
      translated_home_status = "Pronajatý byt"
    default:
      translated_home_status = "Jiné"
   }

  switch translated_education {
    case "PRIMARY":
      translated_education = "Základní vzdělání"
    case "SECONDARY_PROFESSIONAL":
      translated_education = "Středoškolské"
    case "UNIVERSITY_BACHELOR":
      translated_education = "Vysokoškolské - bakalářský titul"
    case "UNIVERSITY_MASTER":
      translated_education = "Vysokoškolské - magisterský titul"
    default:
      translated_education = "Bez vzdělání"
  }

  if "EMPLOYED" == translated_income_type {
    data_map["work_length"] = "Doba určitá"
    data_map["planning_to_leave_job"] = "Ne"
    data_map["training_job_period"] = "Ne"
  } else if "SELF_EMPLOYED" == translated_income_type {
    data_map["company_code"] = pPluginData["company_number"]
  }

  switch translated_income_type {
    case "EMPLOYED":
      translated_income_type = "Zaměstnaný"
    case "SELF_EMPLOYED":
      translated_income_type = "OSVČ"
    case "MATERNITY_LEAVE":
      translated_income_type = "Mateřsk"
    case "STUDENT":
      translated_income_type = "Student"
    case "PENSION":
      translated_income_type = "Důchodce"
    case "UNEMPLOYED":
      translated_income_type = "Nezaměstnaný"
    default:
      translated_income_type = "Nezaměstnaný"
  }
  data_map["income_type"] = translated_income_type
  data_map["address_different"] = "0"
  data_map["place_of_birth"] = pPluginData["city"]
  data_map["place_of_living"] = "0"
  data_map["nationality"] = "CZE"
  data_map["country_of_birth"] = "CZE"
  data_map["education_degree"] = translated_education
  data_map["income"] = pPluginData["monthly_income"]
  data_map["expenses"] = pPluginData["monthly_expenses"]
  data_map["ip"] = pPluginData["ip_address"]
  data_map["contract_accept"] = "1"
  data_map["marketing"] = 1

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var pers_data = GetMap(ret["personalData"])

  if nil != ret["err"] {
    log.Printf("%s SET_RESPONSE_DATA: ERROR: %v", pPluginData["plugin_log"], ret["ok"])

    return
  } else if strings.Contains(fmt.Sprintf("%v", ret), "registered") {
    pPluginData["redirect_url"] = ret["url"]

    if "klient" == strings.ToLower(GetString(pers_data["addinfo"])) {
      log.Printf("%s SET_RESPONSE_DATA: KLIENT: %v", pPluginData["plugin_log"], ret)

      pPluginData["ofin_type"] = "klient"
      result = true
      pPluginData["sale_status"] = "UNCONFIRMED"
      pPluginData["external_id"] = pers_data["customer_id"]
    } else if "registrator" == strings.ToLower(GetString(pers_data["addinfo"])) {
      log.Printf("%s SET_RESPONSE_DATA: REGISTRATOR: %v", pPluginData["plugin_log"], ret)

      pPluginData["ofin_type"] = "registrator"
      result = true
      pPluginData["sale_status"] = "UNCONFIRMED"
      pPluginData["external_id"] = pers_data["customer_id"]
    } else {
      log.Printf("%s SET_RESPONSE_DATA: DUPLICATE: %v", pPluginData["plugin_log"], ret)

      pPluginData["sale_status"] = "DUPLICATE"
    }
  } else if "1" == GetString(ret["ok"]) {
    log.Printf("%s SET_RESPONSE_DATA: OK: %v", pPluginData["plugin_log"], ret["ok"])

    var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
    pPluginData["redirect_url"] = ret["url"]

    if ("_klient" == plugin_postfix && "klient" == strings.ToLower(GetString(pers_data["addinfo"]))) ||
       ("_registrator" == plugin_postfix && "registrator" == strings.ToLower(GetString(pers_data["addinfo"]))) {
    } else {
      result = true
      pPluginData["sale_status"] = "UNCONFIRMED"
      pPluginData["external_id"] = ret["customer_id"]
    }
  } else {
    log.Printf("%s SET_RESPONSE_DATA: REJECTED: %v", pPluginData["plugin_log"], ret)
  }
  return
}

// ################################################################################################################################################################
func set_response_data_get(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA_GET: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var ret_data_map = GetMap(ret["data"])

  if nil != ret["err"] || 200 != GetInt(ret["status"]) || nil == ret_data_map || nil != ret_data_map["errors"] {
    log.Printf("%s SET_RESPONSE_DATA_GET: REJECTED: %v [ DATA_MAP: %v ]", pPluginData["plugin_log"], ret, ret_data_map)

    return
  }
  log.Printf("%v%v SET_RESPONSE_DATA_GET: STATUS: %v OFIN_TYPE: %v [ DATA_MAP: %v ]%v", YELLOW, pPluginData["plugin_log"], ret["status"], pPluginData["ofin_type"], ret_data_map, NC)

  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  // pPluginData["external_id"] = ret_data_map["credit_id"]
  pPluginData["redirect_url"] = ret_data_map["url"]

  if "klient" == pPluginData["ofin_type"] { // README extra settings
  } else if "registrator" == pPluginData["ofin_type"] {
  }

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

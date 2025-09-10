package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "math/rand"

  . "leadz/utils"
)

type leadplugin string

const codename string = "DINEVO"

var configs_map = map[string]string{
  "api_url": "https://api.lainasto.com/api/",
  "api_key": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
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
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "post"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map, plugin_postfix)


  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any, plugin_postfix string) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var config = GetMapStrings(pPluginData["config"])

  if nil == config {
    log.Printf("%v LEAD_CONFIG_ISNULL", codename)

    return
  }
  var translated_bank_code = GetString(pPluginData["bank_code"])
  var translated_home_status int

  switch translated_bank_code {
    case "12":
      translated_bank_code = "Santander"
    case "04":
      translated_bank_code = "Bankia"
    case "11":
      translated_bank_code = "La Caixa"
    case "02":
      translated_bank_code = "Banco Popular"
    case "09":
      translated_bank_code = "Caixa Catalunya"
    case "13":
      translated_bank_code = "UniCaja"
    case "07":
      translated_bank_code = "Caja España"
    case "08":
      translated_bank_code = "Cajamar"
    case "05":
      translated_bank_code = "BBVA"
    case "01":
      translated_bank_code = "Abanca"
    case "10":
      translated_bank_code = "ING Direct"
    case "03":
      translated_bank_code = "Banc Sabadell"
    case "06":
      translated_bank_code = "Caja badajoz"
    default:
      translated_bank_code = "La Caixa"
  }

  switch GetString(pPluginData["home_status"]) {
    case "HOME_OWNER":
      translated_home_status = 2
    case "CO_OWNED":
      translated_home_status = 5
    case "HOSTEL":
      translated_home_status = 6
    default:
      translated_home_status = 4
  }
  var fields = []Pair{{"FirstName", "first_name"},
                      {"FirstLastName", "last_name"},
                      {"SecondLastName", "last_name_2"},
                      {"AppliedAmount", "requested_amount"},
                      {"LoanDuration", "period"},
                      {"SSN", "identity_card_number"},
                      {"Email", "email"},
                      {"MobilePhone", "cell_phone"},
                      {"Birthday", "birth_date"},
                      {"Street", "street"},
                      {"StreetNumber", "house_number"},
                      {"Zip", "zip"},
                      {"City", "city"},
                      {"Province", "province_name"},
                      {"Gender", "gender"},
                    }
  var request_map = map[string]any{}

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      request_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      request_map[GetString(f.A)] = f.B
    }
  }
  request_map["BankNumber"] = pPluginData["iban"]
  request_map["ResidenceType"] = translated_home_status
  request_map["Bank"] = translated_bank_code
  request_map["ApplicantId"] = fmt.Sprintf("%v", strings.Replace(GetString(pPluginData["uid"]), "-", "", -1))
  request_map["CreateTime"] = time.Now().Format("2006-01-02T15:04:05")
  request_map["CallbackUrl"] = "https://www.PARTNER_NAME.cz"
  request_map["Floor"] = rand.Intn(10) + 1
  request_map["apiKey"] = config["api_key" + plugin_postfix]

  if "M" ==  request_map["Gender"] {
    request_map["Gender"] = 1
  } else {
    request_map["Gender"] = 2
  }
  data_map["request"] = request_map

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  var info_map = GetMap(ret["info"])

  if nil == info_map || "accepted" != strings.ToLower(GetString(info_map["status"])) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

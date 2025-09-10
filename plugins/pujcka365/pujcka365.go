package main

import (
  "fmt"
  "log"
  "strings"
  "net/url"

  . "leadz/utils"
)

type leadplugin string

const  codename string = "PUJCKA365"

var configs_map = map[string]string {
  "api_url": "https://www.pujcka365.cz/cz/registrace_PARTNER_NAME",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED"}},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      register_lead,
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
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "register"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  var urlencoded_data = url.Values{}

  for k, v := range data_map {
    if nil == v {
      urlencoded_data.Set(k, "")
    } else {
      urlencoded_data.Set(k, GetString(v))
    }
  }
  pPluginData["urlencoded_data"] = urlencoded_data

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
/*
FIXME TODO *** html parsing TO ADD ***
  // =============================
  fo, err := os.Create(fmt.Sprintf("./response.%v.%v.%v.html", codename, plugin_postfix, pPluginData["uid"]))

  if nil != err {
    log.Printf("%s RESPONSE_WRITE_ERROR: %v\n", pPluginData["plugin_log"], err)
    pPluginData["description"] = err

    return status
  } else {
    defer fo.Close()

    fo.Write(body)
  }
  // =============================
*/

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_bank_code = GetString(pPluginData["bank_code"])
  var fields = []Pair{{"personal_name", "first_name"},
                      {"personal_surname", "last_name"},
                      {"id_number", "identity_card_number"},
                      {"personal_code", "birth_number"},
                      {"email", "email"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }
  switch translated_bank_code {
    case "0800":
      translated_bank_code = "Česká spořitelna a.s."
    case "2010":
      translated_bank_code = "Fio banka, a.s."
    case "5500":
      translated_bank_code = "Raiffeisenbank"
    default:
      translated_bank_code = "Česká spořitelna a.s."
    }
  var cell_phone = GetString(pPluginData["cell_phone"])

  if len(cell_phone) > 8 {
    data_map["phone"]= fmt.Sprintf("%v %v %v", cell_phone[:3], cell_phone[3:6], cell_phone[6:])
  } else {
    data_map["phone"]= cell_phone
  }
  data_map["personal_bank"] = translated_bank_code
  data_map["about_us_feedback"] = "registrace+PARTNER_NAME"
  data_map["checkbox1"] = 1
  data_map["checkbox2"] = 1
  data_map["register_outer"] = "Zažádat o půjčku"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret {
    return
  }
  var body = fmt.Sprintf("%v", ret)

  if code > 299 {
    if strings.Contains(body, "Register") {
      pPluginData["sale_status"] = "DUPLICATE"
    }
    return
  }

  if strings.Contains(strings.ToLower(body), "rejected") {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

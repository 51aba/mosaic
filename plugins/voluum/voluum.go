package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "VOLUUM"

var configs_map = map[string]string {
  "api_url": "",
  "redirect_url": "",
  "username": "",
  "password": "",
}
var plugin_vars = []string{"pdredirect", "clredirect"}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "pdredirect": []func(map[string]any, map[string]string) (bool){
      pd_action,
    },
    "clredirect": []func(map[string]any, map[string]string) (bool){
      cl_action,
    },
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
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func pd_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  if "" == config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = fmt.Sprintf("https://track.lspujcka.cz/81484708-796f-445d-9696-ebbbcad4c6ab?lead_id=%v&aff_id=%v&cascade=pd", pPluginData["uid"], pPluginData["affiliate_name"])
  } else {
    if strings.Contains(config["redirect_url" + plugin_postfix], "%v") {
      pPluginData["redirect_url"] = fmt.Sprintf(config["redirect_url" + plugin_postfix], pPluginData["uid"], pPluginData["affiliate_name"])
    } else {
      pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
    }
  }
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Printf("%v%v%v", YELLOW, pPluginData["description"], NC)

  return true
}

// ################################################################################################################################################################
func cl_action(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  if "" == config["redirect_url" + plugin_postfix] {
    pPluginData["redirect_url"] = fmt.Sprintf("https://track.lspujcka.cz/81484708-796f-445d-9696-ebbbcad4c6ab?lead_id=%v&aff_id=%v&cascade=cl", pPluginData["uid"], pPluginData["affiliate_name"])
  } else {
    if strings.Contains(config["redirect_url" + plugin_postfix], "%v") {
      pPluginData["redirect_url"] = fmt.Sprintf(config["redirect_url" + plugin_postfix], pPluginData["uid"], pPluginData["affiliate_name"])
    } else {
      pPluginData["redirect_url"] = config["redirect_url" + plugin_postfix]
    }
  }
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["description"] = fmt.Sprintf("%v -= REDIRECT =- REDIRECT_URL: %v SALE_STATUS: %v", pPluginData["plugin_log"], pPluginData["redirect_url"], pPluginData["sale_status"])
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
  log.Printf("%v%v%v", YELLOW, pPluginData["description"], NC)

  return true
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": fmt.Sprintf("?token=%v", config["token" + plugin_postfix])}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
                      /*
                      // {"employer", "employer"},
                      // {"employerPosition", "job_title"},
                      // {"employerPhoneNumber", "work_phone"},
                      */
  var fields = []Pair{{"nin", "birth_number"},
                      {"name", "first_name"},
                      {"surname", "last_name"},
                      {"idCardNumber", "identity_card_number"},
                      {"phoneNumber", "cell_phone"},
                      {"email", "email"},
                      {"street", "street"},
                      {"streetNumber", "house_number"},
                      {"city", "city"},
                      {"zipCode", "zip"},
                      {"netIncome", "monthly_income"},
                      {"accountNumber", "bank_account_number"},
                      {"amount", "requested_amount"},
                      {"costs", "monthly_expenses"},
                     }

  for _, f := range fields {
    if nil != pPluginData[GetString(f.B)] {
      data_map[GetString(f.A)] = pPluginData[GetString(f.B)]
    } else {
      data_map[GetString(f.A)] = f.B
    }
  }

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])

  if "_cps" == plugin_postfix {
    if nil != ret["id"] {
      pPluginData["external_id"] = ret["id"]
      result = true
      pPluginData["sale_status"] = "UNCONFIRMED"
    }
  } else {
    log.Printf("%v%v SET_RESPONSE_DATA: DATA: %v%v", YELLOW, pPluginData["plugin_log"], ret["data"], NC)

    if nil == ret["data"] {
      return
    }
    var data_map = GetMap(ret["data"])

    if nil == data_map {
      return
    }
    var rez_val = data_map["result"]
    log.Printf("%v SET_RESPONSE_DATA: RESULT: %v [%v]", pPluginData["plugin_log"], data_map["result"], rez_val)

    if "accepted" != rez_val {
      return
    }
    result = true
    pPluginData["sale_status"] = "UNCONFIRMED"
  } // *** NON CPS ***

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

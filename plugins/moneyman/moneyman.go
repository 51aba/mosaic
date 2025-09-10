package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "github.com/leekchan/timeutil"

  . "leadz/utils"
)

type leadplugin string

const codename string = "MONEYMAN"

var configs_map = map[string]string {
  "api_url": "https://moneyman.es/secure/rest/api/partner/PARTNER_NAME/",
  "username": "",
  "password": "",
}
var plugin_vars = []string{"redirect"}
var validators_map = map[string][]map[string]any {
  "redirect": {
    {},},
  "": {
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 23, "param2": 65},
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 500},
    {"field": "income_type", "func": "DisallowedValuesValidator", "param1": []string{"STATIONARY",}},},
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
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": "register"}
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = data_map

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  data_map["source"] = "PARTNER_NAME"
  data_map["affiliate_id"] = "100036"
  data_map["lead_id"] = pPluginData["uid"]
  data_map["credit"] = map[string]any{
                        "amount": pPluginData["requested_amount"],
                        "period": pPluginData["period"],
                       }
  data_map["contacts"] = map[string]any{
                        "cell_phone": pPluginData["cell_phone"],
                        "email": pPluginData["email"],
                       }
  data_map["personal_data"] = map[string]any{
                        "first_name": pPluginData["first_name"],
                        "last_name": pPluginData["last_name"],
                        "last_name_2": pPluginData["last_name_2"],
                        "gender": pPluginData["gender"],
                        "birth_date": pPluginData["birth_date"],
                        "identity_card_number": pPluginData["identity_card_number"],
                        "marital_status": pPluginData["marital_status"],
                        "dependant_count": pPluginData["dependant_count"],
                       }
  data_map["address"] = map[string]any{
                        "city": pPluginData["city"],
                        "street": pPluginData["street"],
                        "house_number": pPluginData["house_number"],
                        "house_ext": pPluginData["house_ext"],
                        "zip": GetInt(pPluginData["zip"]),
                        "home_phone": pPluginData["home_phone"],
                        "home_status": pPluginData["home_status"],
                       }
  var td = timeutil.Timedelta{Days: time.Duration(30)}

  data_map["employment"] = map[string]any{
                        "education": pPluginData["education"],
                        "income_type": pPluginData["income_type"],
                        "monthly_income": pPluginData["monthly_income"],
                        "specialization": pPluginData["specialization"],
                        // "employer_phone": pPluginData["employer_phone"],
                        "next_pay_date": time.Now().Add(td.Duration()).Format("02-01-2006"),
                        "expenses": pPluginData["monthly_expenses"],
                        "loan_purpose": pPluginData["loan_purpose"],
                      }
  data_map["bank_account_number"] = pPluginData["bank_account_number"]

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }

  if nil != ret["errors"] || strings.Contains(fmt.Sprintf("%v", ret), "rejected") {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = ret["partnerUrl"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

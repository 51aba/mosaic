package main

import (
  "fmt"
  "log"
  "strings"

  . "leadz/utils"
)

type leadplugin string

const codename string = "EXPRESSCASH"

var configs_map = map[string]string{
  "api_url": "https://api.expresscash.cz/Generic/GenericLeadIn.asmx",
  "token": "",
  "CampaignActivityCode": "5611",
}
var plugin_vars = []string{"4470", "5611", }
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "4470": []func(map[string]any, map[string]string) (bool){
      register_lead,},
    "5611": []func(map[string]any, map[string]string) (bool){
      register_lead,},
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
  var plugin_postfix = P_set_postfix(codename, pPluginData, plugin_vars)
  var config map[string]string = P_init_named(codename, pPluginData, configs_map, plugin_postfix)

  pPluginData["headers"] = map[string]string{"Content-Type": "text/xml; charset=utf-8", "SOAPAction": "https://api.expresscash.cz/Generic/GenericLeadIn.asmx/processLead",
                                       "Authorization": fmt.Sprintf("Token %v", config["token" + plugin_postfix])}
  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": ""}

  pPluginData["string_data"] = translate(pPluginData)
  pPluginData["CampaignActivityCode"] = config["CampaignActivityCode" + plugin_postfix]

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any) (ret string) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var translated_income_type = GetString(pPluginData["income_type"])
  var translated_home_status = GetString(pPluginData["home_status"])
  var translated_birth_number = GetString(pPluginData["birth_number"])

  /*
  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%s/%s", translated_birth_number[:6], translated_birth_number[6:])
  }
  */
  /*
  var translated_marital_status = GetString(pPluginData["marital_status"])

  switch translated_marital_status {
    case "MARRIED":
      translated_marital_status = "Vdaná/ženatý"
    case "DIVORCED":
      translated_marital_status = "Rozvedená/ý"
    case "SINGLE":
      translated_marital_status = "Svobodná/ý"
    case "PARTNERSHIP":
      translated_marital_status = "Ve společné domácnosti"
    case "WIDOWED":
      translated_marital_status = "Ovdovělá/ý"
    default:
      translated_marital_status = "Svobodná/ý"
  }o*/

  switch translated_income_type {
    case "EMPLOYED":
      translated_income_type = "Zaměstnaný/á"
    case "SELF_EMPLOYED":
      translated_income_type = "Podnikatel / OSVČ"
    case "MATERNITY_LEAVE":
      translated_income_type = "Mateřská / Rodičovská dovolená"
    case "STUDENT":
      translated_income_type = "Student"
    case "PENSION":
      translated_income_type = "Starobní důchod"
    case "UNEMPLOYED":
      translated_income_type = "Nezaměstnaný/á"
    default:
      translated_income_type = "Ostatní"
  }

  switch translated_home_status {
    case "HOME_OWNER":
      translated_home_status = "Vlastní"
    case "CO_OWNED":
      translated_home_status = "Družstevní"
    case "HOSTEL":
      translated_home_status = "Městský úřad"
    default:
      translated_home_status = "Ostatní"
  }

  ret = fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
                     <soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
                       <soap:Body>
                         <processLead xmlns="https://api.expresscash.cz/Generic/GenericLeadIn.asmx">
                           <requestForm>
                             <firstName>%v</firstName>
                             <surname>%v</surname>
                             <personalNumber>%v</personalNumber>
                             <phone>%v</phone>
                             <email>%v</email>
                             <city>%v</city>
                             <street>%v</street>
                             <postalCode>%v</postalCode>
                             <personalInfoProcessingConsent>true</personalInfoProcessingConsent>
                             <marketingConsent>true</marketingConsent>
                             <personalInfoHandoverConsent>true</personalInfoHandoverConsent>
                             <loanValue>%v</loanValue>
                             <instalmentCount>%v</instalmentCount>
                             <instalmentValue>%v</instalmentValue>
                             <formType></formType>
                             <IncomeSource>%v</IncomeSource>
                             <isCzechNationality>true</isCzechNationality>
                             <personalID>%v</personalID>
                             <usableIncome>%v</usableIncome>
                             <livingStatus>%v</livingStatus>
                             <isCarOwned>true</isCarOwned>
                             <extUID>%v</extUID>
                           </requestForm>
                         </processLead>
                       </soap:Body>
                     </soap:Envelope>`,
                     GetString(pPluginData["first_name"]),
                     GetString(pPluginData["last_name"]),
                     translated_birth_number,
                     "+420" + GetString(pPluginData["cell_phone"]),
                     GetString(pPluginData["email"]),
                     GetString(pPluginData["city"]),
                     GetString(pPluginData["street"]),
                     GetString(pPluginData["zip"]),
                     GetInt(pPluginData["requested_amount"]),
                     1,
                     GetInt(pPluginData["requested_amount"]),
                     translated_income_type,
                     GetString(pPluginData["identity_card_number"]),
                     GetInt(pPluginData["monthly_income"]),
                     translated_home_status,
                     strings.Replace(GetString(pPluginData["uid"]), "-", "", -1))

                     /*
                     GetString(pPluginData["dependent_children"]),
                             <childrenCount>%v</childrenCount>
                     translated_marital_status,
                             <familyStatus>%v</familyStatus>
                     */
  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  var xml_map = GetMap(ret["Envelope"])

  if nil == xml_map {
    xml_map = GetMap(ret["xml"])
  }
  log.Printf("%v%v SET_RESPONSE_DATA: XSL_MAP_ISNULL: %v%v", CYAN, pPluginData["plugin_log"], nil == xml_map, NC)

  if nil == xml_map {
    return
  }
  var result_str string = GetString(GetMap(GetMap(GetMap(xml_map["Body"])["processLeadResponse"])["processLeadResult"])["result"])
  log.Printf("%v SET_RESPONSE_DATA: RESULT: %v", pPluginData["plugin_log"], result_str)

  if "accept" != strings.ToLower(result_str) {
    if "duplicate" == strings.ToLower(result_str) {
      pPluginData["sale_status"] = "DUPLICATE"
    }

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

package main

import (
  "fmt"
  "log"
  "strings"
  "encoding/base64"

  . "leadz/utils"
)

type leadplugin string

const codename string = "PROFICREDIT"

var configs_map = map[string]string {
  "api_url": "https://communication.proficredit.cz:4400/App/Communication.svc",
  "login": "",
  "password": "",
}
var plugin_vars = []string{"razdva", "5611", "mail",}
var validators_map = map[string][]map[string]any {
  "razdva": {
    {"field": "insolvency", "func": "AllowedValuesValidator", "param1": "NO",},
    {"field": "birth_number", "func": "InsolvencyValidator"},
    {"field": "home_status", "func": "DisallowedValuesValidator", "param1": []string{"MINISTRY",}},},
  "5611": {
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 15000},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED", "EMPLOYED", "PENSION", "PART_TIME_EMPLOYMENT",}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 25, "param2": 65},
    {"field": "birth_number", "func": "InsolvencyValidator",},},
  "mail": {
    {"field": "monthly_income", "func": "MinMoneyValidator", "param1": 15000},
    {"field": "income_type", "func": "AllowedValuesValidator", "param1": []string{"SELF_EMPLOYED", "EMPLOYED", "PENSION", "PART_TIME_EMPLOYMENT",}},
    {"field": "birth_date", "func": "AgeRangeValidator", "param1": 25, "param2": 65},
    {"field": "birth_number", "func": "InsolvencyValidator",},},
  "": {
    {"field": "insolvency", "func": "AllowedValuesValidator", "param1": "NO",},
    {"field": "birth_number", "func": "InsolvencyValidator"},
    {"field": "home_status", "func": "DisallowedValuesValidator", "param1": []string{"MINISTRY",}},},
}
var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "razdva": []func(map[string]any, map[string]string) (bool){
      register_lead,
    },
    "5611": []func(map[string]any, map[string]string) (bool){
      register_lead,
    },
    "mail": []func(map[string]any, map[string]string) (bool){
      register_lead,
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
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": ""}
  var basic = base64.StdEncoding.EncodeToString([]byte(config["login" + plugin_postfix] + ":" + config["password" + plugin_postfix]))

  pPluginData["headers"] = map[string]string{"Content-Type": "application/soap+xml; charset=utf-8", "Authorization": fmt.Sprintf("Basic %v", basic),
                                       "SOAPAction": "http://communication.proficredit.cz/2017/01/Communication/Action"}
  pPluginData["string_data"] = translate(pPluginData)

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any) (data_str string) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_birth_number = GetString(pPluginData["birth_number"])

  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%s/%s", translated_birth_number[:6], translated_birth_number[6:])
  }
  data_str = fmt.Sprintf(`<LeadsImportReq xmlns="http://proficredit.cz/Request/LeadsImport/2018/01">
                           <Contacts>
                            <Contact>
                             <LeadID>%v</LeadID>
                             <CampaignActivityCode>%v</CampaignActivityCode>
                             <FirstName>%v</FirstName>
                             <LastName>%v</LastName>
                             <BirthNumber>%v</BirthNumber>
                             <PhoneNumber>%v</PhoneNumber>
                             <Email>%v</Email>
                             <County>%v</County>
                             <City>%v</City>
                             <Street>%v</Street>
                             <HouseNumber>%v</HouseNumber>
                             <PostalCode>%v</PostalCode>
                             <ContactPriorityType>OnlineLoan</ContactPriorityType>
                            </Contact>
                           </Contacts>
                          </LeadsImportReq>`,
                          strings.Replace(GetString(pPluginData["uid"]), "-", "", -1),
                          pPluginData["CampaignActivityCode"],
                          pPluginData["first_name"],
                          pPluginData["last_name"],
                          translated_birth_number,
                          pPluginData["cell_phone"],
                          pPluginData["email"],
                          GetString(pPluginData["district"]),
                          pPluginData["city"],
                          pPluginData["street"],
                          pPluginData["house_number"],
                          pPluginData["zip"],)
  log.Printf("%s TRANSLATOR_DATA: %v\n", pPluginData["plugin_log"], data_str)

  data_str = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
                      <env:Envelope xmlns:env="http://www.w3.org/2003/05/soap-envelope"
                        xmlns:ns1="http://communication.proficredit.cz/ActionRequest/2017/01"
                        xmlns:ns2="http://communication.proficredit.cz/2017/01"
                        xmlns:ns3="http://www.w3.org/2005/08/addressing">
                        <env:Header>
                          <ns3:Action>http://communication.proficredit.cz/2017/01/Communication/Action</ns3:Action>
                          <ns3:To>https://communication.proficredit.cz:4400/App/Communication.svc</ns3:To>
                          <ns3:MessageID>urn:uuid:%v</ns3:MessageID>
                        </env:Header>
                        <env:Body>
                          <ns2:Action><ns2:request>
                            <ns1:Module>Leads</ns1:Module><ns1:Name>LeadsImport</ns1:Name>
                            <ns1:Data>%v</ns1:Data>
                            <ns1:DoFake>0</ns1:DoFake>
                          </ns2:request></ns2:Action>
                        </env:Body>
                      </env:Envelope>`,
                     GetString(pPluginData["uid"]),
                     base64.StdEncoding.EncodeToString([]byte(data_str)),)

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    return
  }
  var xml_map = GetMap(ret["xml"])

  log.Printf("%s SET_RESPONSE_DATA: %v", pPluginData["plugin_log"], xml_map["code"])

  if "0" != GetString(xml_map["code"]) && "accepted" != strings.ToLower(GetString(xml_map["description"])) {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

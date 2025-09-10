package main

import (
  "fmt"
  "log"
  "time"
  "strings"

  . "leadz/utils"
)

type leadplugin string

// openssl pkcs12 -in partner_name.pfx -nocerts -out key.pem -nodes
// openssl pkcs12 -in partner_name.pfx -nokeys -out client.pem
const codename string = "PROVIDENT"

var configs_map = map[string]string {
  "api_url": "https://ows.provident.cz:43008/osc_placebo/service.asmx",
  "login": "",
  "password": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
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
  pPluginData["headers"] = map[string]string{"Content-Type": "text/xml; charset=utf-8", "SOAPAction": "http://providentoz.cz/GetScoringResponse"}
  pPluginData["certificate_path"] = fmt.Sprintf("plugins/%v/assets/", strings.ToLower(codename))

  return P_senddata(codename, pPluginData, configs_map, plugin_vars, is_paused, send_map)
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var call_config = map[string]any{"command": ""}

  pPluginData["string_data"] = translate(pPluginData)

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any) (data_str string) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var translated_bank_account_number string = strings.Replace(strings.Replace(GetString(pPluginData["bank_account_number"]), "000000-", "", -1), "-", "", -1)
  var translated_birth_number string = GetString(pPluginData["birth_number"])

  if len(translated_birth_number) > 6 {
    translated_birth_number = fmt.Sprintf("%s/%s", translated_birth_number[:6], translated_birth_number[6:])
  }

  if len(translated_bank_account_number) > 9 {
    translated_bank_account_number = translated_bank_account_number[:9]
  }
  data_str = fmt.Sprintf(
              `<?xml version="1.0" encoding="UTF-8"?><SOAP-ENV:Envelope xmlns:ns0="http://providentoz.cz/" xmlns:ns1="http://schemas.xmlsoap.org/soap/envelope/"
                xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
                <SOAP-ENV:Header/><ns1:Body><ns0:GetScoringResponse>
                <ns0:CountryCode>CZ</ns0:CountryCode>
                <ns0:CampaignCode>1901</ns0:CampaignCode>
                <ns0:AffiliateCode>1501</ns0:AffiliateCode>
                <ns0:Title></ns0:Title>
                <ns0:RequestedLoanAmount>%v</ns0:RequestedLoanAmount>
                <ns0:RequestedTerm>%v</ns0:RequestedTerm>
                <ns0:Forename>%v</ns0:Forename>
                <ns0:Surname>%v</ns0:Surname>
                <ns0:PersonalNumber>%v</ns0:PersonalNumber>
                <ns0:EmploymentStatus>01</ns0:EmploymentStatus>
                <ns0:HouseNumber>%v</ns0:HouseNumber>
                <ns0:StreetName>%v</ns0:StreetName>
                <ns0:Town>%v</ns0:Town>
                <ns0:Region>%v</ns0:Region>
                <ns0:PostalCode>%v</ns0:PostalCode>
                <ns0:UIRADRID></ns0:UIRADRID>
                <ns0:DateOfBirth></ns0:DateOfBirth>
                <ns0:Gender>%v</ns0:Gender>
                <ns0:MobilePhoneNumber>%v</ns0:MobilePhoneNumber>
                <ns0:HomeTelNumber></ns0:HomeTelNumber>
                <ns0:EmailAddress>%v</ns0:EmailAddress>
                <ns0:BankAccountIndicator>%v</ns0:BankAccountIndicator>
                <ns0:Nationality>0</ns0:Nationality>
                <ns0:AppCode></ns0:AppCode>
                <ns0:traffic_type>affiliate</ns0:traffic_type>
                <ns0:utm_campaign>STANDART</ns0:utm_campaign>
                <ns0:utm_content>%v</ns0:utm_content>
                <ns0:utm_medium>affiliate</ns0:utm_medium>
                <ns0:utm_source>partner_name</ns0:utm_source>
                <ns0:utm_term>STANDART</ns0:utm_term>
                <ns0:RequestedProductCode>4</ns0:RequestedProductCode>
                <ns0:MarketingCommunicationApproval>1</ns0:MarketingCommunicationApproval>
                <ns0:ConsentDate>%v</ns0:ConsentDate>
                <ns0:ConsentSource>https://finero.cz/formular/</ns0:ConsentSource>
                <ns0:PersonalDataProcessingApproval>1</ns0:PersonalDataProcessingApproval></ns0:GetScoringResponse>
                </ns1:Body></SOAP-ENV:Envelope>`,
                GetInt(pPluginData["requested_amount"]),
                GetInt(pPluginData["period"]),
                GetString(pPluginData["first_name"]),
                GetString(pPluginData["last_name"]),
                GetString(pPluginData["birth_number"]),
                GetString(pPluginData["house_number"]),
                GetString(pPluginData["street"]),
                GetString(pPluginData["city"]),
                GetString(pPluginData["district"]),
                GetString(pPluginData["zip"]),
                GetString(pPluginData["gender"]),
                GetString(pPluginData["cell_phone"]),
                GetString(pPluginData["email"]),
                translated_bank_account_number,
                GetString(pPluginData["uid"]),
                time.Now().Format("2006-01-02 15:04:05 -0700"),)

  log.Printf("%s TRANSLATOR_DATA: %v\n", pPluginData["plugin_log"], data_str)
  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])

  return
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], ret == nil, code > 299, NC)

    return
  }
  var envelope_map = GetMap(ret["Envelope"])
  log.Printf("%s SET_RESPONSE_DATA: [1] %v", pPluginData["plugin_log"], envelope_map)

  if nil == envelope_map {
    return
  }
  Pretty(envelope_map)
  /*
  "Body": {
    "GetScoringResponseResponse": {
      "GetScoringResponseResult": {
        "ApplicationReferenceNumber": "32799235",
        "ApplicationStatus": "6",
        "StatusDescription": "Contact accepted",
        "ErrorCode": "0"
      },
  */
  var body_map = GetMap(envelope_map["Body"])
  log.Printf("%s SET_RESPONSE_DATA: [2] %v", pPluginData["plugin_log"], body_map )
  var response_map = GetMap(body_map["GetScoringResponseResponse"])
  log.Printf("%s SET_RESPONSE_DATA: [3] %v", pPluginData["plugin_log"], response_map )
  var result_map = GetMap(response_map["GetScoringResponseResult"])
  log.Printf("%s SET_RESPONSE_DATA: [i4] %v", pPluginData["plugin_log"], result_map)

  if "Contact accepted" != result_map["StatusDescription"] {
    if strings.Contains(GetString(result_map["StatusDescription"]), "Duplicity") {
      pPluginData["sale_status"] = "DUPLICATE"
    }

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = result_map["ApplicationReferenceNumber"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

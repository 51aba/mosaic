package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
  "github.com/bytedance/sonic"

  . "leadz/utils"
)

type leadplugin string

const codename string = "FERRATUM"

var configs_map = map[string]string {
  "api_url": "https://lending-api-ext.sit.ferratum.com",
  "auth_url": "https://auth-server-ext.sit.ferratum.com",
  // "status_url": "https://lending-api-ext.sit.ferratum.com",
  "client_id": "",
  "client_secret": "",
  "hmac": "",
}
var plugin_vars = []string{}
var validators_map = map[string][]map[string]any {
  "": {
    {},},
}

var send_map = SEND_MAP {
  false: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      get_auth_token, register_lead,},
  },
  true: VAR_FUNCS_MAP {
    "": []func(map[string]any, map[string]string) (bool){
      get_auth_token, check_sms,},
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
func get_auth_token(pPluginData map[string]any, config map[string]string) (result bool) {
  config["api_url_default"] = config["api_url"]
  config["api_url"] = config["auth_url"]

  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var call_config = map[string]any{"command": "oauth/token?grant_type=client_credentials"}
  var basic string = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", config["client_id" + plugin_postfix], config["client_secret" + plugin_postfix])))

  pPluginData["string_data"] = ""
  pPluginData["headers"] = map[string]string{"Authorization": fmt.Sprintf("Basic %v", basic)}

  return P_check_unique(pPluginData, call_config, set_response_data_token)
}

// ################################################################################################################################################################
func prepare_data(pPluginData map[string]any, config map[string]string, data []byte, command string) {
  var hash = sha256.Sum256(data)
  var hash_encoded = base64.StdEncoding.EncodeToString(hash[:])
  var headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %v", config["auth_token"])}

  log.Printf("%s DIGEST: DATA: %v HASH: %v", pPluginData["plugin_log"], string(data), hash_encoded)

  headers["Host"] = strings.Replace(config["api_url"], "https://", "", -1)
  headers["Timestamp"] = time.Now().UTC().Format(time.RFC1123)
  headers["X-Request-ID"] = fmt.Sprintf("%v", time.Now().UTC().UnixNano() / 10000)
  headers["Digest"] = fmt.Sprintf("SHA-256=%v", hash_encoded)
  headers["Content-type"] = "application/json"

  create_signature(pPluginData, headers, command)

  pPluginData["headers"] = headers
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  config["api_url"] = config["api_url_default"]

  var command string = "api/v1/loan-applications"
  var call_config = map[string]any{"command": command}
  var data_map = map[string]any{}
  var data []byte
  var err error

  translate(pPluginData, data_map)

  data, err = sonic.Marshal(data_map)

  if nil != err {
    log.Printf("%v%v REGISTER_LEAD: MARSHAL_DATA_ERROR: %v%v", RED, pPluginData["plugin_log"], err, NC)
  }
  prepare_data(pPluginData, config, data, command)

  pPluginData["string_data"] = string(data)

  result = P_register_lead(pPluginData, call_config, set_response_data)

  return
}

// ################################################################################################################################################################
func check_sms(pPluginData map[string]any, config map[string]string) (result bool) {
  config["api_url"] = config["api_url_default"]

  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])

  //=================================================================================================
  var headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %v", config["auth_token"])}
  var command string = fmt.Sprintf("api/v1/applications/%v/status", pPluginData["external_id"])
  var call_config = map[string]any{"command": command}

  headers["Host"] = strings.Replace(config["api_url" + plugin_postfix], "https://", "", -1)
  headers["Content-type"] = "application/json"
  headers["Timestamp"] = time.Now().UTC().Format(time.RFC1123)
  headers["X-Request-ID"] = fmt.Sprintf("%v", time.Now().UTC().UnixNano() / 10000)

  create_signature(pPluginData, headers, command)

  pPluginData["headers"] = headers
  pPluginData["map_data"] = map[string]any{}

  call_config["description"] = "GET_STATUS"

  if result = P_register_lead_get(pPluginData, call_config, set_response_data_get_status); ! result {
    return
  }
  headers["Timestamp"] = time.Now().UTC().Format(time.RFC1123)
  headers["X-Request-ID"] = fmt.Sprintf("%v", time.Now().UTC().UnixNano() / 10000)

  command = fmt.Sprintf("api/v1/applications/%v/offers", pPluginData["external_id"])
  create_signature(pPluginData, headers, command)

  pPluginData["headers"] = headers

  call_config["command"] = command
  call_config["description"] = "GET_OFFERS"

  result = P_register_lead_get(pPluginData, call_config, set_response_data_get_offers)
  //=================================================================================================

  return
}

// ################################################################################################################################################################
func create_signature(pPluginData map[string]any, headers map[string]string, command string) {
  var config = GetMapStrings(pPluginData["config"])

  if nil == config {
    log.Printf("%v LEAD_CONFIG_ISNULL: %v", pPluginData["plugin_log"], codename)

    return
  }
  var headers_data []byte
  var op_method string = "POST"

  if _, hok := headers["Digest"]; ! hok {
    op_method = "GET"
    headers_data = []byte(fmt.Sprintf("timestamp: %v\nhost: %v\n(request-target): get /%v", headers["Timestamp"], headers["Host"], command))
  } else {
    headers_data = []byte(fmt.Sprintf("timestamp: %v\ndigest: %v\nhost: %v\n(request-target): post /%v", headers["Timestamp"], headers["Digest"], headers["Host"], command))
  }
  var h = hmac.New(sha256.New, []byte(config["hmac"]))

  h.Write(headers_data)
  var hash_encoded = base64.StdEncoding.EncodeToString(h.Sum(nil))

  log.Println()
  log.Printf("%v SIGNATURE: PAYLOAD:\n%v HASH: %v\n\n", pPluginData["plugin_log"], string(headers_data), hash_encoded)

  if _, hok := headers["Digest"]; ! hok {
    headers["Signature"] = fmt.Sprintf(`keyId="%v", algorithm="hmac-sha256", headers="timestamp host (request-target)", signature="%v"`, config["client_id"], hash_encoded)
  } else {
    headers["Signature"] = fmt.Sprintf(`keyId="%v", algorithm="hmac-sha256", headers="timestamp digest host (request-target)", signature="%v"`, config["client_id"], hash_encoded)
  }
  pPluginData["description"] = fmt.Sprintf("REQUEST: %v: %s HEADERS: %v", op_method, pPluginData["plugin_log"], headers)
  pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
}

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])
  var translated_home_status (string) = GetString(pPluginData["home_status"])
  var translated_income_type = GetString(pPluginData["income_type"])
  var person = map[string]any{
                "names": []map[string]any{
                          {"primary": true,
                           "firstName": GetString(pPluginData["first_name"]),
                           "lastName": GetString(pPluginData["last_name"])}},
                "birth": map[string]any{
                      "date": GetString(pPluginData["birth_date"])},
                "household": map[string]any{
                      "partnerDwellType": translated_home_status,},
                "employments": []map[string]any{
                          {"primary": true,
                           "partnerType": translated_income_type,
                           "employer": map[string]any{
                                       "name": GetString(pPluginData["employer"]),
                                       "from": ""}}},
                "financialData": map[string]any{
                      "incomes": []map[string]any{
                          {"key": "GROSS_INCOME",
                           "value": map[string]any{
                                   "amount": GetInt(pPluginData["monthly_income"]),
                                   "currency": "CZK"},
                           "period": "MONTH"}},
                      "expenses": []map[string]any{
                          {"key": "TOTAL_EXPENSES",
                           "value": map[string]any{
                                    "amount": int(GetInt(pPluginData["monthly_expenses"])),
                                    "currency": "CZK"},
                           "period": "MONTH"}},
                      "otherLoans": []map[string]any{
                          {"amount": map[string]any{
                                    "amount": 0,
                                    "currency": "CZK"}}}},
              }
  data_map["country"] = "CZ"
  data_map["language"] = "cs"
  data_map["debtors"] = []map[string]any {{
                                    "type": "NATURAL_PERSON",
                                    "role":1,
                                    "person": map[string]any{
                                              "formattedName": GetString(pPluginData["first_name"]) + GetString(pPluginData["last_name"]),
                                              "addresses": []map[string]any{{
                                                            "type": "RESIDENCE",
                                                            "primary": true,
                                                            "country": "CZ",
                                                            "components": []map[string]any {
                                                                          {"type": "ADDRESS_LINE_1",
                                                                           "value": GetString(pPluginData["street"])},
                                                                          {"type": "CITY",
                                                                           "value": GetString(pPluginData["city"])},
                                                                          {"type": "ZIP_CODE",
                                                                           "value": GetString(pPluginData["zip"])},
                                                                          }}},
                                              "identities": []map[string]any{
                                                            {"type": "SSN",
                                                             "primary": true,
                                                             "value": GetString(pPluginData["birth_number"])},
                                                            {"type" : "IDENTITY_PROVIDER",
                                                             "primary" : false,
                                                             "value": GetString(pPluginData["identity_card_number"])}},
                                              "emails": []map[string]any{
                                                            {"primary": true,
                                                             "value": GetString(pPluginData["email"])}},
                                              "phoneNumbers": []map[string]any{
                                                            {"type": "MOBILE",
                                                             "primary": true,
                                                             "value": "+420" + GetString(pPluginData["cell_phone"])}},
                                              "naturalPerson": person,
                                                             }}}
  data_map["applyForLoan"] = map[string]any {
                            "amount": map[string]any{
                                      "amount": GetInt(pPluginData["requested_amount"]),
                                      "currency": "CZK"},
                            "term": map[string]any{
                                      "term": GetInt(pPluginData["period"]),
                                      "termUnit": "MONTH"},
                            "partnerPurpose": "CAR",
                            "preferences": map[string]any{
                                      "oneOfferRequired": true,
                                      "exactMatchOnly": false,
                                      "affiliateMode": "HAND_OVER",
                                      "contractRequired": false},
                            "loanType": "CREDIT_LIMIT"}

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])

  return
}

// ################################################################################################################################################################
func set_response_data_token(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }
  log.Printf("%v SET_RESPONSE_DATA_TOKEN: %v", pPluginData["plugin_log"], ret["access_token"])
  var config = GetMapStrings(pPluginData["config"])

  config["auth_token"] = GetString(ret["access_token"])
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
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil == ret["applicationId"] {
    log.Printf("%s SET_RESPONSE_DATA: MISSING_APPLICATION_ID_ERROR: %v", pPluginData["plugin_log"], ret)

    return
  }
  var callbackOperation = GetMap(ret["callbackOperation"])
  var external_id string = GetString(ret["applicationId"])

  if "" == external_id {
    log.Printf("%s SET_RESPONSE_DATA: APPLICATION_ID_ISNULL_ERROR: %v [ %v ]", pPluginData["plugin_log"], ret["applicationId"], external_id)

    return
  }
  log.Printf("%s SET_RESPONSE_DATA: %v [ %v ]\n%v", pPluginData["plugin_log"], external_id, ret["applicationId"], callbackOperation)
  var plugin_postfix = GetString(pPluginData["plugin_postfix"])
  var config map[string]string = GetMapStrings(pPluginData["config"])
  result = true

  if "" != GetString(config["form_context" + plugin_postfix]) {
    pPluginData["form_context"] = config["form_context" + plugin_postfix]
  } else {
    pPluginData["form_context"] = `{"items": [{"id": 0, "name": "waiting_step", "label": "Waiting Step"}]}`
  }
  pPluginData["sale_status"] = "PAUSED"
  pPluginData["external_id"] = external_id

  return
}

// ################################################################################################################################################################
func set_response_data_get_status(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil == ret["applicationStatus"] {
    log.Printf("%v%v SET_RESPONSE_DATA: STATUS_ISNULL_ERROR:\nRET: %v%v", RED, pPluginData["plugin_log"], ret, NC)

    return
  }
  var status_array = []string{"APPROVED", "PAY_OUT"}

  if ! ArrayStringContains(status_array, GetString(ret["applicationStatus"])) {
    log.Printf("%v%v SET_RESPONSE_DATA: STATUS_REJECTED_ERROR: %v%v", RED, pPluginData["plugin_log"], ret["applicationStatus"], NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"

  return
}

// ################################################################################################################################################################
func set_response_data_get_offers(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], nil == ret, code > 299, NC)

    return
  }

  if nil == ret["offers"] {
    log.Printf("%v%v SET_RESPONSE_DATA: OFFERS_ISNULL_ERROR:\nRET: %v%v", RED, pPluginData["plugin_log"], ret, NC)

    return
  }
  var offers_array = GetArray(ret["offers"])

  if nil  == offers_array || 0 == len(offers_array) {
    log.Printf("%v%v SET_RESPONSE_DATA: OFFERS_ARRAY_ISNULL_ERROR: %v%v", RED, pPluginData["plugin_log"], ret["offers"], NC)

    return
  }
  var offer_map = GetMap(offers_array[0])

  if nil == offer_map || "" == GetString(offer_map["offerUrl"]) {
    log.Printf("%v%v SET_RESPONSE_DATA: OFFER_URL_ISNULL_ERROR: %v%v", RED, pPluginData["plugin_log"], offers_array, NC)

    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["redirect_url"] = GetString(offer_map["offerUrl"])

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

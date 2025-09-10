package main

import (
  "fmt"
  "log"
  "time"
  "strings"
  "github.com/leekchan/timeutil"
  /*
  "golang.org/x/net/html"
  */

  . "leadz/utils"
)

type leadplugin string

const codename string = "NOVACREDIT"

var requested_sum_variants = []int{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 11000, 12000, 13000, 14000, 15000, 16000, 17000, 18000, 19000, 20000}
var requested_sum_variants_fl2 = []int{5000, 10000, 15000, 20000, 25000, 30000, 35000, 40000, 45000, 50000, 55000, 60000}
var period_variants_fl2 = []int{3, 6, 9, 12}

var configs_map = map[string]string {
  "api_url": "https://www.novacredit.cz/loan_calculator/external_form_json",
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
func prepare_data() (ret string) {
  /*

  FIXME *** prepare_data TO ADD ***
        senddata = {}

        senddata["ncId"] = self.ncid
        senddata["action"] = "send-form"
        senddata["form_type"] = data.get("form_type", 2)

        try:
            del data["form_type"]
        except:
            pass

        for k, v in data.iteritems():
            if k in ("loan_cost", "loan_return"):
                senddata["%s" % k] = v
            try:
                senddata["fields[%s]" % k] = unicode(v).encode("utf8")
            except:
                senddata["fields[%s]" % k] = ""
                self.log_warning("WRONG_VALUE: {} => {}".format(k, v))

        if headers:
            response = self._post(url, data, headers=headers)
        else:
            url = url + '?' + urlencode(senddata)
            response = self._post(url)
*/
return
}

// ################################################################################################################################################################
func register_lead(pPluginData map[string]any, config map[string]string) (result bool) {
  var data_map = map[string]any{}

  translate(pPluginData, data_map)

  pPluginData["map_data"] = map[string]any{}
  var call_config = map[string]any{"command": prepare_data()}

  return P_register_lead(pPluginData, call_config, set_response_data)
}

// ################################################################################################################################################################
/*
FIXME *** html parsing TO ADD ***
func call_api(command string, description string, pPluginData map[string]any, data_map map[string]any, plugin_request_timeout time.Duration) (status bool) {
  // =============================
  // var fo *os.File
  // fo, err = os.Create(fmt.Sprintf("./response.%v.%v.%v.html", codename, plugin_postfix, pPluginData["uid"]))
  // if nil != err {
  //   log.Printf("%s RESPONSE_WRITE_ERROR: %v\n", pPluginData["plugin_log"], err)
  //   pPluginData["description"] = err
  //   return status
  // } else {
  //   defer fo.Close()
  //   fo.Write(body)
  // }
  // =============================

  if ! strings.Contains(string(body), "stdClass Object") {
    pPluginData["description"] = fmt.Sprintf("%s RESPONSE_CONTENT_ERROR", pPluginData["plugin_log"])
    pPluginData["sale_logs"] = append(GetArray(pPluginData["sale_logs"]), pPluginData["description"])
    log.Println(pPluginData["description"])

    pPluginData["sale_status"] = "REJECTED"

    return status
  }
  bodylen = len(body)

  if bodylen > 512 {
    bodylen = 512
  }
  var tkn = html.NewTokenizer(strings.NewReader(string(body)))
  var isPre bool

  for {
      var tt = tkn.Next()

      switch {
        case tt == html.ErrorToken:
        case tt == html.StartTagToken:
          var t = tkn.Token()
          isPre = "pre" == t.Data
        case tt == html.TextToken:
          var t = tkn.Token()
          if isPre {
            if strings.Contains(t.Data, "errors") {
              ret["errors"] = t.Data
            }
          }
      }
  }
}
*/

// ################################################################################################################################################################
func translate(pPluginData map[string]any, data_map map[string]any) {
  log.Printf("%v TRANSLATE: STARTED", pPluginData["plugin_log"])

  var plugin_postfix string = GetString(pPluginData["plugin_postfix"])
  var requested_amount int
  var loan_cost int
  var loan_return string

  if strings.Contains(plugin_postfix, "fl2") {
    requested_amount = FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants_fl2)
    var mons int = FindClosest(int(GetFloat(pPluginData["requested_amount"]) / 30.5), period_variants_fl2)

    loan_cost = func(mons int, amount int) (ret int) {
      switch (mons) {
        case 3:
          ret = int(amount / 5000) * 2308
         case 6:
          ret = int(amount / 5000) * 3838
         case 9:
           ret = int(amount / 5000) * 4432
         default:
           ret = int(amount / 5000) * 4912
       }
       return
      }(mons, requested_amount)
    var td = timeutil.Timedelta{Days: time.Duration(int(float64(mons) * 30.5))}

    loan_return = time.Now().Add(td.Duration()).Format("02-01-2006")
  } else {
    requested_amount = FindClosest(GetInt(pPluginData["requested_amount"]), requested_sum_variants)
    var td = timeutil.Timedelta{Days: time.Duration(25)}

    loan_return = time.Now().Add(td.Duration()).Format("02-01-2006")
  }

  var amount string = fmt.Sprintf("%v", requested_amount)
  var fee string = fmt.Sprintf("%v", loan_cost)
  var loan_total string = fmt.Sprintf("%v", requested_amount + loan_cost)

  var account_number string = fmt.Sprintf("%v", pPluginData["bank_account_number"])
  var home_status string = GetString(pPluginData["home_status"])
  var translated_home_status int = 0

  switch home_status {
    case "HOME_OWNER":
      translated_home_status = 1
    default:
      translated_home_status = 0
  }
                      /*
                      // {"employer", "employer"},
                      // {"employer_telephone", "employer_phone"},
                      // {"profession", "job_title"},
                      */
  var fields = []Pair{{"name", "first_name"},
                      {"last_name", "last_name"},
                      {"pin", "birth_number"},
                      {"dni", "identity_card_number"},
                      {"idn", "identity_card_number"},
                      {"city", "city"},
                      {"street", "street"},
                      {"email", "email"},
                      {"zip_code", "zip"},
                      {"mobile", "cell_phone"},
                      {"monthly_net_income", "monthly_income"},
                      {"monthly_expenses", "monthly_expenses_amount"},
                      {"remote_addr", "ip_address"},
                     }

  for _, f := range fields {
    var key = GetString(f.A)

    if nil != pPluginData[GetString(f.B)] {
      data_map[key] = pPluginData[GetString(f.B)]
    } else {
      data_map[key] = f.B
    }
  }
  var config = GetMapStrings(pPluginData["config"])

  data_map["ncId"] = config["api_key"]
  data_map["action"] = "send-form"
  data_map["form_type"] = 2
  data_map["importe"] = amount
  data_map["loan_amount"] = amount
  data_map["loan_cost"] = fee
  data_map["loan_total"] = loan_total
  data_map["loan_return"] = loan_return
  data_map["account_number"] = account_number
  data_map["property_owner"] = translated_home_status
  data_map["motor_vehicle_owner"] = 0
  data_map["terms"] = 1
  data_map["terms_installment"] = 1
  data_map["payment_method"] = "BT"
  data_map["form_submit"] = "1"

  log.Printf("%v TRANSLATE: COMPLETED", pPluginData["plugin_log"])
}

// ################################################################################################################################################################
func set_response_data(code int, pPluginData map[string]any, ret map[string]any, command string) (result bool) {
  if nil == pPluginData || nil == ret || code > 299 {
    log.Printf("%v%v SET_RESPONSE_DATA: INPUT_ERROR: RET_ISNULL: %v CODE_ISERROR: %v%v", RED, pPluginData["plugin_log"], ret == nil, code > 299, NC)

    return
  }
  var data_str string = GetString(ret["data"])

  if strings.Contains(data_str, "errors") {
    log.Printf("%v%v SET_RESPONSE_DATA: DATA_ERRORS%v", RED, pPluginData["plugin_log"], NC)

    return
  }

  if strings.Contains(data_str, "stdClass Object") {
  }

  if nil != ret["errors"] {
    return
  }
  result = true
  pPluginData["sale_status"] = "UNCONFIRMED"
  pPluginData["external_id"] = ret["lcrfid"]

  return
}

// ################################################################################################################################################################
var LeadPlugin leadplugin

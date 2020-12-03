package main

import (
        "bytes"
        "log"
        "os"
        "os/signal"
        "syscall"
        "time"
        "strconv"
        "net/http"
        "sync/atomic"

        "go.uber.org/ratelimit"
)

const ACOUNT_INFO_JSON = `{"account_info":{"i_product":3,"activation_date":"2009-09-18","iso_639_1":"en","iso_4217":"USD","batch_name":"BE_Batch_101","out_date_format":"YYYY-MM-DD","i_account":55,"opening_balance":0.0,"password":"sht4og","has_custom_fields":1,"customer_bill_suspended":0,"blocked":"N","id":"11198700001","out_date_time_format":"YYYY-MM-DD HH24:MI:SS","last_usage":"2011-04-04 10:00:06","h323_password":"6goiwbnm","bill_status":"O","time_zone_name":"Europe/Prague","i_lang":"en","life_time":null,"login":"11198700001","i_role":6,"idle_days":3524,"first_usage":"2009-09-25","balance":3.49865,"is_active":1,"in_date_format":"YYYY-MM-DD","cust_bill_suspension_delayed":0,"i_time_zone":113,"customer_name":"BE_Customer_001","assigned_addons":[],"billing_model":1,"customer_bill_status":"O","customer_blocked":"N","out_time_format":"HH24:MI:SS","password_lifetime":278803363,"service_features":[{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"unified_messaging","attributes":[{"effective_values":["0"],"name":"mailbox_limit","values":["10"]},{"effective_values":["N"],"name":"fax_only_mode","values":["N"]}],"locks":["user"]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"endpoint_redirect","attributes":[],"locks":[]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"rtpp_level","attributes":[],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"individual_routing_plan","attributes":[{"effective_values":[null],"name":"i_routing_plan","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"legal_intercept","attributes":[],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"external_voicemail","attributes":[{"effective_values":[null],"name":"access_number","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"cnam_lookup","attributes":[],"locks":["user"]},{"locked":0,"flag_value":"Y","effective_flag_value":"Y","name":"call_wait_limit","attributes":[],"locks":["user"]},{"locked":0,"flag_value":"7","effective_flag_value":"7","name":"default_action","attributes":[{"effective_values":["30"],"name":"timeout","values":["30"]}],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"phonebook","attributes":[{"effective_values":[null],"name":"favorite_allowed_patterns","values":[null]},{"effective_values":[null],"name":"favorite_change_lock_days","values":[null]},{"effective_values":[null],"name":"max_favorites","values":[null]},{"effective_values":["N"],"name":"enable_abbrev_dial","values":["N"]},{"effective_values":["1"],"name":"abbrev_dial","values":["1"]},{"effective_values":["N"],"name":"favorite_change_allowed","values":["N"]}],"locks":["user"]},{"locked":0,"flag_value":"/","effective_flag_value":"N","name":"sip_static_contact","attributes":[{"effective_values":["N"],"name":"use_tcp","values":["N"]},{"effective_values":[null],"name":"user","values":[null]},{"effective_values":[null],"name":"port","values":[null]},{"effective_values":[null],"name":"host","values":[null]}],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"clir","attributes":[{"effective_values":["N"],"name":"blocked","values":["N"]},{"effective_values":[null],"name":"clir_note","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"auto_attendant","attributes":[{"effective_values":[null],"name":"auto_attendant_note","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"msg_service_policy","attributes":[],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"voice_service_policy","attributes":[{"effective_values":[null],"name":"id","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"distinctive_ring_vpn","attributes":[],"locks":[]},{"locked":0,"flag_value":"N","effective_flag_value":null,"name":"netaccess_sessions","attributes":[{"effective_values":["1"],"name":"max_sessions","values":["1"]}],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":"Y","name":"sip_dynamic_registration","attributes":[],"locks":["user"]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"cli","attributes":[{"effective_values":["N"],"name":"display_number_allow_external","values":[]},{"effective_values":[null],"name":"centrex","values":[null]},{"effective_values":[null],"name":"display_number","values":[null]},{"effective_values":["Y"],"name":"display_number_check","values":["Y"]},{"effective_values":["N"],"name":"display_name_override","values":["N"]},{"effective_values":["A"],"name":"attest","values":["A"]},{"effective_values":[null],"name":"account_group","values":[null]},{"effective_values":[null],"name":"display_name","values":[null]}],"locks":[]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"voice_roaming_protection","attributes":[],"locks":[]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"voice_fup","attributes":[],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"user_ivr_application","attributes":[{"effective_values":[null],"name":"max_numbers","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"voice_authentication","attributes":[{"effective_values":[null],"name":"pin","values":[null]}],"locks":[]},{"locked":0,"flag_value":"^","effective_flag_value":"Y","name":"music_on_hold","attributes":[{"effective_values":["1"],"name":"i_moh","values":["1"]}],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"emergency","attributes":[{"effective_values":[null],"name":"emergency_administrative_unit","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"Y","effective_flag_value":"Y","name":"clip","attributes":[],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":null,"name":"netaccess_static_ip","attributes":[{"effective_values":[],"name":"routed_network","values":[]},{"effective_values":[null],"name":"address","values":[null]},{"effective_values":[null],"name":"netmask","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"forward_mode","attributes":[{"effective_values":[null],"name":"max_forwards","values":[null]},{"effective_values":["N"],"name":"dtmf_control","values":["N"]}],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"voice_pass_through","attributes":[{"effective_values":[null],"name":"outgoing_access_number","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"call_barring","attributes":[{"effective_values":[],"name":"call_barring_rules","values":[]}],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":null,"name":"iptv","attributes":[{"effective_values":[],"name":"service_packages","values":[]},{"effective_values":[],"name":"channel_packages","values":[]},{"effective_values":[null],"name":"activation_pin","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"^","effective_flag_value":"N","name":"cli_trust","attributes":[{"effective_values":["N"],"name":"accept_caller","values":["N"]},{"effective_values":["N"],"name":"supply_caller","values":["N"]}],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":null,"name":"wifi_speed_limit","attributes":[{"effective_values":[null],"name":"tx_rate","values":[null]},{"effective_values":[null],"name":"rx_rate","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"call_processing","attributes":[],"locks":["user"]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"ringback_tone","attributes":[{"effective_values":[null],"name":"i_ringback_tone","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"sim_calls_limit","attributes":[],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":null,"name":"netaccess_policy","attributes":[{"effective_values":[null],"name":"access_policy","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"/","effective_flag_value":"N","name":"voice_dialing","attributes":[{"effective_values":["N"],"name":"translate_cli_out","values":["N"]},{"effective_values":[null],"name":"i_dial_rule","values":[null]},{"effective_values":["N"],"name":"translate_cli_in","values":["N"]},{"effective_values":["N"],"name":"translate_cld_in","values":["N"]}],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":"N","name":"cps_account_limit","attributes":[{"effective_values":["normal"],"name":"burst_method","values":["normal"]},{"effective_values":[null],"name":"cps_max","values":[null]}],"locks":["user"]},{"locked":0,"flag_value":"/","effective_flag_value":"N","name":"voice_location","attributes":[{"effective_values":[null],"name":"allow_roaming","values":[null]},{"effective_values":[null],"name":"primary_location","values":[null]},{"effective_values":[null],"name":"primary_location_data","values":[null]},{"effective_values":[null],"name":"emergency_administrative_unit","values":[null]}],"locks":[]},{"locked":0,"flag_value":"~","effective_flag_value":null,"name":"sms_routing","attributes":[{"effective_values":["N"],"name":"update_netnumber_db","values":["N"]}],"locks":["user"]},{"locked":0,"flag_value":"N","effective_flag_value":null,"name":"conf_enabled","attributes":[{"effective_values":["64"],"name":"max_participants","values":["64"]}],"locks":["user"]},{"locked":0,"flag_value":"N","effective_flag_value":"N","name":"associated_number","attributes":[{"effective_values":[null],"name":"redirect_number","values":[null]}],"locks":[]},{"locked":0,"flag_value":"Y","effective_flag_value":"Y","name":"lan_name","attributes":[{"effective_values":["en"],"name":"iso_639_1","values":["en"]}],"locks":[]}],"ecommerce_enabled":"N","i_account_role":1,"i_batch":7,"in_time_format":"HH24:MI:SS","first_usage_time":"2009-09-25 02:00:00","password_timestamp":"2012-01-26 12:47:51","expiration_date":"2023-01-13","service_flags":" ^^^~NNN ^YN~ ~ N","included_services":[],"issue_date":"2009-09-18","control_number":1,"i_customer":10,"i_acl":155,"inactivity_expire_time":null,"credit_limit":null,"product_name":"BE_Product_001"}}`


func main() {
        if len(os.Args) < 2 {
                log.Printf("use %s URL rps", os.Args[0])
                return
        }

        service_url := os.Args[1]
        rps := 0
        if len(os.Args) > 2 {
                r, err := strconv.Atoi( os.Args[2] )
                if err != nil {
                        log.Printf("cannot parse rps")
                        return
                }
                rps = r

        }

        client := &http.Client{}

        finished := false
        os_chan := make(chan os.Signal)
        signal.Notify(os_chan, os.Interrupt, syscall.SIGTERM)
        go func() {
                <-os_chan
                log.Println("\n Ctrl+C pressed in Terminal")
                finished = true
        }()

        var cnt uint64
        var errors uint64
        err_cnt := make(map[error]uint64)

        go (func() {
            var prev_c uint64
            var prev_e uint64
            for {
                time.Sleep(1 * time.Minute)
                local_cnt := atomic.LoadUint64(&cnt)
                local_err := atomic.LoadUint64(&errors)
                log.Printf("Send RPS: %.2f (%d -> %d)\n", float64(local_cnt - prev_c) / 60.0, prev_c, local_cnt)
                if local_err > prev_e {
                        log.Printf("  Errors: %d", local_err - prev_e)
                        for k, v := range err_cnt {
                                log.Printf("    %s : %s", k, v)
                        }
                }
                prev_c = local_cnt
                prev_e = local_err
            }
        })()

        log.Printf("Start sending requests with %d Kb payload to %s", len(ACOUNT_INFO_JSON) / 1024, service_url)
        
        var rl ratelimit.Limiter
        if rps > 0 {
                log.Printf("Generating load with RPS %d\n", rps)
                rl = ratelimit.New( rps )
        }

        for {
                if finished {
                        break
                }

                if rl != nil {
                        rl.Take()
                }

                req, err := http.NewRequest("POST", service_url, bytes.NewBufferString(ACOUNT_INFO_JSON))
                req.Header.Add("Content-Type", "application/json")
                now := time.Now().UnixNano()
                req.Header.Add("X-Queue-Time", strconv.FormatInt( now, 10 ))
                req.Header.Add("X-Request-Time", strconv.FormatInt( now, 10 ))

                resp, err := client.Do(req)
                if err != nil {
                        log.Printf("Error sending request: %s", err)
                        atomic.AddUint64(&errors, 1)
                        err_cnt[err]++
                        continue
                }
                resp.Body.Close()

                atomic.AddUint64(&cnt, 1)
        }

        log.Printf("Finished");

        return
}

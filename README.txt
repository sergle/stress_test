Benchmark Teporal.io performance

generator - creates workflows with piece of data

worker
  workflow_SingleSendEvent - read the event data, calls the [local] activity to send data to external http service
  activity_SendHTTPTime    - perform POST request with test payload

ext_app - http service mock, measure incoming requests rate

ext_app_test - basic HTTP load generator for ext_app, to check max available rate

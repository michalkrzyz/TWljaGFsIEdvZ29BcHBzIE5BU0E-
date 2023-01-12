*** Settings ***
Library    DateTime
Library    JSONLibrary
Library    RequestsLibrary

*** Variables ***
${INVALID_DATE_FORMAT}    10-11-2020
${FROM_DATE}    2022-10-10
${TO_DATE}    2022-10-20

${URL_COLLECTOR_PORT}    8090
${URL_COLLECTOR_URL}    http://localhost:${URL_COLLECTOR_PORT}

${LAST_HTTP_RESPONSE}

*** Keywords ***
Client GETs health
    ${response} =    GET    ${URL_COLLECTOR_URL}/health    expected_status=Any
    Set test variable    ${LAST_HTTP_RESPONSE}    ${response}

Client GETs pictures with correct date range parameters
    ${parameters} =    Create dictionary    from=${FROM_DATE}    to=${TO_DATE}
    GET pictures    ${parameters}

Client GETs pictures with missing 'from' parameter
    ${parameters} =    Create dictionary    to=${TO_DATE}
    GET pictures    ${parameters}

Client GETs pictures with missing 'to' parameter
    ${parameters} =    Create dictionary    from=${FROM_DATE}
    GET pictures    ${parameters}

Client GETs pictures with invalid 'from' parameter format
    ${parameters} =    Create dictionary    from=${INVALID_DATE_FORMAT}    to=${TO_DATE}
    GET pictures    ${parameters}

Client GETs pictures with invalid 'to' parameter format
    ${parameters} =    Create dictionary    from=${FROM_DATE}    to=${INVALID_DATE_FORMAT}
    GET pictures    ${parameters}

Client GETs pictures with 'to' parameter before 'from' parameter
    ${parameters} =    Create dictionary    from=${TO_DATE}    to=${FROM_DATE}
    GET pictures    ${parameters}

GET pictures
    [Arguments]    ${params}
    ${response} =    GET
    ...    ${URL_COLLECTOR_URL}/pictures
    ...    params=${params}
    ...    expected_status=Any
    Set test variable    ${LAST_HTTP_RESPONSE}    ${response}

Client POSTs pictures
    ${parameters} =    Create dictionary    from=${TO_DATE}    to=${FROM_DATE}
    POST pictures    ${parameters}

POST pictures
    [Arguments]    ${params}
    ${response} =    POST
    ...    ${URL_COLLECTOR_URL}/pictures
    ...    params=${params}
    ...    expected_status=Any
    Set test variable    ${LAST_HTTP_RESPONSE}    ${response}

HTTP response shall be ${expected_status_code}
    Should be equal as integers
    ...    ${LAST_HTTP_RESPONSE.status_code}
    ...    ${expected_status_code}
    ...    msg=Unexpected HTTP status: ${LAST_HTTP_RESPONSE.status_code}.

Json with error detail message about missing '${parameter}' parameter shall be received
    ${msg} =    Get value from json    ${LAST_HTTP_RESPONSE.json()}    error    fail_on_empty=True
    Should be equal as strings    ${msg[0]}    Missing '${parameter}' URL Query parameter

Json with error detail message about invalid '${parameter}' parameter shall be received
    ${msg} =    Get value from json    ${LAST_HTTP_RESPONSE.json()}    error    fail_on_empty=True
    Should be equal as strings    ${msg[0]}    Invalid range parameters: '${parameter}' parameter format is invalid. Should be YYYY-MM-DD.

Json with error detail message about invalid range shall be received
    ${msg} =    Get value from json    ${LAST_HTTP_RESPONSE.json()}    error    fail_on_empty=True
    Should be equal as strings    ${msg[0]}    Invalid range parameters: 'from' date is not before 'to'

Response json shall return valid number of URLs
    ${time_delta} =    Subtract date from date
    ...    ${TO_DATE}
    ...    ${FROM_DATE}
    ...    result_format=timedelta
    ...    exclude_millis=True
    ...    date1_format=%Y-%m-%d
    ...    date2_format=%Y-%m-%d
    ${number_of_days} =    Set variable    ${time_delta.days}
    ${number_of_days} =    Evaluate    ${number_of_days} + 1
    ${msg} =    Get value from json    ${LAST_HTTP_RESPONSE.json()}    urls    fail_on_empty=True
    Length should be    ${msg[0]}    ${number_of_days}

Nasa API shall be available
    ${parameters} =    Create dictionary    api_key=DEMO_KEY    date=2019-12-06
    ${response} =    GET
    ...    https://api.nasa.gov/planetary/apod
    ...    params=${parameters}
    HTTP response shall be 200

*** Test Cases ***
GET /health shall be successful
    When Client GETs health
    Then HTTP response shall be 200

Get /pictures from correct range of dates shall be successful
    When Client GETs pictures with correct date range parameters
    Then HTTP response shall be 200
     And Response json shall return valid number of URLs

GET /pictures with missing 'from' URL query parameter shall respond with error
    When Client GETs pictures with missing 'from' parameter
    Then HTTP response shall be 422
     And Json with error detail message about missing 'from' parameter shall be received

GET /pictures with missing 'to' URL query parameter shall respond with error
    When Client GETs pictures with missing 'to' parameter
    Then HTTP response shall be 422
     And Json with error detail message about missing 'to' parameter shall be received

GET /pictures with invalid format of 'from' URL query parameter shall respond with error
    When Client GETs pictures with invalid 'from' parameter format
    Then HTTP response shall be 400
     And Json with error detail message about invalid 'from' parameter shall be received

GET /pictures with invalid format of 'to' URL query parameter shall respond with error
    When Client GETs pictures with invalid 'to' parameter format
    Then HTTP response shall be 400
     And Json with error detail message about invalid 'to' parameter shall be received

GET /pictures with 'to' before 'from' parameter shall respond with error
    When Client GETs pictures with 'to' parameter before 'from' parameter
    Then HTTP response shall be 400
     And Json with error detail message about invalid range shall be received

POST /pictures shall respond with method not allowed error
    When Client POSTs pictures
    Then HTTP response shall be 405

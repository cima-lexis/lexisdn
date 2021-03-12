#!/bin/bash

TOKEN=`curl -s --location --request POST 'https://testauth.cimafoundation.org/auth/realms/webdrops/protocol/openid-connect/token' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'client_id=webdrops' \
    --data-urlencode 'grant_type=password' \
    --data-urlencode 'password=^8J*ITws38Cd4b5Cg*g%iSni!KqMPH' \
    --data-urlencode 'username=andrea.parodi@cimafoundation.org' | jq -r '.access_token' -
`
echo login done
ID_LIST=`curl -s --location --request GET 'http://webdrops.cimafoundation.org/app/sensors/list/TERMOMETRO/?from=202103092330&to=202103100030&stationgroup=DewetraWorld%25WunderEurope' \
    --header "Authorization: Bearer $TOKEN" | jq '[.[] | .id]'
`
echo "{\"sensors\": $ID_LIST}" | \
curl -s --location --request POST 'http://webdrops.cimafoundation.org/app/sensors/data/TERMOMETRO/?from=202103092330&to=202103100030&aggr=60' \
    --header "Authorization: Bearer $TOKEN" \
    --header 'Content-Type: application/json' \
    --data @- \
| jq . > TERMOMETRO.json
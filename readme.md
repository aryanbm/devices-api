##Test API
####GetDevice:
`
curl -X GET https://xhu35w9l46.execute-api.us-east-1.amazonaws.com/dev/api/devices/id1 
`

------------

####CreateDevice:
`
curl -X POST https://xhu35w9l46.execute-api.us-east-1.amazonaws.com/dev/api/devices/ -H "Content-Type: application/json" --data-binary @- <<DATA
{
    "id": "/devices/id1",
    "deviceModel": "/devicesmodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
}
DATA
`

------------

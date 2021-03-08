# Test API
#### GetDevice:

```curl -X GET https://<api-gateway-url>/api/devices/id1 ```



#### CreateDevice:
```
curl -X POST https://<api-gateway-url>/api/devices/ -H "Content-Type: application/json" --data-binary @- <<DATA
{
    "id": "/devices/id1",
    "deviceModel": "/devicesmodels/id1",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
}
DATA
```

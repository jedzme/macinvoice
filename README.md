
# MacInvoice Parser


## Webservice API Endpoints


### Download and Parse
#### This will parse the CSV that will be downloaded from an external source: HTTP/REST source, etc.

#### Sample Request
##### Method: POST
##### Content-Type: application/json
```

{
    "url": "https://my.callmanager.tel/portlet/export/voip/autoconfigDeviceList.html?6578706f7274=1&vs=7&d-3147565-e=1",
    "name": "CALLAMANGER",
    "cookie": "JSESSIONID=node0e5v7s15l7nf5nj0p5t7fn3tz219300.node0",
    "authorization": "Basic <authorization key>"
}
```

#### Response
```
{
    "message": "success"
}
```

### Upload and Parse
#### This will parse the CSV that will be uploaded locally through multipart/form-data.
####  Sample Request
##### Method: POST
##### Content-Type: multipart/form-data
| KEY    | VALUE                    | Type |
|--------|--------------------------|------|
| server | callamanger              | text |
| file   | autoconfigDeviceList.csv | file |
|        |                          |      |

####  Response
```
{
    "message": "success"
}
```
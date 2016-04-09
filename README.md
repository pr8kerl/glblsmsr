# glblsmsr

A cmd line client to send SMS messages via [SMS Global](http://smsglobal.com) REST API.

Flags:

```
  -d, --debug=false: enable debug
  -f, --from="smsglobal": string identifier of who is sending the msg
  -m, --message="": the message to be sent
  -n, --number="": the mobile phone number to send to
```

## Build

* set GOROOT in go.env to point to your go installation
* source go.env
* make update (to get required dependencies)
* make

## Run

To use this you need an API key and Secret available from your smsglobal account.
Set the following two environment variables:

```
GLBLKEY=<your api key>
GLBLSCRT=<your secret>
```

Once set, run the executable:

```
$ ./glblsmsr -n 61555555555 -m "schön Tag noch"
OK 201 Created : /v1/sms/1225870233/
$
```

You can also pipe the message in via stdin if you like

```
$ echo "schön Tag noch" | ./glblsmsr -n 6155555555
201 Created /v1/sms/1759071884
```

If you need to use a proxy, set the HTTP_PROXY environment variable appropriately.


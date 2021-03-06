# infradb
A simple API for sqlite in GO</br>

# Examples:</br>
> go build</br>
> infradb -db hosts.db -port 8080 </br>

At this point the server is listening over port 8080.</br>
Specify an IP or hostname using the -ip flag. Default is localhost. </br>

> curl -i -GET http://localhost:8080/ 
>
> curl -i -GET http://localhost:8080/all

## Insert a new host into the DB:

> curl -i -X POST http://localhost:8080/host -d '</br>
{</br>
  &emsp;  "Id": 1, </br>
  &emsp;  "hostname": "ovios",</br>
  &emsp;  "ip": "192.168.12.234",</br>
  &emsp;  "os": "ovios linux 3.12",</br>
  &emsp;  "kernel": "linux 3.x",</br>
  &emsp;  "environment": "PROD",</br>
  &emsp;  "is_vm": true</br>
}'

## Update a host by ID:
> curl -i -X PUT http://localhost:8080/host/1 -d '</br>
{</br>
&emsp;  "Id": 1,</br>
&emsp;  "hostname": "ovios",</br>
&emsp;  "ip": "192.168.1.102",</br>
&emsp;  "os": "ovios linux 3.12",</br>
&emsp;  "kernel": "linux 5.x",</br>
&emsp;  "environment": "PROD",</br>
&emsp;  "is_vm": true</br>
}'

## Return an entry by ID:
> curl -i -GET http://localhost:8080/host/1</br>
HTTP/1.1 200 OK</br>
Date: Thu, 10 Feb 2022 17:24:28 GMT</br>
Content-Length: 156</br>
Content-Type: text/plain; charset=utf-8</br>
{</br>
&emsp;  "id": 1,</br>
&emsp;  "hostname": "ovios",</br>
&emsp;  "ip": "192.168.1.102",</br>
&emsp;  "os": "ovios linux 3.12",</br>
&emsp;  "kernel": "linux 5.x",</br>
&emsp;  "environment": "PROD",</br>
&emsp;  "is_vm": true</br>
}</br>

## Delete an entry by ID:

> curl -i -X DELETE http://localhost:8080/host/1


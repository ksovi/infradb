# infradb
A simple API for sqlite in GO
Examples:
> go build
> infradb.exe -db hosts.db -port 8080 

At this point the server is listening to port 8080.

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
> curl -i -GET http://localhost:8080/host/1

## Delete an entry by ID:

> curl -i -X DELETE http://localhost:8080/host/1


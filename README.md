# check_azure
Nagios/Icinga like check to monitor Microsoft azure objects

To use this check you need some keys which are a little hidden if you are new to the subject:
- clientId
- clientSecret
- subscriptionId
- tenantId

The first two are identifying your application. You have to register the app because it is using Microsoft OAuth service. Here is a short documentation: https://msdn.microsoft.com/de-de/library/bb676626.aspx
The third one is the key which identifies your azure abo. If your on the azure website it is basically everywhere, here is one place to find it: https://blogs.msdn.microsoft.com/mschray/2015/05/13/getting-your-azure-guid-subscription-id/
The last one also not too hard to find, here is one way to do it: http://merill.net/2015/01/how-to-get-the-azure-ad-tenant-id-without-powershell/

At last you have to give your new registered application the right to read the data in azure.

Now the hard stuff is over ;)

Currently supported are
- Subscriptions
- Classic Virtual Machines

 - CPU
 - Network
 - Disk

more to come if Microsoft opens up... The REST API is very limited at the moment.

Be careful there is a direct context between subscriptionId -> resourceGroup -> name. You could have multiple subscriptionIds, where every id contains multiple groups with multiple objects.

```
go run check_azure.go --clientId XXX --clientSecret XXX --subscriptionId XXX --tenantId XXX mode cvm cpu --resourceGroup test --name Ubuntu
OK - Percentage CPU last checked: 2016-11-23 09:15:00 +0100 CET|'Percentage CPU'=0.177125%;80;90;; 

go run check_azure.go --clientId XXX --clientSecret XXX --subscriptionId XXX --tenantId XXX mode cvm network --resourceGroup test --name Ubuntu
OK - Network|'Network In'=26125B;100000;200000;; 'Network Out'=16444B;100000;200000;; 
Network In last checked: 2016-11-23 09:15:00 +0100 CET
Network Out last checked: 2016-11-23 09:15:00 +0100 CET

go run check_azure.go --clientId XXX --clientSecret XXX --subscriptionId XXX --tenantId XXX mode cvm disk --resourceGroup test --name Ubuntu
Critical - Disk|'Disk Read Bytes/sec'=0Bps;100;200;; 'Disk Write Bytes/sec'=887.562252Bps;100;200;; 
Disk Read Bytes/sec last checked: 2016-11-23 09:15:00 +0100 CET
Disk Write Bytes/sec last checked: 2016-11-23 09:15:00 +0100 CET
```

This page is very useful to search the API: https://resources.azure.com/

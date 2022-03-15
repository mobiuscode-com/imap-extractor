
# Imap Extractor

This is a very basic tool to scan an email-inbox for certain content.

The parameters are passed through a configuration file in json format, an example can be seen in below.
The path to said configuration shall be passed as first parameter of the program.

Disclaimer: This tool was just a quick implementation for our CI needs, which might also be useful for others.


##Usage:
To call the program use:
```
imap-extractor <path to config json>
```

For detailed description of config json see below. Configured regex needs to contain at least one group.
Result will be content of group of first match found in the email inbox.
Additional unnamed groups may be used inside the regex.

##Parameters:

The following configurations are to be included in the configuration:
* imap-host: Url on which the email host can be reached
* imap-port: Port on which the IMAP protocol is offered (Usually 143 or 993)
* username: Login username for email host
* password: Password for above username
* from-filter: Name filter for email origin
* regexp: Regex for which the emails will be scanned, including a group for the result


On execution this program will go through the inbox of given email address from newest to oldest.
Once any match with the given regex is found, the content of the first capturing group in the regex will be returned. 


##Config Example:
```json
{
  "imap-host": "mobiuscode.de/",
  "imap-port": 993,
  "username": "user@mobiuscode.de",
  "password": "tryOutPanicMode",
  "from-filter": "boss@mobiuscode.de",
  "regexp": "please see important thing below:(?:[\\s]+)([\\S]+)(?:[\\s]+)"
}
```
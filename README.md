
## Imap Extractor

This is a very basic tool to scan an email-inbox for certain content.

The parameters are passed through a configuration file in json format, an example can be seen in "example-config.json".
The path to said configuration shall be passed as first parameter of the program.

The following configurations are to be included in the configuration:
* imap-host: Url on which the email host can be reached
* imap-port: Port on which the IMAP protocol is offered
* username: Login username for email host
* password: Password for above username
* from-filter: Name filter for email origin
* regexp: Regex for which the emails will be scanned, including a group for the result


On execution this program will go through the inbox of given email address from newest to oldest.
Once any match with the given regex is found, the content of the first capturing group in the regex will be returned. 
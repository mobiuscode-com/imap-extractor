
# Imap Extractor

This is a very basic tool to scan an email-inbox for certain content.

The parameters are passed through a configuration file in json format, an example can be seen in below.
The path to said configuration shall be passed as first parameter of the program.

Disclaimer: This tool was just a quick implementation for our CI needs, which might also be useful for others.


## Usage:
To install the program use:
```
go install github.com/mobiuscode-de/imap-extractor@latest
```
or download a binary from the latest [github-release](https://github.com/mobiuscode-de/imap-extractor/releases).


To call the program use:
```
imap-extractor <path to config json>
```

For detailed description of config json see below. Configured regex needs to contain at least one group.
Result will be content of group of first match found in the email inbox.
Additional unnamed groups may be used inside the regex.

## Usage example:

Given you have an email looking like this:
```
Dear Email recepient,
I am sending you this email with random content. 
It contains a super important code later on though.

For future reference please see:

ImportantCode42

Make sure you keep that important code.
Kind regards
```

A configuration to extract the code detailed in the email could look something like this:

```json
{
  "imap-host": "mobiuscode.de",
  "imap-port": 993,
  "username": "user@mobiuscode.de",
  "password": "$EMAIL_PW",
  "from-filter": "boss@mobiuscode.de",
  "regexp": "please see:(?:[\\s]+)([\\S]+)(?:[\\s]+)"
}
```

Note: Password is extracted from environment variable in this case, which is also recommended for usage of the tool in a CI environment.

Above configuration can then be passed to the tool to go through the email inbox configured:

```
imap-extractor imap-config.json
> ImportantCode42
```

The tool will prioritize the latest emails found and will only return the first match.



## Parameters:

The following configurations are to be included in the configuration:
* imap-host: Url on which the email host can be reached
* imap-port: Port on which the IMAP protocol is offered (Usually 143 or 993)
* username: Login username for email host
* password: Password for above username
* from-filter: Name filter for email origin
* regexp: Regex for which the emails will be scanned, including a group for the result

Any value might be preceded by a $ Symbol to indicate that its value shall be fetched by an environment variable. 

On execution this program will go through the inbox of given email address from newest to oldest.
Once any match with the given regex is found, the content of the first capturing group in the regex will be returned. 


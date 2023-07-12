# cloudflare API Routines
## Introduction
This repositary mainly contains routines that access the cloudflare api to changecloudflare records

### cfTokenFile

Each token is stored in a yaml file with additional information.  

## Programs

### listAccount

### listAccounts

### listDomainsShort

A program that lists all the domains (zones) controlled by the user and stores the domains as a yaml file 
with the name cfDomainsShort.yaml in the folder zoneDir.  

### listDomains

### creTokenFile

A program that creates a cfToken file and stores the file in the folder token.  

### readCfToken

A program that reads a cfToken file and prints the output.  

### verifyToken

A program that verifies the validity of a cloudflare api token.  

### creDnsToken

A Program that creates a DnsToken and stores the token in a token file.  

### listTokens

A program that lists all Cloudflare token and prints the details of each token's scope.  

### listUserInfo

### addDnsRec


### delDnsRec

### listDnsREc

### findApiPerms



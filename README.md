# Apartment notifier
## Introduction
Apartment hunting in Amsterdam can be very difficult

This is a simple html parser for the website pararius.com where the user gets notified of newly published apartment data via OS level notifications

## Dependencies
1. Install go
2. Install redis server

## Limitations
I've hardcoded the below:
- redis credentials, change in code(`redis.go`) if required
- Apartment searching URLs
- Parsing only supported for pararius website search page

# Run
```
go mod download
go run .
```

# What this does

1. Parse existing page
2. Get the top apartment - ap1
3. If ap1 not the same as the top apartment(initially empty string)
	1. Notify user of new apartments
	2. Update the latest apartment as the top apartment
4. Repeat 1 - 3
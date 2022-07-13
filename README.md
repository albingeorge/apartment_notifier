# Apartment notifier
## Introduction
Apartment hunting in Amsterdam can be very difficult

This is a simple html parser for the website pararius.com where the user gets notified of newly published apartment data via OS level notifications

## Dependencies
1. Install go
2. Install redis server

## Limitations
I've hardcoded the redis credentials, change in code(`redis.go`) if required

# Run
`go run .`

# What this does
1. Parse existing page
2. Get the top apartment - ap1
3. Compare ap1 to the top apartment(initially empty string)
3. If not same
    3.1 Notify user of new apartments
    3.2 Update the latest apartment as the top apartment
5. Repeat 1 - 5
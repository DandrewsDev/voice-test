This repo is for the purpose of testing and building out the ability to pull audio data from CS2 Demos.

There are currently 2 files, 
1. One written by me for pulling audio data for testing. pull-audio.go
2. And voice-test.go written by a second dev attempting to pull and process the audio data into a .wav audio file.

Converting the audio is not currently working.


Getting setup will require having golang installed, and getting the dependencies. 
`go get ./...`
Then files can be run `go run pull-audio.go`

You'll need a local demo file for testing. One is provided here https://s3.dandrews.net/public/1-48050850-ad9e-4497-bf8b-d20a7d7cb232.dem
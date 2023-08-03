# Hydration Calculator
Small browser-based tool to calculate how much water a streamer should've drunk based on their Twitch stream uptime.

To start, run `go run main.go` 
By default it will listen on port `8080`

To get a streamer's hydration calculation, add their channel username as a GET request through the URL with the parameter `channel`

Example:
`localhost:8080/?channel=CHANNELNAME`

You can easily host it using Google's app engine or other cloud host that suports Go. That allows you to use it as an overlay or as a chat command. In my case I use Fossabot's custom api feature to display it as a chat message.

The code uses a template to codify the response. You can customize it to say whaterver you want. The available variables are: </br>
`{{user}}` The channel username </br>
`{{uptime}}` The channel's current livestream duration </br>
`{{ML}}` The final value in millilitres </br>
`{{OZ}}` The final value in ounces </br>

# Weather API

There's some code that calls Skorpstuff API and replaces the city name to not show my hometown in chat when pulling the command. Normally it would be easier to pull directly from the original API, however this approach was easier to handle the info.

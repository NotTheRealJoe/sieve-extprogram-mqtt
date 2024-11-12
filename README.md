# sieve-extprogram-mqtt
This program reads an email from STDIN and sends it in a structured format to a configured MQTT channel.

## Publish format
Emails are published to the MQTT channel in JSON structured like this:
```json
{
    "headers": {
        "From": "example@example.com",
        "To": "you@example.com"
    },
    "bodyParts": {
        "text/html": "<html><body><h1>Hello world!</h1></body></html>",
        "text/plain": "Hello world!"
    }
}
```
- `headers` is a mapping from the header name to the corresponding value. If the original message had a header specified
  more than once, all the values are joined together separated by `,`.
- `bodyParts` contains all the content types available in the original message, keyed by the type's MIME name.

## Configuration
Copy config.example.json to `config.json` in the working directory where you will run the code from, or to
`/etc/sieve-extprogram-mqtt/config.json`. Set the properties listed in the example file.

`clientIdPrefix` will be appended with a 6 character random string, generated when sieve-extprogram-mqtt starts up. This
can help distinguish between different runtimes of the program.

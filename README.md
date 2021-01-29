# telegram-to-raindrop

## Development

1. Configure a reverse-proxy to expose our server to the Internet. I usually use
   [localhost.run](https://localhost.run) using this command:

   ```sh
   ssh -R 80:localhost:8080 localhost.run
   ```

   The server will reply with the hostname to be used as the Telegram webhook
   URL.

1. Run the local server using this command:
   
   ```sh
   export TELEGRAM_TOKEN="<token>"
   export TELEGRAM_HOOK_URL="<url returned from the SSH session>"
   export RAINDROP_TOKEN="<your Raindrop.io API token>"
   
   make run
   ```
<p align="center">
  <img alt="Evilginx2 Logo" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/evilginx2-logo-512.png" height="160" />
  <p align="center">
    <img alt="Evilginx2 Title" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/evilginx2-title-black-512.png" height="60" />
  </p>
</p>

# Evilginx 3.0

**Evilginx** is a man-in-the-middle attack framework used for phishing login credentials along with session cookies, which in turn allows to bypass 2-factor authentication protection.

This tool is a successor to [Evilginx](https://github.com/kgretzky/evilginx), released in 2017, which used a custom version of nginx HTTP server to provide man-in-the-middle functionality to act as a proxy between a browser and phished website.
Present version is fully written in GO as a standalone application, which implements its own HTTP and DNS server, making it extremely easy to set up and use.

<p align="center">
  <img alt="Screenshot" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/screen.png" height="320" />
</p>

## Disclaimer

I am very much aware that Evilginx can be used for nefarious purposes. This work is merely a demonstration of what adept attackers can do. It is the defender's responsibility to take such attacks into consideration and find ways to protect their users against this type of phishing attacks. Evilginx should be used only in legitimate penetration testing assignments with written permission from to-be-phished parties.

## Evilginx Mastery Training Course

If you want everything about reverse proxy phishing with **Evilginx** - check out my [Evilginx Mastery](https://academy.breakdev.org/evilginx-mastery) course!

<p align="center">
  <a href="https://academy.breakdev.org/evilginx-mastery"><img alt="Evilginx Mastery" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/evilginx_mastery.jpg" height="320" /></a>
</p>

Learn everything about the latest methods of phishing, using reverse proxying to bypass Multi-Factor Authentication. Learn to think like an attacker, during your red team engagements, and become the master of phishing with Evilginx.

Grab it here:
https://academy.breakdev.org/evilginx-mastery

## Official Gophish integration

If you'd like to use Gophish to send out phishing links compatible with Evilginx, please use the official Gophish integration with Evilginx 3.3.
You can find the custom version here in the forked repository: [Gophish with Evilginx integration](https://github.com/kgretzky/gophish/)

If you want to learn more about how to set it up, please follow the instructions in [this blog post](https://breakdev.org/evilginx-3-3-go-phish/)

## Write-ups

If you want to learn more about reverse proxy phishing, I've published extensive blog posts about **Evilginx** here:

[Evilginx 2.0 - Release](https://breakdev.org/evilginx-2-next-generation-of-phishing-2fa-tokens)

[Evilginx 2.1 - First Update](https://breakdev.org/evilginx-2-1-the-first-post-release-update/)

[Evilginx 2.2 - Jolly Winter Update](https://breakdev.org/evilginx-2-2-jolly-winter-update/)

[Evilginx 2.3 - Phisherman's Dream](https://breakdev.org/evilginx-2-3-phishermans-dream/)

[Evilginx 2.4 - Gone Phishing](https://breakdev.org/evilginx-2-4-gone-phishing/)

[Evilginx 3.0](https://breakdev.org/evilginx-3-0-evilginx-mastery/)

[Evilginx 3.2](https://breakdev.org/evilginx-3-2/)

[Evilginx 3.3](https://breakdev.org/evilginx-3-3-go-phish/)

## Token Modification

Evilginx now supports real-time modification of authentication tokens in HTTP requests. This feature allows you to intercept and modify tokens (such as JWT, Bearer tokens, API keys) before they are sent to the target server.

### Features

- **Token Pattern Matching**: Supports regex patterns to identify different token types (JWT, Bearer, Basic Auth, custom tokens)
- **Multiple Actions**: Replace, add prefix/suffix, or remove tokens entirely
- **Multiple Locations**: Modify tokens in HTTP headers, cookies, query parameters, and request body
- **Rule Management**: Enable/disable rules individually, list all rules, clear rules
- **Preset Configurations**: Quick setup for common token types

### Usage

#### Basic Commands

```bash
# Enable token modification globally
tokens enable

# Disable token modification globally
tokens disable

# Show current status
tokens status

# List all configured rules
tokens list

# Clear all rules
tokens clear
```

#### Quick Setup with Presets

```bash
# Add prefix to JWT tokens
tokens preset jwt "modified_"

# Add prefix to Bearer tokens  
tokens preset bearer "test_"
```

#### Advanced Rule Configuration

```bash
# Add custom rule
# Syntax: tokens add rule <pattern> <action> <value> <description> <location> <target>
tokens add rule "^Bearer\s+(.+)$" "replace" "Bearer modified_$1" "Modify Bearer tokens" "header" "Authorization"

# Enable/disable specific rules
tokens rule enable "Modify Bearer tokens"
tokens rule disable "Modify Bearer tokens"
```

#### Rule Actions

- **replace**: Replace the entire token with a new value (supports regex groups with $1, $2, etc.)
- **prefix**: Add text before the token
- **suffix**: Add text after the token  
- **remove**: Remove the token completely
- **extract**: Extract specific parts using regex groups

#### Token Locations

- **header**: HTTP headers (e.g., Authorization, X-API-Key)
- **cookie**: Cookie values
- **query**: URL query parameters
- **body**: Request body content

### Examples

#### Example 1: Modify JWT Tokens
```bash
# Enable token modification
tokens enable

# Add rule to prefix JWT tokens with "test_"
tokens preset jwt "test_"

# All JWT tokens in Authorization header will be prefixed with "test_"
# Before: Authorization: eyJhbGciOiJIUzI1NiIs...
# After:  Authorization: test_eyJhbGciOiJIUzI1NiIs...
```

#### Example 2: Replace API Keys
```bash
# Replace API keys with a test key
tokens add rule "^(.+)$" "replace" "test-api-key-12345" "Replace API keys" "header" "X-API-Key"
```

#### Example 3: Modify Cookie Session Tokens
```bash
# Add prefix to all cookie values
tokens add rule "^(.+)$" "prefix" "modified_" "Modify session cookies" "cookie" ""
```

### Security Considerations

- Token modification happens before requests reach the target server
- Original tokens are logged and can be captured through Evilginx's session management
- Use responsibly and only for authorized penetration testing
- Rules are applied in the order they were added
- Disabled rules are skipped during processing

## Help

In case you want to learn how to install and use **Evilginx**, please refer to online documentation available at:

https://help.evilginx.com

## Support

I DO NOT offer support for providing or creating phishlets. I will also NOT help you with creation of your own phishlets. Please look for ready-to-use phishlets, provided by other people.

## License

**evilginx2** is made by Kuba Gretzky ([@mrgretzky](https://twitter.com/mrgretzky)) and it's released under BSD-3 license.

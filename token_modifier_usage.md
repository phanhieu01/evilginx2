# Token Modifier Usage Guide

## Overview

The Token Modifier is a new module in Evilginx2 that allows real-time modification of authentication tokens in HTTP requests before they are sent to the target server. This is useful for testing different token scenarios and bypass mechanisms.

## Quick Start

1. Start Evilginx2 normally
2. Enable token modification:
   ```
   tokens enable
   ```
3. Add rules for token modification:
   ```
   tokens preset bearer "test_"
   ```
4. Check status:
   ```
   tokens status
   ```

## Command Reference

### Basic Commands

- `tokens` - Show current status and configured rules
- `tokens enable` - Enable token modification globally
- `tokens disable` - Disable token modification globally
- `tokens status` - Show current status
- `tokens list` - List all configured rules
- `tokens clear` - Remove all rules

### Preset Commands

- `tokens preset jwt <prefix>` - Add prefix to JWT tokens
- `tokens preset bearer <prefix>` - Add prefix to Bearer tokens

### Advanced Commands

```
tokens add rule <pattern> <action> <value> <description> <location> <target>
```

Parameters:
- `pattern`: Regular expression to match tokens
- `action`: What to do (replace, prefix, suffix, remove)
- `value`: Value to use for the action
- `description`: Human readable description
- `location`: Where to find tokens (header, cookie, query, body)
- `target`: Specific target (header name, cookie name, query param name)

### Rule Management

- `tokens rule enable <description>` - Enable a specific rule
- `tokens rule disable <description>` - Disable a specific rule

## Examples

### Example 1: Basic Bearer Token Modification

```bash
# Enable token modification
tokens enable

# Add prefix to Bearer tokens
tokens preset bearer "test_"

# Now all Bearer tokens will be prefixed with "test_"
# Before: Authorization: Bearer abc123
# After:  Authorization: Bearer test_abc123
```

### Example 2: JWT Token Prefix

```bash
# Add prefix to JWT tokens  
tokens preset jwt "modified_"

# JWT tokens will be prefixed
# Before: Authorization: eyJhbGciOiJIUzI1NiIs...
# After:  Authorization: modified_eyJhbGciOiJIUzI1NiIs...
```

### Example 3: API Key Replacement

```bash
# Replace any X-API-Key header with a test key
tokens add rule "^(.+)$" "replace" "test-api-key-12345" "Replace API keys" "header" "X-API-Key"
```

### Example 4: Query Parameter Modification

```bash
# Add prefix to access_token query parameter
tokens add rule "^(.+)$" "prefix" "modified_" "Modify access tokens" "query" "access_token"
```

### Example 5: Cookie Modification

```bash
# Add suffix to all cookie values
tokens add rule "^(.+)$" "suffix" "_test" "Test cookie modification" "cookie" ""
```

## Rule Actions

- **replace**: Replace entire token with new value (supports regex groups $1, $2, etc.)
- **prefix**: Add text before the token
- **suffix**: Add text after the token
- **remove**: Remove the token completely
- **extract**: Extract parts using regex groups

## Token Locations

- **header**: HTTP headers (Authorization, X-API-Key, etc.)
- **cookie**: Cookie values
- **query**: URL query parameters
- **body**: Request body content (future enhancement)

## Regular Expression Examples

```bash
# Match Bearer tokens
^Bearer\s+(.+)$

# Match JWT tokens
^(ey[A-Za-z0-9+/=]+\.[A-Za-z0-9+/=]+\.[A-Za-z0-9+/=]*)$

# Match any non-empty value
^(.+)$

# Match specific API key format
^[A-Za-z0-9]{32}$
```

## Best Practices

1. **Test First**: Always test rules on a test environment first
2. **Use Descriptions**: Give meaningful descriptions to rules for easy management
3. **Rule Order**: Rules are applied in the order they were added
4. **Enable/Disable**: Use rule enable/disable for temporary changes
5. **Clear Rules**: Use `tokens clear` to reset all rules

## Integration with Phishlets

The token modifier works automatically with any phishlet. Once enabled, it will process all HTTP requests going through the proxy and apply the configured rules.

## Troubleshooting

- **Rules not working**: Check if token modification is enabled with `tokens status`
- **Wrong pattern**: Test regex patterns carefully, they must match the entire token
- **Order matters**: Rules are applied in sequence, later rules can override earlier ones
- **Case sensitive**: Regular expressions are case-sensitive by default

## Security Notes

- Token modification happens in real-time during the phishing session
- Original tokens are still captured and logged by Evilginx2
- Modified tokens may cause authentication failures on the target server
- Use responsibly and only for authorized penetration testing

## Example Session

```
: tokens enable
[SUCCESS] Token modification enabled

: tokens preset bearer "test_"
[SUCCESS] Bearer token prefix rule added: test_

: tokens list
Token modification rules:
1. enabled [replace] ^Bearer\s+(.+)$ -> Bearer test_$1
   Description: Add prefix to Bearer tokens

: tokens status
Token modification status: enabled
```

This completes the basic usage guide for the Token Modifier functionality.
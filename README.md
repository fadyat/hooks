## Hooks between services

### Gitlab to:

- **Asana**
    ```text
    Merge request:
    ---

    Updates:
    - <Last Commit> task field
    - <Message> task field
    - Creating a custom fields <Last Commit> and <Message> if they don't exist
    - Works for multiple passed tasks

    How: finding in the last message string like
    - #|refs|https://app.asana.com/#/#/#
    - refs|https://app.asana.com/#/#/#
    
    
    ```

## Configuration

```text
// Put in .env file in the root of the project

// asana access token for editing tasks custom fields
ASANA_API_KEY=<your api key>

// secret tokens that will be used to verify the webhook
GITLAB_SECRET_TOKENS=<list of tokens> 
```

## Build

```bash
make up
```
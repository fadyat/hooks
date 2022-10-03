## Hooks between services

### Gitlab to:

- **Asana**
    ```text
    Merge request:
  
    ---
    Description:
    ---
    Updates:
    - <Last Commit> task field
    - <Message> task field
    - Works for multiple passed tasks

    How it works: finds a line in the last commit message, for example:
    - #|refs|https://app.asana.com/#/#/#
    - refs|https://app.asana.com/#/#/#

    ---
    How to use:
    ---
    - Setup a service (check the example below)
    - Launch service (could use ngrok for local testing)
    - Setup a webhook in Gitlab:
        - URL: https://<your service url>/api/v1/asana/merge
        - Secret Token: <your-secret-token>
        - Trigger: Merge Request Events
    - Add a asana task url to the last commit message
    - Create pull request
    - Check the webhook logs
    - Check the task in Asana
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

## Documentation

```text
When service is running: http://localhost:80/swagger/index.html
In project: /api/docs/swagger.yaml

Library: https://github.com/swaggo/gin-swagger
Recreate: make swag
```
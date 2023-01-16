## Hooks between services

### Workflow

Updates:

- **Last Commit** task field
- **Message** task field
- Works for multiple passed tasks

Pass a task to a commit message with the following syntax:

```text
- #|ref|https://app.asana.com/#/#/#
- ref|https://app.asana.com/#/#/#
- #|ref|<asana-task-id>
- ref|<asana-task-id>
```

### Gitlab integration

#### How to use:

```text
- Set up a service (configuration section)
- Launch service (could use ngrok for local testing)
- Set up a webhook in Gitlab:
    - URL: https://<service>/api/v1/asana/push
    - Secret Token: <your-secret-token>
    - Trigger: Push events
```

> Also may use merge request handler:
> ```text
> Supported merge hook actions: open, update, merge
> - URL: https://<service>/api/v1/asana/merge
> - Trigger: Merge request events
> ```


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
When service is running: https://<service>:80/swagger/index.html
In project: /api/docs/swagger.yaml

Library: https://github.com/swaggo/gin-swagger
Recreate: make swag
```

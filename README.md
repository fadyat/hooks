## Hooks between services

### Workflow

Updates:

- **Last Commit** task field
- **Message** task field
- Works for multiple passed tasks

As a separator can be used one of the following characters: `|`, `:`, `-`, `_`, `=`.

In examples below separator is `|`.

### V1

Pass a task to a commit message with the following syntax:

```text
- #|ref|https://app.asana.com/#/#/#
- ref|https://app.asana.com/#/#/#

Not abstract examples:
- complete|ref|https://app.asana.com/1/2/3
- ref|https://app.asana.com/1/2/3
```

### V2

Pass a asana task_id to a branch name with the following syntax:

```text
- asana|#

Not abstract examples:
- asana|123456789
- asana_123456789
```

### Gitlab integration

#### How to use:

```text
- Set up a service (configuration section)
- Launch service (could use ngrok for local testing)
- Set up a webhook in Gitlab:

    v1. Supports marking in the commit message
    - URL: https://<service>/api/v1/asana/push

    v2. Name branch with asana-<task_id>
    - URL: https://<service>/api/v2/asana/push

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

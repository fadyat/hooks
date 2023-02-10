## Hooks between services

### Workflow

Updates:

- **Last Commit** task field
- **Message** task field
- Works for multiple passed tasks

As a separator can be used one of the following characters: `|`, `:`, `-`, `_`, `=`.

Pass a task to a `commit message` with the following syntax or `name a branch` with the following syntax:

```text
Pattern:
- asana|<task_id>
- ref|<task_id>

Not abstract examples:
- asana|123456789
- asana_123456789
- ref=123456789
```

### Some interesting cases:

- If task passed in branch and in commit message, all tasks will be updated.
- When merging a branch with following pattern, message will be created in the task.

### Gitlab integration

#### How to use:

```text
- Set up a service (configuration section)
- Launch service (could use ngrok for local testing)
- Set up a webhook in Gitlab:
    * URL: https://<service>/api/v1/asana/push
    * Secret Token: <your-secret-token>
    * Trigger: Push events
```

## Configuration

```text
// Put in .env file in the root of the project

// asana access token for editing tasks custom fields
ASANA_API_KEY=<your api key>

// secret tokens that will be used to verify the webhook
GITLAB_SECRET_TOKENS=<list of tokens> 
```

## Documentation

```text
When service is running: https://<service>:80/swagger/index.html
```

Recreate:
```bash
make swag
```

## Hooks between services

### Workflow

Pass a task to a `commit message` with the following syntax or `name a branch` with the following syntax:

As a separator can be used one of the following characters: `|`, `:`, `-`, `_`, `=`.

```text
Pattern:
- asana|<task_id>
- ref|<task_id>

Not abstract examples:
- asana|123456789
- asana_123456789
- ref=123456789
```

Works for multiple passed tasks.

### Some interesting cases:

- If task passed in branch and in commit message, all tasks will be updated.
- When merging a branch with following pattern, message will be created in the task.

### Gitlab integration

#### How to use:

```text
- Set up a service (configuration section)
- Launch service (could use ngrok for local testing)
- Set up a webhook in Gitlab
    * URL: <endpoint>
    * Secret Token: <your-secret-token>
    * Trigger: Push events
```

#### Endpoints:

```text
- POST /api/v1/asana/push
  * last commit info to a task. (to custom field or comment)
- POST /api/v1/gitlab/merge
  * binding a short link of the asana task to the description of MR
- POST /api/v1/asana/merge
  * send a message to the task about the merge of the MR
```

### Configuration

- Put in `.env` file in the root of the project, or set up environment variables.

```text
// asana access token for editing tasks custom fields
ASANA_API_KEY=<your api key>

// secret tokens that will be used to verify the webhook
GITLAB_SECRET_TOKENS=<list of tokens> 

// gitlab api key for updating the merge request description
// sure, that generated token has access to the project!
GITLAB_API_KEY=<your api key>
```

### Feature flags:

```text
// getting task mentions from commit message
IS_COMMIT_MENTIONS_ENABLED=<true|false> // default: false

// make some blured logs when server is started
IS_REPRESENT_SECRETS_ENABLED=<true|false> // default: false
```

### Documentation

- `/swagger/index.html` for swagger docs

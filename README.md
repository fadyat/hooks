## Hooks between services

### When triggered?

[From gitlab docs:](https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html#merge-request-events)

```text
Merge request events are triggered when:

- A new merge request is created.
- An existing merge request is updated, approved (by all required approvers), unapproved, merged, or closed.
- An individual user adds or removes their approval to an existing merge request.
- A commit is added in the source branch.
- All threads are resolved on the merge request.
```

### Which actions are supported?

```text
Supported actions:
- open
- update
- merge
```

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
    - #|ref|https://app.asana.com/#/#/#
    - ref|https://app.asana.com/#/#/#
    - #|ref|<asana-task-id>
    - ref|<asana-task-id>

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

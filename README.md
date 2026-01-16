> Since strava webhook is supported by [statistics-for-strava](https://github.com/robiningelbrecht/statistics-for-strava), and I am changing stack to Kubernetes. This repo's `main` branch is for Kubernetes using. Check `docker` branch if deployed by docker compose.

# Strava Webhook

A Strava webhook listener that automatically updates activity data by triggering [statistics-for-strava](https://github.com/robiningelbrecht/statistics-for-strava) imports on new activities.

## Links

-   DockerHub: [kyxie/strava-webhook general | Docker Hub](https://hub.docker.com/repository/docker/kyxie/strava-webhook/general)
-   GitHub: [GitHub - Kyxie/strava-webhook](https://github.com/Kyxie/strava-webhook)

## How to use

- Chinese toturial: [使用Docker部署Strava数据分析面板 | Kunyang's Blog](https://kyxie.me/zh/blog/bike/strava/#同步)

- Modify your `.env`，for k8s users, save these information to `secret.yaml`

  ```python
  ### Stats for Strava
  # The client id of your Strava app.
  STRAVA_CLIENT_ID=YOUR_CLIENT_ID
  # The client secret of your Strava app.
  STRAVA_CLIENT_SECRET=YOUR_CLIENT_SECRET
  # You will need to obtain this token the first time you launch the app. 
  # Leave this unchanged for now until the app tells you otherwise.
  # Do not use the refresh token displayed on your Strava API settings page, it will not work.
  STRAVA_REFRESH_TOKEN=YOUR_REFRESH_TOKEN_OBTAINED_AFTER_AUTH_FLOW
  # Valid timezones can found under TZ Identifier column here: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List
  TZ=America/Toronto
  
  # The UID and GID to create/own files managed by statistics-for-strava
  PUID=1000
  PGID=1000
  
  ### Webhook
  # You should create a verify token for registering webhook
  STRAVA_VERIFY_TOKEN=YOUR_VERIFY_TOKEN
  # Let strava know where to callback, it's your service's endpoint
  STRAVA_CALLBACK_URL=YOUR_STRAVA_CALLBACK_URL
  ATHELETE_ID=YOUR_ATHELETE_ID
  SUBSCRIPTION_ID=YOUR_SUBSCRIPTION_ID
  APP=strava
  NAMESPACE=strava
  ```

- Create a `verify token` by using command below，then save this token to `.env`

  ```bash
  openssl rand -hex 16
  ```

- This `STRAVA_CALLBACK_URL` is `stats_for_strava's url + /webhook`, for example, if my `stats_for_strava` service is deployed on `https://strava.kyxie.me`, then `STRAVA_CALLBACK_URL=https://strava.kyxie.me/webhook`

- `ATHELETE_ID` can be found at your profile page url, e.g. `https://www.strava.com/athletes/12345`, in this case, 12345 is your `ATHELETE_ID`

- For safety, we should add `SUBSCRIPTION_ID`, which can be found after you register the webhook. If you forget this value, simply run `curl http://localhost:8001/subscription/status`, the value with key is `id` is the subscription id

- Here is the flow for why we need those parameters, you can use a mermaid reader like: https://mermaid.live/ to read it:

  ```mermaid
  sequenceDiagram
      autonumber
      actor User as You
      participant Strava-Webhook
      participant Strava
      participant Stats for Strava
  
      rect rgb(230, 240, 255)
          Note over User, Strava: Phase 1: Registration - One time
          
          User->>GoApp: 1. curl POST /subscription/register
          Note right of User: Trigger register
          
          GoApp->>Strava: 2. Request for estabilishing the subscription<br/>(with callback_url + verify token in request body)
          
          Note right of Strava: Strava receive the request and start to check<br/>your URL is valid or not
          
          Strava->>GoApp: 3. GET /webhook?hub.verify_token=your_varify_token
          
          alt Verify Token match?
              GoApp-->>Strava: 4. Return hub.challenge (200 OK)
              Note right of GoApp: Successful
              Strava-->>GoApp: 5. Return JSON (Contains subscription id)
              Note right of GoApp: Now we have the subscription id <br/> we don't use verify token any more
          else Not match
              GoApp-->>Strava: 403 Forbidden
              Note right of Strava: Failed
          end
      end
  
      rect rgb(255, 245, 230)
          Note over User, Strava: Phase 2: Runtime
          
          User->>Strava: 6. New activity
          Note right of User: Strava cloud generate event
          
          Strava->>GoApp: 7. POST /webhook (JSON Payload)
          Note right of Strava: Now header has no verify token!<br/>Only subscription id in payload <br/> To prevent anyone from calling your /webhook endpoint, <br/> we will check subscription id here
          
          rect rgb(255, 230, 230)
              Note over GoApp: Security Check
              GoApp->>GoApp: Check 1:  owner id
              GoApp->>GoApp: Check 2: subscription id
          end
          
          alt ID Match
              GoApp-->>Strava: 8. Return 200 OK
              GoApp->>PHPApp: 9. kubectl exec
              
              Note right of PHPApp: Import and building Stats for strava
              PHPApp-->>Strava: 10. API GET /activities
              
          else ID Not Match
              GoApp-->>Strava: 200 OK
              Note right of GoApp: Drop
          end
      end
  ```

- If you deployed by cloudflared, `Bot Fight Mode` should be disabled, go to Cloudflare → Security → Bots → Bot Fight Mode, and turn it off

- Register

  ```bash
  curl -X POST http://localhost:8001/subscription/register
  ```

  -   Similarly unregister

      ```bash
      curl -X POST http://localhost:8001/subscription/unregister
      ```

  -   Check Strava webhook subscription status

      ```bash
      curl http://localhost:8001/subscription/status
      ```

- Since the `/webhook` endpoint must be publicly accessible for Strava, and Bot Fight Mode is also disabled, it's best to add rate limiting rules through Cloudflare WAF for protection
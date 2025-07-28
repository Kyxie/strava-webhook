# Strava Webhook

A Strava webhook listener that automatically updates activity data by triggering [statistics-for-strava](https://github.com/robiningelbrecht/statistics-for-strava) imports on new activities.

## Links

-   DockerHub: [kyxie/strava-webhook general | Docker Hub](https://hub.docker.com/repository/docker/kyxie/strava-webhook/general)
-   GitHub: [GitHub - Kyxie/strava-webhook](https://github.com/Kyxie/strava-webhook)

## How to use

- Chinese toturial: [使用Docker部署Strava数据分析面板 | Kunyang's Blog](https://kyxie.me/zh/blog/bike/strava/#同步)

- Modify your `.env`

  ```python
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
  
  # Webhook
  STRAVA_VERIFY_TOKEN=YOUR_VERIFY_TOKEN
  STRAVA_CALLBACK_URL=YOUR_STRAVA_CALLBACK_URL
  ```

- Create a `verify token` by using command below，then save this token to `.env`

  ```bash
  openssl rand -hex 16
  ```

- This `STRAVA_CALLBACK_URL` is `stats_for_strava's url + /webhook`, for example, if my `stats_for_strava` service is deployed on `https://strava.kyxie.me`, then `STRAVA_CALLBACK_URL=https://strava.kyxie.me/webhook`

- Update `docker-compose.yml`

  ```yaml
  services:
    app:
      image: robiningelbrecht/strava-statistics:latest
      container_name: strava
      restart: unless-stopped
      volumes:
        - ./config:/var/www/config/app
        - ./build:/var/www/build
        - ./storage/database:/var/www/storage/database
        - ./storage/files:/var/www/storage/files
      env_file: ./.env
      ports:
        - 8000:8080
      networks:
        - cloudflared
  
    # New
    webhook:
      image: kyxie/strava-webhook:latest
      container_name: strava-webhook
      restart: unless-stopped
      env_file: ./.env
      volumes:
        - /var/run/docker.sock:/var/run/docker.sock
      ports:
        - 8001:8001
      networks:
        - cloudflared
  
  networks:
    cloudflared:
      external: true
  ```

- If you deployed by cloudflared, `Bot Fight Mode` should be disabled, go to Cloudflare → Security → Bots → Bot Fight Mode, and turn it off

- Register

  ```bash
  curl -X POST http://localhost:8001/subscription/register
  ```

  -   Similarily unregister

      ```bash
      curl -X POST http://localhost:8001/subscription/unregister
      ```

  -   Check Strava webhook subscription status

      ```bash
      curl http://localhost:8001/subscription/status
      ```

- Since the `/webhook` endpoint must be publicly accessible for Strava, and Bot Fight Mode is also disabled, it's best to add rate limiting rules through Cloudflare WAF for protection
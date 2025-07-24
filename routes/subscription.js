const express = require('express');
const axios = require('axios');
const router = express.Router();

const {
  STRAVA_CLIENT_ID,
  STRAVA_CLIENT_SECRET,
  STRAVA_CALLBACK_URL,
  STRAVA_VERIFY_TOKEN
} = process.env;

const apiUrl = 'https://www.strava.com/api/v3/push_subscriptions';

router.get('/status', async (req, res) => {
  try {
    const response = await axios.get(apiUrl, {
      params: {
        client_id: STRAVA_CLIENT_ID,
        client_secret: STRAVA_CLIENT_SECRET
      }
    });
    res.json(response.data);
  } catch (err) {
    res.status(500).json({ error: err.response?.data || err.message });
  }
});

router.post('/register', async (req, res) => {
  try {
    const check = await axios.get(apiUrl, {
      params: {
        client_id: STRAVA_CLIENT_ID,
        client_secret: STRAVA_CLIENT_SECRET
      }
    });

    if (check.data.length > 0) {
      const id = check.data[0].id;
      console.log(`[Subscription] Existing webhook found (ID: ${id}), deleting`);
      await axios.delete(`${apiUrl}/${id}`, {
        params: {
          client_id: STRAVA_CLIENT_ID,
          client_secret: STRAVA_CLIENT_SECRET
        }
      });
    }

    const response = await axios.post(apiUrl, null, {
      params: {
        client_id: STRAVA_CLIENT_ID,
        client_secret: STRAVA_CLIENT_SECRET,
        callback_url: STRAVA_CALLBACK_URL,
        verify_token: STRAVA_VERIFY_TOKEN
      }
    });

    console.log('[Subscription] Webhook registered');
    res.json(response.data);
  } catch (err) {
    console.error('[Subscription] Registration failed:', err.response?.data || err.message);
    res.status(500).json({ error: err.response?.data || err.message });
  }
});

router.post('/unregister', async (req, res) => {
  try {
    const check = await axios.get(apiUrl, {
      params: {
        client_id: STRAVA_CLIENT_ID,
        client_secret: STRAVA_CLIENT_SECRET
      }
    });

    if (check.data.length === 0) {
      return res.json({ message: 'No active webhook subscription' });
    }

    const id = check.data[0].id;
    await axios.delete(`${apiUrl}/${id}`, {
      params: {
        client_id: STRAVA_CLIENT_ID,
        client_secret: STRAVA_CLIENT_SECRET
      }
    });

    console.log(`[Subscription] Webhook ID ${id} unregistered`);
    res.json({ message: `Webhook ID ${id} unregistered` });
  } catch (err) {
    res.status(500).json({ error: err.response?.data || err.message });
  }
});

module.exports = router;

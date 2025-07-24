const express = require('express');
const router = express.Router();
const STRAVA_VERIFY_TOKEN = process.env.STRAVA_VERIFY_TOKEN;
const { exec } = require('child_process');

router.get('/', (req, res) => {
  const mode = req.query['hub.mode'];
  const token = req.query['hub.verify_token'];
  const challenge = req.query['hub.challenge'];

  if (mode === 'subscribe' && token === STRAVA_VERIFY_TOKEN) {
    console.log('[Webhook] Verified successfully');
    
    let response = {};
    response["hub.challenge"] = challenge;
    
    return res.status(200).json(response);
  } else {
    console.warn('[Webhook] Verification failed');
    return res.status(403).send('Forbidden');
  }
});

router.post('/', (req, res) => {
  const event = req.body;

  if (event.object_type === 'activity' && event.event_type === 'create') {
    console.log(`[Webhook] Received new activity: ${event.object_id}`);

    const runCommand = (cmd) => {
      return new Promise((resolve, reject) => {
        exec(cmd, (err, stdout, stderr) => {
          if (err) {
            console.error(`[Webhook] Command failed: ${cmd}\n${stderr}`);
            return reject(err);
          }
          console.log(`[Webhook] Command succeeded: ${cmd}\n${stdout}`);
          resolve(stdout);
        });
      });
    };

    (async () => {
      try {
        await runCommand('docker exec strava bin/console app:strava:import-data');
        await runCommand('docker exec strava bin/console app:strava:build-files');
        console.log('[Webhook] ✅ Strava update complete.');
      } catch (err) {
        console.error('[Webhook] ❌ Strava update failed.');
      }
    })();
  }

  res.sendStatus(200);
});

module.exports = router;

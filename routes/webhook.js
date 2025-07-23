const express = require('express');
const router = express.Router();
const VERIFY_TOKEN = process.env.VERIFY_TOKEN;

router.get('/', (req, res) => {
  const mode = req.query['hub.mode'];
  const token = req.query['hub.verify_token'];
  const challenge = req.query['hub.challenge'];

  if (mode === 'subscribe' && token === VERIFY_TOKEN) {
    console.log('[Webhook] Verified successfully');
    return res.json({ hub: { challenge } });
  } else {
    console.warn('[Webhook] Verification failed');
    return res.status(403).send('Forbidden');
  }
});

router.post('/', (req, res) => {
  const event = req.body;

  if (event.object_type === 'activity' && event.event_type === 'create') {
    console.log(`[Webhook] Received new activity: ${event.object_id}`);

    const { exec } = require('child_process');
    exec('./update.sh', (err, stdout, stderr) => {
      if (err) {
        console.error('[Webhook] Error executing update.sh:', stderr);
        return;
      }
      console.log('[Webhook] update.sh output:', stdout);
    });
  }

  res.sendStatus(200);
});

module.exports = router;

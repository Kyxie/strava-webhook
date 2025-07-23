const express = require('express');
const bodyParser = require('body-parser');
const dotenv = require('dotenv');
const { exec } = require('child_process');

// Load env
dotenv.config();

const app = express();
const port = process.env.PORT || 8001;

app.use(bodyParser.json());

// Webhook verification endpoint
app.get('/webhook', (req, res) => {
  const mode = req.query['hub.mode'];
  const token = req.query['hub.verify_token'];
  const challenge = req.query['hub.challenge'];

  if (mode === 'subscribe' && token === process.env.STRAVA_VERIFY_TOKEN) {
    console.log('✅ Webhook verify successful');
    res.status(200).json({ 'hub.challenge': challenge });
  } else {
    console.warn('❌ Webhook verify failed.');
    res.sendStatus(403);
  }
});

// Webhook event handler
app.post('/webhook', (req, res) => {
  const event = req.body;

  if (event.object_type === 'activity' && event.aspect_type === 'create') {
    console.log('🚴 New activity uploaded, running update.sh');

    // Execute the update script
    const scriptPath = process.env.UPDATE_SCRIPT_PATH || './update.sh';
    exec(`sh ${scriptPath}`, (err, stdout, stderr) => {
      if (err) {
        console.error(`❌ Script execution failed: ${stderr}`);
      } else {
        console.log(`✅ Script executed successfully:\n${stdout}`);
      }
    });
  }

  res.sendStatus(200);
});

// Start the server
app.listen(port, () => {
  console.log(`🚀 Webhook running, listening on port ${port}`);
});

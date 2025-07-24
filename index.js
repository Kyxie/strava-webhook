const express = require('express');
const bodyParser = require('body-parser');

const app = express();
const port = process.env.PORT || 8001;

app.use(bodyParser.json());

app.use('/webhook', require('./routes/webhook'));
app.use('/subscription', require('./routes/subscription'));

app.listen(port, () => {
  console.log(`Webhook running, listening on port ${port}`);
});

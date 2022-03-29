## A. Risk Identification

### Identifiy assets (e.g. web application)

The application contains multiple assets that potentially could be vulnerable:

- API/server (go application)
- Webclient (react.js application, served with nginx)
- Graphana
- Prometheus
- Cadvisor
- Elastic Search
- Kibana
- Filebeat
- MySQL database

### Identify threat sources (e.g. SQL injection)

1. SQL injection on web client
1. XSS on web client
1. Getting hands on some of our secrets
1. DDoS on VPS

### Construct risk scenarios (e.g. Attacker performs SQL injection on web application to download sensitive user data)

1. Attacker performs sql injection to download or destroy data from the database
1. Attacker inputs javascript in an input field and accesses data of another user
1. Attacker is able to socially engineer a group member to expose a secret.
1. Attacker uses DDoS crash or halt our server or database.

## B. Risk Analysis

### Determine likelihood

Likelihoods: Certain, Likely, Possible, Unlikely, Rare

### Determine impact

Severities: Insignificant, Negligible, Marginal, Critical, Catastrophic

### Use a Risk Matrix to prioritize risk of scenarios

1. Catastrophic, Unlikely
1. Critical, Possible
1. Critical, Rare
1. Marginal, Possible

### Discuss what are you going to do about each of the scenarios

1. Fix injections and restore backups
1. Say sorry to the user and fix injections
1. Give the exposed group member a security course and change all secrets
1. Restart the server. Put a some DDoS protection in front, like a firewall or a CloudFlare.

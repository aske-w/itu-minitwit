# SLA for ITU-Minitwit

## KPI and metrics

| Metric      | Measurement | Description |
| ----------- | ----------- | ----- |
| Total Requests Attemps | TRA |  Is the total number of API request to the service running on port 8080|
| Failed Requests Attemps | FRA | Is the set of all requests that either does not return, returns an error code or does not return within 120 seconds|
| Availability      | MUT       | Monthley uptime in percentage. Calculated by: ```MUT = (TRA - FRA) / TRA``` |
If the above criterias are not met, the customer can make a formal issue under Github issues.

## SLA details

| Service      |  SLA-target | Measurement |
| ----------- | ----------- | ------ |
|   ITU-Minitwit   |    < 95%       | MUT |

## Responses and Responsibilities

### Customer responsibilities

In case the service does not meet the stated SLA-target or finds the service unsatisfactory, the customer can create a Github issue

### Service provider responsibilities

We will provide support for the customer regarding issues. Below is the expected response time:

| Work days      |  Weekends |
| ----------- | ----------- |
|   24 hours   |    48 hours      |

Furthermore we will notify the customer of scheduled and unscheduled service maintaince.

## Exclusions

The SLA does not apply to unauthoized use of the ITU-Miniwit service, during holidays.

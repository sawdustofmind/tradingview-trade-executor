#/bin/bash

#STAGE_HOST=0.0.0.0
STAGE_HOST=5.223.46.44

USER_URL=http://${STAGE_HOST}:8080
ADMIN_URL=http://${STAGE_HOST}:8081
WEBHOOK_URL=http://${STAGE_HOST}:80

# ADMIN
curl -XPOST ${ADMIN_URL}/v1/admin/login \
-H "Content-Type: application/json" \
-d '{"username":"don", "password":"frengoodfren"}'

ADMIN_TOKEN=f4efcf4b2112608fe89ebd8ef40ddd2648822f5df890a0d635d1f0c8ea7b03c9

curl -XPOST ${ADMIN_URL}/v1/admin/portfolio \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
--data-binary @- << EOF
{
  "name": "DCA_WIP",
  "description": "DCA WIP",
  "year_pnl": "47",
  "avg_delay": "Instant",
  "risk_level": "Medium",
  "strategy_type": "dca",
  "dcalevels": 7,
  "leverage": 2,
  "cycle_investment_percent": "100",
  "holdings": [
    {
      "coin": "XRP",
      "percent": "60"
    },
    {
      "coin": "SOL",
      "percent": "40"
    }
  ]
}
EOF

curl -XPUT ${ADMIN_URL}/v1/admin/portfolio \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
--data-binary @- << EOF
{
  "id": 1,
  "holdings": [
    {
      "coin": "XRP",
      "percent": "100"
    }
  ]
}
EOF

curl -XDELETE ${ADMIN_URL}/v1/admin/portfolio \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"id":1}'

curl -XPUT ${ADMIN_URL}/v1/admin/portfolio \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"id":1, "description":"DCA WIP stat"}'

curl ${ADMIN_URL}/v1/admin/portfolio/list \
-H "Authorization: Bearer${ADMIN_TOKEN}"

curl -XPOST ${ADMIN_URL}/v1/admin/fren \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"name":"Jabami Yumeko", "description": "fond of the risk!", "portfolio_ids":[1]}'

curl -XPUT ${ADMIN_URL}/v1/admin/fren \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"id":1, "description":"fond of the risk! would you play with me?"}'

curl -XDELETE ${ADMIN_URL}/v1/admin/fren \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"id":1}'

curl ${ADMIN_URL}/v1/admin/fren/list \
-H "Authorization: Bearer${ADMIN_TOKEN}"

curl -XPOST ${ADMIN_URL}/v1/admin/generate_invite_token \
-H "Authorization: Bearer${ADMIN_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"count":2}'


curl ${ADMIN_URL}/v1/admin/customer/list \
-H "Authorization: Bearer${ADMIN_TOKEN}"

# USER
curl -XPOST ${USER_URL}/v1/user/register \
-H "Content-Type: application/json" \
-d '{"username":"loolmi", "password":"lol", "invite_token":"45d69e2d0e3d25c71a52ef827476cdc83a7b53b3321f80fca199e4296d140ab8"}'

curl -XPOST ${USER_URL}/v1/user/login \
-H "Content-Type: application/json" \
-d '{"username":"loolmi", "password":"lol"}'

# PUT token from reply above to this env
USER_TOKEN=7c88a8dfc9f5b2b07ea36a915b7bda00fa4980982afa523718ef8e527769cc45

curl ${USER_URL}/v1/user/portfolio/list \
-H "Authorization: Bearer${USER_TOKEN}"

curl -XPUT ${USER_URL}/v1/user/settings \
-H "Authorization: Bearer${USER_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"bybit_api_key":"1nNgYiFMO7daNbGPOA","bybit_api_secret":"d8ULhhYl8YUXJtjIMO63D8Canh2vuglbHFAx"}'

curl ${USER_URL}/v1/user/settings \
-H "Authorization: Bearer${USER_TOKEN}"

curl -XPOST ${USER_URL}/v1/user/subscribe \
-H "Authorization: Bearer${USER_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"portfolio_id":13,"is_test":true,"amount":"1000"}'

curl ${USER_URL}/v1/user/subscriptions/list \
-H "Authorization: Bearer${USER_TOKEN}"

curl -XPOST ${USER_URL}/v1/user/unsubscribe \
-H "Authorization: Bearer${USER_TOKEN}" \
-H "Content-Type: application/json" \
-d '{"subscription_id":8}'

curl ${USER_URL}/v1/user/actions/list \
-H "Authorization: Bearer${USER_TOKEN}"

curl ${USER_URL}/v1/user/trades/list \
-H "Authorization: Bearer${USER_TOKEN}"

### WEBHOOK
curl -XPOST ${WEBHOOK_URL}/v1/tv \
-d '{"strategy_name":"Filter3Heatmap","action":"open","exchange":"BYBIT","symbol":"BTCUSDT"}'

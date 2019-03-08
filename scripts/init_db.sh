
#!/bin/bash
set -e

# get local IP.
localIp=`python ip.py`
curl -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${localIp}:admin_port_placeholer/migrate/v3/migrate/community/0

echo ""

#!/bin/bash

declare arr=()
INSTANCES=""

for arg do
	arr+=${arg}
	INSTANCES="${INSTANCES} ${arg}"
done

LBNDNS="phbk-1290848312.us-east-1.elb.amazonaws.com"
LBNAME="phbk"
echo "aws elb register-instances-with-load-balancer --load-balancer-name ${LBNAME} --instances ${INSTANCES}"
aws elb register-instances-with-load-balancer --load-balancer-name ${LBNAME} --instances ${INSTANCES}

# see the results
# aws elb describe-load-balancers --output json
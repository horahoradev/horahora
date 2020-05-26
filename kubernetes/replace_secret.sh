# Credit goes to https://stackoverflow.com/questions/53852007/kubectl-pod-fails-to-pull-down-an-aws-ecr-image
# Created a script that pulls the token from AWS-ECR
ACCOUNT=$1
REGION=us-west-1
SECRET_NAME=${REGION}-ecr-registry
EMAIL=$2

#
#

TOKEN=`aws ecr --region=$REGION get-authorization-token --output text --query authorizationData[].authorizationToken | base64 -d | cut -d: -f2`

#
#  Create or replace registry secret
#

kubectl delete secret --ignore-not-found $SECRET_NAME
kubectl create secret docker-registry $SECRET_NAME \
 --docker-server=https://${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"

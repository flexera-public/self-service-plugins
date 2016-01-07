#! /bin/bash

ACCT=80263

if [[ ! -f .rsc ]]; then
  echo "Please init rsc using rsc -c .rsc setup"
  exit 1
fi

cat="$1"
if [[ ! -f "$cat" ]]; then
  echo "Looks like $cat does not exist"
  exit 1
fi

echo "Launching..."
cat=`cat "$cat"`
exec_href=`rsc -c .rsc --xh Location ss create manager/projects/$ACCT/executions "source=$cat"`
if [[ "$?" != 0 ]]; then exit $?; fi
echo "Execution href: $exec_href"
re='/([0-9a-z]*)$'
[[ "$exec_href" =~ $re ]] || exit 1
exec_id=${BASH_REMATCH[1]}

echo "Waiting..."
while true; do
  status=`rsc -c .rsc --x1 .status ss show $exec_href`
  echo Status: $status
  if [[ "$status" == failed ]]; then
    # Print error
    echo "Error..."
    rsc -c .rsc --x1 'object:has(.category:val("error")).message:contains("Problem:")' ss index "/api/manager/projects/$ACCT/notifications" "filter[]=execution_id==$exec_id"
    echo
    exit 1
  fi
  sleep 2
done

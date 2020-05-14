# source
# https://askubuntu.com/questions/983063/start-a-screen-session-and-run-a-script-without-attaching-to-it

ServiceName='pikopos'

# get occurences of running apps
NumOfOccurences=$(ps -ef | grep -o "./service_${ServiceName}" | wc -l)

# if not running
if [ $NumOfOccurences -ne 2 ]; then
  echo "Restarting service: ${ServiceName}..."
  screen -X -S "session_${ServiceName}" quit
  screen -dmS "session_${ServiceName}"
  screen -S "session_${ServiceName}" -p 0 -X stuff "./service_${ServiceName}\n"
  echo "Service ${ServiceName} started ✔️"
fi

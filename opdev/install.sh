os=$(go env GOOS)
echo $DLOG_DIALS
dials=$DLOG_DIALS
if [ "$dials" = "" ];then
  if [ $os = "windows" ];then
    echo 'setx DLOG_DIALS "server:abc123456@192.168.61.231:22,server:abc123456@192.168.61.234:22"'
    setx DLOG_DIALS "server:abc123456@192.168.61.231:22,server:abc123456@192.168.61.234:22"
  else
    echo '"export DLOG_DIALS=server:abc123456@192.168.61.231:22,server:abc123456@192.168.61.234:22" >> ~/.bash_profile'
    echo "export DLOG_DIALS=server:abc123456@192.168.61.231:22,server:abc123456@192.168.61.234:22" >> ~/.bash_profile
  fi
fi

echo "go install ./dlog"
go install ./dlog

echo "go install ./devd"
go install ./devd
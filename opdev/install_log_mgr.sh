os=$(go env GOOS)
if [ $os = "windows" ];then
    echo "日志目录环境变量TLOGDIR未设置"
    exit 1
else if [ "$TLOGDIR" = "" ];then
    echo '"export TLOGDIR=/var/log/gzjjyz/" >> ~/.bash_profile'
    echo "export TLOGDIR=/var/log/gzjjyz/" >> ~/.bash_profile
fi

if [ "$BACKUP_LOG_HOST" = "" ];then

fi
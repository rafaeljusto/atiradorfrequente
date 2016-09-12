#!/bin/sh

# PROVIDE: restaf
# REQUIRE: network_ipv6 ip6addrctl
# KEYWORD: shutdown

. /etc/rc.subr

restaf_enable=${restaf_enable-"NO"}

name="restaf"
rcvar=`set_rcvar`
command="/usr/local/atiradorfrequente/rest.af/rest.af"

pidfile=/usr/local/atiradorfrequente/rest.af/rest.af.pid

load_rc_config $name

start_cmd=rest_start

rest_run()
{
    sleep 10
    echo "Starting ${name}."
    cd /usr/local/atiradorfrequente/rest.af && ${command}
}

rest_start()
{
    pid=`check_process ${command}`
    if [ -z $pid ]; then
        rest_run &
    else
        echo "${name} already running? (pid = $pid)"
    fi
}

run_rc_command "$1"

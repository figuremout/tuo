#!/bin/bash
PROJ_PATH=/home/hotel-agency/log
INTERVAL=10 # seconds
logs=(server)

# Var start,end,sum are for get scan efficiency
#start=`date +%s%N`
#sum=0
for log in ${logs[*]}; do
    lines=`grep ^ERROR ${PROJ_PATH}/${log}*.log`
    #n=(`wc -l ./log/${log}*.log`)
    #sum=`expr ${sum} + ${n[0]}`
    # Example grep line:
    # log_route:ERROR 2021-07-25 14:09:41 server.py:41 * msg
    cnt=`echo "$lines" | awk -v INTERVAL=${INTERVAL} '
        BEGIN{
            #year = strftime("%Y")
            now = systime()
            last_interval = now - INTERVAL
            cnt = 0
        }
    
        {
            split($2, year_month_day, "-")
            split($3, hour_minute_sec, ":")
    
            year = year_month_day[1]
            month = year_month_day[2]
            day = year_month_day[3]
            hour = hour_minute_sec[1]
            minute = hour_minute_sec[2]
            second = hour_minute_sec[3]
    
            log_time = mktime(year" "month" "day" "hour" "minute" "second)
            if (log_time > last_interval && log_time <= now)
                cnt += 1
        }
    
        END{
            print cnt
        }
    '`
    echo ${log}_error_cnt:${cnt}
done
#end=`date +%s%N`
#expr \( $end - $start \) / 1000000
#echo $sum

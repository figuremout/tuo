# 编译
```shell
make all VERSION=0.0.1
```

# 统计行数
`find . -type f | xargs cat | wc -l`

# metrics
## cpu
_start:time
_stop:time
_time:time

_measurement:string = cpu
cpu:string = cpu-total/cpu0/cpu1/...
host:string = 058d63a9b310
_field:string = usage_guest/usage_guest_nice/usage_idle/usage_iowait/usage_irq/usage_nice/usage_softirq/usage_steal/usage_system/usage_user
_value:

phisical cpu count: 1
logical cpu count: 2
{"cpu":0,"vendorId":"GenuineIntel","family":"6","model":"85","stepping":7,"physicalId":"0","coreId":"0","cores":1,"modelName":"Intel(R) Xeon(R) Platinum 8269CY CPU @ 2.50GHz","mhz":2500.002,"cacheSize":36608,"flags":["fpu","vme","de","pse","tsc","msr","pae","mce","cx8","apic","sep","mtrr","pge","mca","cmov","pat","pse36","clflush","mmx","fxsr","sse","sse2","ss","ht","syscall","nx","pdpe1gb","rdtscp","lm","constant_tsc","rep_good","nopl","xtopology","nonstop_tsc","pni","pclmulqdq","ssse3","fma","cx16","pcid","sse4_1","sse4_2","x2apic","movbe","popcnt","tsc_deadline_timer","aes","xsave","avx","f16c","rdrand","hypervisor","lahf_lm","abm","3dnowprefetch","invpcid_single","kaiser","fsgsbase","tsc_adjust","bmi1","hle","avx2","smep","bmi2","erms","invpcid","rtm","mpx","avx512f","avx512dq","rdseed","adx","smap","clflushopt","clwb","avx512cd","avx512bw","avx512vl","xsaveopt","xsavec","xgetbv1","xsaves","arat"],"microcode":"0x1"}
{"cpu":1,"vendorId":"GenuineIntel","family":"6","model":"85","stepping":7,"physicalId":"0","coreId":"0","cores":1,"modelName":"Intel(R) Xeon(R) Platinum 8269CY CPU @ 2.50GHz","mhz":2500.002,"cacheSize":36608,"flags":["fpu","vme","de","pse","tsc","msr","pae","mce","cx8","apic","sep","mtrr","pge","mca","cmov","pat","pse36","clflush","mmx","fxsr","sse","sse2","ss","ht","syscall","nx","pdpe1gb","rdtscp","lm","constant_tsc","rep_good","nopl","xtopology","nonstop_tsc","pni","pclmulqdq","ssse3","fma","cx16","pcid","sse4_1","sse4_2","x2apic","movbe","popcnt","tsc_deadline_timer","aes","xsave","avx","f16c","rdrand","hypervisor","lahf_lm","abm","3dnowprefetch","invpcid_single","kaiser","fsgsbase","tsc_adjust","bmi1","hle","avx2","smep","bmi2","erms","invpcid","rtm","mpx","avx512f","avx512dq","rdseed","adx","smap","clflushopt","clwb","avx512cd","avx512bw","avx512vl","xsaveopt","xsavec","xgetbv1","xsaves","arat"],"microcode":"0x1"}
cpu total percent in 1 sec: [0.9950248752112774]
cpu per percents in 1 sec: [0.9999999992724042 0.9900990096014213]
cpu total used seconds since power on: [{"cpu":"cpu-total","user":22383.9,"system":37644.5,"idle":14140982.5,"nice":0.1,"iowait":8048.5,"irq":0.0,"softirq":237.5,"steal":0.0,"guest":0.0,"guestNice":0.0}]
cpu total used seconds sum up: 1.4209297e+07
cpu per used seconds since power on: [{"cpu":"cpu0","user":11445.0,"system":18514.9,"idle":7070984.2,"nice":0.1,"iowait":4044.8,"irq":0.0,"softirq":110.1,"steal":0.0,"guest":0.0,"guestNice":0.0} {"cpu":"cpu1","user":10938.9,"system":19129.6,"idle":7069998.2,"nice":0.1,"iowait":4003.7,"irq":0.0,"softirq":127.5,"steal":0.0,"guest":0.0,"guestNice":0.0}]
cpu per used seconds sum up: 7.1050990200000005e+06
cpu per used seconds sum up: 7.1041979799999995e+06

## load
load avg: {"load1":0.16,"load5":0.05,"load15":0.01}
load misc: {"procsTotal":266,"procsCreated":225981,"procsRunning":2,"procsBlocked":0,"ctxt":24021867054}

## disk
_measurement:string = disk
device:string = vda1
fstype:string = ext4
host:string = 058d63a9b310
mode:string = ro/rw
path:string = /etc/telegraf/telegraf.conf or /etc/hostname
_field:string = free/inodes_free/inodes_total/inodes_used/total/used/used_percent/

_measurement:string = diskio
host:string = 058d63a9b310
name:string = vda/vda1
_field:string = io_time/iops_in_progress/merged_reads/merged_writes/read_bytes/read_time/reads/weighted_io_time/write_bytes/write_time/writes/

disk iocouters: map[vda:{"readCount":10180778,"mergedReadCount":1272,"writeCount":3290544,"mergedWriteCount":3020590,"readBytes":376695337984,"writeBytes":41734811648,"readTime":205839992,"writeTime":12097944,"iopsInProgress":0,"ioTime":5116540,"weightedIO":217945448,"name":"vda","serialNumber":"","label":""} vda1:{"readCount":10180545,"mergedReadCount":1272,"writeCount":3252506,"mergedWriteCount":3020590,"readBytes":376689583104,"writeBytes":41734811648,"readTime":205839376,"writeTime":12096764,"iopsInProgress":0,"ioTime":5115576,"weightedIO":217950444,"name":"vda1","serialNumber":"","label":""}]

disk all partitions: [{"device":"sysfs","mountpoint":"/sys","fstype":"sysfs","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"proc","mountpoint":"/proc","fstype":"proc","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"udev","mountpoint":"/dev","fstype":"devtmpfs","opts":["rw","nosuid","relatime"]} {"device":"devpts","mountpoint":"/dev/pts","fstype":"devpts","opts":["rw","nosuid","noexec","relatime"]} {"device":"tmpfs","mountpoint":"/run","fstype":"tmpfs","opts":["rw","nosuid","noexec","relatime"]} {"device":"/dev/vda1","mountpoint":"/","fstype":"ext4","opts":["rw","relatime"]} {"device":"securityfs","mountpoint":"/sys/kernel/security","fstype":"securityfs","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"tmpfs","mountpoint":"/dev/shm","fstype":"tmpfs","opts":["rw","nosuid","nodev"]} {"device":"tmpfs","mountpoint":"/run/lock","fstype":"tmpfs","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"tmpfs","mountpoint":"/sys/fs/cgroup","fstype":"tmpfs","opts":["ro","nosuid","nodev","noexec"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/systemd","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"pstore","mountpoint":"/sys/fs/pstore","fstype":"pstore","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/pids","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/net_cls,net_prio","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/freezer","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/cpu,cpuacct","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/devices","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/memory","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/cpuset","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/perf_event","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"cgroup","mountpoint":"/sys/fs/cgroup/blkio","fstype":"cgroup","opts":["rw","nosuid","nodev","noexec","relatime"]} {"device":"systemd-1","mountpoint":"/proc/sys/fs/binfmt_misc","fstype":"autofs","opts":["rw","relatime"]} {"device":"mqueue","mountpoint":"/dev/mqueue","fstype":"mqueue","opts":["rw","relatime"]} {"device":"hugetlbfs","mountpoint":"/dev/hugepages","fstype":"hugetlbfs","opts":["rw","relatime"]} {"device":"debugfs","mountpoint":"/sys/kernel/debug","fstype":"debugfs","opts":["rw","relatime"]} {"device":"tracefs","mountpoint":"/sys/kernel/debug/tracing","fstype":"tracefs","opts":["rw","relatime"]} {"device":"tmpfs","mountpoint":"/run/user/0","fstype":"tmpfs","opts":["rw","nosuid","nodev","relatime"]} {"device":"overlay","mountpoint":"/var/lib/docker/overlay2/e63782feb076da8d79c9e0fd417b5964a5ea01bd3b353963be4c5699de13ecdc/merged","fstype":"overlay","opts":["rw","relatime"]} {"device":"nsfs","mountpoint":"/run/docker/netns/4c2e8d9aa419","fstype":"nsfs","opts":["rw","bind"]} {"device":"overlay","mountpoint":"/var/lib/docker/overlay2/cfc284a729b6207030c6743d419012a62e20513cc48a3737e6e824e5c1a31f6a/merged","fstype":"overlay","opts":["rw","relatime"]} {"device":"overlay","mountpoint":"/var/lib/docker/overlay2/594bc08fac9a062347604c324932e79b19d7eb0fc2af0799dfd20411830077db/merged","fstype":"overlay","opts":["rw","relatime"]} {"device":"nsfs","mountpoint":"/run/docker/netns/71c39420a17f","fstype":"nsfs","opts":["rw","bind"]} {"device":"nsfs","mountpoint":"/run/docker/netns/23b3a3961dfc","fstype":"nsfs","opts":["rw","bind"]}]

disk phisical partitions: [{"device":"/dev/vda1","mountpoint":"/","fstype":"ext4","opts":["rw","relatime"]}]

disk usage: {"path":"/tuo","fstype":"ext2/ext3","total":21001531392,"free":5652459520,"used":14258659328,"usedPercent":71.61154245951494,"inodesTotal":1310720,"inodesUsed":400847,"inodesFree":909873,"inodesUsedPercent":30.582199096679684}

## 

_measurement:kernel
host:string = 058d63a9b310
_field:string = boot_time/context_switches/entropy_avail/interrupts/processes_forked/

_measurement:mem
host:string = 058d63a9b310
_field:string = active/available/available_percent/buffered/cached/commit_limit/committed_as/dirty/free/high_free/high_total/huge_page_size/huge_pages_free/huge_pages_total/inactive/low_free/low_total/mapped/page_tables/shared/slab/sreclaimable/sunreclaim/swap_cached/swap_free/swap_total/total/used/used_percent/vmalloc_chunk/vmalloc_total/vmalloc_used/write_back/write_back_tmp/

_measurement:string = processes
host:string = 058d63a9b310
_field:string = blocked/dead/idle/paging/running/sleeping/stopped/total/total_threads/unknown/zombies/

_measurement:string = swap
host:string = 058d63a9b310
_field:string = free/in/out/total/used/used_percent

_measurement:string = system
host:string = 058d63a9b310
_field:string = load1/load15/load5/n_cpus/n_users/uptime/uptime_format



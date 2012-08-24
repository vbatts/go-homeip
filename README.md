
Overview
==================
This is a simple webapp to work around my home ISP and the dynamic lease on DHCP.

Have a simple database storing hostname and the IP from the request, with a PUT.
and return the corresponding IP for hostname, with a GET.

The home machine has a crontab, like:
    0 */1 * * * curl -q -O/dev/null -X PUT http://homeip.myhost.com/ip/$(hostname)


Building
==================

    go build app.go



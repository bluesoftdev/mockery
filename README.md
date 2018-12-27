This project has moved to [gitlab](https://github.com/bluesoftdev/mockery),
if you are not on the gitlab site, please go there instead of using this
github repo.

# Mockery
Mockery is a go library that enables programmers to create mock http
servers for the purpose of testing their integrations in isolation.  It
is particularly good at doing performance testing since one instance can
handle a very large number of tps.  I have tested a basic mockery
handling 100,000 tps without using more than 20% CPU on an 8 core
system.
FROM redis:6.2-alpine
CMD ["redis-server", "--include /usr/local/etc/redis/redis.conf"]
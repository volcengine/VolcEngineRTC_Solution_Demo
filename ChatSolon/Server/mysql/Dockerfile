FROM mysql:5.7

ENV MYSQL_ALLOW_EMPTY_PASSWORD yes

LABEL version = "1.1" description="rtc_demo mysql server"

COPY setup.sh /mysql/setup.sh
COPY rtc_demo.sql /mysql/rtc_demo.sql

CMD ["sh","/mysql/setup.sh"]


FROM python:3

EXPOSE 9784

ADD tplink_smartplug.py /

RUN pip install prometheus-client

CMD [ "python", "./tplink_smartplug.py" ]
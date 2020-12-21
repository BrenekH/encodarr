FROM ubuntu:20.04

WORKDIR /usr/src/app

COPY . .

ENV TZ=America/Boise

RUN apt-get update -qq && DEBIAN_FRONTEND="noninteractive" apt-get install -qq -y mediainfo python3-pip

RUN python3 -m pip install --no-cache-dir -r requirements.txt

EXPOSE 5000
CMD ["python3", "main.py"]

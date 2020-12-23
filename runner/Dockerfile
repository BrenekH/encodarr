FROM ubuntu:20.04

ENV TZ=America/Boise

RUN apt-get update -qq && DEBIAN_FRONTEND="noninteractive" apt-get install -qq -y tzdata ffmpeg python3-pip

WORKDIR /usr/src/app

COPY requirements.txt ./

RUN python3 -m pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python3", "main.py"]

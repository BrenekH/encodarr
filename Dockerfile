FROM ubuntu:20.04

WORKDIR /usr/src/app

COPY . .

RUN apt-get update -qq && apt-get install -qq -y software-properties-common && add-apt-repository -y ppa:stebbins/handbrake-releases && apt-get update -qq && apt-get install -qq -y handbrake-cli mediainfo python3-pip

RUN python3 -m pip install --no-cache-dir -r requirements.txt

EXPOSE 5000
CMD ["python3", "main.py"]

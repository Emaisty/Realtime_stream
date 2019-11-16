# Realtime Stream

Project create to Inno final contest.
Created by @Emaisty, @alexnurin, @fokurly and @kodzhimbo.(t.me)

1)
```
$sudo apt-get install libjpeg-turbo8-dev libavcodec-dev libavformat-dev
$mkdir Screenstreamer
$cd ScreenStreamer
$git clone https://github.com/Emaisty/Realtime_stream.git src/ScreenStreamer
$cd src/ScreenStreamer
(set paths)
$. ./dev.sh
(build)
$go build cmd/mjpeg/mjpeg.go
$./mjpeg
```
When, you can go to ip and port (add /mjpeg), which located in configuration.mjpeg.yml

(Standart "0.0.0.0:8080/mjpeg")
But, if you want to connect from other device, you should change ip of server to ip of your main machine. 
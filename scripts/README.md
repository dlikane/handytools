## Install

```shell
pip install librosa soundfile numpy
```

## Run 

```shell
echo Extract audio from file
ffmpeg -i input.mov -vn -ac 1 -ar 44100 -f wav audio.wav

echo Beat detection
py ../scripts/detect_beats.py audio.wav beats.json

echo Extract images: 
mkdir frames
ffmpeg -i input.mov -qscale:v 2 frames/frame_%05d.jpg

echo Run effect: input: beats.json ./frames/*.jpg output ./processed/*.jpg
go run ../cmd/effect

echo Put it back together: input: ./processed/*.jpg input.mov(audio) output: output.mov
ffmpeg -framerate 25 -i processed/frame_%05d.jpg -i input.mov -map 0:v -map 1:a -c:v libx264 -c:a copy -shortest output.mov


```
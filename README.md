<h1 align="left">Golang based telegram bot that translate indonesia to japanase with voicevox engine</h1>

###

<h3 align="left">How to run?</h3>

###

<p align="left">CPU</p>

###

<p align="left">docker pull voicevox/voicevox_engine:cpu-ubuntu20.04-latest<br>docker run --rm -it -p '127.0.0.1:50021:50021' voicevox/voicevox_engine:cpu-ubuntu20.04-latest</p>

###

<p align="left">GPU</p>

###

<p align="left">docker pull voicevox/voicevox_engine:nvidia-ubuntu20.04-latest<br>docker run --rm --gpus all -p '127.0.0.1:50021:50021' voicevox/voicevox_engine:nvidia-ubuntu20.04-latest</p>

###

<p align="left">Config</p>

###

<p align="left">TELEGRAM_BOT_API_KEY =</p>

###

<p align="left">Run command</p>

###

<p align="left">go run main.go</p>

###

## Example running bot

<div align="center">
    <img src="https://github.com/eldhoral/translator-telegram-bot/blob/main/Screen_Recording_20230915_171924_Telegram.gif" width="200px"</img> 
</div>


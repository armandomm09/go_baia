const { Client, LocalAuth } = require('whatsapp-web.js');
const fs = require('fs');
const qrcode = require('qrcode-terminal');
const axios = require('axios');
const FormData = require('form-data');



  var mediaContador = 0

  async function sendGPTAudio(filePath) {
    if (!fs.existsSync(filePath)) {
        console.log("Archivo no encontrado:", filePath);
        return;
    }

    const formData = new FormData();
    formData.append('audio', fs.createReadStream(filePath));

    try {
        const response = await axios.post("http://localhost:8888/baia/askGPT/audio/", formData, {
            headers: {
                ...formData.getHeaders()
            }
        });
        if (response.status !== 201) {
            console.log("Error fetching:", response.statusText);
            console.log("Response data:", response.data);
            return "Error fetchingg"
        } else {
            console.log("File uploaded successfully");
            console.log("Response data:", response.data);
            return response.data["Answer"]
        }
    } catch (error) {
        return "catched Error" + error.toString()
        console.log("Error uploading file:", error.message);
          }
}

  async function sendGPTMessage(mensaje) {
    const response = await fetch("http://localhost:8888/baia/askGPT/text/question", {
        method: 'POST',
        body: JSON.stringify({ // Convert data to JSON string
            "question": mensaje
        }),
        headers: { // Set Content-Type header for JSON data
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        console.log("Error asking gpt: " + response.status);
        return "Hubo un error"
    } else {
        const responseData = await response.json(); // Parse JSON response
        console.log(responseData.toString());
        return responseData["Answer"] // Print the parsed JSON data
    }
}

const client = new Client({
    authStrategy: new LocalAuth(),
    webVersion: "2.2412.54",
    webVersionCache: {
        type: "remote",
        remotePath:
            "https://raw.githubusercontent.com/wppconnect-team/wa-version/main/html/2.2412.54.html",
    },
});

client.on('message', async message => {
    console.log(message.from)
    if(message.from === "5212721976963@c.us"){
              console.log(message.body)
        console.log(message.from)
      if(message.hasMedia){
        console.log(message.hasMedia)
        const msgmedia =  await message.downloadMedia()
        console.log(msgmedia.filename)
        const mediaLocalPath = "../audios/base64EncodedMedia/" + "audioNum" + mediaContador.toString()
        fs.writeFile(
            mediaLocalPath,
            msgmedia.data,
            "base64",
            function (err) {
              if (err) {
                console.log(err);
              }
            }
          );
          const oggAudioPath = `../audios/mediaInOgg/audio${mediaContador}.ogg`
          fs.writeFileSync(oggAudioPath, Buffer.from(msgmedia.data.replace(`data:audio/ogg; codecs=opus;base64,`, ''), 'base64'));
        var answer = await sendGPTAudio(oggAudioPath)
        message.reply(answer)
        mediaContador++

    } else {
        message.reply(await sendGPTMessage(message.body))
        }
        }
    if(message.body === "!ping"){
        message.reply("pong")
    }

});

client.on('ready', () => {
    console.log('Client is ready!');

});

client.on('qr', qr => {
    qrcode.generate(qr, { small: true });
});

client.initialize();

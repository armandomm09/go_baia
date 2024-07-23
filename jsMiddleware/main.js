// archivo principal
const { Client, LocalAuth, MessageMedia } = require('whatsapp-web.js');
const fs = require('fs');
const qrcode = require('qrcode-terminal');
// const qrcode = require('qrcode');
const axios = require('axios');
const FormData = require('form-data');

var mediaContador = 0;

const client = new Client({
    authStrategy: new LocalAuth(),
    webVersion: "2.2412.54",
    webVersionCache: {
        type: "remote",
        remotePath:
            "https://raw.githubusercontent.com/wppconnect-team/wa-version/main/html/2.2412.54.html",
    },
});

async function sendGPTAudio(filePath, senderID) {
    if (!fs.existsSync(filePath)) {
        console.log("Archivo no encontrado:", filePath);
        return;
    }

    const formData = new FormData();
    formData.append('audio', fs.createReadStream(filePath));
    formData.append('senderID', senderID); // Agrega senderID directamente al formData

    try {
        const response = await axios.post("http://localhost:8888/baia/askGPT/audio/", formData, {
            headers: {
                ...formData.getHeaders()
            }
        });

        if (response.status !== 201) {
            console.log("Error fetching:", response.statusText);
            console.log("Response data:", response.data);
            return "Error fetching";
        } else {
            console.log("File uploaded successfully");
            console.log("Response data:", response.data);
            return response.data["Answer"];
        }
    } catch (error) {
        console.log("Error uploading file:", error.message);
        return "catched Error: " + error.toString();
    }
}


async function sendGPTMessage(mensaje, senderID) {
    const response = await fetch("http://localhost:8888/baia/askGPT/text/question", {
        method: 'POST',
        body: JSON.stringify({ // Convert data to JSON string
            "question": mensaje,
            "senderID": senderID
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
        return responseData["messages"] // Print the parsed JSON data
    }
}

client.on('message', async message => {
    console.log(message.from)
    if (message.from === "5212223201384@c.us") {
        console.log(message.body)
        console.log(message.from)
        if (message.hasMedia) {
            console.log(message.hasMedia)
            const msgmedia = await message.downloadMedia()
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
            var answer = await sendGPTAudio(oggAudioPath, message.from)
            message.reply(answer)
            mediaContador++

        } else {
            const responseMessages = await sendGPTMessage(message.body, message.from)
            const messageFrom = message.from
            console.log(JSON.stringify(responseMessages))
            for (let i = 0; i < responseMessages.length; i++) {
                if (responseMessages[i]["isImage"]) {
                    let imageToSend = await MessageMedia.fromUrl(responseMessages[i]["response"])
                    await client.sendMessage(messageFrom, imageToSend)
                } else {
                    console.log(JSON.stringify(responseMessages[i]))
                    console.log("Message send")
                    await client.sendMessage(messageFrom, responseMessages[i]['response'])
                }
            }
        }
    }
    if (message.body === "!ping") {
        message.reply("pong")
    }

});

client.on('ready', async () => {
    console.log('Client is ready!');
    await client.sendMessage('5212223201384@c.us', "hola")

});

client.on('qr', qr => {
    console.log("Generating qr...")
    qrcode.generate(qr, { small: true });
});

// client.on('qr', qr => {
//     // qrcode.generate(qr, { small: true });
//     console.log('QR RECEIVED', qr);

//     // Genera el código QR y guárdalo como PNG
//     qrcode.toFile('{clientName}.png', qr, {
//         color: {
//             dark: '#000000',  // Color de los puntos
//             light: '#FFFFFF'  // Color de fondo
//         }
//     }, function (err) {
//         if (err) throw err;
//         console.log('QR code saved as qr-code.png');
//     });
// });

client.initialize();
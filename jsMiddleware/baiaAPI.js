// chatgpt.js
const fs = require('fs');
const axios = require('axios');
const FormData = require('form-data');

export default async function sendGPTAudio(filePath) {
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
            return "Error fetching"
        } else {
            console.log("File uploaded successfully");
            console.log("Response data:", response.data);
            return response.data["Answer"]
        }
    } catch (error) {
        console.log("Error uploading file:", error.message);
        return "catched Error" + error.toString()
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

module.exports = {
    sendGPTAudio,
    sendGPTMessage
};



// client.on('message', async message => {
//     console.log(message.from)
//     if(message.from === "5212721976963@c.us"){
//               console.log(message.body)
//         console.log(message.from)
//       if(message.hasMedia){
//         console.log(message.hasMedia)
//         const msgmedia =  await message.downloadMedia()
//         console.log(msgmedia.filename)
//         const mediaLocalPath = "../audios/base64EncodedMedia/" + "audioNum" + mediaContador.toString()
//         fs.writeFile(
//             mediaLocalPath,
//             msgmedia.data,
//             "base64",
//             function (err) {
//               if (err) {
//                 console.log(err);
//               }
//             }
//           );
//           const oggAudioPath = `../audios/mediaInOgg/audio${mediaContador}.ogg`
//           fs.writeFileSync(oggAudioPath, Buffer.from(msgmedia.data.replace(`data:audio/ogg; codecs=opus;base64,`, ''), 'base64'));
//         var answer = await sendGPTAudio(oggAudioPath)
//         message.reply(answer)
//         mediaContador++

//     } else {
//         message.reply(await sendGPTMessage(message.body))
//         }
//         }
//     if(message.body === "!ping"){
//         message.reply("pong")
//     }

// });

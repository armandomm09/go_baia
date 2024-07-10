// archivo principal
const { Client, LocalAuth, Location, MessageMedia } = require('whatsapp-web.js');
const fs = require('fs');
const qrcode = require('qrcode-terminal');
// const { sendGPTAudio, sendGPTMessage } = require('./chatgpt');

var mediaContador = 0;

const client = new Client({
    authStrategy: new LocalAuth()
});

client.on('message', async message => {
    
    
    console.log(message.from);
    console.log(message.body);
    console.log(message.type);
    // if((message.from === '5212222150794@c.us' || message.from === '5212223201384@c.us') && message.hasMedia){
    //     const newMedia = await message.downloadMedia()
    //     client.sendMessage('5212222150794@c.us', newMedia)

    // }
    if((message.from === '5212721873974@c.us' ) && message.type === 'location'){
        console.log("LOCATIONNN *********** \n")

        var newLatitude = message.location.latitude
        var newLongitude = message.location.longitude
        console.log(newLatitude)
        console.log(newLongitude)
        var newLocation = new Location(newLatitude, newLongitude);
        // message.reply(greenwichLocation)
        client.sendMessage('5212721976963@c.us', newLocation)

    }
    if(message.body === "!ping"){
        message.reply('pong');
    }
});

client.on('ready', async () => {
    console.log('Client is ready!');
    // client.sendMessage(`@5212221882222`, )
    var newLocation = new Location(40.7128, 74.0060);

    client.sendMessage('5212721976963@c.us', newLocation)
    
});

client.on('qr', qr => {
    qrcode.generate(qr, { small: true });
});

client.initialize();

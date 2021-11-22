let hls = new Hls();
let alreadyOpened = false;

function turnSong(path) {
    var audio = document.querySelector('audio');

    if (Hls.isSupported()) {
        hls.loadSource(path);
        hls.attachMedia(audio);
    }

    if (!alreadyOpened) {
        plyr.setup(audio);
        alreadyOpened = true;
    }
    plyr.source = audio;
}

window.onload = function (e) {
   console.log("did");
    $.ajax({
        url: "http://localhost:8081/api/v1/getallsongs",
    })
        .done(function( data ) {
            for (let i = 0; i < data.songs.length; i++) {
                console.log(data.songs[i]);
                let el = `<button onclick="turnSong('${data.songs[i].path}')">${data.songs[i].name}</button>`;
                $("#buttons").append(el);
            }
        });
}

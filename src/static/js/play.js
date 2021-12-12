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
    $.ajax({
        url: "http://localhost:8081/api/v1/getallsongs",
    })
        .done(function( data ) {
            for (let i = 0; i < data.songs.length; i++) {
                let el = `<p class='songStyle' onclick="turnSong('${data.songs[i].path}')">${data.songs[i].band}: ${data.songs[i].name}</p>`;
                $("#buttons").append(el);
            }
        });
}

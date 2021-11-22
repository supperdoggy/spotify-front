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
        url: "http://localhost:8080/allsongs",
    })
        .done(function( data ) {
            console.log(data);
        });
}

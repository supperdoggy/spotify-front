
// now it works omg
function send() {
    let name = $("#name")[0].value;
    let album = $("#album")[0].value;
    let band = $("#band")[0].value;
    let release = $("#release")[0].value;
    let file = $("#file")[0].files[0];

    if (name === "" || album === "" || band === "" || release === "" || file === undefined) {
        alert("you need to fill all the fields!");
        return
    }
    getBase64(file).then( function (data) {
        let req = JSON.stringify({
            "name": name,
            "album": album,
            "band": band,
            "release_date": release,
            'song_data': data
        })
        $.ajax({
            url: "http://localhost:8081/api/v1/song",
            method: "post",
            processData: false,
            data: req,
            crossDomain: true,
            dataType: 'json',
            contentType: "application/json; charset=utf-8",
        }).done(function( data ) {
            console.log(data);
        });
    });

}

function toTimestamp(strDate){
    var datum = Date.parse(strDate);
    return datum/1000;
}

function getBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}

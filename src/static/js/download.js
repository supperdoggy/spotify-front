
// now it works omg
function send() {
    console.log("click");
    $("#file")[0].files[0].text().then(function(text) {
        console.log(text);
        let data = JSON.stringify({'song_data': text})
        $.ajax({
            url: "http://localhost:8081/api/v1/song",
            method: "post",
            processData: false,
            data: data,
            crossDomain: true,
            dataType: 'json',
            contentType: "application/json; charset=utf-8",
        })
            .done(function( data ) {
                console.log(data);
            });
    });

}
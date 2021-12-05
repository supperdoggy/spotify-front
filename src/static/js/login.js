
let login = true;

function activateReg() {
    $("#sign_up").attr("class", "active");
    $("#sign_in").removeClass('active');
    login = false;

    // add fields

}

function activateLog() {
    $("#sign_in").attr("class", "active");
    $("#sign_up").removeClass('active');
    login = true;

    // remove fields
}

function send() {
    let email = $("#email")[0].value;
    let password = $("#password")[0].value
    let obj = {
        "email": email,
        "password": password,
        "login": login
    }

    $.ajax({
        url: "http://localhost:8081/api/v1/auth",
        method: "post",
        type: "post",
        dataType: "json",
        contentType: "application/json; charset=utf-8",
        data: JSON.stringify(obj),
    }).done(function (data) {
        console.log(data);
    });

}


let login = true;

function activateReg() {
    // add fields
    $("#first_name").attr("style", "visibility:visible");
    $("#last_name").attr("style", "visibility:visible");

    $("#sign_up").attr("class", "active");
    $("#sign_in").removeClass('active');
    login = false;
}

function activateLog() {
    // remove fields
    $("#first_name").attr("style", "visibility:hidden");
    $("#last_name").attr("style", "visibility:hidden");

    $("#sign_in").attr("class", "active");
    $("#sign_up").removeClass('active');
    login = true;
}

function send() {
    eraseCookie("t");
    let email = $("#email")[0].value;
    let password = $("#password")[0].value;
    let first_name = $("#first_name")[0].value;
    let last_name = $("#last_name")[0].value;

    let obj = {
        "email": email,
        "password": password,
        "login": login,
        "first_name": first_name,
        "last_name": last_name
    };

    $.ajax({
        url: "http://localhost:8081/api/v1/auth",
        method: "post",
        type: "post",
        dataType: "json",
        contentType: "application/json; charset=utf-8",
        data: JSON.stringify(obj),
    }).done(function (data) {
        createCookie("t", data.token, 31);
        window.location.href = "http://localhost:8081/";
    });

}

function createCookie(name, value, days) {
    var expires;

    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toGMTString();
    } else {
        expires = "";
    }
    document.cookie = encodeURIComponent(name) + "=" + encodeURIComponent(value) + expires + "; path=/";
}

function readCookie(name) {
    var nameEQ = encodeURIComponent(name) + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === ' ')
            c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0)
            return decodeURIComponent(c.substring(nameEQ.length, c.length));
    }
    return null;
}

function eraseCookie(name) {
    createCookie(name, "", -1);
}

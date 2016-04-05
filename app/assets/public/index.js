(function() {
    var snackbar = document.querySelector("#snackbar");
    var showErr = function(err) {
        snackbar.style.backgroundColor = '#F44336';
        snackbar.MaterialSnackbar.showSnackbar({message: err});
    };
    var showMsg = function(msg) {
        snackbar.style.backgroundColor = '#323232';
        snackbar.MaterialSnackbar.showSnackbar({message: msg});
    };

    var ajax = function(url, cb) {
        var xh = new XMLHttpRequest();
        xh.onreadystatechange = function() {
            if (xh.readyState == XMLHttpRequest.DONE) {
                switch (xh.status) {
                case 200:
                    cb(xh.responseText);
                    break;
                default:
                    showErr(xh.responseText)
                    break;
                }
            }
        };
        xh.open("GET", url, true);
        xh.send();
    }

    var result = document.querySelector("#result");
    var search = document.querySelector("#search");
    search.addEventListener("keyup", function() {
        if (search.value === "") {
            result.innerHTML = "";
            return;
        }
        ajax("/foods/" + search.value, function(data) {
            result.innerHTML = data;
        });
    });
})()

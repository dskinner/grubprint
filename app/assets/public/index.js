(function() {
    var snackbar = document.querySelector("#snackbar");
    var showErr = function(err) {
        snackbar.style.backgroundColor = "#F44336";
        snackbar.MaterialSnackbar.showSnackbar({message: err});
    };
    var showMsg = function(msg) {
        snackbar.style.backgroundColor = "#323232";
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

    var detail = document.querySelector("#detail");

    var details = function(id) {
        ajax("/nutrients/"+id, function(data) {
            detail.innerHTML = data;
        })
    };

    var result = document.querySelector("#result");
    var search = document.querySelector("#search");
    search.addEventListener("keyup", function() {
        result.style.animation = "fadeout 250ms forwards";
        if (search.value === "") {
            result.innerHTML = "";
            return;
        }
        ajax("/foods/" + search.value, function(data) {
            result.innerHTML = data;
            result.style.animation = "fadein 250ms";
            var l = document.querySelectorAll("ul li");
            for (var i = 0; i < l.length; i++) {
                l[i].addEventListener("click", function() {
                    details(this.getAttribute("data-id"));
                });
            };
        });
    });
})()

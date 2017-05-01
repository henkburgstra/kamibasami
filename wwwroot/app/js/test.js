var Test = function() {
    this.url = null;
    this.path = null;

    this.submit = function(url, path) {
        var alerts = new alertContainer('alerts');
        alerts.clear();
        var data = {
            url: url,
            path: path
        };
        var xmlhttp = new XMLHttpRequest();

        xmlhttp.onreadystatechange = function () {
            if (xmlhttp.readyState == 4) {
                if (xmlhttp.status == 200) {
                    var result = JSON.parse(xmlhttp.responseText);
                    if (result.status != 200) {
                        alerts.add('danger', '[Test.submit] Error submitting URL. Error: ' + 
                            result.status + ', ' +
                            result.error);                    
                    }
                } else {
                    alerts.add('danger', '[Test.submit] Error submitting URL. HTTP-status: ' + xmlhttp.status);
                }
            }
        };

        xmlhttp.open("POST", '/api/webpage');
        xmlhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        xmlhttp.send(JSON.stringify(data));
        
    }

    this.init = function() {
        var self = this;
        this.url = document.getElementById('url');
        this.path = document.getElementById('path');
        var ok = document.getElementById('btn-ok');
        ok.addEventListener('click', function(evt){
            self.submit(self.url.value, self.path.value);
        });
    };
};

document.addEventListener('DOMContentLoaded', function () {
    var test = new Test();
    test.init();
});